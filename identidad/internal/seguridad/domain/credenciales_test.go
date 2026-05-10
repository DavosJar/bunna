package domain

import (
	"testing"
	"time"
)

func TestNuevaCredencialesUsuario(t *testing.T) {
	// Arrange & Act
	cred := NuevaCredencialesUsuario("user123", "hash_password")

	// Assert
	if cred.usuarioID != "user123" {
		t.Errorf("expected usuarioID 'user123', got '%s'", cred.usuarioID)
	}
	if cred.passwordHash != "hash_password" {
		t.Errorf("expected passwordHash 'hash_password', got '%s'", cred.passwordHash)
	}
	if !cred.activo {
		t.Error("expected activo to be true by default")
	}
	if cred.correoVerificado {
		t.Error("expected correoVerificado to be false by default")
	}
	if cred.intentosFallidos != 0 {
		t.Errorf("expected intentosFallidos 0, got %d", cred.intentosFallidos)
	}
	if !cred.bloqueadoHasta.IsZero() {
		t.Error("expected bloqueadoHasta to be zero time by default")
	}
}

func TestNuevaCredencialesUsuarioDesdeBD(t *testing.T) {
	// Arrange
	bloqueadoHasta := time.Now().Add(10 * time.Minute)

	// Act
	cred := NuevaCredencialesUsuarioDesdeBD(
		"user456",
		"hash_bd",
		false,
		true,
		3,
		bloqueadoHasta,
	)

	// Assert
	if cred.usuarioID != "user456" {
		t.Errorf("expected usuarioID 'user456', got '%s'", cred.usuarioID)
	}
	if cred.passwordHash != "hash_bd" {
		t.Errorf("expected passwordHash 'hash_bd', got '%s'", cred.passwordHash)
	}
	if cred.activo {
		t.Error("expected activo to be false")
	}
	if !cred.correoVerificado {
		t.Error("expected correoVerificado to be true")
	}
	if cred.intentosFallidos != 3 {
		t.Errorf("expected intentosFallidos 3, got %d", cred.intentosFallidos)
	}
	if cred.bloqueadoHasta != bloqueadoHasta {
		t.Error("expected bloqueadoHasta to match provided value")
	}
}

func TestVerificarPassword_Correcto(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash_correcto")

	// Act
	resultado := cred.VerificarPassword("hash_correcto")

	// Assert
	if !resultado {
		t.Error("expected VerificarPassword to return true for matching hash")
	}
}

func TestVerificarPassword_Incorrecto(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash_correcto")

	// Act
	resultado := cred.VerificarPassword("hash_incorrecto")

	// Assert
	if resultado {
		t.Error("expected VerificarPassword to return false for non-matching hash")
	}
}

func TestMarcarIntentoFallido_IncrementaContador(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	ahora := time.Now()

	// Act
	cred.MarcarIntentoFallido(ahora)
	cred.MarcarIntentoFallido(ahora)
	cred.MarcarIntentoFallido(ahora)

	// Assert
	if cred.intentosFallidos != 3 {
		t.Errorf("expected intentosFallidos 3, got %d", cred.intentosFallidos)
	}
	if !cred.bloqueadoHasta.IsZero() {
		t.Error("expected bloqueadoHasta to remain zero after 3 intentos")
	}
}

func TestMarcarIntentoFallido_BloqueaDespues5Intentos(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	ahora := time.Now()

	// Act - marcar 5 intentos fallidos
	for i := 0; i < 5; i++ {
		cred.MarcarIntentoFallido(ahora)
	}

	// Assert
	if cred.intentosFallidos != 5 {
		t.Errorf("expected intentosFallidos 5, got %d", cred.intentosFallidos)
	}
	if cred.bloqueadoHasta.IsZero() {
		t.Error("expected bloqueadoHasta to be set after 5 intentos")
	}
	// Verificar que el bloqueo es aproximadamente 15 minutos desde ahora
	expectedMin := ahora.Add(14 * time.Minute)
	expectedMax := ahora.Add(16 * time.Minute)
	if cred.bloqueadoHasta.Before(expectedMin) || cred.bloqueadoHasta.After(expectedMax) {
		t.Errorf("expected bloqueadoHasta to be ~15 minutes from now, got %v", cred.bloqueadoHasta)
	}
}

func TestResetearIntentos_LimpiaBloqueoyContador(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	ahora := time.Now()
	for i := 0; i < 5; i++ {
		cred.MarcarIntentoFallido(ahora)
	}
	if cred.intentosFallidos != 5 {
		t.Fatal("expected intentosFallidos to be 5 before reset")
	}
	if cred.bloqueadoHasta.IsZero() {
		t.Fatal("expected bloqueadoHasta to be set before reset")
	}

	// Act
	cred.ResetearIntentos()

	// Assert
	if cred.intentosFallidos != 0 {
		t.Errorf("expected intentosFallidos 0 after reset, got %d", cred.intentosFallidos)
	}
	if !cred.bloqueadoHasta.IsZero() {
		t.Error("expected bloqueadoHasta to be zero after reset")
	}
}

func TestEstaBloqueado_DentroDeTiempo(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	ahora := time.Now()
	futuro := ahora.Add(5 * time.Minute)

	// Act - asignar un bloqueo en el futuro
	cred.bloqueadoHasta = futuro
	resultado := cred.EstaBloqueado(ahora)

	// Assert
	if !resultado {
		t.Error("expected EstaBloqueado to return true when bloqueadoHasta is in the future")
	}
}

func TestEstaBloqueado_FueradeTiempo(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	ahora := time.Now()

	// bloqueadoHasta es tiempo cero por defecto, que está en el pasado
	resultado := cred.EstaBloqueado(ahora)

	// Assert
	if resultado {
		t.Error("expected EstaBloqueado to return false when bloqueadoHasta is zero or in the past")
	}
}

func TestEstaBloqueado_DespuesDeTiempoBloqueo(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	ahora := time.Now()
	pasado := ahora.Add(-5 * time.Minute)

	// Act - asignar un bloqueo en el pasado
	cred.bloqueadoHasta = pasado
	resultado := cred.EstaBloqueado(ahora)

	// Assert
	if resultado {
		t.Error("expected EstaBloqueado to return false when bloqueadoHasta is in the past")
	}
}

func TestVerificarCorreo_MarcaCorreoVerificado(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	if cred.correoVerificado {
		t.Fatal("expected correoVerificado to be false initially")
	}

	// Act
	cred.VerificarCorreo()

	// Assert
	if !cred.correoVerificado {
		t.Error("expected correoVerificado to be true after VerificarCorreo")
	}
}

func TestDesactivar_CambiaEstadoActivo(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	if !cred.activo {
		t.Fatal("expected activo to be true initially")
	}

	// Act
	cred.Desactivar()

	// Assert
	if cred.activo {
		t.Error("expected activo to be false after Desactivar")
	}
}

func TestActivar_CambiaEstadoActivo(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuarioDesdeBD("user1", "hash", false, false, 0, time.Time{})
	if cred.activo {
		t.Fatal("expected activo to be false initially")
	}

	// Act
	cred.Activar()

	// Assert
	if !cred.activo {
		t.Error("expected activo to be true after Activar")
	}
}

func TestActivar_SiYaEstaActivo(t *testing.T) {
	// Arrange
	cred := NuevaCredencialesUsuario("user1", "hash")
	if !cred.activo {
		t.Fatal("expected activo to be true initially")
	}

	// Act
	cred.Activar()

	// Assert
	if !cred.activo {
		t.Error("expected activo to remain true after Activar")
	}
}
