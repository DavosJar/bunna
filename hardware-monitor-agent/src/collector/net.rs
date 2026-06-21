use crate::types::{InterfaceMetrics, NetMetrics};
use anyhow::{Context, Result};
use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader};

pub struct NetCollector {
    /// Previous tick values: interface name -> (rx_bytes, tx_bytes)
    previous: HashMap<String, (u64, u64)>,
    /// Collection interval in milliseconds (used for rate calculation)
    interval_ms: u64,
}

impl NetCollector {
    pub fn new(interval_ms: u64) -> Self {
        Self {
            previous: HashMap::new(),
            interval_ms,
        }
    }

    pub fn collect(&mut self) -> Result<NetMetrics> {
        let file = File::open("/host/proc/net/dev")
            .or_else(|_| File::open("/proc/net/dev"))
            .context("Could not open /proc/net/dev")?;

        let reader = BufReader::new(file);
        let mut interfaces = Vec::new();
        let mut current: HashMap<String, (u64, u64)> = HashMap::new();

        let interval_secs = self.interval_ms as f64 / 1000.0;

        for line in reader.lines() {
            let line = line?;

            // Skip header lines (they have no colon)
            if !line.contains(':') {
                continue;
            }

            let parts: Vec<&str> = line.split_whitespace().collect();
            if parts.len() < 10 {
                continue;
            }

            // Interface name has a trailing colon, e.g. "eth0:" or "  lo:"
            let name = parts[0].trim_end_matches(':').to_string();
            let rx_bytes: u64 = parts[1].parse().unwrap_or(0);
            let tx_bytes: u64 = parts[9].parse().unwrap_or(0);

            // Calculate rates from previous values if available
            let (rx_per_sec, tx_per_sec) = if let Some(&(prev_rx, prev_tx)) = self.previous.get(&name)
            {
                if interval_secs > 0.0 {
                    let rx_delta = rx_bytes.saturating_sub(prev_rx);
                    let tx_delta = tx_bytes.saturating_sub(prev_tx);
                    (
                        rx_delta as f64 / interval_secs,
                        tx_delta as f64 / interval_secs,
                    )
                } else {
                    (0.0, 0.0)
                }
            } else {
                // First reading: no rate yet
                (0.0, 0.0)
            };

            // Store current values for next tick
            current.insert(name.clone(), (rx_bytes, tx_bytes));

            interfaces.push(InterfaceMetrics {
                name,
                received_bytes: rx_bytes,
                transmitted_bytes: tx_bytes,
                received_bytes_per_sec: rx_per_sec,
                transmitted_bytes_per_sec: tx_per_sec,
            });
        }

        // Swap previous state for next iteration
        self.previous = current;

        Ok(NetMetrics { interfaces })
    }
}
