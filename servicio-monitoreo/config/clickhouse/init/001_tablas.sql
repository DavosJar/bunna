CREATE DATABASE IF NOT EXISTS bunna_monitoreo;

USE bunna_monitoreo;

CREATE TABLE IF NOT EXISTS telemetria_api (
    _v UInt8 DEFAULT 1,
    log_type LowCardinality(String),
    level LowCardinality(String),
    timestamp DateTime64(3),
    trace_id String,
    span_id String,
    service_name LowCardinality(String),
    environment LowCardinality(String),
    method LowCardinality(String),
    path String,
    status_code UInt16,
    duration_ms Float64,
    client_ip String,
    user_agent String,
    content_length Int64,
    _insertado_en DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (toStartOfHour(timestamp), method, status_code)
TTL timestamp + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS telemetria_negocio (
    _v UInt8 DEFAULT 1,
    log_type LowCardinality(String),
    level LowCardinality(String),
    timestamp DateTime64(3),
    trace_id String,
    span_id String,
    service_name LowCardinality(String),
    environment LowCardinality(String),
    use_case LowCardinality(String),
    command String,
    result LowCardinality(String),
    user_id String,
    details String,
    duration_usecase_ms Float64,
    _insertado_en DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (toDate(timestamp), use_case, result)
TTL timestamp + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS telemetria_bd (
    _v UInt8 DEFAULT 1,
    log_type LowCardinality(String),
    level LowCardinality(String),
    timestamp DateTime64(3),
    trace_id String,
    span_id String,
    service_name LowCardinality(String),
    environment LowCardinality(String),
    operation LowCardinality(String),
    table LowCardinality(String),
    duration_ms Float64,
    rows_affected UInt32,
    error_sql_state String,
    query_hash String,
    _insertado_en DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (toDate(timestamp), table, operation)
TTL timestamp + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS hardware_metricas (
    _v UInt8 DEFAULT 1,
    node_id String,
    timestamp DateTime64(3),
    interval_ms UInt32,
    cpu_usage_percent Float64,
    cpu_cores UInt16,
    ram_total_mb UInt64,
    ram_used_mb UInt64,
    ram_available_mb UInt64,
    ram_usage_percent Float64,
    disco_mount String,
    disco_total_gb Float64,
    disco_used_gb Float64,
    disco_available_gb Float64,
    disco_usage_percent Float64,
    interfaz_name String,
    interfaz_received_bytes UInt64,
    interfaz_transmitted_bytes UInt64,
    interfaz_received_bytes_per_sec Float64,
    interfaz_transmitted_bytes_per_sec Float64,
    contenedor_id String,
    contenedor_cpu_shares UInt64,
    contenedor_memory_limit_mb UInt64,
    _insertado_en DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (toStartOfHour(timestamp), node_id)
TTL timestamp + INTERVAL 7 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS hardware_alertas (
    _v UInt8 DEFAULT 1,
    node_id String,
    timestamp DateTime64(3),
    metric LowCardinality(String),
    severity LowCardinality(String),
    value Float64,
    threshold Float64,
    message String,
    previous_state String,
    event_type LowCardinality(String),
    _insertado_en DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (toDate(timestamp), node_id, severity)
TTL timestamp + INTERVAL 90 DAY
SETTINGS index_granularity = 8192;
