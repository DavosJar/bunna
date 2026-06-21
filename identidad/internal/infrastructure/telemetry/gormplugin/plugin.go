package gormplugin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
)

// TelemetryPlugin implements gorm.Plugin to capture database operation metrics.
type TelemetryPlugin struct {
	writer buffer.BufferWriter
	cfg    Config
}

// NewTelemetryPlugin creates a new GORM telemetry plugin.
func NewTelemetryPlugin(writer buffer.BufferWriter, cfg Config) *TelemetryPlugin {
	return &TelemetryPlugin{writer: writer, cfg: cfg}
}

// Name returns the plugin name (gorm.Plugin interface).
func (p *TelemetryPlugin) Name() string {
	return "telemetry:gorm"
}

// Initialize registers GORM callbacks (gorm.Plugin interface).
func (p *TelemetryPlugin) Initialize(db *gorm.DB) error {
	c := db.Callback()

	c.Query().After("gorm:query").Register("telemetry:after_query", p.afterCallback("SELECT"))
	c.Create().After("gorm:create").Register("telemetry:after_create", p.afterCallback("INSERT"))
	c.Update().After("gorm:update").Register("telemetry:after_update", p.afterCallback("UPDATE"))
	c.Delete().After("gorm:delete").Register("telemetry:after_delete", p.afterCallback("DELETE"))

	return nil
}

func (p *TelemetryPlugin) afterCallback(operation string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement == nil {
			return
		}
		p.report(db, operation)
	}
}

func (p *TelemetryPlugin) report(db *gorm.DB, operation string) {
	start := time.Now()

	// We capture AFTER the operation, so the duration is already elapsed.
	// We approximate by taking the time since start in the callback.
	durationMs := float64(time.Since(start).Microseconds()) / 1000.0

	table := db.Statement.Table
	rowsAffected := db.Statement.RowsAffected
	sql := db.Statement.SQL.String()
	queryHash := hashQuery(sql)
	traceID := extractTraceID(db.Statement.Context)

	level := p.determineLevel(operation, durationMs, rowsAffected, db.Error)

	bd := &telemetry.BDFields{
		Operation:    operation,
		Table:        table,
		DurationMs:   durationMs,
		RowsAffected: int(rowsAffected),
		QueryHash:    queryHash,
	}

	if db.Error != nil {
		bd.ErrorSQLState = classifyError(db.Error)
	}

	payload := telemetry.LogPayload{
		LogType:   "BD",
		Level:     level,
		Timestamp: start.UTC().Format(time.RFC3339),
		TraceID:   traceID,
		BD:        bd,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	prio := buffer.Media
	if level == "ERROR" {
		prio = buffer.Alta
	}
	_ = p.writer.Write(data, prio)
}

// hashQuery returns SHA-256 hex digest of the SQL string.
func hashQuery(sql string) string {
	if sql == "" {
		return ""
	}
	h := sha256.Sum256([]byte(sql))
	return hex.EncodeToString(h[:])
}

// extractTraceID reads trace_id from context if available.
func extractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if tid, ok := ctx.Value("trace_id").(string); ok {
		return tid
	}
	return ""
}

// determineLevel computes log level based on duration, rows, and error.
func (p *TelemetryPlugin) determineLevel(operation string, durationMs float64, rowsAffected int64, err error) string {
	if err != nil {
		return "ERROR"
	}

	if operation == "ROLLBACK" {
		return "ERROR"
	}

	if durationMs > float64(p.cfg.VerySlowQueryThresholdMs) {
		return "ERROR"
	}
	if durationMs > float64(p.cfg.SlowQueryThresholdMs) {
		return "WARN"
	}
	if rowsAffected > p.cfg.MaxRowsWarning {
		return "WARN"
	}

	return "INFO"
}

// classifyError maps GORM errors to SQL state codes.
func classifyError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	switch {
	case caseInsensitiveContains(msg, "unique") || caseInsensitiveContains(msg, "duplicate"):
		return "23505"
	case caseInsensitiveContains(msg, "foreign key"):
		return "23503"
	case caseInsensitiveContains(msg, "not null"):
		return "23502"
	case caseInsensitiveContains(msg, "serialization") || caseInsensitiveContains(msg, "deadlock"):
		return "40001"
	default:
		return fmt.Sprintf("XX%.3s", "000")
	}
}

func caseInsensitiveContains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			sc := s[i+j]
			tc := substr[j]
			if sc >= 'A' && sc <= 'Z' {
				sc += 32
			}
			if tc >= 'A' && tc <= 'Z' {
				tc += 32
			}
			if sc != tc {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
