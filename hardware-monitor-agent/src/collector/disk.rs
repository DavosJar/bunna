use crate::types::DiskMetrics;
use anyhow::{Result, Context};
use std::collections::HashMap;
use std::ffi::CString;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::mem::MaybeUninit;

/// Filesystem types considered real/physical storage.
const VALID_FILESYSTEMS: &[&str] = &["ext4", "xfs", "btrfs", "zfs", "fuseblk"];

/// Mount point prefixes to exclude (Docker bind mounts, etc.).
const MOUNT_EXCLUSIONS: &[&str] = &[
    "/etc/hostname",
    "/etc/hosts",
    "/etc/resolv.conf",
    "/usr/local/cargo/",
];

pub struct DiskCollector;

impl DiskCollector {
    pub fn new() -> Self {
        Self
    }

    pub fn collect(&self) -> Result<Vec<DiskMetrics>> {
        let file = File::open("/host/proc/mounts")
            .or_else(|_| File::open("/proc/mounts"))
            .context("Could not open /proc/mounts")?;

        let reader = BufReader::new(file);
        // Collect disk entries paired with their filesystem ID for later dedup.
        let mut entries: Vec<(DiskMetrics, u64)> = Vec::new();

        for line in reader.lines() {
            let line = match line {
                Ok(l) => l,
                Err(_) => continue,
            };

            // /proc/mounts format: device mount_point fstype options dump pass
            let parts: Vec<&str> = line.split_whitespace().collect();
            if parts.len() < 3 {
                continue;
            }

            let fstype = parts[2];
            if !VALID_FILESYSTEMS.contains(&fstype) {
                continue;
            }

            let mount_point = parts[1];

            // Skip excluded mount-point prefixes (Docker bind mounts, etc.)
            if MOUNT_EXCLUSIONS
                .iter()
                .any(|excl| mount_point.starts_with(excl))
            {
                continue;
            }

            // Call statvfs(2) to get filesystem statistics
            let mount_c = match CString::new(mount_point.as_bytes()) {
                Ok(c) => c,
                Err(_) => continue,
            };

            let mut stat = MaybeUninit::<libc::statvfs>::uninit();
            let ret = unsafe { libc::statvfs(mount_c.as_ptr(), stat.as_mut_ptr()) };
            if ret != 0 {
                // Skip mounts where statvfs fails (e.g., permissions, disconnected)
                continue;
            }
            let stat = unsafe { stat.assume_init() };

            let fsid = stat.f_fsid as u64;
            let block_size = stat.f_frsize as u64;
            let total_blocks = stat.f_blocks as u64;
            let free_blocks = stat.f_bfree as u64;
            let avail_blocks = stat.f_bavail as u64;

            let total_bytes = total_blocks.saturating_mul(block_size);
            let free_bytes = free_blocks.saturating_mul(block_size);
            let avail_bytes = avail_blocks.saturating_mul(block_size);

            let gb_div = 1024.0_f64 * 1024.0 * 1024.0;
            let total_gb = total_bytes as f64 / gb_div;
            let avail_gb = avail_bytes as f64 / gb_div;
            let used_bytes = total_bytes.saturating_sub(free_bytes);
            let used_gb = used_bytes as f64 / gb_div;

            let usage_percent = if total_bytes == 0 {
                0.0
            } else {
                (used_bytes as f64 / total_bytes as f64) * 100.0
            };

            entries.push((
                DiskMetrics {
                    mount: mount_point.to_string(),
                    total_gb,
                    used_gb,
                    available_gb: avail_gb,
                    usage_percent,
                },
                fsid,
            ));
        }

        // Deduplicate by filesystem ID (f_fsid), keeping the shortest mount point.
        // Docker bind mounts of the same ext4 filesystem all share the same f_fsid
        // so this eliminates duplicates while preferring the shortest path (e.g. /app > /app/target).
        let mut fs_map: HashMap<u64, DiskMetrics> = HashMap::new();
        for (disk, fsid) in entries {
            use std::collections::hash_map::Entry;
            match fs_map.entry(fsid) {
                Entry::Occupied(mut occupied) => {
                    if disk.mount.len() < occupied.get().mount.len() {
                        occupied.insert(disk);
                    }
                }
                Entry::Vacant(vacant) => {
                    vacant.insert(disk);
                }
            }
        }

        let mut disks: Vec<DiskMetrics> = fs_map.into_values().collect();
        // Sort by mount point for deterministic output.
        disks.sort_by(|a, b| a.mount.cmp(&b.mount));

        Ok(disks)
    }
}
