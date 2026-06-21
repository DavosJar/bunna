//! Módulo de consumo de Kafka.
//!
//! Contiene el `ManejadorKafka` que recibe mensajes crudos de los topics
//! y los rutea al `BatchInserter` correspondiente según el tipo.

pub mod error;

use std::sync::Arc;

use rdkafka::Message;
use tracing::warn;

use crate::sink::insertador::BatchInserter;
use crate::types::*;
use error::ErrorConsumidor;

/// Manejador de mensajes de Kafka.
///
/// Recibe un mensaje crudo, determina su topic, lo parsea al tipo
/// correspondiente y lo envía al `BatchInserter`.
pub struct ManejadorKafka {
    insertador: Arc<BatchInserter>,
}

impl ManejadorKafka {
    /// Crea un nuevo manejador asociado al insertador.
    pub fn nuevo(insertador: Arc<BatchInserter>) -> Self {
        Self { insertador }
    }

    /// Procesa un mensaje recibido de Kafka.
    ///
    /// # Parámetros
    ///
    /// * `mensaje` – Mensaje crudo de rdkafka.
    /// * `topic`  – Topic del que proviene el mensaje.
    pub async fn manejar_mensaje(
        &self,
        mensaje: &rdkafka::message::BorrowedMessage<'_>,
    ) -> Result<(), ErrorConsumidor> {
        let topic = mensaje
            .topic()
            .to_string();

        let payload = mensaje
            .payload()
            .ok_or_else(|| ErrorConsumidor::Otro("Mensaje sin payload".into()))?;

        let texto = std::str::from_utf8(payload)
            .map_err(|e| ErrorConsumidor::Otro(format!("Payload no es UTF-8: {}", e)))?;

        tracing::debug!("Mensaje recibido de topic '{}': {} bytes", topic, texto.len());

        match topic.as_str() {
            t if t.ends_with("telemetry")
                || t == "telemetry"
                || t.ends_with("telemetria") =>
            {
                self.procesar_telemetria(texto).await
            }
            t if t.ends_with("hardware.metrics")
                || t.ends_with("hardware_metricas") =>
            {
                self.procesar_hardware_metricas(texto).await
            }
            t if t.ends_with("hardware.alerts")
                || t.ends_with("hardware_alertas") =>
            {
                self.procesar_hardware_alertas(texto).await
            }
            otro => {
                warn!("Topic desconocido, se ignora: {}", otro);
                Ok(())
            }
        }
    }

    // ----------------------------------------------------------------
    // Procesadores específicos
    // ----------------------------------------------------------------

    /// Parsea un `MensajeTelemetria` y lo rutea según `log_type`.
    async fn procesar_telemetria(&self, json: &str) -> Result<(), ErrorConsumidor> {
        let msg: MensajeTelemetria = serde_json::from_str(json)?;

        match msg.log_type.as_str() {
            "API" => {
                if let Some(ref api) = msg.api {
                    self.insertador.insertar_fila(
                        "telemetria_api",
                        vec![
                            "log_type".into(),
                            "level".into(),
                            "timestamp".into(),
                            "trace_id".into(),
                            "span_id".into(),
                            "service_name".into(),
                            "environment".into(),
                            "method".into(),
                            "path".into(),
                            "status_code".into(),
                            "duration_ms".into(),
                            "client_ip".into(),
                            "user_agent".into(),
                            "content_length".into(),
                        ],
                        vec![
                            msg.log_type.clone(),
                            msg.level.clone(),
                            msg.timestamp.clone(),
                            msg.trace_id.clone(),
                            msg.span_id.clone(),
                            msg.service_name.clone(),
                            msg.environment.clone(),
                            api.method.clone(),
                            api.path.clone(),
                            api.status_code.to_string(),
                            api.duration_ms.to_string(),
                            api.client_ip.clone(),
                            api.user_agent.clone(),
                            api.content_length.to_string(),
                        ],
                    );
                } else {
                    warn!("log_type=API pero el campo api es None");
                }
            }
            "NEGOCIO" => {
                if let Some(ref neg) = msg.negocio {
                    let command_json = serde_json::to_string(&neg.command)
                        .unwrap_or_default();
                    let details_json = serde_json::to_string(&neg.details)
                        .unwrap_or_default();

                    self.insertador.insertar_fila(
                        "telemetria_negocio",
                        vec![
                            "log_type".into(),
                            "level".into(),
                            "timestamp".into(),
                            "trace_id".into(),
                            "span_id".into(),
                            "service_name".into(),
                            "environment".into(),
                            "use_case".into(),
                            "command".into(),
                            "result".into(),
                            "user_id".into(),
                            "details".into(),
                            "duration_usecase_ms".into(),
                        ],
                        vec![
                            msg.log_type.clone(),
                            msg.level.clone(),
                            msg.timestamp.clone(),
                            msg.trace_id.clone(),
                            msg.span_id.clone(),
                            msg.service_name.clone(),
                            msg.environment.clone(),
                            neg.use_case.clone(),
                            command_json,
                            neg.result.clone(),
                            neg.user_id.clone(),
                            details_json,
                            neg.duration_usecase_ms.to_string(),
                        ],
                    );
                } else {
                    warn!("log_type=NEGOCIO pero el campo negocio es None");
                }
            }
            "BD" => {
                if let Some(ref bd) = msg.bd {
                    self.insertador.insertar_fila(
                        "telemetria_bd",
                        vec![
                            "log_type".into(),
                            "level".into(),
                            "timestamp".into(),
                            "trace_id".into(),
                            "span_id".into(),
                            "service_name".into(),
                            "environment".into(),
                            "operation".into(),
                            "table".into(),
                            "duration_ms".into(),
                            "rows_affected".into(),
                            "error_sql_state".into(),
                            "query_hash".into(),
                        ],
                        vec![
                            msg.log_type.clone(),
                            msg.level.clone(),
                            msg.timestamp.clone(),
                            msg.trace_id.clone(),
                            msg.span_id.clone(),
                            msg.service_name.clone(),
                            msg.environment.clone(),
                            bd.operation.clone(),
                            bd.table.clone(),
                            bd.duration_ms.to_string(),
                            bd.rows_affected.to_string(),
                            bd.error_sql_state.clone().unwrap_or_default(),
                            bd.query_hash.clone(),
                        ],
                    );
                } else {
                    warn!("log_type=BD pero el campo bd es None");
                }
            }
            otro => {
                warn!("log_type desconocido: {}", otro);
            }
        }

        Ok(())
    }

    /// Parsea una `InstantaneaHardware` y la aplanada en filas para
    /// hardware_metricas (una fila por disco, interfaz y contenedor).
    async fn procesar_hardware_metricas(&self, json: &str) -> Result<(), ErrorConsumidor> {
        let snap: InstantaneaHardware = serde_json::from_str(json)?;

        let columnas_base = vec![
            "node_id".into(),
            "timestamp".into(),
            "interval_ms".into(),
            "cpu_usage_percent".into(),
            "cpu_cores".into(),
            "ram_total_mb".into(),
            "ram_used_mb".into(),
            "ram_available_mb".into(),
            "ram_usage_percent".into(),
        ];

        let valores_base: Vec<String> = vec![
            snap.node_id.clone(),
            snap.timestamp.clone(),
            snap.interval_ms.to_string(),
            snap.cpu.usage_percent.to_string(),
            snap.cpu.cores.to_string(),
            snap.ram.total_mb.to_string(),
            snap.ram.used_mb.to_string(),
            snap.ram.available_mb.to_string(),
            snap.ram.usage_percent.to_string(),
        ];

        // Fila base con métricas de CPU y RAM
        self.insertador
            .insertar_fila("hardware_metricas", columnas_base.clone(), valores_base.clone());

        // Una fila por disco
        for disco in &snap.disks {
            let mut cols = columnas_base.clone();
            let mut vals = valores_base.clone();
            cols.extend_from_slice(&[
                "disco_mount".into(),
                "disco_total_gb".into(),
                "disco_used_gb".into(),
                "disco_available_gb".into(),
                "disco_usage_percent".into(),
            ]);
            vals.extend_from_slice(&[
                disco.mount.clone(),
                disco.total_gb.to_string(),
                disco.used_gb.to_string(),
                disco.available_gb.to_string(),
                disco.usage_percent.to_string(),
            ]);
            self.insertador.insertar_fila("hardware_metricas", cols, vals);
        }

        // Una fila por interfaz de red
        for iface in &snap.net.interfaces {
            let mut cols = columnas_base.clone();
            let mut vals = valores_base.clone();
            cols.extend_from_slice(&[
                "interfaz_name".into(),
                "interfaz_received_bytes".into(),
                "interfaz_transmitted_bytes".into(),
                "interfaz_received_bytes_per_sec".into(),
                "interfaz_transmitted_bytes_per_sec".into(),
            ]);
            vals.extend_from_slice(&[
                iface.name.clone(),
                iface.received_bytes.to_string(),
                iface.transmitted_bytes.to_string(),
                iface.received_bytes_per_sec.to_string(),
                iface.transmitted_bytes_per_sec.to_string(),
            ]);
            self.insertador.insertar_fila("hardware_metricas", cols, vals);
        }

        // Una fila por contenedor
        for cont in &snap.containers {
            let mut cols = columnas_base.clone();
            let mut vals = valores_base.clone();
            cols.extend_from_slice(&[
                "contenedor_id".into(),
                "contenedor_cpu_shares".into(),
                "contenedor_memory_limit_mb".into(),
            ]);
            vals.extend_from_slice(&[
                cont.container_id.clone(),
                cont.cpu_shares.to_string(),
                cont.memory_limit_mb.to_string(),
            ]);
            self.insertador.insertar_fila("hardware_metricas", cols, vals);
        }

        Ok(())
    }

    /// Parsea una `AlertaHardware` y la inserta como fila.
    async fn procesar_hardware_alertas(&self, json: &str) -> Result<(), ErrorConsumidor> {
        let alerta: AlertaHardware = serde_json::from_str(json)?;

        self.insertador.insertar_fila(
            "hardware_alertas",
            vec![
                "node_id".into(),
                "timestamp".into(),
                "metric".into(),
                "severity".into(),
                "value".into(),
                "threshold".into(),
                "message".into(),
                "previous_state".into(),
                "event_type".into(),
            ],
            vec![
                alerta.node_id,
                alerta.timestamp,
                alerta.metric,
                alerta.severity,
                alerta.value.to_string(),
                alerta.threshold.to_string(),
                alerta.message,
                alerta.previous_state,
                alerta.event_type,
            ],
        );

        Ok(())
    }
}
