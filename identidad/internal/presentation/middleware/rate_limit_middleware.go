package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitConfig struct {
	MaxRequests int
	Ventana     time.Duration
}

type ipCounter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	config   rateLimitConfig
}

func newIPCounter(cfg rateLimitConfig) *ipCounter {
	return &ipCounter{
		requests: make(map[string][]time.Time),
		config:   cfg,
	}
}

func (c *ipCounter) allow(ip string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-c.config.Ventana)

	times := c.requests[ip]
	var valid []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= c.config.MaxRequests {
		c.requests[ip] = valid
		return false
	}

	c.requests[ip] = append(valid, now)
	return true
}

// NuevoRateLimitMiddleware crea un middleware que limita requests por IP.
// Usa c.ClientIP() que respeta X-Forwarded-For cuando hay proxy configurado.
func NuevoRateLimitMiddleware(maxRequests int, ventana time.Duration) gin.HandlerFunc {
	counter := newIPCounter(rateLimitConfig{
		MaxRequests: maxRequests,
		Ventana:     ventana,
	})

	return func(c *gin.Context) {
		ip := c.GetHeader("X-Client-IP")
		if ip == "" {
			ip = c.ClientIP()
		}

		log.Printf("[RATE_LIMIT_DEBUG] ClientIP=%q | X-Client-IP=%q | RemoteAddr=%q | X-Forwarded-For=%q | X-Real-IP=%q | Path=%s",
			ip,
			c.GetHeader("X-Client-IP"),
			c.Request.RemoteAddr,
			c.GetHeader("X-Forwarded-For"),
			c.GetHeader("X-Real-IP"),
			c.Request.URL.Path,
		)

		if !counter.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Demasiadas solicitudes. Intente más tarde.",
				"codigo":  "RATE_LIMIT_EXCEDIDO",
			})
			return
		}
		c.Next()
	}
}
