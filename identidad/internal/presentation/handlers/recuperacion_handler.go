package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	presentation "github.com/davosjar/bunna/services/identidad/shared/presentation"
)

// ── Solicitar Recuperación ─────────────────────────────────────────────────────

type SolicitarRecuperacionInput struct {
	Body   dto.SolicitarRecuperacionRequest
	RealIP string `header:"X-Real-IP"`
}

type SolicitarRecuperacionOutput struct {
	Body presentation.ApiResponse[dto.SolicitarRecuperacionResponse]
}

type SolicitarRecuperacionHandler struct {
	facade facades.RecuperacionFacade
}

func NewSolicitarRecuperacionHandler(facade facades.RecuperacionFacade) *SolicitarRecuperacionHandler {
	return &SolicitarRecuperacionHandler{facade: facade}
}

func (h *SolicitarRecuperacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "solicitar-recuperacion",
		Method:      http.MethodPost,
		Path:        "/api/v1/recuperacion/solicitar",
		Summary:     "Solicitar recuperación de contraseña",
		Description: "Envía un enlace de recuperación al correo electrónico proporcionado.",
		Tags:        []string{"Recuperación"},
	}, h.handle)
}

func (h *SolicitarRecuperacionHandler) handle(ctx context.Context, input *SolicitarRecuperacionInput) (*SolicitarRecuperacionOutput, error) {
	resp, err := h.facade.SolicitarRecuperacion(ctx, facades.ComandoSolicitarRecuperacion{
		Email:    input.Body.Correo,
		IPOrigen: input.RealIP,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &SolicitarRecuperacionOutput{}
	out.Body = presentation.NewApiResponse(dto.SolicitarRecuperacionResponse{Mensaje: resp.Mensaje})
	return out, nil
}

// ── Validar Token de Recuperación ──────────────────────────────────────────────

type ValidarTokenRecuperacionInput struct {
	Body dto.ValidarTokenRecuperacionRequest
}

type ValidarTokenRecuperacionOutput struct {
	Body presentation.ApiResponse[dto.ValidarTokenRecuperacionResponse]
}

type ValidarTokenRecuperacionHandler struct {
	facade facades.RecuperacionFacade
}

func NewValidarTokenRecuperacionHandler(facade facades.RecuperacionFacade) *ValidarTokenRecuperacionHandler {
	return &ValidarTokenRecuperacionHandler{facade: facade}
}

func (h *ValidarTokenRecuperacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "validar-token-recuperacion",
		Method:      http.MethodPost,
		Path:        "/api/v1/recuperacion/validar",
		Summary:     "Validar token de recuperación",
		Description: "Valida si un token de recuperación es válido y devuelve el ID del usuario asociado.",
		Tags:        []string{"Recuperación"},
	}, h.handle)
}

func (h *ValidarTokenRecuperacionHandler) handle(ctx context.Context, input *ValidarTokenRecuperacionInput) (*ValidarTokenRecuperacionOutput, error) {
	resp, err := h.facade.ValidarTokenRecuperacion(ctx, facades.ComandoValidarTokenRecuperacion{
		Token: input.Body.Token,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ValidarTokenRecuperacionOutput{}
	out.Body = presentation.NewApiResponse(dto.ValidarTokenRecuperacionResponse{
		UsuarioID: resp.UsuarioID,
		Valido:    resp.Valido,
	})
	return out, nil
}

// ── Confirmar Restablecimiento ─────────────────────────────────────────────────

type ConfirmarRecuperacionInput struct {
	Body dto.ConfirmarRecuperacionRequest
}

type ConfirmarRecuperacionOutput struct {
	Body presentation.ApiResponse[dto.ConfirmarRecuperacionResponse]
}

type ConfirmarRecuperacionHandler struct {
	facade facades.RecuperacionFacade
}

func NewConfirmarRecuperacionHandler(facade facades.RecuperacionFacade) *ConfirmarRecuperacionHandler {
	return &ConfirmarRecuperacionHandler{facade: facade}
}

func (h *ConfirmarRecuperacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "confirmar-recuperacion",
		Method:      http.MethodPost,
		Path:        "/api/v1/recuperacion/confirmar",
		Summary:     "Confirmar restablecimiento de contraseña",
		Description: "Restablece la contraseña usando el token de recuperación.",
		Tags:        []string{"Recuperación"},
	}, h.handle)
}

func (h *ConfirmarRecuperacionHandler) handle(ctx context.Context, input *ConfirmarRecuperacionInput) (*ConfirmarRecuperacionOutput, error) {
	resp, err := h.facade.ConfirmarRecuperacion(ctx, facades.ComandoConfirmarRecuperacion{
		Token:         input.Body.Token,
		NuevaPassword: input.Body.NuevaPassword,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ConfirmarRecuperacionOutput{}
	out.Body = presentation.NewApiResponse(dto.ConfirmarRecuperacionResponse{Mensaje: resp.Mensaje})
	return out, nil
}
