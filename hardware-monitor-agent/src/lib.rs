//! hardware-monitor-agent library

pub mod config;
pub mod types;
pub mod collector;
pub mod aggregator;
pub mod guardian;
pub mod kafka_publisher;

pub use config::Config;
pub use types::{Snapshot, AlertEvent, Thresholds};
pub use aggregator::Aggregator;
pub use guardian::Guardian;
pub use kafka_publisher::KafkaPublisher;
