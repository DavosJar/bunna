package tenant

import (
	"testing"
	"time"
)

func TestNuevoTenantValido(t *testing.T) {
	ten, err := NuevoTenant("id-1", "Finca La Esperanza", "finca-la-esperanza")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ten.Nombre() != "Finca La Esperanza" {
		t.Errorf("Expected nombre 'Finca La Esperanza', got '%s'", ten.Nombre())
	}
	if ten.Slug() != "finca-la-esperanza" {
		t.Errorf("Expected slug 'finca-la-esperanza', got '%s'", ten.Slug())
	}
	if !ten.EstaActivo() {
		t.Error("Expected tenant activo=true")
	}
}

func TestNuevoTenantSinNombre(t *testing.T) {
	_, err := NuevoTenant("id-1", "", "finca-la-esperanza")
	if err != ErrNombreRequerido {
		t.Errorf("Expected ErrNombreRequerido, got %v", err)
	}
}

func TestNuevoTenantSinSlug(t *testing.T) {
	_, err := NuevoTenant("id-1", "Finca La Esperanza", "")
	if err != ErrSlugRequerido {
		t.Errorf("Expected ErrSlugRequerido, got %v", err)
	}
}

func TestNuevoTenantSlugConMayusculas(t *testing.T) {
	_, err := NuevoTenant("id-1", "Finca", "Finca-La-Esperanza")
	if err != ErrSlugInvalido {
		t.Errorf("Expected ErrSlugInvalido, got %v", err)
	}
}

func TestNuevoTenantSlugConEspacios(t *testing.T) {
	_, err := NuevoTenant("id-1", "Finca", "finca la esperanza")
	if err != ErrSlugInvalido {
		t.Errorf("Expected ErrSlugInvalido, got %v", err)
	}
}

func TestNuevoTenantSlugConCaracteresEspeciales(t *testing.T) {
	_, err := NuevoTenant("id-1", "Finca", "finca@esperanza")
	if err != ErrSlugInvalido {
		t.Errorf("Expected ErrSlugInvalido, got %v", err)
	}
}

func TestNuevoTenantSlugConNumeros(t *testing.T) {
	ten, err := NuevoTenant("id-1", "Finca 2024", "finca-2024")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ten.Slug() != "finca-2024" {
		t.Errorf("Expected slug 'finca-2024', got '%s'", ten.Slug())
	}
}

func TestNuevoTenantDesdeBD(t *testing.T) {
	ahora := time.Now()
	ten := NuevoTenantDesdeBD("id-1", "Finca", "finca", false, ahora, ahora)
	if ten == nil {
		t.Fatal("Expected Tenant, got nil")
	}
	if ten.EstaActivo() {
		t.Error("Expected tenant inactivo")
	}
}

func TestActivarTenantInactivo(t *testing.T) {
	ahora := time.Now()
	ten := NuevoTenantDesdeBD("id-1", "Finca", "finca", false, ahora, ahora)
	err := ten.Activar()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !ten.EstaActivo() {
		t.Error("Expected tenant activo=true")
	}
}

func TestActivarTenantYaActivo(t *testing.T) {
	ten, _ := NuevoTenant("id-1", "Finca", "finca")
	err := ten.Activar()
	if err != ErrTenantYaActivo {
		t.Errorf("Expected ErrTenantYaActivo, got %v", err)
	}
}

func TestDesactivarTenantActivo(t *testing.T) {
	ten, _ := NuevoTenant("id-1", "Finca", "finca")
	err := ten.Desactivar()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ten.EstaActivo() {
		t.Error("Expected tenant activo=false")
	}
}

func TestDesactivarTenantYaInactivo(t *testing.T) {
	ahora := time.Now()
	ten := NuevoTenantDesdeBD("id-1", "Finca", "finca", false, ahora, ahora)
	err := ten.Desactivar()
	if err != ErrTenantYaInactivo {
		t.Errorf("Expected ErrTenantYaInactivo, got %v", err)
	}
}

func TestNuevoTenantCreadoActivo(t *testing.T) {
	ten, _ := NuevoTenant("id-1", "Finca", "finca")
	if !ten.EstaActivo() {
		t.Error("Todo tenant nuevo debe crearse activo")
	}
}

func TestNuevoTenantFechasNoVacias(t *testing.T) {
	ten, _ := NuevoTenant("id-1", "Finca", "finca")
	if ten.FechaCreacion().IsZero() {
		t.Error("FechaCreacion no debe ser zero")
	}
	if ten.FechaActualizacion().IsZero() {
		t.Error("FechaActualizacion no debe ser zero")
	}
}

func TestNuevaMembresiaValida(t *testing.T) {
	m, err := NuevaMembresia("usuario-1", "tenant-1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if m.UsuarioID() != "usuario-1" {
		t.Errorf("Expected usuarioID 'usuario-1', got '%s'", m.UsuarioID())
	}
	if m.TenantID() != "tenant-1" {
		t.Errorf("Expected tenantID 'tenant-1', got '%s'", m.TenantID())
	}
}

func TestNuevaMembresiaUsuarioVacio(t *testing.T) {
	_, err := NuevaMembresia("", "tenant-1")
	if err == nil {
		t.Error("Expected error for usuarioID vacío")
	}
}