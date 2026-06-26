package middleware

import (
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/infrastructure/telemetry"
)

type Config struct {
	MaxDurationWarning time.Duration `json:"maxDurationWarning"`
	MaxDurationError   time.Duration `json:"maxDurationError"`
	TraceHeaderNames   []string      `json:"traceHeaderNames"`
	Service            telemetry.ServiceInfo
}
