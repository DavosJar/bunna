use std::collections::HashMap;
use crate::types::{AlertEvent, EventType, Severity, Snapshot, MetricType};
use crate::config::Config;
use chrono::Utc;
use tokio::sync::mpsc;
use tracing::{info, warn};

pub enum AlertState {
    Normal,
    AlertSent,
}

pub struct Guardian {
    config: Config,
    state: HashMap<MetricType, AlertState>,
    snapshot_rx: mpsc::Receiver<Snapshot>,
    alert_tx: mpsc::Sender<AlertEvent>,
}

impl Guardian {
    pub fn new(
        config: Config,
        snapshot_rx: mpsc::Receiver<Snapshot>,
        alert_tx: mpsc::Sender<AlertEvent>,
    ) -> Self {
        Self {
            config,
            state: HashMap::new(),
            snapshot_rx,
            alert_tx,
        }
    }

    pub async fn run(mut self) {
        info!("Guardian iniciado");
        while let Some(snapshot) = self.snapshot_rx.recv().await {
            self.evaluate_cpu(&snapshot).await;
            self.evaluate_ram(&snapshot).await;
            self.evaluate_disks(&snapshot).await;
        }
    }

    async fn evaluate_cpu(&mut self, snapshot: &Snapshot) {
        let val = snapshot.cpu.usage_percent;
        self.evaluate_metric(
            MetricType::Cpu,
            "cpu",
            val,
            self.config.cpu_warn,
            self.config.cpu_critical,
            format!("CPU usage at {:.1}% exceeds CRITICAL threshold of {:.1}%", val, self.config.cpu_critical),
            snapshot,
        ).await;
    }

    async fn evaluate_ram(&mut self, snapshot: &Snapshot) {
        let val = snapshot.ram.usage_percent;
        self.evaluate_metric(
            MetricType::Ram,
            "ram",
            val,
            self.config.ram_warn,
            self.config.ram_critical,
            format!("RAM usage at {:.1}% exceeds CRITICAL threshold of {:.1}%", val, self.config.ram_critical),
            snapshot,
        ).await;
    }

    async fn evaluate_disks(&mut self, snapshot: &Snapshot) {
        for (i, disk) in snapshot.disks.iter().enumerate() {
            let val = disk.usage_percent;
            self.evaluate_metric(
                MetricType::Disk(i),
                &format!("disk:{}", disk.mount),
                val,
                self.config.disk_warn,
                self.config.disk_critical,
                format!("Disk usage on {} at {:.1}% exceeds CRITICAL threshold of {:.1}%", disk.mount, val, self.config.disk_critical),
                snapshot,
            ).await;
        }
    }

    async fn evaluate_metric(
        &mut self,
        metric_type: MetricType,
        metric_name: &str,
        value: f64,
        warn_threshold: f64,
        critical_threshold: f64,
        message: String,
        snapshot: &Snapshot,
    ) {
        let current_state = self.state.entry(metric_type).or_insert(AlertState::Normal);

        // 1. Caso CRITICAL
        if value >= critical_threshold {
            // Control de inundación: Si ya enviamos alerta, no hacemos nada
            if let AlertState::AlertSent = current_state {
                return;
            }

            // Primera vez que entra en CRITICAL
            warn!("CRÍTICO: {}", message);
            let alert = AlertEvent {
                node_id: snapshot.node_id.clone(),
                timestamp: Utc::now(),
                metric: metric_name.to_string(),
                severity: Severity::Critical,
                value,
                threshold: critical_threshold,
                message,
                previous_state: "normal".to_string(),
                event_type: EventType::Alert,
            };
            let _ = self.alert_tx.send(alert).await;
            *current_state = AlertState::AlertSent;
            return;
        }

        // 2. Caso WARN (solo log, según REQ-FUNC-007)
        if value >= warn_threshold {
            warn!("ADVERTENCIA: {} usage at {:.1}% (threshold {:.1}%)", metric_name, value, warn_threshold);
            return;
        }

        // 3. Caso RECOVERY (Si llegamos aquí es porque value < warn_threshold)
        // Solo enviamos recovery si antes estábamos en estado de alerta
        if let AlertState::Normal = current_state {
            return;
        }

        info!("RECUPERADO: {} volvió a la normalidad: {:.1}%", metric_name, value);
        let recovery = AlertEvent {
            node_id: snapshot.node_id.clone(),
            timestamp: Utc::now(),
            metric: metric_name.to_string(),
            severity: Severity::Info,
            value,
            threshold: critical_threshold,
            message: format!("{} usage returned to normal: {:.1}%", metric_name, value),
            previous_state: "critical".to_string(),
            event_type: EventType::Recovery,
        };
        let _ = self.alert_tx.send(recovery).await;
        *current_state = AlertState::Normal;
    }
}
