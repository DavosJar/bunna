//! Tipos para mensajes de telemetría (API, NEGOCIO, BD).
//!
//! Todos los structs permiten `non_snake_case` para mantener
//! nombres en español legibles.

use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;

/// Mensaje de telemetría proveniente de Kafka (topic `telemetry`).
///
/// El campo `log_type` determina cuál de los tres detalles opcionales
/// (`api`, `negocio`, `bd`) viene poblado.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct MensajeTelemetria {
    /// Tipo de log: "API" | "NEGOCIO" | "BD"
    pub log_type: String,
    pub level: String,
    pub timestamp: String,
    pub trace_id: String,
    pub span_id: String,
    pub service_name: String,
    pub environment: String,

    // Cada variante es opcional; se interpreta según `log_type`.
    pub api: Option<DatosApi>,
    pub negocio: Option<DatosNegocio>,
    pub bd: Option<DatosBd>,
}

/// Datos específicos de una llamada a API REST.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct DatosApi {
    pub method: String,
    pub path: String,
    pub status_code: u16,
    pub duration_ms: f64,
    pub client_ip: String,
    pub user_agent: String,
    pub content_length: i64,
}

/// Datos específicos de un caso de uso de negocio.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct DatosNegocio {
    pub use_case: String,
    pub command: HashMap<String, Value>,
    pub result: String,
    pub user_id: String,
    pub details: HashMap<String, Value>,
    pub duration_usecase_ms: f64,
}

/// Datos específicos de una operación de base de datos.
#[derive(Debug, Clone, Serialize, Deserialize)]
#[allow(non_snake_case)]
pub struct DatosBd {
    pub operation: String,
    pub table: String,
    pub duration_ms: f64,
    pub rows_affected: i64,
    pub error_sql_state: Option<String>,
    pub query_hash: String,
}
