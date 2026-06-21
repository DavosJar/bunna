#[cfg(test)]
mod tests {
    use super::*;
    use std::env;

    #[test]
    fn test_config_from_env() {
        // Set env vars
        env::set_var("NODE_ID", "test-node");
        env::set_var("KAFKA_BROKERS", "localhost:9092");
        env::set_var("INTERVAL_MS", "5000");
        env::set_var("CPU_WARN_PERCENT", "70.0");
        env::set_var("CPU_CRITICAL_PERCENT", "90.0");
        env::set_var("RAM_WARN_PERCENT", "75.0");
        env::set_var("RAM_CRITICAL_PERCENT", "95.0");
        env::set_var("DISK_WARN_PERCENT", "80.0");
        env::set_var("DISK_CRITICAL_PERCENT", "98.0");

        let cfg = Config::from_env().expect("Config should load");
        assert_eq!(cfg.node_id, "test-node");
        assert_eq!(cfg.kafka_brokers, "localhost:9092");
        assert_eq!(cfg.interval_ms, 5000);
        assert!((cfg.cpu_warn - 70.0).abs() < f64::EPSILON);
        assert!((cfg.cpu_critical - 90.0).abs() < f64::EPSILON);
        assert!((cfg.ram_warn - 75.0).abs() < f64::EPSILON);
        assert!((cfg.ram_critical - 95.0).abs() < f64::EPSILON);
        assert!((cfg.disk_warn - 80.0).abs() < f64::EPSILON);
        assert!((cfg.disk_critical - 98.0).abs() < f64::EPSILON);
    }
}
