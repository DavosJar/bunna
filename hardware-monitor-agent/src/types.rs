use serde::{Deserialize, Serialize};
use chrono::{DateTime, Utc};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Snapshot {
    pub node_id: String,
    pub timestamp: DateTime<Utc>,
    pub interval_ms: u64,
    pub cloud_provider: String,
    pub node_hostname: String,
    pub cpu: CpuMetrics,
    pub ram: RamMetrics,
    pub disks: Vec<DiskMetrics>,
    pub net: NetMetrics,
    pub containers: Vec<ContainerMetrics>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct CpuMetrics {
    pub usage_percent: f64,
    pub cores: u32,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct RamMetrics {
    pub total_mb: u64,
    pub used_mb: u64,
    pub available_mb: u64,
    pub usage_percent: f64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct DiskMetrics {
    pub mount: String,
    pub total_gb: f64,
    pub used_gb: f64,
    pub available_gb: f64,
    pub usage_percent: f64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct NetMetrics {
    pub interfaces: Vec<InterfaceMetrics>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct InterfaceMetrics {
    pub name: String,
    pub received_bytes: u64,
    pub transmitted_bytes: u64,
    pub received_bytes_per_sec: f64,
    pub transmitted_bytes_per_sec: f64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct ContainerMetrics {
    pub container_id: String,
    pub cpu_shares: u64,
    pub memory_limit_mb: u64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(rename_all = "lowercase")]
pub enum EventType {
    Alert,
    Recovery,
}

use std::fmt;


impl fmt::Display for Severity {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{:?}", self)
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum Severity {
    Warn,
    Critical,
    Info,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Thresholds {
    pub cpu_warn: f64,
    pub cpu_critical: f64,
    pub ram_warn: f64,
    pub ram_critical: f64,
    pub disk_warn: f64,
    pub disk_critical: f64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct AlertEvent {
    pub node_id: String,
    pub timestamp: DateTime<Utc>,
    pub metric: String,
    pub severity: Severity,
    pub value: f64,
    pub threshold: f64,
    pub message: String,
    pub previous_state: String,
    pub event_type: EventType,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum MetricType {
    Cpu,
    Ram,
    Disk(usize), // Index in the disks vector
}
