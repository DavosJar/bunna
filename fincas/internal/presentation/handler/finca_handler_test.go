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
	"github.com/davosjar/bunna/services/fincas/internal/presentation/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFincasFacade is a mock implementation of facades.FincasFacade
type MockFincasFacade struct {
	mock.Mock
}

func (m *MockFincasFacade) Registrar(ctx context.Context, auth *application.AuthContext, req dto.RegistrarFincaRequest) (*dto.FincaResponse, error) {
	args := m.Called(ctx, auth, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FincaResponse), args.Error(1)
}

func (m *MockFincasFacade) Desactivar(ctx context.Context, auth *application.AuthContext, fincaID string, req dto.DesactivarFincaRequest) (*dto.EstadoCambioResponse, error) {
	args := m.Called(ctx, auth, fincaID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.EstadoCambioResponse), args.Error(1)
}

func (m *MockFincasFacade) Listar(ctx context.Context, auth *application.AuthContext) ([]dto.FincaResponse, error) {
	args := m.Called(ctx, auth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.FincaResponse), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func authMiddlewareMock() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mock auth context for tests
		c.Set(middleware.ClaveAuthContext, &application.AuthContext{
			UsuarioID: "user-123",
			TenantID:  "tenant-1",
			Permisos:  []string{"ADMIN"},
		})
		c.Next()
	}
}

func TestFincaHandler_Registrar(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(m *MockFincasFacade)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Success",
			requestBody: dto.RegistrarFincaRequest{
				Nombre:      "Mi Finca",
				Ubicacion:   "Antioquia",
				Descripcion: "Finca cafetera",
			},
			setupMock: func(m *MockFincasFacade) {
				m.On("Registrar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), mock.AnythingOfType("dto.RegistrarFincaRequest")).
					Return(&dto.FincaResponse{
						ID:          "finca-1",
						Nombre:      "Mi Finca",
						Ubicacion:   "Antioquia",
						Descripcion: "Finca cafetera",
						Estado:      "ACTIVO",
					}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "Invalid Body",
			requestBody:    "invalid json",
			setupMock:      func(m *MockFincasFacade) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Facade Error",
			requestBody: dto.RegistrarFincaRequest{
				Nombre: "Finca Error",
			},
			setupMock: func(m *MockFincasFacade) {
				m.On("Registrar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), mock.AnythingOfType("dto.RegistrarFincaRequest")).
					Return(nil, errors.New("business error")).Once()
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockFincasFacade)
			tt.setupMock(mockFacade)

			h := handler.NewFincaHandler(mockFacade)
			router := setupTestRouter()
			router.Use(authMiddlewareMock())
			router.POST("/fincas", h.Registrar)

			var bodyBytes []byte
			if b, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(b)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/fincas", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockFacade.AssertExpectations(t)
		})
	}
}

func TestFincaHandler_Desactivar(t *testing.T) {
	tests := []struct {
		name           string
		fincaID        string
		requestBody    interface{}
		setupMock      func(m *MockFincasFacade)
		expectedStatus int
	}{
		{
			name:    "Success",
			fincaID: "finca-1",
			requestBody: dto.DesactivarFincaRequest{
				Confirmar: true,
			},
			setupMock: func(m *MockFincasFacade) {
				m.On("Desactivar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), "finca-1", mock.AnythingOfType("dto.DesactivarFincaRequest")).
					Return(&dto.EstadoCambioResponse{
						ID:     "finca-1",
						Estado: "INACTIVO",
					}, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Body",
			fincaID:        "finca-1",
			requestBody:    "invalid json",
			setupMock:      func(m *MockFincasFacade) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Facade Error",
			fincaID: "finca-1",
			requestBody: dto.DesactivarFincaRequest{
				Confirmar: true,
			},
			setupMock: func(m *MockFincasFacade) {
				m.On("Desactivar", mock.Anything, mock.AnythingOfType("*application.AuthContext"), "finca-1", mock.AnythingOfType("dto.DesactivarFincaRequest")).
					Return(nil, errors.New("cannot deactivate")).Once()
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockFincasFacade)
			tt.setupMock(mockFacade)

			h := handler.NewFincaHandler(mockFacade)
			router := setupTestRouter()
			router.Use(authMiddlewareMock())
			router.POST("/fincas/:id/desactivar", h.Desactivar)

			var bodyBytes []byte
			if b, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(b)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/fincas/"+tt.fincaID+"/desactivar", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockFacade.AssertExpectations(t)
		})
	}
}
