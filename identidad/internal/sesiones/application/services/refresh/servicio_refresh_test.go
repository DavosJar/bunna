package refresh_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/sesiones/application/services/refresh"
	sesiones_domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad_domain "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	usuario_domain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

// ── Mocks ────────────────────────────────────────────────────────────────────

type mockUnitOfWork struct {
	sesionRepo    sesiones_domain.SesionRepositorio
	tokenServicio sesiones_domain.TokenServicio
}

func (m *mockUnitOfWork) Transaccional(ctx context.Context, fn func(tx sesiones_domain.UnitOfWork) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return fn(m)
}
func (m *mockUnitOfWork) SesionRepositorio() sesiones_domain.SesionRepositorio             { return m.sesionRepo }
func (m *mockUnitOfWork) TokenServicio() sesiones_domain.TokenServicio                      { return m.tokenServicio }
func (m *mockUnitOfWork) CredencialesRepositorio() seguridad_domain.CredencialesRepositorio { return nil }
func (m *mockUnitOfWork) UsuarioRepositorio() usuario_domain.UsuarioRepositorio             { return nil }
func (m *mockUnitOfWork) EncriptacionServicio() seguridad_domain.EncriptacionServicio       { return nil }
func (m *mockUnitOfWork) GeneradorID() shared_domain.GeneradorID                            { return nil }

// mockTokenServicio
type mockTokenServicio struct {
	claimsValidos *sesiones_domain.TokenClaims
	errValidar    error
	failAccess    bool
	failRefresh   bool
}

func (m *mockTokenServicio) GenerarAccessToken(usuarioID, sesionID string, claims *rbac.UsuarioClaims) (string, time.Time, error) {
	if m.failAccess {
		return "", time.Time{}, errors.New("fallo access token")
	}
	return "nuevo-access-token", time.Now().Add(15 * time.Minute), nil
}
func (m *mockTokenServicio) GenerarRefreshToken(usuarioID, sesionID string) (string, time.Time, error) {
	if m.failRefresh {
		return "", time.Time{}, errors.New("fallo refresh token")
	}
	return "nuevo-refresh-token", time.Now().Add(24 * time.Hour), nil
}
func (m *mockTokenServicio) ValidarAccessToken(token string) (*sesiones_domain.TokenClaims, error) {
	return nil, nil
}
func (m *mockTokenServicio) ValidarRefreshToken(token string) (*sesiones_domain.TokenClaims, error) {
	return m.claimsValidos, m.errValidar
}
func (m *mockTokenServicio) HashearToken(token string) string { return "hash:" + token }

// mockSesionRepo
type mockSesionRepo struct {
	sesionPorHash     *sesiones_domain.Sesion
	errPorHash        error
	errActualizar     error
	invalidadasUserID string
}

func (m *mockSesionRepo) Crear(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	return s, nil
}
func (m *mockSesionRepo) Actualizar(ctx context.Context, s *sesiones_domain.Sesion) (*sesiones_domain.Sesion, error) {
	if m.errActualizar != nil {
		return nil, m.errActualizar
	}
	return s, nil
}
func (m *mockSesionRepo) ObtenerPorID(ctx context.Context, id string) (*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) ObtenerPorRefreshTokenHash(ctx context.Context, hash string) (*sesiones_domain.Sesion, error) {
	return m.sesionPorHash, m.errPorHash
}
func (m *mockSesionRepo) ListarActivasPorUsuarioID(ctx context.Context, usuarioID string, ahora time.Time) ([]*sesiones_domain.Sesion, error) {
	return nil, nil
}
func (m *mockSesionRepo) InvalidarTodasPorUsuarioID(ctx context.Context, usuarioID string) error {
	m.invalidadasUserID = usuarioID
	return nil
}
func (m *mockSesionRepo) Eliminar(ctx context.Context, id string) error { return nil }

// ── Helpers ───────────────────────────────────────────────────────────────────

// claimsValidos retorna claims de prueba para un token válido.
func claimsValidos() *sesiones_domain.TokenClaims {
	return &sesiones_domain.TokenClaims{
		UsuarioID: "user-id-1",
		SesionID:  "sesion-id-1",
		Tipo:      "refresh",
		Expira:    time.Now().Add(24 * time.Hour),
	}
}

// sesionActiva crea una sesión activa de prueba.
// El refreshTokenHash debe coincidir con HashearToken("refresh-token-valido") = "hash:refresh-token-valido".
func sesionActiva() *sesiones_domain.Sesion {
	ahora := time.Now()
	s := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1",
		"user-id-1",
		"hash:access-token",
		"hash:refresh-token-valido",
		sesiones_domain.EstadoActiva,
		"192.168.1.1",
		ahora,
		ahora,
		ahora.Add(15*time.Minute),
		ahora.Add(24*time.Hour),
		ahora,
		0,
	)
	return s
}

// configDefault retorna una configuración sin límites para tests básicos.
func configDefault() refresh.ConfigRefresh {
	return refresh.ConfigRefresh{MaxRefrescos: 0, TimeoutAbsoluto: 0}
}

func uowValido(sesionRepo *mockSesionRepo, tokenSvc *mockTokenServicio) *mockUnitOfWork {
	return &mockUnitOfWork{sesionRepo: sesionRepo, tokenServicio: tokenSvc}
}

// ── Tests ────────────────────────────────────────────────────────────────────

// Escenario 1: Refresh exitoso
func TestRefresh_Exitoso(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{sesionPorHash: sesionActiva()},
		&mockTokenServicio{claimsValidos: claimsValidos()},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	resp, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{
		RefreshToken: "refresh-token-valido",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("tokens vacíos en respuesta exitosa")
	}
	if resp.UsuarioID != "user-id-1" {
		t.Errorf("usuarioID incorrecto: %v", resp.UsuarioID)
	}
}

// Escenario 2: Múltiples refrescos incrementan contador
func TestRefresh_MultiplesRefrescos(t *testing.T) {
	ahora := time.Now()
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh-token-valido",
		sesiones_domain.EstadoActiva, "", ahora, ahora,
		ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora, 2,
	)
	uow := uowValido(
		&mockSesionRepo{sesionPorHash: sesion},
		&mockTokenServicio{claimsValidos: claimsValidos()},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	resp, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{
		RefreshToken: "refresh-token-valido",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("esperaba tokens en respuesta")
	}
}

// Escenario 3: Refresh token vacío
func TestRefresh_TokenVacio(t *testing.T) {
	svc := refresh.NuevoServicioRefresh(&mockUnitOfWork{}, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: ""})
	if !errors.Is(err, refresh.ErrRefreshTokenRequerido) {
		t.Errorf("esperaba ErrRefreshTokenRequerido, got %v", err)
	}
}

// Escenario 4: Token JWT expirado
func TestRefresh_TokenExpirado(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{},
		&mockTokenServicio{errValidar: errors.New("token expirado")},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "token-expirado"})
	if !errors.Is(err, refresh.ErrTokenInvalido) {
		t.Errorf("esperaba ErrTokenInvalido, got %v", err)
	}
}

// Escenario 5: Token mal formado
func TestRefresh_TokenMalFormado(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{},
		&mockTokenServicio{errValidar: errors.New("token mal formado")},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "no-es-un-jwt"})
	if !errors.Is(err, refresh.ErrTokenInvalido) {
		t.Errorf("esperaba ErrTokenInvalido, got %v", err)
	}
}

// Escenario 6: Firma inválida
func TestRefresh_FirmaInvalida(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{},
		&mockTokenServicio{errValidar: errors.New("firma inválida")},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "token-alterado"})
	if !errors.Is(err, refresh.ErrTokenInvalido) {
		t.Errorf("esperaba ErrTokenInvalido, got %v", err)
	}
}

// Escenario 7: Sesión revocada
func TestRefresh_SesionRevocada(t *testing.T) {
	ahora := time.Now()
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh-token-valido",
		sesiones_domain.EstadoRevocada, "", ahora, ahora,
		ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora, 0,
	)
	uow := uowValido(
		&mockSesionRepo{sesionPorHash: sesion},
		&mockTokenServicio{claimsValidos: claimsValidos()},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "refresh-token-valido"})
	if !errors.Is(err, refresh.ErrSesionNoValida) {
		t.Errorf("esperaba ErrSesionNoValida, got %v", err)
	}
}

// Escenario 8: Sesión expirada
func TestRefresh_SesionExpirada(t *testing.T) {
	ahora := time.Now()
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh-token-valido",
		sesiones_domain.EstadoExpirada, "", ahora, ahora,
		ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora, 0,
	)
	uow := uowValido(
		&mockSesionRepo{sesionPorHash: sesion},
		&mockTokenServicio{claimsValidos: claimsValidos()},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "refresh-token-valido"})
	if !errors.Is(err, refresh.ErrSesionNoValida) {
		t.Errorf("esperaba ErrSesionNoValida, got %v", err)
	}
}

// Escenario 9 y 12: Hash no existe en BD → detección de robo
func TestRefresh_DeteccionRobo(t *testing.T) {
	sesionRepo := &mockSesionRepo{
		errPorHash: errors.New("no encontrado"),
	}
	uow := uowValido(sesionRepo, &mockTokenServicio{claimsValidos: claimsValidos()})
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "refresh-token-rotado"})
	if !errors.Is(err, refresh.ErrTokenInvalido) {
		t.Errorf("esperaba ErrTokenInvalido, got %v", err)
	}
	if sesionRepo.invalidadasUserID != "user-id-1" {
		t.Errorf("esperaba invalidar sesiones de user-id-1, got '%v'", sesionRepo.invalidadasUserID)
	}
}

// Escenario 10: Excede límite de refrescos
func TestRefresh_LimiteRefrescosAlcanzado(t *testing.T) {
	ahora := time.Now()
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh-token-valido",
		sesiones_domain.EstadoActiva, "", ahora, ahora,
		ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora, 10,
	)
	uow := uowValido(
		&mockSesionRepo{sesionPorHash: sesion},
		&mockTokenServicio{claimsValidos: claimsValidos()},
	)
	cfg := refresh.ConfigRefresh{MaxRefrescos: 10, TimeoutAbsoluto: 0}
	svc := refresh.NuevoServicioRefresh(uow, cfg)
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "refresh-token-valido"})
	if !errors.Is(err, refresh.ErrLimiteRefrescosAlcanzado) {
		t.Errorf("esperaba ErrLimiteRefrescosAlcanzado, got %v", err)
	}
}

// Escenario 11: Timeout absoluto de sesión
func TestRefresh_TimeoutAbsoluto(t *testing.T) {
	hace8dias := time.Now().Add(-8 * 24 * time.Hour)
	sesion := sesiones_domain.NuevaSesionDesdeBD(
		"sesion-id-1", "user-id-1", "hash:access", "hash:refresh-token-valido",
		sesiones_domain.EstadoActiva, "", hace8dias, hace8dias,
		time.Now().Add(15*time.Minute), time.Now().Add(24*time.Hour), hace8dias, 0,
	)
	uow := uowValido(
		&mockSesionRepo{sesionPorHash: sesion},
		&mockTokenServicio{claimsValidos: claimsValidos()},
	)
	cfg := refresh.ConfigRefresh{MaxRefrescos: 0, TimeoutAbsoluto: 7 * 24 * time.Hour}
	svc := refresh.NuevoServicioRefresh(uow, cfg)
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "refresh-token-valido"})
	if !errors.Is(err, refresh.ErrSesionAbsolutaExpirada) {
		t.Errorf("esperaba ErrSesionAbsolutaExpirada, got %v", err)
	}
}

// Escenario 14: Fallo al persistir sesión actualizada → rollback
func TestRefresh_FalloAlActualizar(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{
			sesionPorHash: sesionActiva(),
			errActualizar: errors.New("fallo bd"),
		},
		&mockTokenServicio{claimsValidos: claimsValidos()},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "refresh-token-valido"})
	if err == nil {
		t.Error("esperaba error al fallar la actualización de sesión")
	}
}

// Escenario 15: Fallo al generar access token
func TestRefresh_FalloAccessToken(t *testing.T) {
	uow := uowValido(
		&mockSesionRepo{sesionPorHash: sesionActiva()},
		&mockTokenServicio{claimsValidos: claimsValidos(), failAccess: true},
	)
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{RefreshToken: "refresh-token-valido"})
	if !errors.Is(err, refresh.ErrErrorGenerandoTokens) {
		t.Errorf("esperaba ErrErrorGenerandoTokens, got %v", err)
	}
}
// Escenario 13: Post-detección de robo no quedan sesiones activas
func TestRefresh_SinSesionesActivasPostDeteccion(t *testing.T) {
	sesionRepo := &mockSesionRepo{
		errPorHash: errors.New("no encontrado"),
	}
	uow := uowValido(sesionRepo, &mockTokenServicio{claimsValidos: claimsValidos()})
	svc := refresh.NuevoServicioRefresh(uow, configDefault())
	_, err := svc.Ejecutar(context.Background(), refresh.ComandoRefresh{
		RefreshToken: "refresh-token-rotado",
	})
	if !errors.Is(err, refresh.ErrTokenInvalido) {
		t.Errorf("esperaba ErrTokenInvalido, got %v", err)
	}
	// verificar que se intentó invalidar todas las sesiones del usuario
	if sesionRepo.invalidadasUserID == "" {
		t.Error("esperaba que se invalidaran las sesiones del usuario tras detección de robo")
	}
}