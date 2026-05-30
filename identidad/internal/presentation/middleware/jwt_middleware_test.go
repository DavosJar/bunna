package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	"github.com/gin-gonic/gin"
)

// ── Mock TokenServicio ────────────────────────────────────────────────────────

type mockTokenServicio struct {
	claims *sesiones_domain.TokenClaims
	err    error
}

func (m *mockTokenServicio) GenerarAccessToken(usuarioID, sesionID string, tenantID string, rol string) (string, time.Time, error) {
	return "", time.Time{}, nil
}
func (m *mockTokenServicio) GenerarRefreshToken(usuarioID, sesionID string) (string, time.Time, error) {
	return "", time.Time{}, nil
}
func (m *mockTokenServicio) ValidarAccessToken(token string) (*sesiones_domain.TokenClaims, error) {
	return m.claims, m.err
}
func (m *mockTokenServicio) ValidarRefreshToken(token string) (*sesiones_domain.TokenClaims, error) {
	return nil, nil
}
func (m *mockTokenServicio) HashearToken(token string) string { return token }

// ── Helper ────────────────────────────────────────────────────────────────────

func setupMiddlewareRouter(tokenSvc sesiones_domain.TokenServicio) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.JWTMiddleware(tokenSvc))
	router.GET("/protegido", func(c *gin.Context) {
		usuarioID, _ := c.Get(middleware.ClaveUsuarioID)
		sesionID, _ := c.Get(middleware.ClaveSesionID)
		c.JSON(http.StatusOK, gin.H{
			"usuario_id": usuarioID,
			"sesion_id":  sesionID,
		})
	})
	return router
}

func doMiddlewareRequest(router *gin.Engine, authHeader string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ── Tests ─────────────────────────────────────────────────────────────────────

// AC-PRES-008: sin header Authorization → 401
func TestJWTMiddleware_SinHeader_Retorna401(t *testing.T) {
	router := setupMiddlewareRouter(&mockTokenServicio{
		claims: &sesiones_domain.TokenClaims{UsuarioID: "uid", SesionID: "sid"},
	})
	w := doMiddlewareRequest(router, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got %d", w.Code)
	}
}

// AC-PRES-008: formato inválido (sin Bearer) → 401
func TestJWTMiddleware_FormatoInvalido_Retorna401(t *testing.T) {
	router := setupMiddlewareRouter(&mockTokenServicio{
		claims: &sesiones_domain.TokenClaims{UsuarioID: "uid", SesionID: "sid"},
	})
	w := doMiddlewareRequest(router, "Token abc123")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got %d", w.Code)
	}
}

// AC-PRES-008: token inválido → 401
func TestJWTMiddleware_TokenInvalido_Retorna401(t *testing.T) {
	router := setupMiddlewareRouter(&mockTokenServicio{
		err: errors.New("token inválido"),
	})
	w := doMiddlewareRequest(router, "Bearer token-invalido")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got %d", w.Code)
	}
}

// AC-PRES-008: token válido → 200 y claims inyectados en contexto
func TestJWTMiddleware_TokenValido_Retorna200(t *testing.T) {
	router := setupMiddlewareRouter(&mockTokenServicio{
		claims: &sesiones_domain.TokenClaims{
			UsuarioID: "usuario-id-1",
			SesionID:  "sesion-id-1",
		},
	})
	w := doMiddlewareRequest(router, "Bearer token-valido")
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got %d", w.Code)
	}
	if !containsStr(w.Body.String(), "usuario-id-1") {
		t.Errorf("esperaba usuarioID en respuesta, got %s", w.Body.String())
	}
}

// AC-PRES-008: token válido inyecta sesionID en contexto
func TestJWTMiddleware_TokenValido_InyectaSesionID(t *testing.T) {
	router := setupMiddlewareRouter(&mockTokenServicio{
		claims: &sesiones_domain.TokenClaims{
			UsuarioID: "usuario-id-1",
			SesionID:  "sesion-id-1",
		},
	})
	w := doMiddlewareRequest(router, "Bearer token-valido")
	if !containsStr(w.Body.String(), "sesion-id-1") {
		t.Errorf("esperaba sesionID en respuesta, got %s", w.Body.String())
	}
}

// AC-PRES-008: solo token vacío después de Bearer → 401
func TestJWTMiddleware_BearerSinToken_Retorna401(t *testing.T) {
	router := setupMiddlewareRouter(&mockTokenServicio{
		err: errors.New("token inválido"),
	})
	w := doMiddlewareRequest(router, "Bearer ")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got %d", w.Code)
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
