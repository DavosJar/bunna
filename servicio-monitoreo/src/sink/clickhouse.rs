//! Cliente HTTP directo hacia ClickHouse.
//!
//! Construye consultas INSERT ... FORMAT TabSeparated y las envía
//! mediante POST al puerto HTTP de ClickHouse (por defecto 8123).

use std::collections::HashMap;

use tracing::{debug, error};

use crate::consumidor::error::ErrorConsumidor;

/// Cliente HTTP para ClickHouse.
///
/// Se conecta a la interfaz HTTP de ClickHouse usando autenticación
/// básica con usuario y contraseña.
#[derive(Debug, Clone)]
#[allow(non_snake_case)]
pub struct ClienteClickHouse {
    /// URL base (ej. "http://localhost:8123").
    pub url: String,
    pub usuario: String,
    pub contrasena: String,
    pub base_datos: String,

    /// Cliente HTTP reutilizable.
    #[allow(dead_code)]
    cliente_http: reqwest::Client,
}

impl ClienteClickHouse {
    /// Crea una nueva conexión a ClickHouse.
    ///
    /// # Argumentos
    ///
    /// * `url` – URL del servidor ClickHouse (ej. "http://clickhouse:8123").
    /// * `usuario` – Nombre de usuario.
    /// * `contrasena` – Contraseña.
    /// * `base_datos` – Base de datos por defecto.
    pub fn nuevo(url: String, usuario: String, contrasena: String, base_datos: String) -> Self {
        let cliente_http = reqwest::Client::builder()
            .user_agent("servicio-monitoreo/0.1.0")
            .build()
            .expect("Error al construir cliente HTTP");

        Self {
            url,
            usuario,
            contrasena,
            base_datos,
            cliente_http,
        }
    }

    /// Inserta un lote de filas en una tabla de ClickHouse.
    ///
    /// Construye una consulta `INSERT INTO tabla (col1, col2, ...) FORMAT TabSeparated`
    /// y envía los datos crudos en el cuerpo del POST.
    ///
    /// # Argumentos
    ///
    /// * `tabla`   – Nombre de la tabla (incluyendo base de datos si aplica).
    /// * `lote`    – Mapa de columna → lista de valores.  Todas las listas
    ///               deben tener la misma longitud (una entrada por fila).
    pub async fn insertar(
        &self,
        tabla: &str,
        lote: &HashMap<String, Vec<String>>,
    ) -> Result<u64, ErrorConsumidor> {
        if lote.is_empty() {
            return Ok(0);
        }

        // Determinar el número de filas y el orden de columnas.
        let num_filas: usize = lote.values().next().map(|v| v.len()).unwrap_or(0);
        if num_filas == 0 {
            return Ok(0);
        }

        let columnas: Vec<&String> = lote.keys().collect();
        // Usamos siempre el mismo orden para consistencia.
        let mut columnas_ordenadas: Vec<&String> = columnas.clone();
        columnas_ordenadas.sort();

        // Construir el cuerpo en formato TabSeparated.
        let mut cuerpo = String::with_capacity(num_filas * 128);
        for fila_idx in 0..num_filas {
            for (col_idx, col) in columnas_ordenadas.iter().enumerate() {
                if col_idx > 0 {
                    cuerpo.push('\t');
                }
                if let Some(valores) = lote.get(*col) {
                    if fila_idx < valores.len() {
                        // Escapar caracteres especiales para TabSeparated.
                        let valor_escapado = escapar_tab_separated(&valores[fila_idx]);
                        cuerpo.push_str(&valor_escapado);
                    }
                }
            }
            cuerpo.push('\n');
        }

        let columnas_str: Vec<String> = columnas_ordenadas
            .iter()
            .map(|c| (*c).clone())
            .collect();

        let query = format!(
            "INSERT INTO {}.{} ({}) FORMAT TabSeparated",
            self.base_datos,
            tabla,
            columnas_str.join(", ")
        );

        debug!(
            "Insertando {} filas en {}.{} ({} columnas)",
            num_filas, self.base_datos, tabla, columnas_ordenadas.len()
        );

        let url_con_query = format!("{}/?query={}", self.url, urlencoding(&query));

        let response = self
            .cliente_http
            .post(&url_con_query)
            .basic_auth(&self.usuario, Some(&self.contrasena))
            .header("Content-Type", "text/plain; charset=UTF-8")
            .body(cuerpo)
            .send()
            .await
            .map_err(|e| ErrorConsumidor::ClickHouse(format!("Error HTTP: {}", e)))?;

        let status = response.status();
        let body = response
            .text()
            .await
            .unwrap_or_else(|_| "<sin cuerpo>".to_string());

        if !status.is_success() {
            error!(
                "ClickHouse respondió con {}: {}",
                status,
                &body[..body.len().min(512)]
            );
            return Err(ErrorConsumidor::ClickHouse(format!(
                "HTTP {}: {}",
                status,
                &body[..body.len().min(512)]
            )));
        }

        Ok(num_filas as u64)
    }
}

/// Escapa un valor para el formato TabSeparated de ClickHouse.
///
/// Reemplaza:
/// - `\` → `\\`
/// - `\t` → `\t` (tabulador literal)
/// - `\n` → `\n` (nueva línea literal)
fn escapar_tab_separated(valor: &str) -> String {
    let mut resultado = String::with_capacity(valor.len());
    for ch in valor.chars() {
        match ch {
            '\\' => resultado.push_str("\\\\"),
            '\t' => resultado.push_str("\\t"),
            '\n' => resultado.push_str("\\n"),
            c => resultado.push(c),
        }
    }
    resultado
}

/// Codifica URL mínima para el parámetro `query`.
fn urlencoding(query: &str) -> String {
    let mut encoded = String::with_capacity(query.len() * 2);
    for byte in query.bytes() {
        match byte {
            b'A'..=b'Z' | b'a'..=b'z' | b'0'..=b'9' | b'_' | b'.' | b'-' | b'~' => {
                encoded.push(byte as char);
            }
            b' ' => encoded.push_str("%20"),
            _ => {
                encoded.push_str(&format!("%{:02X}", byte));
            }
        }
    }
    encoded
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_escapar_tab_separated() {
        assert_eq!(escapar_tab_separated("hola"), "hola");
        assert_eq!(escapar_tab_separated("a\tb"), "a\\tb");
        assert_eq!(escapar_tab_separated("a\nb"), "a\\nb");
        assert_eq!(escapar_tab_separated("a\\b"), "a\\\\b");
        assert_eq!(escapar_tab_separated("a\t\n\\"), "a\\t\\n\\\\");
    }

    #[test]
    fn test_urlencoding() {
        let result = urlencoding("INSERT INTO test VALUES");
        assert!(result.contains("INSERT"));
        assert!(result.contains("%20"));
    }
}
