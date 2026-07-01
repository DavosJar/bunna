package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDiagnosticosFacade is a mock implementation of facades.DiagnosticosFacade
type MockDiagnosticosFacade struct {
	mock.Mock
}

func (m *MockDiagnosticosFacade) SolicitarManual(ctx context.Context, auth *application.AuthContext, muestraID string, req dto.SolicitarDiagnosticoManualRequest) (*dto.SolicitudDiagnosticoResponse, error) {
	args := m.Called(ctx, auth, muestraID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SolicitudDiagnosticoResponse), args.Error(1)
}

func (m *MockDiagnosticosFacade) Aceptar(ctx context.Context, auth *application.AuthContext, diagnosticoID string) (*dto.EstadoCambioResponse, error) {
	args := m.Called(ctx, auth, diagnosticoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EstadoCambioResponse), args.Error(1)
}

func (m *MockDiagnosticosFacade) Rechazar(ctx context.Context, auth *application.AuthContext, diagnosticoID string, req dto.RechazarDiagnosticoRequest) (*dto.EstadoCambioResponse, error) {
	args := m.Called(ctx, auth, diagnosticoID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EstadoCambioResponse), args.Error(1)
}
func (m *MockDiagnosticosFacade) GuardarResultadoManual(ctx context.Context, auth *application.AuthContext, muestraID string, req dto.GuardarResultadoManualRequest) (*dto.DiagnosticoResponse, error) {
	args := m.Called(ctx, auth, muestraID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DiagnosticoResponse), args.Error(1)
}

func TestDiagnosticoHandler_SolicitarManual(t *testing.T) {
	tests := []struct {
		name           string
		muestraID      string
		reqBody        interface{}
		mockSetup      func(m *MockDiagnosticosFacade)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Exito",
			muestraID: "muestra-1",
			reqBody: dto.SolicitarDiagnosticoManualRequest{
				ImageURL: "http://example.com/img.jpg",
			},
			mockSetup: func(m *MockDiagnosticosFacade) {
				m.On("SolicitarManual", mock.Anything, mock.Anything, "muestra-1", mock.AnythingOfType("dto.SolicitarDiagnosticoManualRequest")).
					Return(&dto.SolicitudDiagnosticoResponse{SolicitudID: "diag-1"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"solicitudID":"diag-1"`,
		},
		{
			name:           "CuerpoInvalido",
			muestraID:      "muestra-1",
			reqBody:        "invalid json",
			mockSetup:      func(m *MockDiagnosticosFacade) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"cuerpo inválido"`,
		},
		{
			name:      "ErrorFacade",
			muestraID: "muestra-1",
			reqBody: dto.SolicitarDiagnosticoManualRequest{
				ImageURL: "http://example.com/img.jpg",
			},
			mockSetup: func(m *MockDiagnosticosFacade) {
				m.On("SolicitarManual", mock.Anything, mock.Anything, "muestra-1", mock.AnythingOfType("dto.SolicitarDiagnosticoManualRequest")).
					Return(nil, errors.New("error del facade"))
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `"error":"error del facade"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockDiagnosticosFacade)
			tt.mockSetup(mockFacade)

			h := handler.NewDiagnosticoHandler(mockFacade)
			router := setupTestRouter()
			router.POST("/muestras/:muestraID/diagnosticos/manual", authMiddlewareMock(), h.SolicitarManual)

			var body []byte
			if str, ok := tt.reqBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.reqBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/muestras/"+tt.muestraID+"/diagnosticos/manual", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
			mockFacade.AssertExpectations(t)
		})
	}
}

func TestDiagnosticoHandler_Aceptar(t *testing.T) {
	tests := []struct {
		name           string
		diagnosticoID  string
		mockSetup      func(m *MockDiagnosticosFacade)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Exito",
			diagnosticoID: "diag-1",
			mockSetup: func(m *MockDiagnosticosFacade) {
				m.On("Aceptar", mock.Anything, mock.Anything, "diag-1").
					Return(&dto.EstadoCambioResponse{ID: "diag-1", Estado: "aceptado"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"estado":"aceptado"`,
		},
		{
			name:          "ErrorFacade",
			diagnosticoID: "diag-1",
			mockSetup: func(m *MockDiagnosticosFacade) {
				m.On("Aceptar", mock.Anything, mock.Anything, "diag-1").
					Return(nil, errors.New("error al aceptar"))
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `"error":"error al aceptar"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockDiagnosticosFacade)
			tt.mockSetup(mockFacade)

			h := handler.NewDiagnosticoHandler(mockFacade)
			router := setupTestRouter()
			router.POST("/diagnosticos/:id/aceptar", authMiddlewareMock(), h.Aceptar)

			req, _ := http.NewRequest(http.MethodPost, "/diagnosticos/"+tt.diagnosticoID+"/aceptar", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
			mockFacade.AssertExpectations(t)
		})
	}
}

func TestDiagnosticoHandler_Rechazar(t *testing.T) {
	tests := []struct {
		name           string
		diagnosticoID  string
		reqBody        interface{}
		mockSetup      func(m *MockDiagnosticosFacade)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:          "Exito",
			diagnosticoID: "diag-1",
			reqBody: dto.RechazarDiagnosticoRequest{
				Motivo: "incorrecto",
			},
			mockSetup: func(m *MockDiagnosticosFacade) {
				m.On("Rechazar", mock.Anything, mock.Anything, "diag-1", mock.AnythingOfType("dto.RechazarDiagnosticoRequest")).
					Return(&dto.EstadoCambioResponse{ID: "diag-1", Estado: "rechazado"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"estado":"rechazado"`,
		},
		{
			name:           "CuerpoInvalido",
			diagnosticoID:  "diag-1",
			reqBody:        "invalid json",
			mockSetup:      func(m *MockDiagnosticosFacade) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"cuerpo inválido"`,
		},
		{
			name:          "ErrorFacade",
			diagnosticoID: "diag-1",
			reqBody: dto.RechazarDiagnosticoRequest{
				Motivo: "incorrecto",
			},
			mockSetup: func(m *MockDiagnosticosFacade) {
				m.On("Rechazar", mock.Anything, mock.Anything, "diag-1", mock.AnythingOfType("dto.RechazarDiagnosticoRequest")).
					Return(nil, errors.New("error al rechazar"))
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `"error":"error al rechazar"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockDiagnosticosFacade)
			tt.mockSetup(mockFacade)

			h := handler.NewDiagnosticoHandler(mockFacade)
			router := setupTestRouter()
			router.POST("/diagnosticos/:id/rechazar", authMiddlewareMock(), h.Rechazar)

			var body []byte
			if str, ok := tt.reqBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.reqBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/diagnosticos/"+tt.diagnosticoID+"/rechazar", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
			mockFacade.AssertExpectations(t)
		})
	}
}
