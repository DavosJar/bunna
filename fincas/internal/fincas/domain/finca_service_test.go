package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFincaRepo struct {
	mock.Mock
}

func (m *mockFincaRepo) Crear(ctx context.Context, finca *Finca) error {
	return m.Called(ctx, finca).Error(0)
}

func (m *mockFincaRepo) ObtenerPorID(ctx context.Context, id string) (*Finca, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Finca), args.Error(1)
}

func (m *mockFincaRepo) ListarPorUsuario(ctx context.Context, usuarioID string) ([]Finca, error) {
	args := m.Called(ctx, usuarioID)
	return args.Get(0).([]Finca), args.Error(1)
}

func (m *mockFincaRepo) ListarTodas(ctx context.Context) ([]Finca, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Finca), args.Error(1)
}

func (m *mockFincaRepo) Actualizar(ctx context.Context, finca *Finca) error {
	return m.Called(ctx, finca).Error(0)
}

func (m *mockFincaRepo) Eliminar(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockLoteRepo struct {
	mock.Mock
}

func (m *mockLoteRepo) Crear(ctx context.Context, lote *Lote) error {
	return m.Called(ctx, lote).Error(0)
}

func (m *mockLoteRepo) ObtenerPorID(ctx context.Context, id string) (*Lote, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Lote), args.Error(1)
}

func (m *mockLoteRepo) ListarPorFinca(ctx context.Context, fincaID string) ([]Lote, error) {
	args := m.Called(ctx, fincaID)
	return args.Get(0).([]Lote), args.Error(1)
}

func (m *mockLoteRepo) ContarPorFinca(ctx context.Context, fincaID string) (int, error) {
	args := m.Called(ctx, fincaID)
	return args.Int(0), args.Error(1)
}

func (m *mockLoteRepo) Actualizar(ctx context.Context, lote *Lote) error {
	return m.Called(ctx, lote).Error(0)
}

func (m *mockLoteRepo) Eliminar(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockLoteRepo) EliminarPorFinca(ctx context.Context, fincaID string) error {
	return m.Called(ctx, fincaID).Error(0)
}

func TestFincaService_RegistrarLoteEnFinca_OK(t *testing.T) {
	finca, _ := NewFinca("Mi Finca", "ubicación", "", "user-1")
	finca.ID = uuid.NewString()

	loteRepo := new(mockLoteRepo)
	loteRepo.On("Crear", mock.Anything, mock.AnythingOfType("*domain.Lote")).Return(nil)

	fincaRepo := new(mockFincaRepo)

	svc := NewFincaService(fincaRepo, loteRepo)
	lote, err := svc.RegistrarLoteEnFinca(context.Background(), finca, "Lote 1", 10.5, "desc")

	assert.NoError(t, err)
	assert.NotNil(t, lote)
	assert.Equal(t, finca.ID, lote.FincaID)
	assert.Equal(t, "Lote 1", lote.Nombre)
	assert.Equal(t, 10.5, lote.Area)
	assert.NotEmpty(t, lote.ID)
	assert.False(t, lote.CreatedAt.IsZero())
	assert.Equal(t, lote.CreatedAt, lote.UpdatedAt)
	loteRepo.AssertExpectations(t)
}

func TestFincaService_RegistrarLoteEnFinca_ErrValidacion(t *testing.T) {
	finca, _ := NewFinca("Mi Finca", "ubicación", "", "user-1")
	finca.ID = uuid.NewString()

	loteRepo := new(mockLoteRepo)
	fincaRepo := new(mockFincaRepo)

	svc := NewFincaService(fincaRepo, loteRepo)
	lote, err := svc.RegistrarLoteEnFinca(context.Background(), finca, "ab", 10.5, "")

	assert.Nil(t, lote)
	assert.ErrorIs(t, err, ErrNombreLoteRequerido)
	loteRepo.AssertNotCalled(t, "Crear")
}

func TestFincaService_EliminarFincaConLotes_SinLotes(t *testing.T) {
	finca, _ := NewFinca("Mi Finca", "ubicación", "", "user-1")
	finca.ID = uuid.NewString()

	loteRepo := new(mockLoteRepo)
	loteRepo.On("ContarPorFinca", mock.Anything, finca.ID).Return(0, nil)
	loteRepo.On("EliminarPorFinca", mock.Anything, finca.ID).Return(nil)

	fincaRepo := new(mockFincaRepo)
	fincaRepo.On("Eliminar", mock.Anything, finca.ID).Return(nil)

	svc := NewFincaService(fincaRepo, loteRepo)
	err := svc.EliminarFincaConLotes(context.Background(), finca, false)

	assert.NoError(t, err)
	loteRepo.AssertExpectations(t)
	fincaRepo.AssertExpectations(t)
}

func TestFincaService_EliminarFincaConLotes_Confirmado(t *testing.T) {
	finca, _ := NewFinca("Mi Finca", "ubicación", "", "user-1")
	finca.ID = uuid.NewString()

	loteRepo := new(mockLoteRepo)
	loteRepo.On("ContarPorFinca", mock.Anything, finca.ID).Return(3, nil)
	loteRepo.On("EliminarPorFinca", mock.Anything, finca.ID).Return(nil)

	fincaRepo := new(mockFincaRepo)
	fincaRepo.On("Eliminar", mock.Anything, finca.ID).Return(nil)

	svc := NewFincaService(fincaRepo, loteRepo)
	err := svc.EliminarFincaConLotes(context.Background(), finca, true)

	assert.NoError(t, err)
	loteRepo.AssertExpectations(t)
	fincaRepo.AssertExpectations(t)
}

func TestFincaService_EliminarFincaConLotes_ErrConLotes(t *testing.T) {
	finca, _ := NewFinca("Mi Finca", "ubicación", "", "user-1")
	finca.ID = uuid.NewString()

	loteRepo := new(mockLoteRepo)
	loteRepo.On("ContarPorFinca", mock.Anything, finca.ID).Return(3, nil)

	fincaRepo := new(mockFincaRepo)

	svc := NewFincaService(fincaRepo, loteRepo)
	err := svc.EliminarFincaConLotes(context.Background(), finca, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "3")
	assert.Contains(t, err.Error(), "lote(s)")
	loteRepo.AssertNotCalled(t, "EliminarPorFinca")
	fincaRepo.AssertNotCalled(t, "Eliminar")
}

func TestFincaService_EliminarFincaConLotes_ErrRepo(t *testing.T) {
	finca, _ := NewFinca("Mi Finca", "ubicación", "", "user-1")
	finca.ID = uuid.NewString()

	expectedErr := errors.New("error de BD")

	loteRepo := new(mockLoteRepo)
	loteRepo.On("ContarPorFinca", mock.Anything, finca.ID).Return(0, expectedErr)

	fincaRepo := new(mockFincaRepo)

	svc := NewFincaService(fincaRepo, loteRepo)
	err := svc.EliminarFincaConLotes(context.Background(), finca, false)

	assert.ErrorIs(t, err, expectedErr)
}
