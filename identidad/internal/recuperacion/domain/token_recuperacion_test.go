package recuperacion

import (
	"testing"
	"time"
)

func TestNuevoTokenRecuperacion(t *testing.T) {
	expira := time.Now().Add(time.Hour)
	token := NuevoTokenRecuperacion("id-1", "usuario-1", "token-plano", expira)

	if token.ID() != "id-1" {
		t.Errorf("Expected ID 'id-1', got '%s'", token.ID())
	}
	if token.UsuarioID() != "usuario-1" {
		t.Errorf("Expected UsuarioID 'usuario-1', got '%s'", token.UsuarioID())
	}
	if token.TokenHash() == "token-plano" {
		t.Error("Expected hash, no token en plano")
	}
	if token.TokenHash() == "" {
		t.Error("Expected hash no vacío")
	}
	if token.Usado() {
		t.Error("Expected Usado() = false para token nuevo")
	}
}

func TestTokenRecuperacionEsValidoVigente(t *testing.T) {
	expira := time.Now().Add(time.Hour)
	token := NuevoTokenRecuperacion("id-1", "usuario-1", "token", expira)

	if err := token.EsValido(time.Now()); err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestTokenRecuperacionEsValidoExpirado(t *testing.T) {
	expira := time.Now().Add(-time.Hour)
	token := NuevoTokenRecuperacion("id-1", "usuario-1", "token", expira)

	if err := token.EsValido(time.Now()); err != ErrEnlaceExpirado {
		t.Errorf("Expected ErrEnlaceExpirado, got %v", err)
	}
}

func TestTokenRecuperacionEsValidoYaUsado(t *testing.T) {
	expira := time.Now().Add(time.Hour)
	token := NuevoTokenRecuperacion("id-1", "usuario-1", "token", expira)
	token.Usar(time.Now())

	if err := token.EsValido(time.Now()); err != ErrEnlaceYaUtilizado {
		t.Errorf("Expected ErrEnlaceYaUtilizado, got %v", err)
	}
}

func TestTokenRecuperacionUsar(t *testing.T) {
	expira := time.Now().Add(time.Hour)
	token := NuevoTokenRecuperacion("id-1", "usuario-1", "token", expira)

	ahora := time.Now()
	token.Usar(ahora)

	if !token.Usado() {
		t.Error("Expected Usado() = true después de Usar()")
	}
	if token.UsadoEn() == nil {
		t.Error("Expected UsadoEn() no nil")
	}
}

func TestHashearTokenDeterministico(t *testing.T) {
	hash1 := HashearToken("mi-token")
	hash2 := HashearToken("mi-token")
	if hash1 != hash2 {
		t.Error("Expected mismo hash para mismo token")
	}
}

func TestHashearTokenDistinto(t *testing.T) {
	hash1 := HashearToken("token-a")
	hash2 := HashearToken("token-b")
	if hash1 == hash2 {
		t.Error("Expected hashes distintos para tokens distintos")
	}
}

func TestTokenRecuperacionDesdeBD(t *testing.T) {
	ahora := time.Now()
	expira := ahora.Add(time.Hour)
	hash := HashearToken("token-plano")

	token := NuevoTokenRecuperacionDesdeBD("id-1", "usuario-1", hash, expira, false, ahora, nil)

	if token.TokenHash() != hash {
		t.Errorf("Expected hash '%s', got '%s'", hash, token.TokenHash())
	}
	if token.Usado() {
		t.Error("Expected Usado() = false")
	}
}
