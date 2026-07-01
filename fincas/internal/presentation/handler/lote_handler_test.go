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

// MockLotesFacade is a mock implementation of facades.LotesFacade
type MockLotesFacade struct {
	mock.Mock
}

func (m *MockLotesFacade) Agregar(ctx context.Context, auth *application.AuthContext, fincaID string, req dto.AgregarLoteRequest) (*dto.LoteResponse, error) {
	args := m.Called(ctx, auth, fincaID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoteResponse), args.Error(1)
}

func (m *MockLotesFacade) Eliminar(ctx context.Context, auth *application.AuthContext, loteID string) (*dto.EstadoCambioResponse, error) {
	args := m.Called(ctx, auth, loteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EstadoCambioResponse), args.Error(1)
}

func (m *MockLotesFacade) Listar(ctx context.Context, auth *application.AuthContext, fincaID string) ([]dto.LoteResponse, error) {
	args := m.Called(ctx, auth, fincaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.LoteResponse), args.Error(1)
}

func TestLoteHandler_Agregar(t *testing.T) {
	tests := []struct {
		name           string
		fincaID        string
		requestBody    interface{}
		setupMock      func(m *MockLotesFacade)
		expectedStatus int
	}{
		{
			name:    "Success",
			fincaID: "finca-1",
			requestBody: dto.AgregarLoteRequest{
				Nombre:      "Lote A",
				Area:        10.5,
				Descripcion: "Lote de prueba",
			},
			setupMock: func(m *MockLotesFacade) {
				m.On("Agregar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), "finca-1", mock.AnythingOfType("dto.AgregarLoteRequest")).
					Return(&dto.LoteResponse{
						ID:          "lote-1",
						Nombre:      "Lote A",
						Area:        10.5,
						Descripcion: "Lote de prueba",
						Estado:      "ACTIVO",
					}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Body",
			fincaID:        "finca-1",
			requestBody:    "invalid json",
			setupMock:      func(m *MockLotesFacade) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Facade Error",
			fincaID: "finca-1",
			requestBody: dto.AgregarLoteRequest{
				Nombre: "Lote Error",
			},
			setupMock: func(m *MockLotesFacade) {
				m.On("Agregar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), "finca-1", mock.AnythingOfType("dto.AgregarLoteRequest")).
					Return(nil, errors.New("business error")).Once()
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockLotesFacade)
			tt.setupMock(mockFacade)

			h := handler.NewLoteHandler(mockFacade)
			router := setupTestRouter() // Reuses setupTestRouter from finca_handler_test.go
			router.Use(authMiddlewareMock()) // Reuses authMiddlewareMock from finca_handler_test.go
			router.POST("/fincas/:id/lotes", h.Agregar)

			var bodyBytes []byte
			if b, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(b)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/fincas/"+tt.fincaID+"/lotes", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockFacade.AssertExpectations(t)
		})
	}
}

func TestLoteHandler_Eliminar(t *testing.T) {
	tests := []struct {
		name           string
		loteID         string
		setupMock      func(m *MockLotesFacade)
		expectedStatus int
	}{
		{
			name:   "Success",
			loteID: "lote-1",
			setupMock: func(m *MockLotesFacade) {
				m.On("Eliminar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), "lote-1").
					Return(&dto.EstadoCambioResponse{
						ID:     "lote-1",
						Estado: "ELIMINADO",
					}, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Facade Error",
			loteID: "lote-1",
			setupMock: func(m *MockLotesFacade) {
				m.On("Eliminar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), "lote-1").
					Return(nil, errors.New("cannot delete")).Once()
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockLotesFacade)
			tt.setupMock(mockFacade)

			h := handler.NewLoteHandler(mockFacade)
			router := setupTestRouter()
			router.Use(authMiddlewareMock())
			router.POST("/lotes/:id/eliminar", h.Eliminar)

			req, _ := http.NewRequest(http.MethodPost, "/lotes/"+tt.loteID+"/eliminar", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockFacade.AssertExpectations(t)
		})
	}
}
