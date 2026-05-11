package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/handlers"
)

// ── Mock AuthFacade ───────────────────────────────────────────────────────────

type mockAuthFacade struct {
	registrarResp *facades.RespuestaRegistro
	registrarErr  error
	loginResp     *facades.RespuestaLogin
	loginErr      error
}

func (m *mockAuthFacade) Registrar(ctx context.Context, cmd facades.ComandoRegistro) (*facades.RespuestaRegistro, error) {
	return m.registrarResp, m.registrarErr
}

func (m *mockAuthFacade) Login(ctx context.Context, cmd facades.ComandoLogin) (*facades.RespuestaLogin, error) {
	return m.loginResp, m.loginErr
}

// ── Helper ────────────────────────────────────────────────────────────────────

func setupRouter(facade facades.AuthFacade) (*gin.Engine, huma.API) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := humagin.New(router, huma.DefaultConfig("Test API", "1.0.0"))

	handlers.RegisterHealthHandler(api)
	handlers.NewRegisterHandler(facade).Register(api)
	handlers.NewLoginHandler(facade).Register(api)

	return router, api
}

func doRequest(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// ── Health ────────────────────────────────────────────────────────────────────

// AC-PRES-008 (health): GET /health responde 200
func TestHealthHandler_Responde200(t *testing.T) {
	router, _ := setupRouter(&mockAuthFacade{})
	w := doRequest(router, http.MethodGet, "/health", "")
	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got %d", w.Code)
	}
}

func TestHealthHandler_CuerpoContieneStatusOk(t *testing.T) {
	router, _ := setupRouter(&mockAuthFacade{})
	w := doRequest(router, http.MethodGet, "/health", "")
	if !strings.Contains(w.Body.String(), "ok") {
		t.Errorf("esperaba 'ok' en body, got %s", w.Body.String())
	}
}

// ── Register ──────────────────────────────────────────────────────────────────

// AC-PRES-003: request válido → 201 con ApiResponse
func TestRegisterHandler_Exitoso(t *testing.T) {
	facade := &mockAuthFacade{
		registrarResp: &facades.RespuestaRegistro{
			UsuarioID: "uid-1",
			Correo:    "juan@correo.com",
			Estado:    "NO_VERIFICADO",
		},
	}
	router, _ := setupRouter(facade)

	body := `{"nombre":"Juan","apellido":"Pérez","correo":"juan@correo.com","password":"secreto123"}`
	w := doRequest(router, http.MethodPost, "/api/v1/auth/register", body)

	if w.Code != http.StatusCreated {
		t.Errorf("esperaba 201, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("esperaba campo 'data' en respuesta")
	}
	if data["usuario_id"] != "uid-1" {
		t.Errorf("usuario_id incorrecto: %v", data["usuario_id"])
	}
}

// AC-PRES-003: respuesta incluye _links HATEOAS
func TestRegisterHandler_IncluyeLinks(t *testing.T) {
	facade := &mockAuthFacade{
		registrarResp: &facades.RespuestaRegistro{
			UsuarioID: "uid-1",
			Correo:    "juan@correo.com",
			Estado:    "NO_VERIFICADO",
		},
	}
	router, _ := setupRouter(facade)

	body := `{"nombre":"Juan","apellido":"Pérez","correo":"juan@correo.com","password":"secreto123"}`
	w := doRequest(router, http.MethodPost, "/api/v1/auth/register", body)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if _, ok := resp["_links"]; !ok {
		t.Error("esperaba campo '_links' en respuesta")
	}
}

// AC-PRES-004: error del facade → 400
func TestRegisterHandler_ErrorFacade_Retorna400(t *testing.T) {
	facade := &mockAuthFacade{
		registrarErr: errors.New("correo ya registrado"),
	}
	router, _ := setupRouter(facade)

	body := `{"nombre":"Juan","apellido":"Pérez","correo":"juan@correo.com","password":"secreto123"}`
	w := doRequest(router, http.MethodPost, "/api/v1/auth/register", body)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("esperaba 400, got %d", w.Code)
	}
}

// ── Login ─────────────────────────────────────────────────────────────────────

// AC-PRES-005: credenciales correctas → 200 con tokens
func TestLoginHandler_Exitoso(t *testing.T) {
	facade := &mockAuthFacade{
		loginResp: &facades.RespuestaLogin{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    900,
			TokenType:    "Bearer",
			UsuarioID:    "uid-1",
			SesionID:     "sid-1",
		},
	}
	router, _ := setupRouter(facade)

	body := `{"correo":"juan@correo.com","password":"secreto123"}`
	w := doRequest(router, http.MethodPost, "/api/v1/auth/login", body)

	if w.Code != http.StatusOK {
		t.Errorf("esperaba 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("esperaba campo 'data' en respuesta")
	}
	if data["access_token"] == "" {
		t.Error("esperaba access_token no vacío")
	}
	if data["token_type"] != "Bearer" {
		t.Errorf("esperaba token_type=Bearer, got %v", data["token_type"])
	}
}

// AC-PRES-005: respuesta incluye _links self y refresh
func TestLoginHandler_IncluyeLinks(t *testing.T) {
	facade := &mockAuthFacade{
		loginResp: &facades.RespuestaLogin{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    900,
			TokenType:    "Bearer",
			UsuarioID:    "uid-1",
			SesionID:     "sid-1",
		},
	}
	router, _ := setupRouter(facade)

	body := `{"correo":"juan@correo.com","password":"secreto123"}`
	w := doRequest(router, http.MethodPost, "/api/v1/auth/login", body)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	links, ok := resp["_links"].(map[string]interface{})
	if !ok {
		t.Fatal("esperaba campo '_links' en respuesta")
	}
	if _, ok := links["refresh"]; !ok {
		t.Error("esperaba link 'refresh' en _links")
	}
}

// AC-PRES-006: credenciales incorrectas → 401
func TestLoginHandler_ErrorFacade_Retorna401(t *testing.T) {
	facade := &mockAuthFacade{
		loginErr: errors.New("credenciales inválidas"),
	}
	router, _ := setupRouter(facade)

	body := `{"correo":"juan@correo.com","password":"incorrecta"}`
	w := doRequest(router, http.MethodPost, "/api/v1/auth/login", body)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, got %d", w.Code)
	}
}
