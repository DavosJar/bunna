package generarreporteporlote

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = application.PermisoGenerarReporte

// UseCase orquesta la generación del reporte de clorosis para un lote.
// Es un caso de uso de solo lectura: no publica eventos ni usa Unit of Work.
type UseCase struct {
	loteRepo       fincasdomain.LoteRepositorio
	muestraRepo    diagnosticodomain.MuestraRepositorio
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio
}

// NewUseCase crea un nuevo caso de uso con las dependencias inyectadas.
func NewUseCase(
	loteRepo fincasdomain.LoteRepositorio,
	muestraRepo diagnosticodomain.MuestraRepositorio,
	diagnosticoRepo diagnosticodomain.DiagnosticoRepositorio,
) *UseCase {
	return &UseCase{
		loteRepo:        loteRepo,
		muestraRepo:     muestraRepo,
		diagnosticoRepo: diagnosticoRepo,
	}
}

// Ejecutar valida permisos, carga el lote, consulta muestras y diagnósticos
// aceptados, calcula métricas de clorosis y retorna el reporte completo.
func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	// 1. Validar permisos
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	// 2. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 3. Cargar lote
	lote, err := uc.loteRepo.ObtenerPorID(ctx, cmd.LoteID)
	if err != nil {
		if errors.Is(err, fincasdomain.ErrLoteNoEncontrado) {
			return nil, application.ErrNotFound
		}
		return nil, err
	}

	// 4. Validar tenencia: si el TenantID no coincide → ErrNotFound
	if auth.TenantID != "" && lote.TenantID() != auth.TenantID {
		return nil, application.ErrNotFound
	}

	// 5. Buscar muestras del lote
	filtrosMuestra := []shared.CriterioFiltro{
		{Campo: "loteID", Operador: "=", Valor: cmd.LoteID},
	}
	if auth.TenantID != "" {
		filtrosMuestra = append(filtrosMuestra, shared.CriterioFiltro{
			Campo: "tenantID", Operador: "=", Valor: auth.TenantID,
		})
	}

	muestras, err := uc.muestraRepo.Buscar(ctx, diagnosticodomain.EspecificacionMuestra{
		Filtros: filtrosMuestra,
	}, shared.Paginacion{})
	if err != nil {
		return nil, err
	}

	// 6-7. Recorrer muestras, buscar diagnósticos y calcular métricas
	var (
		conClorosis  int
		sinClorosis  int
		pendientes   int
		zonas        []ZonaAfectada
		muestrasRepo []MuestraReporte
	)

	for _, m := range muestras {
		// Buscar diagnósticos de la muestra
		filtrosDiag := []shared.CriterioFiltro{
			{Campo: "muestraID", Operador: "=", Valor: m.ID()},
		}

		diagnosticos, err := uc.diagnosticoRepo.Buscar(ctx, diagnosticodomain.EspecificacionDiagnostico{
			Filtros: filtrosDiag,
		}, shared.Paginacion{})
		if err != nil {
			return nil, err
		}

		// Identificar el primer diagnóstico ACEPTADO y el primero disponible
		var diagAceptado *diagnosticodomain.Diagnostico
		primerDiagIdx := -1
		for i, d := range diagnosticos {
			if primerDiagIdx == -1 {
				primerDiagIdx = i
			}
			if d.Estado() == diagnosticodomain.EstadoDiagnosticoAceptado {
				diagAceptado = &diagnosticos[i]
				break
			}
		}

		// Métricas: solo cuentan diagnósticos ACEPTADOS
		if diagAceptado != nil {
			res := diagAceptado.ResultadoInferencia()
			if res != nil && res.TieneClorosis() {
				conClorosis++
				ubic := m.Ubicacion()
				zonas = append(zonas, ZonaAfectada{
					Latitud:  ubic.Latitud(),
					Longitud: ubic.Longitud(),
					RadioMts: 2.0,
				})
			} else {
				sinClorosis++
			}
		} else {
			pendientes++
		}

		// Armar MuestraReporte con la información de diagnóstico disponible
		ubicMuestra := m.Ubicacion()
		mr := MuestraReporte{
			ID:       m.ID(),
			Latitud:  ubicMuestra.Latitud(),
			Longitud: ubicMuestra.Longitud(),
		}

		// Usar el ACEPTADO si existe, si no el primero disponible
		var diagParaReporte *diagnosticodomain.Diagnostico
		if diagAceptado != nil {
			diagParaReporte = diagAceptado
		} else if primerDiagIdx != -1 {
			diagParaReporte = &diagnosticos[primerDiagIdx]
		}

		if diagParaReporte != nil {
			mr.DiagnosticoID = diagParaReporte.ID()
			mr.EstadoDiagnostico = string(diagParaReporte.Estado())
			if res := diagParaReporte.ResultadoInferencia(); res != nil {
				mr.ImageURL = res.ImageUrl()
				mr.ImageBase64 = res.ImageBase64()
				tc := res.TieneClorosis()
				mr.TieneClorosis = &tc
				cf := res.Confianza()
				mr.Confianza = &cf
			}
		}

		muestrasRepo = append(muestrasRepo, mr)
	}

	// 8. Calcular métricas de área
	areaAfectada := float64(len(zonas)) * math.Pi * 2.0 * 2.0
	var porcentajeAfectado float64
	if lote.Area() > 0 {
		// El área del lote está en hectáreas → convertir a m²
		porcentajeAfectado = (areaAfectada / (lote.Area() * 10000)) * 100
	}

	// 9. Ensamblar salida
	salida := Salida{
		ID:        lote.ID(),
		Nombre:    lote.Nombre(),
		AreaTotal: lote.Area(),
		Estado:    string(lote.Estado()),
		Muestras:  muestrasRepo,
		Zonas:     zonas,
		Metricas: Metricas{
			TotalMuestras:        len(muestras),
			ConClorosis:          conClorosis,
			SinClorosis:          sinClorosis,
			Pendientes:           pendientes,
			AreaAfectadaEstimada: areaAfectada,
			PorcentajeAfectado:   porcentajeAfectado,
		},
		GeneradoEn: time.Now(),
	}

	return &salida, nil
}
