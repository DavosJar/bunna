//! Punto de entrada del Servicio de Monitoreo.
//!
//! 1. Carga variables de entorno desde `.env`.
//! 2. Inicia el sistema de tracing (formato JSON, filtro vía `RUST_LOG`).
//! 3. Lee la configuración.
//! 4. Arranca el insertador por lotes (BatchInserter).
//! 5. Arranca el consumidor de Kafka.
//! 6. Expone un servidor HTTP con `/metrics` y `/health`.
//! 7. Corre todo con `tokio::select!`.

mod config;
mod consumidor;
mod sink;
mod types;

use std::sync::Arc;

use axum::{
    Router,
    routing::get,
    response::IntoResponse,
};
use futures::StreamExt;
use prometheus::{Encoder, TextEncoder};
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::ClientConfig;
use tracing::{error, info, warn};
use std::net::SocketAddr;

use config::Config;
use consumidor::ManejadorKafka;
use sink::{BatchInserter, ClienteClickHouse};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // 1. Cargar .env (silencioso si no existe)
    dotenvy::dotenv().ok();

    // 2. Iniciar tracing
    let config = Config::desde_env();

    tracing_subscriber::fmt()
        .json()
        .with_env_filter(&config.RUST_LOG)
        .with_target(true)
        .with_current_span(true)
        .init();

    info!("Iniciando Servicio de Monitoreo v{}", env!("CARGO_PKG_VERSION"));

    // 3. Leer configuración (ya se hizo para RUST_LOG, volvemos a tenerla)
    //    (Config::desde_env() ya se invocó arriba; reusamos la misma)
    info!("Configuración cargada desde variables de entorno");

    // 4. Crear cliente ClickHouse y el insertador por lotes
    let cliente_clickhouse = ClienteClickHouse::nuevo(
        config.CLICKHOUSE_URL.clone(),
        config.CLICKHOUSE_USUARIO.clone(),
        config.CLICKHOUSE_CONTRASENA.clone(),
        config.CLICKHOUSE_BD.clone(),
    );
    let insertador = Arc::new(BatchInserter::nuevo(cliente_clickhouse));
    let manejador = Arc::new(ManejadorKafka::nuevo(insertador.clone()));

    // 5. Crear consumidor de Kafka
    let consumidor: StreamConsumer = ClientConfig::new()
        .set("bootstrap.servers", &config.KAFKA_BROKERS)
        .set("group.id", &config.GRUPO_CONSUMIDOR)
        .set("auto.offset.reset", "earliest")
        .set("enable.auto.commit", "true")
        .set("session.timeout.ms", "6000")
        .create()
        .expect("Error al crear consumidor de Kafka");

    let topics: Vec<&str> = vec![
        &config.KAFKA_TOPIC_TELEMETRIA,
        &config.KAFKA_TOPIC_HARDWARE_METRICAS,
        &config.KAFKA_TOPIC_HARDWARE_ALERTAS,
    ];
    consumidor
        .subscribe(&topics)
        .expect("Error al suscribirse a los topics de Kafka");

    info!(
        "Consumidor Kafka suscrito a: {:?} (grupo: {})",
        topics, config.GRUPO_CONSUMIDOR
    );

    // 6. Servidor HTTP con métricas y health check
    let router = Router::new()
        .route("/health", get(health_handler))
        .route("/metrics", get(metrics_handler));

    let addr = SocketAddr::from(([0, 0, 0, 0], config.PUERTO_METRICAS));
    info!("Servidor HTTP escuchando en {}", addr);

    let listener = tokio::net::TcpListener::bind(addr)
        .await
        .expect("Error al bindear puerto de métricas");

    // 7. Ejecutar consumidor y servidor concurrentemente
    let manejador_arc = manejador.clone();
    let consumidor_arc = consumidor;

    let tarea_consumidor = tokio::spawn(async move {
        info!("Bucle de consumidor Kafka iniciado");
        let mut stream = consumidor_arc.stream();
        loop {
            match stream.next().await {
                Some(Ok(mensaje)) => {
                    if let Err(e) = manejador_arc.manejar_mensaje(&mensaje).await {
                        error!("Error procesando mensaje: {}", e);
                    }
                }
                Some(Err(e)) => {
                    error!("Error de Kafka: {:?}", e);
                }
                None => {
                    warn!("Stream de Kafka finalizó inesperadamente");
                    break;
                }
            }
        }
    });

    let tarea_servidor = tokio::spawn(async move {
        axum::serve(listener, router)
            .await
            .expect("Error en servidor HTTP");
    });

    // Esperar a que alguna de las dos tareas termine
    tokio::select! {
        _ = tarea_consumidor => {
            info!("Tarea de consumidor finalizó");
        }
        _ = tarea_servidor => {
            info!("Servidor HTTP finalizó");
        }
    }

    info!("Servicio de Monitoreo detenido");
    Ok(())
}

// ---------------------------------------------------------------------------
// Handlers HTTP
// ---------------------------------------------------------------------------

/// Handler de health check.
async fn health_handler() -> impl IntoResponse {
    axum::response::Json(serde_json::json!({
        "status": "ok",
        "service": "servicio-monitoreo",
        "version": env!("CARGO_PKG_VERSION"),
    }))
}

/// Handler de métricas Prometheus.
async fn metrics_handler() -> impl IntoResponse {
    let encoder = TextEncoder::new();
    let mut buffer = Vec::new();
    match encoder.encode(&prometheus::gather(), &mut buffer) {
        Ok(()) => axum::response::Response::builder()
            .header("Content-Type", "text/plain; charset=utf-8")
            .body(axum::body::Body::from(buffer))
            .unwrap(),
        Err(e) => axum::response::Response::builder()
            .status(500)
            .body(axum::body::Body::from(format!("Error al codificar métricas: {}", e)))
            .unwrap(),
    }
}
