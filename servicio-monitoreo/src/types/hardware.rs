//! Tipos para métricas y alertas de hardware.
//!
//! `InstantaneaHardware` representa una fotografía completa del estado
//! de un nodo en un instante dado.  `AlertaHardware` representa una
//! alerta o recuperación disparada por una condición de threshold.

use serde::{Deserialize, Serialize};

/// Fotografía completa del hardware de un nodo.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct InstantaneaHardware {
    pub node_id: String,
    pub timestamp: String,
    pub interval_ms: u64,

    pub cpu: CpuMetricas,
    pub ram: RamMetricas,
    pub disks: Vec<DiscoMetricas>,
    pub net: RedMetricas,
    pub containers: Vec<ContenedorMetricas>,
}

/// Métricas de CPU.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct CpuMetricas {
    pub usage_percent: f64,
    pub cores: u32,
}

/// Métricas de memoria RAM.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct RamMetricas {
    pub total_mb: u64,
    pub used_mb: u64,
    pub available_mb: u64,
    pub usage_percent: f64,
}

/// Métricas de un disco / punto de montaje.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct DiscoMetricas {
    pub mount: String,
    pub total_gb: f64,
    pub used_gb: f64,
    pub available_gb: f64,
    pub usage_percent: f64,
}

/// Métricas de red (contiene una lista de interfaces).
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct RedMetricas {
    pub interfaces: Vec<InterfazMetricas>,
}

/// Métricas de una interfaz de red individual.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct InterfazMetricas {
    pub name: String,
    pub received_bytes: u64,
    pub transmitted_bytes: u64,
    pub received_bytes_per_sec: f64,
    pub transmitted_bytes_per_sec: f64,
}

/// Métricas de un contenedor Docker / container runtime.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct ContenedorMetricas {
    pub container_id: String,
    pub cpu_shares: f64,
    pub memory_limit_mb: u64,
}

/// Alerta de hardware (disparo o recuperación).
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct AlertaHardware {
    pub node_id: String,
    pub timestamp: String,

    /// Nombre de la métrica que disparó la alerta (ej. "cpu.usage_percent").
    pub metric: String,
    /// Severidad: "critical", "warning", "info".
    pub severity: String,

    /// Valor actual de la métrica.
    pub value: f64,
    /// Threshold que se superó.
    pub threshold: f64,
    /// Mensaje legible de la alerta.
    pub message: String,

    /// Estado anterior del nodo (ej. "ok", "warning").
    pub previous_state: String,
    /// Tipo de evento: "alert" o "recovery".
    pub event_type: String,
}
