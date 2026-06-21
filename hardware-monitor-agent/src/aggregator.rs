use crate::types::{Snapshot, CpuMetrics, RamMetrics, DiskMetrics, NetMetrics};
use tokio::sync::mpsc;
use tracing::info;

pub struct Aggregator {
    metrics_rx: mpsc::Receiver<PartialMetrics>,
    snapshot_tx: mpsc::Sender<Snapshot>,
    guardian_tx: mpsc::Sender<Snapshot>,
    node_id: String,
    interval_ms: u64,
    cloud_provider: String,
    node_hostname: String,
}

pub enum PartialMetrics {
    Cpu(CpuMetrics),
    Ram(RamMetrics),
    Disks(Vec<DiskMetrics>),
    Net(NetMetrics),
    // Agregaremos containers después
}

impl Aggregator {
    pub fn new(
        node_id: String,
        interval_ms: u64,
        cloud_provider: String,
        node_hostname: String,
        metrics_rx: mpsc::Receiver<PartialMetrics>,
        snapshot_tx: mpsc::Sender<Snapshot>,
        guardian_tx: mpsc::Sender<Snapshot>,
    ) -> Self {
        Self {
            node_id,
            interval_ms,
            cloud_provider,
            node_hostname,
            metrics_rx,
            snapshot_tx,
            guardian_tx,
        }
    }

    pub async fn run(mut self) {
        info!("Agregador iniciado");
        
        let mut cpu = None;
        let mut ram = None;
        let mut disks = None;
        let mut net = None;

        while let Some(partial) = self.metrics_rx.recv().await {
            match partial {
                PartialMetrics::Cpu(m) => cpu = Some(m),
                PartialMetrics::Ram(m) => ram = Some(m),
                PartialMetrics::Disks(m) => disks = Some(m),
                PartialMetrics::Net(m) => net = Some(m),
            }

            // REQ-GUD-002: Esperar a todos antes de construir snapshot
            if let (Some(c), Some(r), Some(d), Some(n)) = (cpu.clone(), ram.clone(), disks.clone(), net.clone()) {
                let snapshot = Snapshot {
                    node_id: self.node_id.clone(),
                    timestamp: chrono::Utc::now(),
                    interval_ms: self.interval_ms,
                    cloud_provider: self.cloud_provider.clone(),
                    node_hostname: self.node_hostname.clone(),
                    cpu: c,
                    ram: r,
                    disks: d,
                    net: n,
                    containers: vec![], // TODO
                };

                // Enviamos a guardian para alertas
                let _ = self.guardian_tx.send(snapshot.clone()).await;
                // Enviamos a publisher para métricas periódicas
                let _ = self.snapshot_tx.send(snapshot).await;

                // Reset para el siguiente tick
                cpu = None;
                ram = None;
                disks = None;
                net = None;
            }
        }
    }
}
