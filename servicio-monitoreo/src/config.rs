//! Configuración leída desde variables de entorno.
//!
//! # Variables de entorno esperadas
//!
//! | Variable                  | Obligatoria | Descripción                               |
//! |---------------------------|-------------|-------------------------------------------|
//! | `KAFKA_BROKERS`           | Sí          | Lista de brokers de Kafka                 |
//! | `KAFKA_TOPIC_TELEMETRIA`  | Sí          | Topic para mensajes de telemetría         |
//! | `KAFKA_TOPIC_HARDWARE_METRICAS` | Sí   | Topic para métricas de hardware           |
//! | `KAFKA_TOPIC_HARDWARE_ALERTAS`  | Sí   | Topic para alertas de hardware            |
//! | `GRUPO_CONSUMIDOR`        | Sí          | Grupo de consumidor de Kafka              |
//! | `CLICKHOUSE_URL`          | Sí          | URL de ClickHouse (ej. http://localhost:8123) |
//! | `CLICKHOUSE_USUARIO`      | Sí          | Usuario de ClickHouse                     |
//! | `CLICKHOUSE_CONTRASENA`   | Sí          | Contraseña de ClickHouse                  |
//! | `CLICKHOUSE_BD`           | Sí          | Base de datos de ClickHouse               |
//! | `PUERTO_METRICAS`         | No          | Puerto para servidor HTTP (defecto: 9100) |
//! | `RUST_LOG`                | No          | Nivel de traza (defecto: info)            |

use std::env;

/// Configuración central del servicio de monitoreo.
#[derive(Debug, Clone)]
#[allow(non_snake_case)]
pub struct Config {
    // Kafka
    pub KAFKA_BROKERS: String,
    pub KAFKA_TOPIC_TELEMETRIA: String,
    pub KAFKA_TOPIC_HARDWARE_METRICAS: String,
    pub KAFKA_TOPIC_HARDWARE_ALERTAS: String,
    pub GRUPO_CONSUMIDOR: String,

    // ClickHouse
    pub CLICKHOUSE_URL: String,
    pub CLICKHOUSE_USUARIO: String,
    pub CLICKHOUSE_CONTRASENA: String,
    pub CLICKHOUSE_BD: String,

    // Servidor de métricas
    pub PUERTO_METRICAS: u16,

    // Tracing
    pub RUST_LOG: String,
}

impl Config {
    /// Lee la configuración desde variables de entorno.
    ///
    /// # Panics
    ///
    /// Paniccea si una variable obligatoria no está presente.
    pub fn desde_env() -> Self {
        Self {
            KAFKA_BROKERS: env::var("KAFKA_BROKERS")
                .expect("KAFKA_BROKERS debe estar definida"),
            KAFKA_TOPIC_TELEMETRIA: env::var("KAFKA_TOPIC_TELEMETRIA")
                .expect("KAFKA_TOPIC_TELEMETRIA debe estar definida"),
            KAFKA_TOPIC_HARDWARE_METRICAS: env::var("KAFKA_TOPIC_HARDWARE_METRICAS")
                .expect("KAFKA_TOPIC_HARDWARE_METRICAS debe estar definida"),
            KAFKA_TOPIC_HARDWARE_ALERTAS: env::var("KAFKA_TOPIC_HARDWARE_ALERTAS")
                .expect("KAFKA_TOPIC_HARDWARE_ALERTAS debe estar definida"),
            GRUPO_CONSUMIDOR: env::var("GRUPO_CONSUMIDOR")
                .expect("GRUPO_CONSUMIDOR debe estar definida"),

            CLICKHOUSE_URL: env::var("CLICKHOUSE_URL")
                .expect("CLICKHOUSE_URL debe estar definida"),
            CLICKHOUSE_USUARIO: env::var("CLICKHOUSE_USUARIO")
                .expect("CLICKHOUSE_USUARIO debe estar definida"),
            CLICKHOUSE_CONTRASENA: env::var("CLICKHOUSE_CONTRASENA")
                .expect("CLICKHOUSE_CONTRASENA debe estar definida"),
            CLICKHOUSE_BD: env::var("CLICKHOUSE_BD")
                .expect("CLICKHOUSE_BD debe estar definida"),

            PUERTO_METRICAS: env::var("PUERTO_METRICAS")
                .ok()
                .and_then(|v| v.parse().ok())
                .unwrap_or(9100),

            RUST_LOG: env::var("RUST_LOG").unwrap_or_else(|_| "info".to_string()),
        }
    }
}
