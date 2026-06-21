package middleware

import "time"

type Config struct {
    MaxDurationWarning time.Duration `json:"maxDurationWarning"`
    MaxDurationError   time.Duration `json:"maxDurationError"`
    TraceHeaderNames   []string       `json:"traceHeaderNames"`
}
