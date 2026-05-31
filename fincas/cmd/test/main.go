package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/aceptardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/agregarlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/desactivarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/eliminarlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/generarreporteporlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/listarmuestrasporlote"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/rechazardiagnostico"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarfinca"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/registrarinferencia"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/solicitardiagnosticomanual"
	"github.com/davosjar/bunna/services/fincas/internal/application/usecases/tomarmuestra"
	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
	fincasdomain "github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

// ──────────────────────────────────────────────
// Test helpers
// ──────────────────────────────────────────────

type testResult struct {
	name string
	ok   bool
	msg  string
}

var passed, failed int
var results []testResult

func runTest(name string, fn func() error) {
	err := fn()
	if err != nil {
		failed++
		results = append(results, testResult{name, false, err.Error()})
		fmt.Printf("  ❌ %s: %v\n", name, err)
	} else {
		passed++
		results = append(results, testResult{name, true, ""})
		fmt.Printf("  ✅ %s\n", name)
	}
}

func runTestExpected(name string, fn func() error, expectError bool) {
	err := fn()
	if expectError && err != nil {
		passed++
		results = append(results, testResult{name, true, fmt.Sprintf("(expected error: %v)", err)})
		fmt.Printf("  ✅ %s (expected error: %v)\n", name, err)
	} else if !expectError && err == nil {
		passed++
		results = append(results, testResult{name, true, ""})
		fmt.Printf("  ✅ %s\n", name)
	} else if expectError && err == nil {
		failed++
		results = append(results, testResult{name, false, "expected error but got nil"})
		fmt.Printf("  ❌ %s: expected error but got nil\n", name)
	} else {
		failed++
		results = append(results, testResult{name, false, err.Error()})
		fmt.Printf("  ❌ %s: unexpected error: %v\n", name, err)
	}
}

func resetMocks(
	fincaRepo *mockFincaRepo,
	loteRepo *mockLoteRepo,
	muestraRepo *mockMuestraRepo,
	diagnosticoRepo *mockDiagnosticoRepo,
	candidatoRepo *mockCandidatoRepo,
) {
	fincaRepo.mu.Lock()
	fincaRepo.stores = make(map[string]*fincasdomain.Finca)
	fincaRepo.mu.Unlock()

	loteRepo.mu.Lock()
	loteRepo.stores = make(map[string]*fincasdomain.Lote)
	loteRepo.mu.Unlock()

	muestraRepo.mu.Lock()
	muestraRepo.stores = make(map[string]*diagnosticodomain.Muestra)
	muestraRepo.mu.Unlock()

	diagnosticoRepo.mu.Lock()
	diagnosticoRepo.stores = make(map[string]*diagnosticodomain.Diagnostico)
	diagnosticoRepo.mu.Unlock()

	candidatoRepo.mu.Lock()
	candidatoRepo.stores = make(map[string]*diagnosticodomain.CandidatoReentrenamiento)
	candidatoRepo.mu.Unlock()
}

// ──────────────────────────────────────────────
// Mock FincaRepositorio
// ──────────────────────────────────────────────

type mockFincaRepo struct {
	mu       sync.RWMutex
	stores   map[string]*fincasdomain.Finca
	loteRepo *mockLoteRepo
}

func (m *mockFincaRepo) Crear(_ context.Context, finca *fincasdomain.Finca) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[finca.ID()] = finca
	return nil
}

func (m *mockFincaRepo) ObtenerPorID(_ context.Context, id string) (*fincasdomain.Finca, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.stores[id]
	if !ok {
		return nil, fincasdomain.ErrFincaNoEncontrada
	}
	return f, nil
}

func (m *mockFincaRepo) ListarPorUsuario(_ context.Context, _ string) ([]fincasdomain.Finca, error) {
	return nil, nil
}

func (m *mockFincaRepo) ListarTodas(_ context.Context) ([]fincasdomain.Finca, error) {
	return nil, nil
}

func (m *mockFincaRepo) Listar(_ context.Context, _ fincasdomain.EspecificacionFinca, _ shared.Paginacion) ([]fincasdomain.Finca, error) {
	return nil, nil
}

func (m *mockFincaRepo) Actualizar(_ context.Context, finca *fincasdomain.Finca) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[finca.ID()] = finca
	return nil
}

func (m *mockFincaRepo) Eliminar(_ context.Context, _ string) error {
	return nil
}

func (m *mockFincaRepo) ContarLotes(_ context.Context, fincaID string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.loteRepo == nil {
		return 0, nil
	}
	m.loteRepo.mu.RLock()
	defer m.loteRepo.mu.RUnlock()
	count := 0
	for _, l := range m.loteRepo.stores {
		if l.FincaID() == fincaID {
			count++
		}
	}
	return count, nil
}

// ──────────────────────────────────────────────
// Mock LoteRepositorio
// ──────────────────────────────────────────────

type mockLoteRepo struct {
	mu     sync.RWMutex
	stores map[string]*fincasdomain.Lote
}

func (m *mockLoteRepo) Crear(_ context.Context, lote *fincasdomain.Lote) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[lote.ID()] = lote
	return nil
}

func (m *mockLoteRepo) ObtenerPorID(_ context.Context, id string) (*fincasdomain.Lote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.stores[id]
	if !ok {
		return nil, fincasdomain.ErrLoteNoEncontrado
	}
	return l, nil
}

func (m *mockLoteRepo) ListarPorFinca(_ context.Context, _ string) ([]fincasdomain.Lote, error) {
	return nil, nil
}

func (m *mockLoteRepo) Listar(_ context.Context, _ fincasdomain.EspecificacionLote, _ shared.Paginacion) ([]fincasdomain.Lote, error) {
	return nil, nil
}

func (m *mockLoteRepo) Actualizar(_ context.Context, lote *fincasdomain.Lote) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[lote.ID()] = lote
	return nil
}

func (m *mockLoteRepo) Eliminar(_ context.Context, _ string) error {
	return nil
}

// ──────────────────────────────────────────────
// Mock MuestraRepositorio
// ──────────────────────────────────────────────

type mockMuestraRepo struct {
	mu     sync.RWMutex
	stores map[string]*diagnosticodomain.Muestra
}

func (m *mockMuestraRepo) Crear(_ context.Context, muestra *diagnosticodomain.Muestra) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[muestra.ID()] = muestra
	return nil
}

func (m *mockMuestraRepo) ObtenerPorID(_ context.Context, id string) (*diagnosticodomain.Muestra, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	muestra, ok := m.stores[id]
	if !ok {
		return nil, diagnosticodomain.ErrMuestraNoEncontrada
	}
	return muestra, nil
}

func (m *mockMuestraRepo) ListarPorDiagnostico(_ context.Context, _ string) ([]diagnosticodomain.Muestra, error) {
	return nil, nil
}

func (m *mockMuestraRepo) Actualizar(_ context.Context, _ *diagnosticodomain.Muestra) error {
	return nil
}

func (m *mockMuestraRepo) Eliminar(_ context.Context, _ string) error {
	return nil
}

func (m *mockMuestraRepo) Buscar(_ context.Context, especificacion diagnosticodomain.EspecificacionMuestra, _ shared.Paginacion) ([]diagnosticodomain.Muestra, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var loteIDFilter, tenantIDFilter string
	hasLoteFilter := false
	hasTenantFilter := false

	for _, f := range especificacion.Filtros {
		if f.Campo == "loteID" {
			loteIDFilter = fmt.Sprintf("%v", f.Valor)
			hasLoteFilter = true
		}
		if f.Campo == "tenantID" {
			tenantIDFilter = fmt.Sprintf("%v", f.Valor)
			hasTenantFilter = true
		}
	}

	var result []diagnosticodomain.Muestra
	for _, m := range m.stores {
		if hasLoteFilter && m.LoteID() != loteIDFilter {
			continue
		}
		if hasTenantFilter && m.TenantID() != tenantIDFilter {
			continue
		}
		result = append(result, *m)
	}
	return result, nil
}

// ──────────────────────────────────────────────
// Mock DiagnosticoRepositorio
// ──────────────────────────────────────────────

type mockDiagnosticoRepo struct {
	mu     sync.RWMutex
	stores map[string]*diagnosticodomain.Diagnostico
}

func (m *mockDiagnosticoRepo) Crear(_ context.Context, diagnostico *diagnosticodomain.Diagnostico) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[diagnostico.ID()] = diagnostico
	return nil
}

func (m *mockDiagnosticoRepo) ObtenerPorID(_ context.Context, id string) (*diagnosticodomain.Diagnostico, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.stores[id]
	if !ok {
		return nil, diagnosticodomain.ErrDiagnosticoNoEncontrado
	}
	return d, nil
}

func (m *mockDiagnosticoRepo) ListarPorFinca(_ context.Context, _ string) ([]diagnosticodomain.Diagnostico, error) {
	return nil, nil
}

func (m *mockDiagnosticoRepo) Actualizar(_ context.Context, diagnostico *diagnosticodomain.Diagnostico) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[diagnostico.ID()] = diagnostico
	return nil
}

func (m *mockDiagnosticoRepo) Eliminar(_ context.Context, _ string) error {
	return nil
}

func (m *mockDiagnosticoRepo) Buscar(_ context.Context, especificacion diagnosticodomain.EspecificacionDiagnostico, _ shared.Paginacion) ([]diagnosticodomain.Diagnostico, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var muestraIDFilter string
	hasMuestraFilter := false

	for _, f := range especificacion.Filtros {
		if f.Campo == "muestraID" {
			muestraIDFilter = fmt.Sprintf("%v", f.Valor)
			hasMuestraFilter = true
		}
	}

	var result []diagnosticodomain.Diagnostico
	for _, d := range m.stores {
		if hasMuestraFilter && d.MuestrasId() != muestraIDFilter {
			continue
		}
		result = append(result, *d)
	}
	return result, nil
}

// ──────────────────────────────────────────────
// Mock CandidatoReentrenamientoRepositorio
// ──────────────────────────────────────────────

type mockCandidatoRepo struct {
	mu     sync.RWMutex
	stores map[string]*diagnosticodomain.CandidatoReentrenamiento
}

func (m *mockCandidatoRepo) Crear(_ context.Context, candidato *diagnosticodomain.CandidatoReentrenamiento) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[candidato.ID()] = candidato
	return nil
}

func (m *mockCandidatoRepo) ObtenerPorDiagnosticoID(_ context.Context, diagnosticoID string) (*diagnosticodomain.CandidatoReentrenamiento, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, c := range m.stores {
		if c.DiagnosticoID() == diagnosticoID {
			return c, nil
		}
	}
	return nil, diagnosticodomain.ErrCandidatoNoEncontrado
}

func (m *mockCandidatoRepo) ListarPendientes(_ context.Context, _ int) ([]diagnosticodomain.CandidatoReentrenamiento, error) {
	return nil, nil
}

// ──────────────────────────────────────────────
// Mock GeneradorID
// ──────────────────────────────────────────────

type mockGeneradorID struct {
	counter int64
}

func (g *mockGeneradorID) NextID(_ context.Context) (string, error) {
	next := atomic.AddInt64(&g.counter, 1)
	return fmt.Sprintf("id-%04d", next), nil
}

// ──────────────────────────────────────────────
// Mock EventPublisher
// ──────────────────────────────────────────────

type mockPublisher struct {
	mu      sync.Mutex
	events  []any
}

func (p *mockPublisher) Publish(_ context.Context, _ string, event any) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, event)
	return nil
}

// ──────────────────────────────────────────────
// Mock UnitOfWorkDiagnostico
// ──────────────────────────────────────────────

type mockUoW struct {
	diagnosticoRepo *mockDiagnosticoRepo
	candidatoRepo   *mockCandidatoRepo
}

func (u *mockUoW) Transaccional(_ context.Context, fn func(application.UnitOfWorkDiagnostico) error) error {
	return fn(u)
}

func (u *mockUoW) DiagnosticoRepo() diagnosticodomain.DiagnosticoRepositorio {
	return u.diagnosticoRepo
}

func (u *mockUoW) CandidatoRepo() diagnosticodomain.CandidatoReentrenamientoRepositorio {
	return u.candidatoRepo
}

// ──────────────────────────────────────────────
// main
// ──────────────────────────────────────────────

func main() {
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("  TEST FUNCIONAL DE CASOS DE USO - FINCAS")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println()

	ctx := context.Background()
	now := time.Now()

	// ── Create shared mocks ──────────────────────

	fincaRepo := &mockFincaRepo{stores: make(map[string]*fincasdomain.Finca)}
	loteRepo := &mockLoteRepo{stores: make(map[string]*fincasdomain.Lote)}
	fincaRepo.loteRepo = loteRepo

	muestraRepo := &mockMuestraRepo{stores: make(map[string]*diagnosticodomain.Muestra)}
	diagnosticoRepo := &mockDiagnosticoRepo{stores: make(map[string]*diagnosticodomain.Diagnostico)}
	candidatoRepo := &mockCandidatoRepo{stores: make(map[string]*diagnosticodomain.CandidatoReentrenamiento)}
	generador := &mockGeneradorID{counter: 0}
	publisher := &mockPublisher{}

	uow := &mockUoW{
		diagnosticoRepo: diagnosticoRepo,
		candidatoRepo:   candidatoRepo,
	}

	// ── Create use cases ─────────────────────────

	registrarFincaUC := registrarfinca.NewUseCase(fincaRepo, generador, publisher)
	desactivarFincaUC := desactivarfinca.NewUseCase(fincaRepo, fincasdomain.NewFincaService(), generador, publisher)
	agregarLoteUC := agregarlote.NewUseCase(fincaRepo, loteRepo, generador, publisher)
	eliminarLoteUC := eliminarlote.NewUseCase(loteRepo, generador, publisher)
	tomarMuestraUC := tomarmuestra.NewUseCase(loteRepo, muestraRepo, generador, publisher)
	listarMuestrasUC := listarmuestrasporlote.NewUseCase(loteRepo, muestraRepo)
	solicitarDiagnosticoUC := solicitardiagnosticomanual.NewUseCase(muestraRepo, generador, publisher)
	registrarInferenciaUC := registrarinferencia.NewUseCase(muestraRepo, diagnosticoRepo, generador, publisher)
	aceptarDiagnosticoUC := aceptardiagnostico.NewUseCase(diagnosticoRepo, generador, publisher)
	rechazarDiagnosticoUC := rechazardiagnostico.NewUseCase(diagnosticoRepo, candidatoRepo, uow, generador, publisher)
	generarReporteUC := generarreporteporlote.NewUseCase(loteRepo, muestraRepo, diagnosticoRepo)

	// ── Auth contexts ────────────────────────────

	authFull := &application.AuthContext{
		UsuarioID: "user-1",
		TenantID:  "tenant-1",
		Permisos: []string{
			"CREAR_FINCA",
			"DESACTIVAR_FINCA",
			"CREAR_LOTE",
			"ELIMINAR_LOTE",
			"CREAR_MUESTRA",
			"VER_MUESTRAS",
			"SOLICITAR_DIAGNOSTICO",
			"ACEPTAR_DIAGNOSTICO",
			"RECHAZAR_DIAGNOSTICO",
			"GENERAR_REPORTE",
		},
	}
	authNoPermisos := &application.AuthContext{
		UsuarioID: "user-2",
		TenantID:  "tenant-1",
		Permisos:  []string{},
	}

	_ = now // used in test setups below

	// ══════════════════════════════════════════════
	//  4.1 RegistrarFinca
	// ══════════════════════════════════════════════

	fmt.Println("─── 4.1 RegistrarFinca ───────────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	runTest("Normal: nombre valido y ubicacion", func() error {
		cmd := registrarfinca.Command{
			Nombre:      "Mi Finca",
			Ubicacion:   "Ubicacion",
			Descripcion: "desc",
		}
		_, err := registrarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	})

	runTestExpected("Edge: nombre vacio", func() error {
		cmd := registrarfinca.Command{
			Nombre:      "",
			Ubicacion:   "Ubicacion",
			Descripcion: "desc",
		}
		_, err := registrarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: nombre 2 caracteres", func() error {
		cmd := registrarfinca.Command{
			Nombre:      "AB",
			Ubicacion:   "Ubicacion",
			Descripcion: "desc",
		}
		_, err := registrarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: nombre 201 caracteres", func() error {
		cmd := registrarfinca.Command{
			Nombre:      strings.Repeat("X", 201),
			Ubicacion:   "Ubicacion",
			Descripcion: "desc",
		}
		_, err := registrarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: ubicacion vacia", func() error {
		cmd := registrarfinca.Command{
			Nombre:      "Finca Test",
			Ubicacion:   "",
			Descripcion: "desc",
		}
		_, err := registrarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: descripcion 1001 caracteres", func() error {
		cmd := registrarfinca.Command{
			Nombre:      "Finca Test",
			Ubicacion:   "Ubicacion",
			Descripcion: strings.Repeat("X", 1001),
		}
		_, err := registrarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: sin permiso CREAR_FINCA", func() error {
		cmd := registrarfinca.Command{
			Nombre:      "Mi Finca",
			Ubicacion:   "Ubicacion",
			Descripcion: "desc",
		}
		_, err := registrarFincaUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.2 DesactivarFinca
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.2 DesactivarFinca ───────────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	tenantID := "tenant-1"
	finteg := func(id, nombre, ubicacion, descripcion, usuario string, estado fincasdomain.EstadoFinca) *fincasdomain.Finca {
		f := fincasdomain.NewFincaFromPersistence(id, nombre, ubicacion, descripcion, usuario, &tenantID, estado, time.Now(), time.Now())
		return f
	}

	// Pre-populate fincas
	fincaSinLotes := finteg("finca-sin-lotes", "Sin Lotes", "Ubicacion", "desc", "user-1", fincasdomain.FincaActiva)
	fincaConLotes := finteg("finca-con-lotes", "Con Lotes", "Ubicacion", "desc", "user-1", fincasdomain.FincaActiva)
	fincaYaPendiente := finteg("finca-pendiente", "Ya Pendiente", "Ubicacion", "desc", "user-1", fincasdomain.FincaPendienteEliminar)
	_ = fincaYaPendiente

	fincaRepo.stores["finca-sin-lotes"] = fincaSinLotes
	fincaRepo.stores["finca-con-lotes"] = fincaConLotes
	fincaRepo.stores["finca-pendiente"] = fincaYaPendiente

	// Add a lote to finca-con-lotes
	loteParaFinca := fincasdomain.NewLoteFromPersistence(
		"lote-de-finca", "finca-con-lotes", tenantID, "Lote 1", 100.0, "lote desc",
		fincasdomain.LoteActivo, time.Now(), time.Now(),
	)
	loteRepo.stores["lote-de-finca"] = loteParaFinca

	runTest("Normal: finca sin lotes", func() error {
		cmd := desactivarfinca.Command{FincaID: "finca-sin-lotes", Confirmar: false}
		_, err := desactivarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	})

	runTest("Normal: finca con lotes + Confirmar=true", func() error {
		cmd := desactivarfinca.Command{FincaID: "finca-con-lotes", Confirmar: true}
		_, err := desactivarFincaUC.Ejecutar(ctx, authFull, cmd)
		return err
	})

	runTestExpected("Edge: finca con lotes + Confirmar=false", func() error {
		cmd := desactivarfinca.Command{FincaID: "finca-con-lotes", Confirmar: false}
		_, err := desactivarFincaUC.Ejecutar(ctx, authFull, cmd)
		if err == nil {
			return fmt.Errorf("expected error but got nil")
		}
		return err
	}, true)

	runTestExpected("Edge: finca no existe", func() error {
		cmd := desactivarfinca.Command{FincaID: "no-existe", Confirmar: false}
		_, err := desactivarFincaUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	// Use a fresh finca in PENDIENTE state for this test
	runTestExpected("Edge: finca ya en PENDIENTE_ELIMINACION", func() error {
		// Insert a fresh finca directly in PENDIENTE state
		freshPendiente := finteg("finca-pendiente-fresh", "Fresh", "Ubi", "desc", "user-1", fincasdomain.FincaPendienteEliminar)
		fincaRepo.mu.Lock()
		fincaRepo.stores["finca-pendiente-fresh"] = freshPendiente
		fincaRepo.mu.Unlock()

		cmd := desactivarfinca.Command{FincaID: "finca-pendiente-fresh", Confirmar: false}
		_, err := desactivarFincaUC.Ejecutar(ctx, authFull, cmd)
		if err == nil {
			return fmt.Errorf("expected error but got nil")
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso DESACTIVAR_FINCA", func() error {
		cmd := desactivarfinca.Command{FincaID: "finca-sin-lotes", Confirmar: false}
		_, err := desactivarFincaUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.3 AgregarLote
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.3 AgregarLote ───────────────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate fincas
	fincaActiva := finteg("finca-activa", "Activa", "Ubicacion", "desc", "user-1", fincasdomain.FincaActiva)
	fincaPendiente := finteg("finca-pend-agregar", "Pendiente", "Ubicacion", "desc", "user-1", fincasdomain.FincaPendienteEliminar)
	fincaRepo.stores["finca-activa"] = fincaActiva
	fincaRepo.stores["finca-pend-agregar"] = fincaPendiente

	runTest("Normal: finca activa", func() error {
		cmd := agregarlote.Command{
			FincaID:     "finca-activa",
			Nombre:      "Lote 1",
			Area:        100.0,
			Descripcion: "lote desc",
		}
		_, err := agregarLoteUC.Ejecutar(ctx, authFull, cmd)
		return err
	})

	runTestExpected("Edge: finca en PENDIENTE_ELIMINACION", func() error {
		cmd := agregarlote.Command{
			FincaID:     "finca-pend-agregar",
			Nombre:      "Lote X",
			Area:        50.0,
			Descripcion: "",
		}
		_, err := agregarLoteUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: finca no existe", func() error {
		cmd := agregarlote.Command{
			FincaID:     "no-existe",
			Nombre:      "Lote X",
			Area:        50.0,
			Descripcion: "",
		}
		_, err := agregarLoteUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: area = 0", func() error {
		cmd := agregarlote.Command{
			FincaID:     "finca-activa",
			Nombre:      "Lote X",
			Area:        0,
			Descripcion: "",
		}
		_, err := agregarLoteUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: nombre vacio", func() error {
		cmd := agregarlote.Command{
			FincaID:     "finca-activa",
			Nombre:      "",
			Area:        50.0,
			Descripcion: "",
		}
		_, err := agregarLoteUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: sin permiso CREAR_LOTE", func() error {
		cmd := agregarlote.Command{
			FincaID:     "finca-activa",
			Nombre:      "Lote 2",
			Area:        100.0,
			Descripcion: "",
		}
		_, err := agregarLoteUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.4 EliminarLote
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.4 EliminarLote ──────────────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate finca + lotes
	fincaParaEliminar := finteg("finca-eliminar", "Finca", "Ubi", "desc", "user-1", fincasdomain.FincaActiva)
	fincaRepo.stores["finca-eliminar"] = fincaParaEliminar

	loteActivo := fincasdomain.NewLoteFromPersistence(
		"lote-activo", "finca-eliminar", tenantID, "Lote Activo", 100, "desc",
		fincasdomain.LoteActivo, time.Now(), time.Now(),
	)
	loteEliminado := fincasdomain.NewLoteFromPersistence(
		"lote-eliminado", "finca-eliminar", tenantID, "Lote Eliminado", 100, "desc",
		fincasdomain.LoteEliminado, time.Now(), time.Now(),
	)
	loteRepo.stores["lote-activo"] = loteActivo
	loteRepo.stores["lote-eliminado"] = loteEliminado

	runTest("Normal: lote activo", func() error {
		cmd := eliminarlote.Command{LoteID: "lote-activo"}
		_, err := eliminarLoteUC.Ejecutar(ctx, authFull, cmd)
		return err
	})

	runTestExpected("Edge: lote ya ELIMINADO", func() error {
		cmd := eliminarlote.Command{LoteID: "lote-eliminado"}
		_, err := eliminarLoteUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: lote no existe", func() error {
		cmd := eliminarlote.Command{LoteID: "no-existe"}
		_, err := eliminarLoteUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso ELIMINAR_LOTE", func() error {
		cmd := eliminarlote.Command{LoteID: "lote-activo"}
		_, err := eliminarLoteUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.5 TomarMuestra
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.5 TomarMuestra ──────────────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate finca + lotes
	fincaMuestra := finteg("finca-muestra", "Finca", "Ubi", "desc", "user-1", fincasdomain.FincaActiva)
	fincaRepo.stores["finca-muestra"] = fincaMuestra

	loteActivoMuestra := fincasdomain.NewLoteFromPersistence(
		"lote-muestra-activo", "finca-muestra", tenantID, "Lote Activo", 100, "desc",
		fincasdomain.LoteActivo, time.Now(), time.Now(),
	)
	loteEliminadoMuestra := fincasdomain.NewLoteFromPersistence(
		"lote-muestra-eliminado", "finca-muestra", tenantID, "Lote Eliminado", 100, "desc",
		fincasdomain.LoteEliminado, time.Now(), time.Now(),
	)
	loteRepo.stores["lote-muestra-activo"] = loteActivoMuestra
	loteRepo.stores["lote-muestra-eliminado"] = loteEliminadoMuestra

	runTest("Normal: lote activo, coordenadas validas", func() error {
		cmd := tomarmuestra.Command{
			LoteID:   "lote-muestra-activo",
			Latitud:  4.711,
			Longitud: -74.072,
		}
		_, err := tomarMuestraUC.Ejecutar(ctx, authFull, cmd)
		return err
	})

	runTestExpected("Edge: lote ELIMINADO", func() error {
		cmd := tomarmuestra.Command{
			LoteID:   "lote-muestra-eliminado",
			Latitud:  4.711,
			Longitud: -74.072,
		}
		_, err := tomarMuestraUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: latitud > 90", func() error {
		cmd := tomarmuestra.Command{
			LoteID:   "lote-muestra-activo",
			Latitud:  91.0,
			Longitud: 0.0,
		}
		_, err := tomarMuestraUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: longitud < -180", func() error {
		cmd := tomarmuestra.Command{
			LoteID:   "lote-muestra-activo",
			Latitud:  0.0,
			Longitud: -181.0,
		}
		_, err := tomarMuestraUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: lote no existe", func() error {
		cmd := tomarmuestra.Command{
			LoteID:   "no-existe",
			Latitud:  4.711,
			Longitud: -74.072,
		}
		_, err := tomarMuestraUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso CREAR_MUESTRA", func() error {
		cmd := tomarmuestra.Command{
			LoteID:   "lote-muestra-activo",
			Latitud:  4.711,
			Longitud: -74.072,
		}
		_, err := tomarMuestraUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.6 ListarMuestrasPorLote
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.6 ListarMuestrasPorLote ──────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate finca, lote and muestras
	fincaListar := finteg("finca-listar", "Finca", "Ubi", "desc", "user-1", fincasdomain.FincaActiva)
	fincaRepo.stores["finca-listar"] = fincaListar

	loteListar := fincasdomain.NewLoteFromPersistence(
		"lote-listar", "finca-listar", tenantID, "Lote Listar", 100, "desc",
		fincasdomain.LoteActivo, time.Now(), time.Now(),
	)
	loteRepo.stores["lote-listar"] = loteListar

	// Create 2 muestras
	ubic1, _ := diagnosticodomain.NewUbicacion(4.711, -74.072)
	muestra1, _ := diagnosticodomain.NewMuestra("muestra-listar-1", "lote-listar", *ubic1, tenantID)
	muestraRepo.stores["muestra-listar-1"] = muestra1

	ubic2, _ := diagnosticodomain.NewUbicacion(4.712, -74.073)
	muestra2, _ := diagnosticodomain.NewMuestra("muestra-listar-2", "lote-listar", *ubic2, tenantID)
	muestraRepo.stores["muestra-listar-2"] = muestra2

	runTest("Normal: lote con 2 muestras", func() error {
		cmd := listarmuestrasporlote.Command{LoteID: "lote-listar"}
		items, err := listarMuestrasUC.Ejecutar(ctx, authFull, cmd)
		if err != nil {
			return err
		}
		if len(items) != 2 {
			return fmt.Errorf("expected 2 items, got %d", len(items))
		}
		return nil
	})

	runTest("Edge: lote sin muestras", func() error {
		// Create a lote with no muestras
		loteSinMuestras := fincasdomain.NewLoteFromPersistence(
			"lote-sin-muestras", "finca-listar", tenantID, "Sin Muestras", 100, "desc",
			fincasdomain.LoteActivo, time.Now(), time.Now(),
		)
		loteRepo.stores["lote-sin-muestras"] = loteSinMuestras

		cmd := listarmuestrasporlote.Command{LoteID: "lote-sin-muestras"}
		items, err := listarMuestrasUC.Ejecutar(ctx, authFull, cmd)
		if err != nil {
			return err
		}
		if len(items) != 0 {
			return fmt.Errorf("expected 0 items, got %d", len(items))
		}
		return nil
	})

	runTestExpected("Edge: lote no existe", func() error {
		cmd := listarmuestrasporlote.Command{LoteID: "no-existe"}
		_, err := listarMuestrasUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso VER_MUESTRAS", func() error {
		cmd := listarmuestrasporlote.Command{LoteID: "lote-listar"}
		_, err := listarMuestrasUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.7 SolicitarDiagnosticoManual
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.7 SolicitarDiagnosticoManual ─────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate finca, lote and muestra
	fincaSolDiag := finteg("finca-soldiag", "Finca", "Ubi", "desc", "user-1", fincasdomain.FincaActiva)
	fincaRepo.stores["finca-soldiag"] = fincaSolDiag

	loteSolDiag := fincasdomain.NewLoteFromPersistence(
		"lote-soldiag", "finca-soldiag", tenantID, "Lote", 100, "desc",
		fincasdomain.LoteActivo, time.Now(), time.Now(),
	)
	loteRepo.stores["lote-soldiag"] = loteSolDiag

	ubicSolDiag, _ := diagnosticodomain.NewUbicacion(4.711, -74.072)
	muestraSolDiag, _ := diagnosticodomain.NewMuestra("muestra-soldiag", "lote-soldiag", *ubicSolDiag, tenantID)
	muestraRepo.stores["muestra-soldiag"] = muestraSolDiag

	runTest("Normal: muestra existe, url HTTPS", func() error {
		cmd := solicitardiagnosticomanual.Command{
			MuestraID: "muestra-soldiag",
			ImageURL:  "https://ejemplo.com/imagen.jpg",
		}
		_, err := solicitarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		return err
	})

	runTestExpected("Edge: imageURL http", func() error {
		cmd := solicitardiagnosticomanual.Command{
			MuestraID: "muestra-soldiag",
			ImageURL:  "http://ejemplo.com/imagen.jpg",
		}
		_, err := solicitarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: muestra no existe", func() error {
		cmd := solicitardiagnosticomanual.Command{
			MuestraID: "no-existe",
			ImageURL:  "https://ejemplo.com/imagen.jpg",
		}
		_, err := solicitarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso SOLICITAR_DIAGNOSTICO", func() error {
		cmd := solicitardiagnosticomanual.Command{
			MuestraID: "muestra-soldiag",
			ImageURL:  "https://ejemplo.com/imagen.jpg",
		}
		_, err := solicitarDiagnosticoUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.8 RegistrarInferencia (no usa auth)
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.8 RegistrarInferencia ────────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate muestra
	ubicInf, _ := diagnosticodomain.NewUbicacion(4.711, -74.072)
	muestraInf, _ := diagnosticodomain.NewMuestra("muestra-inf", "lote-inf", *ubicInf, tenantID)
	muestraRepo.stores["muestra-inf"] = muestraInf

	pasado := time.Now().Add(-1 * time.Hour)

	runTest("Normal: todo valido", func() error {
		cmd := registrarinferencia.Command{
			MuestraID:     "muestra-inf",
			ImageURL:      "https://ejemplo.com/yolo.jpg",
			TieneClorosis: true,
			Confianza:     0.95,
			ProcesadoAt:   pasado,
		}
		_, err := registrarInferenciaUC.Ejecutar(ctx, cmd)
		return err
	})

	runTestExpected("Edge: confianza = 1.5", func() error {
		cmd := registrarinferencia.Command{
			MuestraID:     "muestra-inf",
			ImageURL:      "https://ejemplo.com/yolo.jpg",
			TieneClorosis: true,
			Confianza:     1.5,
			ProcesadoAt:   pasado,
		}
		_, err := registrarInferenciaUC.Ejecutar(ctx, cmd)
		return err
	}, true)

	runTestExpected("Edge: confianza = -0.1", func() error {
		cmd := registrarinferencia.Command{
			MuestraID:     "muestra-inf",
			ImageURL:      "https://ejemplo.com/yolo.jpg",
			TieneClorosis: false,
			Confianza:     -0.1,
			ProcesadoAt:   pasado,
		}
		_, err := registrarInferenciaUC.Ejecutar(ctx, cmd)
		return err
	}, true)

	runTestExpected("Edge: muestra no existe", func() error {
		cmd := registrarinferencia.Command{
			MuestraID:     "no-existe",
			ImageURL:      "https://ejemplo.com/yolo.jpg",
			TieneClorosis: false,
			Confianza:     0.5,
			ProcesadoAt:   pasado,
		}
		_, err := registrarInferenciaUC.Ejecutar(ctx, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: ProcesadoAt futuro", func() error {
		futuro := time.Now().Add(1 * time.Hour)
		cmd := registrarinferencia.Command{
			MuestraID:     "muestra-inf",
			ImageURL:      "https://ejemplo.com/yolo.jpg",
			TieneClorosis: false,
			Confianza:     0.5,
			ProcesadoAt:   futuro,
		}
		_, err := registrarInferenciaUC.Ejecutar(ctx, cmd)
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.9 AceptarDiagnostico
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.9 AceptarDiagnostico ────────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate muestra and diagnostico PENDIENTE
	ubicAcep, _ := diagnosticodomain.NewUbicacion(4.711, -74.072)
	muestraAcep, _ := diagnosticodomain.NewMuestra("muestra-acep", "lote-acep", *ubicAcep, tenantID)
	muestraRepo.stores["muestra-acep"] = muestraAcep

	resAcep, _ := diagnosticodomain.NewResultadoInferencia("https://img.url", true, 0.95, time.Now().Add(-1*time.Hour))
	diagPendiente, _ := diagnosticodomain.NewDiagnostico("diag-pendiente", "INF-001", "muestra-acep", tenantID, resAcep)
	diagnosticoRepo.stores["diag-pendiente"] = diagPendiente

	diagnosticoAceptado, _ := diagnosticodomain.NewDiagnostico("diag-ya-aceptado", "INF-002", "muestra-acep", tenantID, resAcep)
	_ = diagnosticoAceptado.MarcarComoAceptado()
	diagnosticoRepo.stores["diag-ya-aceptado"] = diagnosticoAceptado

	runTest("Normal: diagnostico PENDIENTE -> ACEPTADO", func() error {
		cmd := aceptardiagnostico.Command{DiagnosticoID: "diag-pendiente"}
		salida, err := aceptarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		if err != nil {
			return err
		}
		if salida.Estado != "ACEPTADO" {
			return fmt.Errorf("expected estado ACEPTADO, got %s", salida.Estado)
		}
		return nil
	})

	runTestExpected("Edge: diagnostico ya ACEPTADO", func() error {
		cmd := aceptardiagnostico.Command{DiagnosticoID: "diag-ya-aceptado"}
		_, err := aceptarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: diagnostico no existe", func() error {
		cmd := aceptardiagnostico.Command{DiagnosticoID: "no-existe"}
		_, err := aceptarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso ACEPTAR_DIAGNOSTICO", func() error {
		cmd := aceptardiagnostico.Command{DiagnosticoID: "diag-pendiente"}
		_, err := aceptarDiagnosticoUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.10 RechazarDiagnostico
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.10 RechazarDiagnostico ───────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate muestra and diagnostico PENDIENTE
	ubicRech, _ := diagnosticodomain.NewUbicacion(4.711, -74.072)
	muestraRech, _ := diagnosticodomain.NewMuestra("muestra-rech", "lote-rech", *ubicRech, tenantID)
	muestraRepo.stores["muestra-rech"] = muestraRech

	resRech, _ := diagnosticodomain.NewResultadoInferencia("https://img.url/rech", true, 0.92, time.Now().Add(-1*time.Hour))
	diagPendienteRech, _ := diagnosticodomain.NewDiagnostico("diag-pend-rech", "INF-003", "muestra-rech", tenantID, resRech)
	diagnosticoRepo.stores["diag-pend-rech"] = diagPendienteRech

	// A diag that is already RECHAZADO
	diagYaRechazado, _ := diagnosticodomain.NewDiagnostico("diag-ya-rech", "INF-004", "muestra-rech", tenantID, resRech)
	_ = diagYaRechazado.MarcarComoRechazado()
	diagnosticoRepo.stores["diag-ya-rech"] = diagYaRechazado

	runTest("Normal: diagnostico PENDIENTE -> RECHAZADO + candidato", func() error {
		cmd := rechazardiagnostico.Command{
			DiagnosticoID: "diag-pend-rech",
			Motivo:        "falso positivo manual",
		}
		salida, err := rechazarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		if err != nil {
			return err
		}
		if salida.Estado != "RECHAZADO" {
			return fmt.Errorf("expected estado RECHAZADO, got %s", salida.Estado)
		}
		// Verify candidato was created
		candidato, err := candidatoRepo.ObtenerPorDiagnosticoID(ctx, "diag-pend-rech")
		if err != nil {
			return fmt.Errorf("candidato should exist: %w", err)
		}
		if candidato.DiagnosticoID() != "diag-pend-rech" {
			return fmt.Errorf("candidato references wrong diagnostico")
		}
		if !candidato.TieneClorosis() {
			return fmt.Errorf("candidato should have tieneClorosis=true")
		}
		return nil
	})

	runTestExpected("Edge: diagnostico ya RECHAZADO", func() error {
		cmd := rechazardiagnostico.Command{
			DiagnosticoID: "diag-ya-rech",
			Motivo:        "nuevo rechazo",
		}
		_, err := rechazarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		return err
	}, true)

	runTestExpected("Edge: diagnostico no existe", func() error {
		cmd := rechazardiagnostico.Command{
			DiagnosticoID: "no-existe",
			Motivo:        "motivo",
		}
		_, err := rechazarDiagnosticoUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso RECHAZAR_DIAGNOSTICO", func() error {
		cmd := rechazardiagnostico.Command{
			DiagnosticoID: "diag-pend-rech",
			Motivo:        "motivo",
		}
		_, err := rechazarDiagnosticoUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  4.11 GenerarReportePorLote
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("─── 4.11 GenerarReportePorLote ─────────────")
	resetMocks(fincaRepo, loteRepo, muestraRepo, diagnosticoRepo, candidatoRepo)

	// Pre-populate lote with area=10.0 hectares
	loteReporte := fincasdomain.NewLoteFromPersistence(
		"lote-reporte", "finca-reporte", tenantID, "Lote Reporte", 10.0, "lote para reporte",
		fincasdomain.LoteActivo, time.Now(), time.Now(),
	)
	loteRepo.stores["lote-reporte"] = loteReporte

	// Create 4 muestras
	ubicR1, _ := diagnosticodomain.NewUbicacion(4.711, -74.072)
	m1, _ := diagnosticodomain.NewMuestra("muestra-r1", "lote-reporte", *ubicR1, tenantID)
	muestraRepo.stores["muestra-r1"] = m1

	ubicR2, _ := diagnosticodomain.NewUbicacion(4.712, -74.073)
	m2, _ := diagnosticodomain.NewMuestra("muestra-r2", "lote-reporte", *ubicR2, tenantID)
	muestraRepo.stores["muestra-r2"] = m2

	ubicR3, _ := diagnosticodomain.NewUbicacion(4.713, -74.074)
	m3, _ := diagnosticodomain.NewMuestra("muestra-r3", "lote-reporte", *ubicR3, tenantID)
	muestraRepo.stores["muestra-r3"] = m3

	ubicR4, _ := diagnosticodomain.NewUbicacion(4.714, -74.075)
	m4, _ := diagnosticodomain.NewMuestra("muestra-r4", "lote-reporte", *ubicR4, tenantID)
	muestraRepo.stores["muestra-r4"] = m4

	// Create diagnosticos:
	// m1: ACEPTADO con clorosis
	// m2: ACEPTADO con clorosis
	// m3: ACEPTADO sin clorosis
	// m4: solo PENDIENTE (no ACEPTADO)
	procesado := time.Now().Add(-1 * time.Hour)

	res1, _ := diagnosticodomain.NewResultadoInferencia("https://img/r1", true, 0.95, procesado)
	d1, _ := diagnosticodomain.NewDiagnostico("diag-r1", "INF-R1", "muestra-r1", tenantID, res1)
	_ = d1.MarcarComoAceptado()
	diagnosticoRepo.stores["diag-r1"] = d1

	res2, _ := diagnosticodomain.NewResultadoInferencia("https://img/r2", true, 0.88, procesado)
	d2, _ := diagnosticodomain.NewDiagnostico("diag-r2", "INF-R2", "muestra-r2", tenantID, res2)
	_ = d2.MarcarComoAceptado()
	diagnosticoRepo.stores["diag-r2"] = d2

	res3, _ := diagnosticodomain.NewResultadoInferencia("https://img/r3", false, 0.99, procesado)
	d3, _ := diagnosticodomain.NewDiagnostico("diag-r3", "INF-R3", "muestra-r3", tenantID, res3)
	_ = d3.MarcarComoAceptado()
	diagnosticoRepo.stores["diag-r3"] = d3

	res4, _ := diagnosticodomain.NewResultadoInferencia("https://img/r4", true, 0.70, procesado)
	d4, _ := diagnosticodomain.NewDiagnostico("diag-r4", "INF-R4", "muestra-r4", tenantID, res4)
	// d4 stays PENDIENTE
	diagnosticoRepo.stores["diag-r4"] = d4

	runTest("Normal: lote con muestras y diagnosticos", func() error {
		cmd := generarreporteporlote.Command{LoteID: "lote-reporte"}
		salida, err := generarReporteUC.Ejecutar(ctx, authFull, cmd)
		if err != nil {
			return err
		}
		if salida.Metricas.TotalMuestras != 4 {
			return fmt.Errorf("expected TotalMuestras=4, got %d", salida.Metricas.TotalMuestras)
		}
		if salida.Metricas.ConClorosis != 2 {
			return fmt.Errorf("expected ConClorosis=2, got %d", salida.Metricas.ConClorosis)
		}
		if salida.Metricas.SinClorosis != 1 {
			return fmt.Errorf("expected SinClorosis=1, got %d", salida.Metricas.SinClorosis)
		}
		if salida.Metricas.Pendientes != 1 {
			return fmt.Errorf("expected Pendientes=1, got %d", salida.Metricas.Pendientes)
		}
		expectedArea := 2.0 * math.Pi * 2.0 * 2.0
		if math.Abs(salida.Metricas.AreaAfectadaEstimada-expectedArea) > 0.01 {
			return fmt.Errorf("expected AreaAfectadaEstimada=%.4f, got %.4f", expectedArea, salida.Metricas.AreaAfectadaEstimada)
		}
		return nil
	})

	runTest("Edge: lote sin muestras", func() error {
		loteVacio := fincasdomain.NewLoteFromPersistence(
			"lote-vacio", "finca-vacia", tenantID, "Vacio", 5.0, "",
			fincasdomain.LoteActivo, time.Now(), time.Now(),
		)
		loteRepo.stores["lote-vacio"] = loteVacio

		cmd := generarreporteporlote.Command{LoteID: "lote-vacio"}
		salida, err := generarReporteUC.Ejecutar(ctx, authFull, cmd)
		if err != nil {
			return err
		}
		if salida.Metricas.TotalMuestras != 0 {
			return fmt.Errorf("expected TotalMuestras=0, got %d", salida.Metricas.TotalMuestras)
		}
		if salida.Metricas.ConClorosis != 0 {
			return fmt.Errorf("expected ConClorosis=0, got %d", salida.Metricas.ConClorosis)
		}
		if salida.Metricas.Pendientes != 0 {
			return fmt.Errorf("expected Pendientes=0, got %d", salida.Metricas.Pendientes)
		}
		return nil
	})

	runTestExpected("Edge: lote no existe", func() error {
		cmd := generarreporteporlote.Command{LoteID: "no-existe"}
		_, err := generarReporteUC.Ejecutar(ctx, authFull, cmd)
		if !errors.Is(err, application.ErrNotFound) {
			return fmt.Errorf("expected ErrNotFound, got: %w", err)
		}
		return err
	}, true)

	runTestExpected("Edge: sin permiso GENERAR_REPORTE", func() error {
		cmd := generarreporteporlote.Command{LoteID: "lote-reporte"}
		_, err := generarReporteUC.Ejecutar(ctx, authNoPermisos, cmd)
		if !errors.Is(err, application.ErrForbidden) {
			return fmt.Errorf("expected ErrForbidden, got: %w", err)
		}
		return err
	}, true)

	// ══════════════════════════════════════════════
	//  Summary
	// ══════════════════════════════════════════════

	fmt.Println()
	fmt.Println("══════════════════════════════════════════")
	total := passed + failed
	fmt.Printf("  RESULTADO: %d/%d tests passed\n", passed, total)
	if failed > 0 {
		fmt.Println("  FALLOS DETECTADOS:")
		for _, r := range results {
			if !r.ok {
				fmt.Printf("    ❌ %s: %s\n", r.name, r.msg)
			}
		}
	}
	fmt.Println("══════════════════════════════════════════")

	if failed > 0 {
		os.Exit(1)
	}
}
