package domain_test

import (
	"testing"
	"time"

	domain "github.com/davosjar/bunna/services/identidad/internal/sesiones/domain"
)

func sesionValida(t *testing.T) *domain.Sesion {
	t.Helper()
	ahora := time.Now()
	s, err := domain.NuevaSesion(
		"sesion-id-1",
		"usuario-id-1",
		"hash-access-token",
		"hash-refresh-token",
		"192.168.1.1",
		ahora,
		ahora.Add(15*time.Minute),
		ahora.Add(24*time.Hour),
	)
	if err != nil {
		t.Fatalf("no se esperaba error al crear sesión válida: %v", err)
	}
	return s
}

// --- Escenario 1: Creación exitosa ---
func TestNuevaSesion_CreacionExitosa(t *testing.T) {
	ahora := time.Now()
	s, err := domain.NuevaSesion(
		"id-1", "user-1", "acc-hash", "ref-hash", "10.0.0.1",
		ahora, ahora.Add(15*time.Minute), ahora.Add(24*time.Hour),
	)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if s.Estado() != domain.EstadoActiva {
		t.Errorf("esperaba ACTIVA, got %v", s.Estado())
	}
	if !s.EstaActiva(ahora) {
		t.Error("esperaba EstaActiva=true")
	}
	if s.UsuarioID() != "user-1" {
		t.Errorf("usuarioID incorrecto: %v", s.UsuarioID())
	}
}

// --- Escenario 2: Reconstrucción desde BD ---
func TestNuevaSesionDesdeBD_Reconstruccion(t *testing.T) {
	ahora := time.Now()
	s := domain.NuevaSesionDesdeBD(
		"id-bd", "user-bd", "acc", "ref",
		domain.EstadoRevocada, "1.2.3.4",
		ahora, ahora, ahora.Add(time.Hour), ahora.Add(24*time.Hour),
		ahora, 3,
	)
	if s.Estado() != domain.EstadoRevocada {
		t.Errorf("esperaba REVOCADA, got %v", s.Estado())
	}
	if s.ContadorRefrescos() != 3 {
		t.Errorf("esperaba contadorRefrescos=3, got %v", s.ContadorRefrescos())
	}
}

// --- Escenario 3: usuarioID vacío ---
func TestNuevaSesion_UsuarioIDVacio(t *testing.T) {
	ahora := time.Now()
	_, err := domain.NuevaSesion("id", "", "acc", "ref", "", ahora, ahora.Add(time.Minute), ahora.Add(time.Hour))
	if err != domain.ErrUsuarioIDRequerido {
		t.Errorf("esperaba ErrUsuarioIDRequerido, got %v", err)
	}
}

// --- Escenario 4: refreshTokenHash vacío ---
func TestNuevaSesion_RefreshTokenHashVacio(t *testing.T) {
	ahora := time.Now()
	_, err := domain.NuevaSesion("id", "user", "acc", "", "", ahora, ahora.Add(time.Minute), ahora.Add(time.Hour))
	if err != domain.ErrRefreshTokenHashRequerido {
		t.Errorf("esperaba ErrRefreshTokenHashRequerido, got %v", err)
	}
}

// --- Escenario 5: accessTokenHash vacío ---
func TestNuevaSesion_AccessTokenHashVacio(t *testing.T) {
	ahora := time.Now()
	_, err := domain.NuevaSesion("id", "user", "", "ref", "", ahora, ahora.Add(time.Minute), ahora.Add(time.Hour))
	if err != domain.ErrAccessTokenHashRequerido {
		t.Errorf("esperaba ErrAccessTokenHashRequerido, got %v", err)
	}
}

// --- Escenario 6: Expiración en el pasado ---
func TestNuevaSesion_FechaExpiracionEnPasado(t *testing.T) {
	ahora := time.Now()
	_, err := domain.NuevaSesion("id", "user", "acc", "ref", "", ahora, ahora.Add(-1*time.Minute), ahora.Add(time.Hour))
	if err != domain.ErrFechaExpiracionInvalida {
		t.Errorf("esperaba ErrFechaExpiracionInvalida, got %v", err)
	}
}

// --- Escenario 7: IP vacía permitida ---
func TestNuevaSesion_IPOrigenVaciaPermitida(t *testing.T) {
	ahora := time.Now()
	_, err := domain.NuevaSesion("id", "user", "acc", "ref", "", ahora, ahora.Add(time.Minute), ahora.Add(time.Hour))
	if err != nil {
		t.Errorf("no se esperaba error con IP vacía, got %v", err)
	}
}

// --- Escenario 8: Sesión activa vigente ---
func TestEstaActiva_SesionActivaVigente(t *testing.T) {
	s := sesionValida(t)
	if !s.EstaActiva(time.Now()) {
		t.Error("esperaba EstaActiva=true")
	}
}

// --- Escenario 9: Sesión con fecha expirada, estado NO cambia ---
func TestEstaActiva_FechaExpirada(t *testing.T) {
	ahora := time.Now()
	s, _ := domain.NuevaSesion("id", "user", "acc", "ref", "", ahora, ahora.Add(time.Millisecond), ahora.Add(time.Hour))
	futuro := ahora.Add(time.Minute)
	if s.EstaActiva(futuro) {
		t.Error("esperaba EstaActiva=false cuando el access token expiró")
	}
	if s.Estado() != domain.EstadoActiva {
		t.Error("el estado no debe cambiar automáticamente")
	}
}

// --- Escenario 10: Sesión revocada ---
func TestEstaActiva_SesionRevocada(t *testing.T) {
	s := sesionValida(t)
	s.Revocar()
	if s.EstaActiva(time.Now()) {
		t.Error("esperaba EstaActiva=false para sesión REVOCADA")
	}
}

// --- Escenario 11: MarcarExpirada desde activa ---
func TestMarcarExpirada_DesdeActiva(t *testing.T) {
	s := sesionValida(t)
	err := s.MarcarExpirada()
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if s.Estado() != domain.EstadoExpirada {
		t.Errorf("esperaba EXPIRADA, got %v", s.Estado())
	}
}

// --- Escenario 12: Revocar desde activa ---
func TestRevocar_DesdeActiva(t *testing.T) {
	s := sesionValida(t)
	s.Revocar()
	if s.Estado() != domain.EstadoRevocada {
		t.Errorf("esperaba REVOCADA, got %v", s.Estado())
	}
}

// --- Escenario 13: MarcarExpirada cuando ya está REVOCADA ---
func TestMarcarExpirada_DesdeRevocada(t *testing.T) {
	s := sesionValida(t)
	s.Revocar()
	err := s.MarcarExpirada()
	if err != domain.ErrTransicionEstadoInvalida {
		t.Errorf("esperaba ErrTransicionEstadoInvalida, got %v", err)
	}
}

// --- Escenario 14: Revocar desde expirada ---
func TestRevocar_DesdeExpirada(t *testing.T) {
	s := sesionValida(t)
	_ = s.MarcarExpirada()
	s.Revocar()
	if s.Estado() != domain.EstadoRevocada {
		t.Errorf("esperaba REVOCADA, got %v", s.Estado())
	}
}

// --- Escenario 15: Refresh token vigente ---
func TestRefreshTokenValido_Vigente(t *testing.T) {
	s := sesionValida(t)
	if !s.RefreshTokenValido(time.Now()) {
		t.Error("esperaba RefreshTokenValido=true")
	}
}

// --- Escenario 16: Refresh token expirado ---
func TestRefreshTokenValido_Expirado(t *testing.T) {
	ahora := time.Now()
	s, _ := domain.NuevaSesion("id", "user", "acc", "ref", "", ahora, ahora.Add(time.Hour), ahora.Add(time.Millisecond))
	futuro := ahora.Add(time.Minute)
	if s.RefreshTokenValido(futuro) {
		t.Error("esperaba RefreshTokenValido=false cuando el refresh expiró")
	}
}

// --- Escenario 17: Refresh token en sesión revocada ---
func TestRefreshTokenValido_SesionRevocada(t *testing.T) {
	s := sesionValida(t)
	s.Revocar()
	if s.RefreshTokenValido(time.Now()) {
		t.Error("esperaba RefreshTokenValido=false para sesión REVOCADA")
	}
}

// --- Escenario 18: Refresh token con fecha zero ---
func TestRefreshTokenValido_FechaZero(t *testing.T) {
	ahora := time.Now()
	s := domain.NuevaSesionDesdeBD(
		"id", "user", "acc", "ref",
		domain.EstadoActiva, "",
		ahora, ahora, ahora.Add(time.Hour), time.Time{},
		ahora, 0,
	)
	if s.RefreshTokenValido(ahora) {
		t.Error("esperaba RefreshTokenValido=false con fecha zero")
	}
}

// --- Escenario 19: TokenPair válido ---
func TestNuevoTokenPair_Valido(t *testing.T) {
	ahora := time.Now()
	tp, err := domain.NuevoTokenPair("access-token", "refresh-token", ahora.Add(15*time.Minute), ahora.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if tp.AccessToken() != "access-token" {
		t.Errorf("AccessToken incorrecto: %v", tp.AccessToken())
	}
	if tp.RefreshToken() != "refresh-token" {
		t.Errorf("RefreshToken incorrecto: %v", tp.RefreshToken())
	}
}

// --- Escenario 20: TokenPair con access token vacío ---
func TestNuevoTokenPair_AccessTokenVacio(t *testing.T) {
	ahora := time.Now()
	_, err := domain.NuevoTokenPair("", "refresh", ahora.Add(time.Minute), ahora.Add(time.Hour))
	if err != domain.ErrAccessTokenRequerido {
		t.Errorf("esperaba ErrAccessTokenRequerido, got %v", err)
	}
}

// --- Escenario 21: TokenPair con refresh token vacío ---
func TestNuevoTokenPair_RefreshTokenVacio(t *testing.T) {
	ahora := time.Now()
	_, err := domain.NuevoTokenPair("access", "", ahora.Add(time.Minute), ahora.Add(time.Hour))
	if err != domain.ErrRefreshTokenRequerido {
		t.Errorf("esperaba ErrRefreshTokenRequerido, got %v", err)
	}
}

// --- RegistrarActividad ---
func TestRegistrarActividad(t *testing.T) {
	s := sesionValida(t)
	antes := s.UltimaActividad()
	time.Sleep(time.Millisecond)
	ahora := time.Now()
	s.RegistrarActividad(ahora)
	if !s.UltimaActividad().After(antes) {
		t.Error("esperaba que ultimaActividad se actualizara")
	}
}

// --- TimeoutExcedido cuando sí excede ---
func TestTimeoutExcedido_Excedido(t *testing.T) {
	s := sesionValida(t)
	futuro := time.Now().Add(2 * time.Hour)
	if !s.TimeoutExcedido(futuro, 30*time.Minute) {
		t.Error("esperaba TimeoutExcedido=true")
	}
}

// --- TimeoutExcedido cuando NO excede ---
func TestTimeoutExcedido_NoExcedido(t *testing.T) {
	s := sesionValida(t)
	if s.TimeoutExcedido(time.Now(), 30*time.Minute) {
		t.Error("esperaba TimeoutExcedido=false")
	}
}

// --- RotarTokens exitoso ---
func TestRotarTokens(t *testing.T) {
	s := sesionValida(t)
	ahora := time.Now()
	err := s.RotarTokens("nuevo-acc", "nuevo-ref", ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora)
	if err != nil {
		t.Fatalf("error inesperado al rotar: %v", err)
	}
	if s.AccessTokenHash() != "nuevo-acc" {
		t.Errorf("esperaba nuevo-acc, got %v", s.AccessTokenHash())
	}
	if s.RefreshTokenHash() != "nuevo-ref" {
		t.Errorf("esperaba nuevo-ref, got %v", s.RefreshTokenHash())
	}
	if s.ContadorRefrescos() != 1 {
		t.Errorf("esperaba contadorRefrescos=1, got %v", s.ContadorRefrescos())
	}
}

// --- TokenPair expiraciones ---
func TestNuevoTokenPair_Expiraciones(t *testing.T) {
	ahora := time.Now()
	expAccess := ahora.Add(15 * time.Minute)
	expRefresh := ahora.Add(24 * time.Hour)
	tp, _ := domain.NuevoTokenPair("access", "refresh", expAccess, expRefresh)
	if !tp.ExpiracionAccess().Equal(expAccess) {
		t.Error("ExpiracionAccess incorrecta")
	}
	if !tp.ExpiracionRefresh().Equal(expRefresh) {
		t.Error("ExpiracionRefresh incorrecta")
	}
}

// --- Tests del tester ---
func TestRotarTokens_SesionRevocada(t *testing.T) {
	s := sesionValida(t)
	s.Revocar()
	ahora := time.Now()
	err := s.RotarTokens("nuevo-acc", "nuevo-ref", ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora)
	if err != domain.ErrTransicionEstadoInvalida {
		t.Errorf("esperaba ErrTransicionEstadoInvalida, got %v", err)
	}
}

func TestRotarTokens_SesionExpirada(t *testing.T) {
	s := sesionValida(t)
	_ = s.MarcarExpirada()
	ahora := time.Now()
	err := s.RotarTokens("nuevo-acc", "nuevo-ref", ahora.Add(15*time.Minute), ahora.Add(24*time.Hour), ahora)
	if err != domain.ErrTransicionEstadoInvalida {
		t.Errorf("esperaba ErrTransicionEstadoInvalida, got %v", err)
	}
}

func TestMarcarExpirada_DesdeExpirada(t *testing.T) {
	s := sesionValida(t)
	_ = s.MarcarExpirada()
	err := s.MarcarExpirada()
	if err != domain.ErrTransicionEstadoInvalida {
		t.Errorf("esperaba ErrTransicionEstadoInvalida, got %v", err)
	}
}
