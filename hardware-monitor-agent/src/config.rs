use anyhow::{Context, Result};
use std::env;

#[derive(Debug, Clone)]
pub struct Config {
    pub node_id: String,
    pub kafka_brokers: String,
    pub kafka_topic_metrics: String,
    pub kafka_topic_alerts: String,
    pub interval_ms: u64,
    pub cpu_warn: f64,
    pub cpu_critical: f64,
    pub ram_warn: f64,
    pub ram_critical: f64,
    pub disk_warn: f64,
    pub disk_critical: f64,
    pub cloud_provider: String,
}

impl Config {
    pub fn from_env() -> Result<Self> {
        Ok(Self {
            node_id: env::var("NODE_ID").context("NODE_ID not set")?,
            kafka_brokers: env::var("KAFKA_BROKERS").unwrap_or_else(|_| "localhost:9092".to_string()),
            kafka_topic_metrics: env::var("KAFKA_TOPIC_METRICS")
                .unwrap_or_else(|_| "hardware.metrics".to_string()),
            kafka_topic_alerts: env::var("KAFKA_TOPIC_ALERTS")
                .unwrap_or_else(|_| "hardware.alerts".to_string()),
            interval_ms: env::var("INTERVAL_MS")
                .unwrap_or_else(|_| "10000".to_string())
                .parse()?,
            cpu_warn: env::var("CPU_WARN_PERCENT")
                .unwrap_or_else(|_| "80.0".to_string())
                .parse()?,
            cpu_critical: env::var("CPU_CRITICAL_PERCENT")
                .unwrap_or_else(|_| "90.0".to_string())
                .parse()?,
            ram_warn: env::var("RAM_WARN_PERCENT")
                .unwrap_or_else(|_| "85.0".to_string())
                .parse()?,
            ram_critical: env::var("RAM_CRITICAL_PERCENT")
                .unwrap_or_else(|_| "95.0".to_string())
                .parse()?,
            disk_warn: env::var("DISK_WARN_PERCENT")
                .unwrap_or_else(|_| "85.0".to_string())
                .parse()?,
            disk_critical: env::var("DISK_CRITICAL_PERCENT")
                .unwrap_or_else(|_| "95.0".to_string())
                .parse()?,
            cloud_provider: env::var("CLOUD_PROVIDER").unwrap_or_else(|_| "local".to_string()),
        })
    }
}
