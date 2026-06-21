package gormplugin

// Config holds thresholds for database query telemetry.
type Config struct {
	// SlowQueryThresholdMs: queries taking longer than this are flagged WARN.
	SlowQueryThresholdMs int64 `json:"slow_query_threshold_ms"`
	// VerySlowQueryThresholdMs: queries taking longer than this are ERROR.
	VerySlowQueryThresholdMs int64 `json:"very_slow_query_threshold_ms"`
	// MaxRowsWarning: row counts above this trigger WARN.
	MaxRowsWarning int64 `json:"max_rows_warning"`
	// LongTxThresholdMs: transactions longer than this trigger WARN.
	LongTxThresholdMs int64 `json:"long_tx_threshold_ms"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		SlowQueryThresholdMs:     200,
		VerySlowQueryThresholdMs: 1000,
		MaxRowsWarning:           1000,
		LongTxThresholdMs:        5000,
	}
}
