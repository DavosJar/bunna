use crate::types::CpuMetrics;
use anyhow::{Result, Context};
use std::fs::File;
use std::io::{BufRead, BufReader};

pub struct CpuCollector {
    last_total: u64,
    last_idle: u64,
}

impl CpuCollector {
    pub fn new() -> Self {
        Self {
            last_total: 0,
            last_idle: 0,
        }
    }

    pub fn collect(&mut self) -> Result<CpuMetrics> {
        let (total, idle) = self.read_values()?;

        if self.last_total == 0 {
            // Primera lectura, inicializamos y devolvemos 0
            self.last_total = total;
            self.last_idle = idle;
            return Ok(CpuMetrics {
                usage_percent: 0.0,
                cores: self.count_cores().unwrap_or(1),
            });
        }

        let delta_total = total - self.last_total;
        let delta_idle = idle - self.last_idle;

        self.last_total = total;
        self.last_idle = idle;

        let usage_percent = if delta_total == 0 {
            0.0
        } else {
            100.0 * (1.0 - (delta_idle as f64 / delta_total as f64))
        };

        Ok(CpuMetrics {
            usage_percent,
            cores: self.count_cores().unwrap_or(1),
        })
    }

    fn read_values(&self) -> Result<(u64, u64)> {
        let file = File::open("/host/proc/stat")
            .or_else(|_| File::open("/proc/stat")) // Fallback para dev local
            .context("Could not open /proc/stat")?;
        
        let reader = BufReader::new(file);
        let first_line = reader.lines().next().context("Empty /proc/stat")??;

        // Formato: cpu  user nice system idle iowait irq softirq steal guest guest_nice
        let parts: Vec<&str> = first_line.split_whitespace().collect();
        if parts.is_empty() || parts[0] != "cpu" {
            return Err(anyhow::anyhow!("Invalid format in /proc/stat"));
        }

        let user: u64 = parts[1].parse()?;
        let nice: u64 = parts[2].parse()?;
        let system: u64 = parts[3].parse()?;
        let idle: u64 = parts[4].parse()?;
        let iowait: u64 = parts[5].parse()?;
        let irq: u64 = parts[6].parse()?;
        let softirq: u64 = parts[7].parse()?;
        let steal: u64 = parts[8].parse()?;

        let total = user + nice + system + idle + iowait + irq + softirq + steal;
        
        Ok((total, idle))
    }

    fn count_cores(&self) -> Result<u32> {
        let file = File::open("/host/proc/stat")
            .or_else(|_| File::open("/proc/stat"))?;
        let reader = BufReader::new(file);
        
        let mut count = 0;
        for line in reader.lines() {
            let line = line?;
            if line.starts_with("cpu") && line != "cpu " && !line.starts_with("cpu  ") {
                // Buscamos líneas que empiecen con cpu seguido de un número (ej: cpu0, cpu1)
                let first_word = line.split_whitespace().next().unwrap_or("");
                if first_word.len() > 3 && first_word[3..].chars().all(|c| c.is_ascii_digit()) {
                    count += 1;
                }
            }
        }
        Ok(if count == 0 { 1 } else { count })
    }
}
