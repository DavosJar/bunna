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

// MockMuestrasFacade is a mock implementation of facades.MuestrasFacade
type MockMuestrasFacade struct {
	mock.Mock
}

func (m *MockMuestrasFacade) Tomar(ctx context.Context, auth *application.AuthContext, loteID string, req dto.TomarMuestraRequest) (*dto.MuestraResponse, error) {
	args := m.Called(ctx, auth, loteID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MuestraResponse), args.Error(1)
}

func (m *MockMuestrasFacade) ListarPorLote(ctx context.Context, auth *application.AuthContext, loteID string) ([]dto.MuestraItemResponse, error) {
	args := m.Called(ctx, auth, loteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.MuestraItemResponse), args.Error(1)
}

func TestMuestraHandler_Tomar(t *testing.T) {
	tests := []struct {
		name           string
		loteID         string
		reqBody        interface{}
		mockSetup      func(m *MockMuestrasFacade)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Exito",
			loteID: "lote-1",
			reqBody: dto.TomarMuestraRequest{
				Latitud:  4.5,
				Longitud: -74.0,
			},
			mockSetup: func(m *MockMuestrasFacade) {
				m.On("Tomar", mock.Anything, mock.Anything, "lote-1", mock.AnythingOfType("dto.TomarMuestraRequest")).
					Return(&dto.MuestraResponse{ID: "muestra-1"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"id":"muestra-1"`,
		},
		{
			name:           "CuerpoInvalido",
			loteID:         "lote-1",
			reqBody:        "invalid json",
			mockSetup:      func(m *MockMuestrasFacade) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"cuerpo inválido"`,
		},
		{
			name:   "ErrorFacade",
			loteID: "lote-1",
			reqBody: dto.TomarMuestraRequest{
				Latitud:  4.5,
				Longitud: -74.0,
			},
			mockSetup: func(m *MockMuestrasFacade) {
				m.On("Tomar", mock.Anything, mock.Anything, "lote-1", mock.AnythingOfType("dto.TomarMuestraRequest")).
					Return(nil, errors.New("error del facade"))
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `"error":"error del facade"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockMuestrasFacade)
			tt.mockSetup(mockFacade)

			h := handler.NewMuestraHandler(mockFacade)
			router := setupTestRouter()
			router.POST("/lotes/:loteID/muestras", authMiddlewareMock(), h.Tomar)

			var body []byte
			if str, ok := tt.reqBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.reqBody)
			}

			req, _ := http.NewRequest(http.MethodPost, "/lotes/"+tt.loteID+"/muestras", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
			mockFacade.AssertExpectations(t)
		})
	}
}

func TestMuestraHandler_ListarPorLote(t *testing.T) {
	tests := []struct {
		name           string
		loteID         string
		mockSetup      func(m *MockMuestrasFacade)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Exito",
			loteID: "lote-1",
			mockSetup: func(m *MockMuestrasFacade) {
				m.On("ListarPorLote", mock.Anything, mock.Anything, "lote-1").
					Return([]dto.MuestraItemResponse{{ID: "muestra-1"}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":"muestra-1"`,
		},
		{
			name:   "ErrorFacade",
			loteID: "lote-1",
			mockSetup: func(m *MockMuestrasFacade) {
				m.On("ListarPorLote", mock.Anything, mock.Anything, "lote-1").
					Return(nil, errors.New("error al listar"))
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `"error":"error al listar"`,
		},
		{
			name:   "Vacio",
			loteID: "lote-1",
			mockSetup: func(m *MockMuestrasFacade) {
				m.On("ListarPorLote", mock.Anything, mock.Anything, "lote-1").
					Return([]dto.MuestraItemResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"data":[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockMuestrasFacade)
			tt.mockSetup(mockFacade)

			h := handler.NewMuestraHandler(mockFacade)
			router := setupTestRouter()
			router.GET("/lotes/:loteID/muestras", authMiddlewareMock(), h.ListarPorLote)

			req, _ := http.NewRequest(http.MethodGet, "/lotes/"+tt.loteID+"/muestras", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
			mockFacade.AssertExpectations(t)
		})
	}
}
