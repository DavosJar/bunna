//! Módulo de escritura (sink) hacia ClickHouse.
//!
//! Expone `ClienteClickHouse` para la comunicación HTTP directa y
//! `BatchInserter` para el encolamiento y flush por lotes.

pub mod clickhouse;
pub mod insertador;

pub use clickhouse::ClienteClickHouse;
pub use insertador::BatchInserter;
