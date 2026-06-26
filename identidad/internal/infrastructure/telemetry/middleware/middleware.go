package middleware

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry"
	"github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
	"github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/idgenerator"
	"github.com/gin-gonic/gin"
)

type ctxKey struct{}

type spanCtxKey struct{}

// GetTraceIDFromCtx extracts the trace_id previously stored by the middleware.
// Returns empty string if not found.
func GetTraceIDFromCtx(ctx context.Context) string {
	if v := ctx.Value(ctxKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetSpanIDFromCtx extracts the span_id previously stored by the middleware.
// Returns empty string if not found.
func GetSpanIDFromCtx(ctx context.Context) string {
	if v := ctx.Value(spanCtxKey{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

var idGenerator = idgenerator.NewUUIDv7Generator()

func NewTelemetryMiddleware(writer buffer.BufferWriter, cfg Config) gin.HandlerFunc {
	// ensure default header names if not provided
	if len(cfg.TraceHeaderNames) == 0 {
		cfg.TraceHeaderNames = []string{"X-Trace-ID", "X-Cloud-Trace-Context"}
	}
	return func(c *gin.Context) {
		start := time.Now()
		// extract trace id from headers
		var traceID string
		for _, h := range cfg.TraceHeaderNames {
			if v := c.GetHeader(h); v != "" {
				// take first value before possible delimiters
				traceID = strings.Split(v, ",")[0]
				break
			}
		}
		if traceID == "" {
			if id, err := idGenerator.NextID(context.Background()); err != nil {
				traceID = ""
			} else {
				traceID = id
			}
		}
		// generate spanID; declare variable first
		var spanID string
		if id, err := idGenerator.NextID(context.Background()); err != nil {
			spanID = ""
		} else {
			spanID = id
		}
		// client IP extraction
		clientIP := c.GetHeader("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.GetHeader("X-Real-IP")
		}
		if clientIP == "" {
			clientIP = c.ClientIP()
		}
		// Anonymize: mask last octet for IPv4
		clientIP = anonymizeIP(clientIP)
		// store trace_id and span_id in context
		reqCtx := context.WithValue(c.Request.Context(), ctxKey{}, traceID)
		reqCtx = context.WithValue(reqCtx, spanCtxKey{}, spanID)
		c.Request = c.Request.WithContext(reqCtx)
		// proceed
		c.Next()
		// after request
		latency := time.Since(start)
		status := c.Writer.Status()
		path := c.Request.URL.Path
		ua := c.Request.UserAgent()
		if len(ua) > 128 {
			ua = ua[:128]
		}
		// build payload
		payload := telemetry.LogPayload{
			LogType:     "API",
			Level:       determineLevel(status, latency, cfg),
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			TraceID:     traceID,
			SpanID:      spanID,
			ServiceName: "identidad",
			Environment: "dev",
			API: &telemetry.APIFields{
				Method:        c.Request.Method,
				Path:          path,
				StatusCode:    status,
				DurationMs:    float64(latency.Microseconds()) / 1000.0,
				ClientIP:      clientIP,
				UserAgent:     ua,
				ContentLength: c.Request.ContentLength,
			},
		}
		data, _ := json.Marshal(payload)
		_ = writer.Write(data, buffer.Alta)
	}
}

func determineLevel(status int, latency time.Duration, cfg Config) string {
	if status >= 500 {
		return "ERROR"
	}
	if status >= 400 {
		// 401/403 without brute force detection → INFO (spec says WARN only
		// when brute force pattern detected, which requires external state)
		if status == http.StatusUnauthorized || status == http.StatusForbidden {
			return "INFO"
		}
		return "WARN"
	}
	if latency > cfg.MaxDurationError {
		return "ERROR"
	}
	if latency > cfg.MaxDurationWarning {
		return "WARN"
	}
	return "INFO"
}

func anonymizeIP(ip string) string {
	if parsed := net.ParseIP(ip); parsed != nil {
		if net4 := parsed.To4(); net4 != nil {
			return net4.Mask(net.CIDRMask(24, 32)).String()
		}
		// IPv6: mask first 64 bits
		return ip[:strings.LastIndex(ip, ":")+1] + ":xxxx"
	}
	return ip
}
