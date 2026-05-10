package bcrypt

import (
	"testing"
)

// TestBcryptEncriptacionHashear verifica que Hashear genera un hash válido
func TestBcryptEncriptacionHashear(t *testing.T) {
	encriptacion := NewBcryptEncriptacion(12)

	password := "testPassword123!"
	hash, err := encriptacion.Hashear(password)

	if err != nil {
		t.Fatalf("Error al hashear password: %v", err)
	}

	if hash == "" {
		t.Error("Hash no debería estar vacío")
	}

	if hash == password {
		t.Error("Hash no debería ser igual al password")
	}
}

// TestBcryptEncriptacionHashearGeneraDiferentesHashes verifica que cada hash es único (usa salt)
func TestBcryptEncriptacionHashearGeneraDiferentesHashes(t *testing.T) {
	encriptacion := NewBcryptEncriptacion(12)

	password := "testPassword123!"
	hash1, err1 := encriptacion.Hashear(password)
	hash2, err2 := encriptacion.Hashear(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("Error al hashear password: %v, %v", err1, err2)
	}

	if hash1 == hash2 {
		t.Error("Dos hashes del mismo password no deberían ser idénticos (bcrypt usa salt)")
	}
}

// TestBcryptEncriptacionVerificarPasswordCorrecto verifica que Verificar retorna true con password correcto
func TestBcryptEncriptacionVerificarPasswordCorrecto(t *testing.T) {
	encriptacion := NewBcryptEncriptacion(12)

	password := "testPassword123!"
	hash, err := encriptacion.Hashear(password)
	if err != nil {
		t.Fatalf("Error al hashear password: %v", err)
	}

	if !encriptacion.Verificar(password, hash) {
		t.Error("Verificar debería retornar true con password correcto")
	}
}

// TestBcryptEncriptacionVerificarPasswordIncorrecto verifica que Verificar retorna false con password incorrecto
func TestBcryptEncriptacionVerificarPasswordIncorrecto(t *testing.T) {
	encriptacion := NewBcryptEncriptacion(12)

	password := "testPassword123!"
	hash, err := encriptacion.Hashear(password)
	if err != nil {
		t.Fatalf("Error al hashear password: %v", err)
	}

	if encriptacion.Verificar("wrongPassword", hash) {
		t.Error("Verificar debería retornar false con password incorrecto")
	}
}

// TestBcryptEncriptacionVerificarHashVacio verifica que Verificar retorna false con hash vacío
func TestBcryptEncriptacionVerificarHashVacio(t *testing.T) {
	encriptacion := NewBcryptEncriptacion(12)

	if encriptacion.Verificar("password", "") {
		t.Error("Verificar debería retornar false con hash vacío")
	}
}

// TestBcryptEncriptacionVerificarPasswordVacio verifica que Verificar retorna false con password vacío
func TestBcryptEncriptacionVerificarPasswordVacio(t *testing.T) {
	encriptacion := NewBcryptEncriptacion(12)

	hash, err := encriptacion.Hashear("password")
	if err != nil {
		t.Fatalf("Error al hashear password: %v", err)
	}

	if encriptacion.Verificar("", hash) {
		t.Error("Verificar debería retornar false con password vacío")
	}
}

// TestBcryptEncriptacionCostValido verifica que diferentes costs generan hashes válidos
func TestBcryptEncriptacionCostValido(t *testing.T) {
	password := "testPassword123!"

	// Probar con diferentes costos
	costs := []int{10, 11, 12, 13, 14}

	for _, cost := range costs {
		encriptacion := NewBcryptEncriptacion(cost)

		hash, err := encriptacion.Hashear(password)
		if err != nil {
			t.Fatalf("Error al hashear con cost %d: %v", cost, err)
		}

		if !encriptacion.Verificar(password, hash) {
			t.Errorf("Verificar falló con cost %d", cost)
		}
	}
}
