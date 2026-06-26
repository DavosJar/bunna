package telemetry

// LogPayload is the unified telemetry event structure.
// Every event has common fields + exactly one type-specific section.
type LogPayload struct {
	LogType     string `json:"log_type"`     // "API" | "NEGOCIO" | "BD"
	Level       string `json:"level"`        // "INFO" | "WARN" | "ERROR"
	Timestamp   string `json:"timestamp"`    // ISO8601 UTC
	TraceID     string `json:"trace_id"`     // UUID v7
	SpanID      string `json:"span_id"`      // UUID v7
	ServiceName string `json:"service_name"` // "fincas"
	Environment string `json:"environment"`  // "dev" | "staging" | "production"

	API     *APIFields     `json:"api,omitempty"`
	Negocio *NegocioFields `json:"negocio,omitempty"`
	BD      *BDFields      `json:"bd,omitempty"`
}

// APIFields contains HTTP request/response metrics.
// Never includes sensitive request data (body, query params, auth headers).
type APIFields struct {
	Method        string  `json:"method"`
	Path          string  `json:"path"`
	StatusCode    int     `json:"status_code"`
	DurationMs    float64 `json:"duration_ms"`
	ClientIP      string  `json:"client_ip"`
	UserAgent     string  `json:"user_agent"`
	ContentLength int64   `json:"content_length"`
}

// NegocioFields contains business-level audit information.
// Never includes passwords, tokens, biometric data, or unnecessary PII.
type NegocioFields struct {
	UseCase           string         `json:"use_case"`
	Command           map[string]any `json:"command"`
	Result            string         `json:"result"`
	UserID            string         `json:"user_id"`
	Details           map[string]any `json:"details"`
	DurationUsecaseMs float64        `json:"duration_usecase_ms"`
}

// BDFields contains database query metrics.
// Never includes the raw query (only hash), actual data values, or results.
type BDFields struct {
	Operation     string  `json:"operation"`
	Table         string  `json:"table"`
	DurationMs    float64 `json:"duration_ms"`
	RowsAffected  int     `json:"rows_affected"`
	ErrorSQLState string  `json:"error_sql_state,omitempty"`
	QueryHash     string  `json:"query_hash"`
}
