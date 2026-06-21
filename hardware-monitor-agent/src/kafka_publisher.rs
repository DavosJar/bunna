use crate::types::{Snapshot, AlertEvent};
use rdkafka::producer::{FutureProducer, FutureRecord};
use rdkafka::config::ClientConfig;
use std::time::Duration;
use tokio::sync::mpsc;
use tracing::{error, info};

pub struct KafkaPublisher {
    producer: FutureProducer,
    metrics_topic: String,
    alerts_topic: String,
    snapshot_rx: mpsc::Receiver<Snapshot>,
    alert_rx: mpsc::Receiver<AlertEvent>,
}

impl KafkaPublisher {
    pub fn new(
        brokers: &str,
        metrics_topic: &str,
        alerts_topic: &str,
        snapshot_rx: mpsc::Receiver<Snapshot>,
        alert_rx: mpsc::Receiver<AlertEvent>,
    ) -> anyhow::Result<Self> {
        let producer: FutureProducer = ClientConfig::new()
            .set("bootstrap.servers", brokers)
            .set("message.timeout.ms", "15000")
            // tighter buffering settings
            .set("queue.buffering.max.messages", "100")
            .set("queue.buffering.max.kbytes", "10240")
            .set("queue.buffering.max.ms", "10")
            .set("socket.send.buffer.bytes", "131072")
            .create()?;

        Ok(Self {
            producer,
            metrics_topic: metrics_topic.to_string(),
            alerts_topic: alerts_topic.to_string(),
            snapshot_rx,
            alert_rx,
        })
    }

    pub async fn run(mut self) {
        info!("Publicador de Kafka iniciado");
        loop {
            tokio::select! {
                Some(snapshot) = self.snapshot_rx.recv() => {
                    self.publish_metrics(snapshot).await;
                }
                Some(alert) = self.alert_rx.recv() => {
                    self.publish_alert(alert).await;
                }
                else => break,
            }
        }
    }

    async fn publish_metrics(&self, snapshot: Snapshot) {
        match serde_json::to_string(&snapshot) {
            Ok(json) => {
                let record = FutureRecord::to(&self.metrics_topic)
                    .payload(&json)
                    .key(&snapshot.node_id);
                
                match self.producer.send(record, Duration::from_secs(0)).await {
                    Err((e, _)) => error!("Failed to publish metric to Kafka: {}", e),
                    Ok(_) => info!("Métrica publicada en Kafka: {} @ {}", snapshot.node_id, snapshot.timestamp),
                }
            }
            Err(e) => error!("Failed to serialize snapshot: {}", e),
        }
    }

    async fn publish_alert(&self, alert: AlertEvent) {
        match serde_json::to_string(&alert) {
            Ok(json) => {
                let record = FutureRecord::to(&self.alerts_topic)
                    .payload(&json)
                    .key(&alert.node_id);
                
                if let Err((e, _)) = self.producer.send(record, Duration::from_secs(0)).await {
                    error!("Failed to publish alert to Kafka: {}", e);
                } else {
                    info!("Alerta publicada en Kafka: {} ({})", alert.metric, alert.severity);
                }
            }
            Err(e) => error!("Failed to serialize alert: {}", e),
        }
    }
}
