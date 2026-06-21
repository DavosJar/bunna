//! Insertador por lotes (BatchInserter).
//!
//! Acumula filas en memoria y las descarga periódicamente a ClickHouse
//! usando el `ClienteClickHouse`.  El flush se dispara cada 500 ms o
//! al alcanzar 5000 filas acumuladas, lo que ocurra primero.

use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;

use prometheus::{register_int_counter, IntCounter};
use tokio::sync::mpsc;
use tokio::time::interval;
use tracing::{error, info};

use super::clickhouse::ClienteClickHouse;
use crate::consumidor::error::ErrorConsumidor;

// ---------------------------------------------------------------------------
// Métricas Prometheus
// ---------------------------------------------------------------------------

lazy_static::lazy_static! {
    static ref FILAS_INSERTADAS: IntCounter = register_int_counter!(
        "smon_clickhouse_filas_insertadas_total",
        "Número total de filas insertadas exitosamente en ClickHouse"
    )
    .expect("registro de métrica FILAS_INSERTADAS");

    static ref FILAS_FALLIDAS: IntCounter = register_int_counter!(
        "smon_clickhouse_filas_fallidas_total",
        "Número total de filas que fallaron al insertarse en ClickHouse"
    )
    .expect("registro de métrica FILAS_FALLIDAS");

    static ref LOTES_ENVIADOS: IntCounter = register_int_counter!(
        "smon_clickhouse_lotes_enviados_total",
        "Número total de lotes (batches) enviados a ClickHouse"
    )
    .expect("registro de métrica LOTES_ENVIADOS");

    static ref LOTES_FALLIDOS: IntCounter = register_int_counter!(
        "smon_clickhouse_lotes_fallidos_total",
        "Número total de lotes que fallaron al enviarse a ClickHouse"
    )
    .expect("registro de métrica LOTES_FALLIDOS");
}

// ---------------------------------------------------------------------------
// Tipos internos
// ---------------------------------------------------------------------------

/// Una fila pendiente de inserción: (columnas, valores).
type Fila = (Vec<String>, Vec<String>);

/// Mensaje interno del canal: (nombre_tabla, columnas, valores).
type Mensaje = (String, Vec<String>, Vec<String>);

// ---------------------------------------------------------------------------
// BatchInserter
// ---------------------------------------------------------------------------

/// Insertador por lotes con buffers encolados y flush periódico.
///
/// # Ejemplo
///
/// ```ignore
/// let cliente = ClienteClickHouse::new(url, user, pass, db);
/// let insertador = Arc::new(BatchInserter::nuevo(cliente));
///
/// insertador.insertar_fila("telemetria_api", columnas, valores);
/// ```
pub struct BatchInserter {
    /// Canal de transmisión hacia la tarea de fondo.
    tx: mpsc::UnboundedSender<Mensaje>,
}

impl BatchInserter {
    /// Crea un nuevo `BatchInserter` e inicia la tarea de fondo que
    /// vacía los buffers periódicamente.
    ///
    /// # Argumentos
    ///
    /// * `cliente` – Cliente HTTP para ClickHouse.
    pub fn nuevo(cliente: ClienteClickHouse) -> Self {
        let (tx, rx) = mpsc::unbounded_channel::<Mensaje>();

        let cliente = Arc::new(cliente);
        tokio::spawn(tarea_flujo(cliente, rx));

        info!("BatchInserter iniciado con flush cada 500 ms / 5000 filas");

        Self { tx }
    }

    /// Encola una fila para inserción en la tabla indicada.
    ///
    /// # Argumentos
    ///
    /// * `tabla`    – Nombre de la tabla destino (ej. "telemetria_api").
    /// * `columnas` – Nombres de las columnas en orden.
    /// * `valores`  – Valores de la fila en el mismo orden que `columnas`.
    pub fn insertar_fila(&self, tabla: &str, columnas: Vec<String>, valores: Vec<String>) {
        if let Err(e) = self.tx.send((tabla.to_string(), columnas, valores)) {
            error!("Error al encolar fila para '{}': {}", tabla, e);
            FILAS_FALLIDAS.inc();
        }
    }
}

// ---------------------------------------------------------------------------
// Tarea de fondo
// ---------------------------------------------------------------------------

/// Tarea asíncrona que recibe filas del canal, las acumula en buffers
/// por tabla y las descarga a ClickHouse cada 500 ms o al llegar a 5000.
async fn tarea_flujo(cliente: Arc<ClienteClickHouse>, mut rx: mpsc::UnboundedReceiver<Mensaje>) {
    // Buffer por tabla: nombre_tabla → Vec<Fila>
    let mut buffers: HashMap<String, Vec<Fila>> = HashMap::new();
    let mut total_filas: usize = 0;

    // Temporizador para flush periódico
    let mut tick = interval(Duration::from_millis(500));
    // El primer tick se dispara inmediatamente; lo ajustamos para que no.
    tick.tick().await;

    loop {
        tokio::select! {
            // Recibir una nueva fila desde el canal
            Some((tabla, columnas, valores)) = rx.recv() => {
                buffers
                    .entry(tabla)
                    .or_default()
                    .push((columnas, valores));
                total_filas += 1;
            }

            // Tick de tiempo: forzar flush de lo acumulado
            _ = tick.tick() => {
                if total_filas > 0 {
                    flushear_todo(&cliente, &mut buffers, &mut total_filas).await;
                }
            }
        }

        // Si alcanzamos el límite de filas, flusheamos inmediatamente
        if total_filas >= 5000 {
            flushear_todo(&cliente, &mut buffers, &mut total_filas).await;
        }
    }
}

// ---------------------------------------------------------------------------
// Flush
// ---------------------------------------------------------------------------

/// Vacía todos los buffers acumulados hacia ClickHouse.
async fn flushear_todo(
    cliente: &Arc<ClienteClickHouse>,
    buffers: &mut HashMap<String, Vec<Fila>>,
    total_filas: &mut usize,
) {
    if *total_filas == 0 {
        return;
    }

    info!("Flusheando {} filas a ClickHouse...", *total_filas);

    // Extraer los buffers para no mantener el borrow.
    let tablas: Vec<String> = buffers.keys().cloned().collect();

    for tabla in &tablas {
        if let Some(filas) = buffers.remove(tabla) {
            if filas.is_empty() {
                continue;
            }
            if let Err(e) = flushear_tabla(cliente, tabla, &filas).await {
                error!("Error al flushear '{}': {}. Se pierden {} filas.", tabla, e, filas.len());
                FILAS_FALLIDAS.inc_by(filas.len() as u64);
                LOTES_FALLIDOS.inc();
            }
        }
    }

    *total_filas = 0;
}

/// Envía un lote de filas de una misma tabla a ClickHouse.
async fn flushear_tabla(
    cliente: &Arc<ClienteClickHouse>,
    tabla: &str,
    filas: &[Fila],
) -> Result<(), ErrorConsumidor> {
    if filas.is_empty() {
        return Ok(());
    }

    // Recolectar el conjunto completo de columnas (la primera fila define el
    // esquema; si hay más columnas en filas posteriores se agregan).
    let mut columnas_set: Vec<String> = Vec::new();
    let mut indices: std::collections::HashMap<&str, usize> = std::collections::HashMap::new();
    for (cols, _) in filas {
        for col in cols {
            if !indices.contains_key(col.as_str()) {
                indices.insert(col.as_str(), columnas_set.len());
                columnas_set.push(col.clone());
            }
        }
    }

    // Construir el lote como HashMap<columna → Vec<valor>>
    let mut lote: HashMap<String, Vec<String>> = HashMap::new();
    for col in &columnas_set {
        lote.insert(col.clone(), Vec::with_capacity(filas.len()));
    }

    for (cols, vals) in filas {
        // Inicializar todas las columnas con cadena vacía por defecto
        for valores_col in lote.values_mut() {
            valores_col.push(String::new());
        }
        let fila_idx = lote.values().next().map(|v| v.len()).unwrap_or(1) - 1;
        for (col, val) in cols.iter().zip(vals.iter()) {
            if let Some(valores_col) = lote.get_mut(col) {
                if fila_idx < valores_col.len() {
                    valores_col[fila_idx] = val.clone();
                }
            }
        }
    }

    let num_filas = cliente.insertar(tabla, &lote).await?;

    FILAS_INSERTADAS.inc_by(num_filas);
    LOTES_ENVIADOS.inc();

    info!("{} filas insertadas en '{}'", num_filas, tabla);

    Ok(())
}
