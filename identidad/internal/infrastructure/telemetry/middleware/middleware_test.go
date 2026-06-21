package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/buffer"
)

type testWriter struct {
    count int
    last []byte
}

func (w *testWriter) Write(event []byte, priority buffer.Prioridad) error {
    w.count++
    w.last = event
    return nil
}

func TestTelemetryMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)
    cfg := Config{}
    tw := &testWriter{}
    r := gin.New()
    r.Use(NewTelemetryMiddleware(tw, cfg))
    r.GET("/ping", func(c *gin.Context) {
        // read trace id from context and set header for verification
        if tid, ok := c.Request.Context().Value(ctxKey{}).(string); ok {
            c.Header("X-Trace-ID", tid)
        }
        c.String(http.StatusOK, "pong")
    })
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", w.Code)
    }
    if tw.count != 1 {
        t.Fatalf("expected writer count 1, got %d", tw.count)
    }
    if w.Header().Get("X-Trace-ID") == "" {
        t.Fatalf("trace id not set in response header")
    }
}
