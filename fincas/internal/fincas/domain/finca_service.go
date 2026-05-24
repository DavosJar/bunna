package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type FincaService struct {
	fincaRepo FincaRepositorio
	loteRepo  LoteRepositorio
}

func NewFincaService(fincaRepo FincaRepositorio, loteRepo LoteRepositorio) *FincaService {
	return &FincaService{
		fincaRepo: fincaRepo,
		loteRepo:  loteRepo,
	}
}

func (s *FincaService) RegistrarLoteEnFinca(
	ctx context.Context,
	finca *Finca,
	nombre string,
	area float64,
	descripcion string,
) (*Lote, error) {
	lote, err := NewLote(finca.ID, nombre, area, descripcion)
	if err != nil {
		return nil, err
	}

	lote.ID = uuid.NewString()
	lote.CreatedAt = time.Now()
	lote.UpdatedAt = lote.CreatedAt

	if err := s.loteRepo.Crear(ctx, lote); err != nil {
		return nil, err
	}

	return lote, nil
}

func (s *FincaService) EliminarFincaConLotes(
	ctx context.Context,
	finca *Finca,
	confirmado bool,
) error {
	count, err := s.loteRepo.ContarPorFinca(ctx, finca.ID)
	if err != nil {
		return err
	}

	if count > 0 && !confirmado {
		return ErrFincaConLotes(count)
	}

	if err := s.loteRepo.EliminarPorFinca(ctx, finca.ID); err != nil {
		return err
	}

	if err := s.fincaRepo.Eliminar(ctx, finca.ID); err != nil {
		return err
	}

	return nil
}
