package usuario

import (
	"testing"
)

func TestNuevoCorreoElectronicoValido(t *testing.T) {
	vo, err := NuevoCorreoElectronico("juan@example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if vo == nil {
		t.Fatal("Expected CorreoElectronico, got nil")
	}
	if vo.Direccion() != "juan@example.com" {
		t.Errorf("Expected 'juan@example.com', got '%s'", vo.Direccion())
	}
	if vo.Estado() != PENDIENTE_VERIFICACION {
		t.Errorf("Expected PENDIENTE_VERIFICACION, got %s", vo.Estado())
	}
}

func TestNuevoCorreoElectronicoVacio(t *testing.T) {
	_, err := NuevoCorreoElectronico("")
	if err != ErrCorreoRequerido {
		t.Errorf("Expected ErrCorreoRequerido, got %v", err)
	}
}

func TestNuevoCorreoElectronicoFormatoInvalido(t *testing.T) {
	_, err := NuevoCorreoElectronico("invalido")
	if err != ErrCorreoFormatoInvalido {
		t.Errorf("Expected ErrCorreoFormatoInvalido, got %v", err)
	}
}

func TestNuevoCorreoElectronicoDesdeBD(t *testing.T) {
	vo := NuevoCorreoElectronicoDesdeBD("juan@example.com", VERIFICADO)
	if vo == nil {
		t.Fatal("Expected CorreoElectronico, got nil")
	}
	if vo.Direccion() != "juan@example.com" {
		t.Errorf("Expected 'juan@example.com', got '%s'", vo.Direccion())
	}
	if vo.Estado() != VERIFICADO {
		t.Errorf("Expected VERIFICADO, got %s", vo.Estado())
	}
}

func TestCorreoElectronicoVerificarDesdePendiente(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	err := vo.Verificar()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if vo.Estado() != VERIFICADO {
		t.Errorf("Expected VERIFICADO, got %s", vo.Estado())
	}
}

func TestCorreoElectronicoVerificarDesdeVerificado(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	vo.Verificar()
	err := vo.Verificar()
	if err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("Expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
}

func TestCorreoElectronicoVerificarDesdeEnlaceExpirado(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	vo.MarcarExpirado()
	err := vo.Verificar()
	if err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("Expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
}

func TestCorreoElectronicoMarcarExpiradoDesdePendiente(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	err := vo.MarcarExpirado()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if vo.Estado() != ENLACE_EXPIRADO {
		t.Errorf("Expected ENLACE_EXPIRADO, got %s", vo.Estado())
	}
}

func TestCorreoElectronicoSolicitarReenvioDesdePendiente(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	err := vo.SolicitarReenvio()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if vo.Estado() != REENVIO_SOLICITADO {
		t.Errorf("Expected REENVIO_SOLICITADO, got %s", vo.Estado())
	}
}

func TestCorreoElectronicoSolicitarReenvioDesdeEnlaceExpirado(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	vo.MarcarExpirado()
	err := vo.SolicitarReenvio()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if vo.Estado() != REENVIO_SOLICITADO {
		t.Errorf("Expected REENVIO_SOLICITADO, got %s", vo.Estado())
	}
}

func TestCorreoElectronicoVerificadoEsTerminal(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	vo.Verificar()

	if err := vo.Verificar(); err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("Verificar: expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
	if err := vo.MarcarExpirado(); err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("MarcarExpirado: expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
	if err := vo.SolicitarReenvio(); err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("SolicitarReenvio: expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
}

func TestCorreoElectronicoEstaVerificadoTrue(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	vo.Verificar()
	if !vo.EstaVerificado() {
		t.Error("Expected EstaVerificado() = true")
	}
}

func TestCorreoElectronicoEstaVerificadoFalse(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	if vo.EstaVerificado() {
		t.Error("Expected EstaVerificado() = false")
	}
}

func TestCorreoElectronicoEstaPendienteTrue(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	if !vo.EstaPendiente() {
		t.Error("Expected EstaPendiente() = true")
	}
}

func TestCorreoElectronicoEstaPendienteFalse(t *testing.T) {
	vo, _ := NuevoCorreoElectronico("juan@example.com")
	vo.Verificar()
	if vo.EstaPendiente() {
		t.Error("Expected EstaPendiente() = false")
	}
}