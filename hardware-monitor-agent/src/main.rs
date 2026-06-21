use hardware_monitor_agent::aggregator::PartialMetrics;
use hardware_monitor_agent::collector::{CpuCollector, DiskCollector, NetCollector, RamCollector};
use hardware_monitor_agent::{Aggregator, Config, Guardian, KafkaPublisher};
use tokio::sync::mpsc;
use tracing::{info, error};
use std::time::Duration;
use std::process;
#[tokio::main(flavor = "current_thread")]
async fn main() -> anyhow::Result<()> {
    // Cargar variables de entorno desde .env
    dotenvy::dotenv().ok();

    // Inicializar logs
    tracing_subscriber::fmt::init();

    info!("Iniciando agente de monitoreo de hardware...");
    let pid = process::id();
    info!("PID: {}", pid);

    // 1. Cargar configuración
    let config = Config::from_env()?;
    let interval = Duration::from_millis(config.interval_ms);

    // 2. Crear canales de comunicación (tokio::sync::mpsc)
    // REQ-NF-003: Comunicación interna mediante mpsc
    let (metrics_tx, metrics_rx) = mpsc::channel(32);
    let (snapshot_for_publisher_tx, snapshot_for_publisher_rx) = mpsc::channel(16);
    let (snapshot_for_guardian_tx, snapshot_for_guardian_rx) = mpsc::channel(16);
    let (alert_tx, alert_rx) = mpsc::channel(64);

    // 3. Crear componentes
    let aggregator = Aggregator::new(
        config.node_id.clone(),
        config.interval_ms,
        metrics_rx,
        snapshot_for_publisher_tx,
        snapshot_for_guardian_tx,
    );

    let guardian = Guardian::new(
        config.clone(),
        snapshot_for_guardian_rx,
        alert_tx,
    );

    let publisher = KafkaPublisher::new(
        &config.kafka_brokers,
        "hardware.metrics",
        "hardware.alerts",
        snapshot_for_publisher_rx,
        alert_rx,
    )?;

    // 4. Iniciar tareas en paralelo
    tokio::spawn(async move { aggregator.run().await });
    tokio::spawn(async move { guardian.run().await });
    tokio::spawn(async move { publisher.run().await });

    // 5. Bucle principal de recolección
    let mut cpu_coll = CpuCollector::new();
    let ram_coll = RamCollector::new();
    let disk_coll = DiskCollector::new();
    let mut net_coll = NetCollector::new(config.interval_ms);

    info!("Bucle de monitoreo iniciado (intervalo: {:?})", interval);
    
    let mut ticker = tokio::time::interval(interval);
    loop {
        ticker.tick().await;

        // Recolectar CPU
        match cpu_coll.collect() {
            Ok(m) => { let _ = metrics_tx.send(PartialMetrics::Cpu(m)).await; }
            Err(e) => error!("CPU collection failed: {}", e),
        }

        // Recolectar RAM
        match ram_coll.collect() {
            Ok(m) => { let _ = metrics_tx.send(PartialMetrics::Ram(m)).await; }
            Err(e) => error!("RAM collection failed: {}", e),
        }

        // Recolectar discos
        match disk_coll.collect() {
            Ok(m) => { let _ = metrics_tx.send(PartialMetrics::Disks(m)).await; }
            Err(e) => error!("Disk collection failed: {}", e),
        }

        // Recolectar red
        match net_coll.collect() {
            Ok(m) => { let _ = metrics_tx.send(PartialMetrics::Net(m)).await; }
            Err(e) => error!("Net collection failed: {}", e),
        }
    }
}
