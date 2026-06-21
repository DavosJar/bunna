//! Errores específicos del consumidor de Kafka.

use std::fmt;

/// Errores que puede producir el consumidor de Kafka al procesar mensajes.
#[derive(Debug)]
pub enum ErrorConsumidor {
    /// Error de conexión o lectura de Kafka.
    #[allow(dead_code)]
    Kafka(String),
    /// Error al parsear JSON del mensaje.
    ParseoJson(String),
    /// Error al insertar en ClickHouse.
    ClickHouse(String),
    /// Error genérico / desconocido.
    Otro(String),
}

impl fmt::Display for ErrorConsumidor {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ErrorConsumidor::Kafka(msg) => write!(f, "Error de Kafka: {}", msg),
            ErrorConsumidor::ParseoJson(msg) => write!(f, "Error de parseo JSON: {}", msg),
            ErrorConsumidor::ClickHouse(msg) => write!(f, "Error de ClickHouse: {}", msg),
            ErrorConsumidor::Otro(msg) => write!(f, "Error: {}", msg),
        }
    }
}

impl std::error::Error for ErrorConsumidor {}

impl From<serde_json::Error> for ErrorConsumidor {
    fn from(e: serde_json::Error) -> Self {
        ErrorConsumidor::ParseoJson(e.to_string())
    }
}

impl From<reqwest::Error> for ErrorConsumidor {
    fn from(e: reqwest::Error) -> Self {
        ErrorConsumidor::ClickHouse(e.to_string())
    }
}

impl From<String> for ErrorConsumidor {
    fn from(e: String) -> Self {
        ErrorConsumidor::Otro(e)
    }
}
