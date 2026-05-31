package handler_test

import (
	"context"
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

// MockReportesFacade is a mock implementation of facades.ReportesFacade
type MockReportesFacade struct {
	mock.Mock
}

func (m *MockReportesFacade) GenerarPorLote(ctx context.Context, auth *application.AuthContext, loteID string) (*dto.ReporteLoteResponse, error) {
	args := m.Called(ctx, auth, loteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReporteLoteResponse), args.Error(1)
}

func TestReporteHandler_GenerarPorLote(t *testing.T) {
	tests := []struct {
		name           string
		loteID         string
		mockSetup      func(m *MockReportesFacade)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Exito",
			loteID: "lote-1",
			mockSetup: func(m *MockReportesFacade) {
				m.On("GenerarPorLote", mock.Anything, mock.Anything, "lote-1").
					Return(&dto.ReporteLoteResponse{ID: "lote-1", Nombre: "Buen estado"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"id":"lote-1"`,
		},
		{
			name:   "ErrorFacade",
			loteID: "lote-1",
			mockSetup: func(m *MockReportesFacade) {
				m.On("GenerarPorLote", mock.Anything, mock.Anything, "lote-1").
					Return(nil, errors.New("error al generar"))
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `"error":"error al generar"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFacade := new(MockReportesFacade)
			tt.mockSetup(mockFacade)

			h := handler.NewReporteHandler(mockFacade)
			router := setupTestRouter()
			router.GET("/lotes/:loteID/reportes", authMiddlewareMock(), h.GenerarPorLote)

			req, _ := http.NewRequest(http.MethodGet, "/lotes/"+tt.loteID+"/reportes", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
			mockFacade.AssertExpectations(t)
		})
	}
}
