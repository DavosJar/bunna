package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	presentation "github.com/davosjar/bunna/services/identidad/internal/shared/presentation"
)

// LoginInput es el input del endpoint POST /api/v1/identidad/auth/login.
type LoginInput struct {
	RealIP string `header:"X-Real-IP"`
	Body   dto.LoginRequest
}

// LoginOutput es el output del endpoint POST /api/v1/identidad/auth/login.
type LoginOutput struct {
	Body presentation.ApiResponse[dto.LoginResponse]
}

// LoginHandler maneja el inicio de sesión.
// Sigue CON-PRES-002: no importa domain ni mapper.
type LoginHandler struct {
	facade facades.AuthFacade
}

// NewLoginHandler construye el handler con su facade.
func NewLoginHandler(facade facades.AuthFacade) *LoginHandler {
	return &LoginHandler{facade: facade}
}

// Register registra el endpoint en la API Huma.
func (h *LoginHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "login-usuario",
		Method:      http.MethodPost,
		Path:        "/api/v1/identidad/auth/login",
		Summary:     "Inicio de sesión",
		Description: "Autentica al usuario y devuelve tokens JWT de acceso y refresco.",
		Tags:        []string{"Autenticación"},
	}, h.handle)
}

func (h *LoginHandler) handle(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	resp, err := h.facade.Login(ctx, facades.ComandoLogin{
		Email:    input.Body.Correo,
		Password: input.Body.Password,
		IPOrigen: input.RealIP,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	// CON-PRES-005: links HATEOAS los construye el handler
	links := map[string]presentation.Link{
		"self": {
			Href:   "/api/v1/identidad/usuarios/" + resp.UsuarioID,
			Method: http.MethodGet,
		},
		"refresh": {
			Href:   "/api/v1/identidad/auth/refresh",
			Method: http.MethodPost,
		},
	}

	out := &LoginOutput{}
	out.Body = presentation.NewApiResponseWithLinks(dto.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		UsuarioID:    resp.UsuarioID,
		TenantID:     resp.TenantID,
		Rol:          resp.Rol,
	}, links)

	return out, nil
}
