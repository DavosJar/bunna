package verificacion

import (
	"testing"
	"time"
)

func TestNuevaPruebaVerificacion(t *testing.T) {
	expira := time.Now().Add(24 * time.Hour)
	prueba := NuevaPruebaVerificacion("token-123", expira)

	if !prueba.EstaPendiente() {
		t.Error("Expected EstaPendiente() = true")
	}
	if prueba.SecretoHash() == "" {
		t.Error("Expected hash no vacío")
	}
	if prueba.SecretoHash() == "token-123" {
		t.Error("Expected hash, no token en plano")
	}
}

func TestPruebaVerificacionCoincideCon(t *testing.T) {
	expira := time.Now().Add(24 * time.Hour)
	prueba := NuevaPruebaVerificacion("token-123", expira)

	if !prueba.CoincideCon("token-123") {
		t.Error("Expected CoincideCon('token-123') = true")
	}
	if prueba.CoincideCon("token-incorrecto") {
		t.Error("Expected CoincideCon('token-incorrecto') = false")
	}
}

func TestPruebaVerificacionExpiro(t *testing.T) {
	expirada := NuevaPruebaVerificacion("token", time.Now().Add(-1*time.Hour))
	if !expirada.Expiro(time.Now()) {
		t.Error("Expected Expiro() = true para token expirado")
	}

	vigente := NuevaPruebaVerificacion("token", time.Now().Add(24*time.Hour))
	if vigente.Expiro(time.Now()) {
		t.Error("Expected Expiro() = false para token vigente")
	}
}

func TestPruebaVerificacionVacia(t *testing.T) {
	prueba := PruebaVerificacionVacia()
	if prueba.EstaPendiente() {
		t.Error("Expected EstaPendiente() = false para prueba vacía")
	}
	if !prueba.Expiro(time.Now()) {
		t.Error("Expected Expiro() = true para prueba vacía")
	}
}

func TestNuevaPruebaVerificacionDesdeBD(t *testing.T) {
	hash := HashearToken("mi-token")
	expira := time.Now().Add(24 * time.Hour)
	prueba := NuevaPruebaVerificacionDesdeBD(hash, expira)

	if prueba.SecretoHash() != hash {
		t.Errorf("Expected hash '%s', got '%s'", hash, prueba.SecretoHash())
	}
	if !prueba.CoincideCon("mi-token") {
		t.Error("Expected CoincideCon('mi-token') = true")
	}
}

func TestHashearTokenEsDeterministico(t *testing.T) {
	hash1 := HashearToken("mismo-token")
	hash2 := HashearToken("mismo-token")
	if hash1 != hash2 {
		t.Error("Expected mismo hash para mismo token")
	}
}

func TestHashearTokenDistintosTokens(t *testing.T) {
	hash1 := HashearToken("token-a")
	hash2 := HashearToken("token-b")
	if hash1 == hash2 {
		t.Error("Expected hashes distintos para tokens distintos")
	}
}
