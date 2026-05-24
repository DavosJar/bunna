package domain

import (
	"errors"
)

var (
	ErrDiagnosticoNoEncontrado      = errors.New("diagnostico no encontrado")
	ErrMuestraNoEncontrada          = errors.New("muestra no encontrada")
	ErrImageUrlRequerida            = errors.New("la url de la imagen es requerida")
	ErrConfianzaInvalida            = errors.New("la confianza debe estar entre 0 y 1")
	ErrNombreRequerido              = errors.New("el nombre es requerido")
	ErrMuestrasIdRequerido          = errors.New("el muestrasId es requerido")
	ErrTenantIdRequerido            = errors.New("el tenantID es requerido")
	ErrResultadoInferenciaRequerido = errors.New("el resultadoInferencia es requerido")
	ErrCreatedAtRequerido           = errors.New("el createdAt es requerido")
	ErrUpdatedAtRequerido           = errors.New("el updatedAt es requerido")
	ErrIdRequerido                  = errors.New("el id es requerido")
	ErrLoteIdRequerido              = errors.New("el loteID es requerido")
	ErrUbicacionRequerida           = errors.New("la ubicacion es requerida")
	ErrLatitudInvalida              = errors.New("la latitud es invalida")
	ErrLongitudInvalida             = errors.New("la longitud es invalida")
)
