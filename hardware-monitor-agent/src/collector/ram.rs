use crate::types::RamMetrics;
use anyhow::{Result, Context};
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::collections::HashMap;

pub struct RamCollector;

impl RamCollector {
    pub fn new() -> Self {
        Self
    }

    pub fn collect(&self) -> Result<RamMetrics> {
        let file = File::open("/host/proc/meminfo")
            .or_else(|_| File::open("/proc/meminfo"))
            .context("Could not open /proc/meminfo")?;
        
        let reader = BufReader::new(file);
        let mut stats = HashMap::new();

        for line in reader.lines() {
            let line = line?;
            let parts: Vec<&str> = line.split(':').collect();
            if parts.len() < 2 { continue; }
            
            let key = parts[0].trim();
            // El valor suele venir como "16384 kB"
            let val_str = parts[1].trim().split_whitespace().next().unwrap_or("0");
            let val: u64 = val_str.parse().unwrap_or(0);
            
            stats.insert(key.to_string(), val);
        }

        let total_kb = *stats.get("MemTotal").unwrap_or(&0);
        let available_kb = *stats.get("MemAvailable").unwrap_or(&0);
        let used_kb = total_kb.saturating_sub(available_kb);

        let total_mb = total_kb / 1024;
        let available_mb = available_kb / 1024;
        let used_mb = used_kb / 1024;

        let usage_percent = if total_mb == 0 {
            0.0
        } else {
            (used_mb as f64 / total_mb as f64) * 100.0
        };

        Ok(RamMetrics {
            total_mb,
            used_mb,
            available_mb,
            usage_percent,
        })
    }
}
