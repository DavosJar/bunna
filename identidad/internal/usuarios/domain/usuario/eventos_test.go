package usuario

import (
	"testing"
)

func TestNuevosEventosUsuario(t *testing.T) {
	ev := NuevosEventosUsuario()

	if ev == nil {
		t.Fatal("Expected EventosUsuario instance, got nil")
	}

	if ev.Cantidad() != 0 {
		t.Errorf("Expected 0 eventos, got %d", ev.Cantidad())
	}
}

func TestRegistrarCreacion(t *testing.T) {
	ev := NuevosEventosUsuario()

	ev.RegistrarCreacion("user-123", "test@example.com")

	if ev.Cantidad() != 1 {
		t.Errorf("Expected 1 evento, got %d", ev.Cantidad())
	}

	eventos := ev.Extraer()
	if len(eventos) != 1 {
		t.Fatalf("Expected 1 evento, got %d", len(eventos))
	}

	if eventos[0].Nombre != "UsuarioCreado" {
		t.Errorf("Expected 'UsuarioCreado', got '%s'", eventos[0].Nombre)
	}
}

func TestRegistrarBloqueo(t *testing.T) {
	ev := NuevosEventosUsuario()

	ev.RegistrarBloqueo("user-123")

	if ev.Cantidad() != 1 {
		t.Errorf("Expected 1 evento, got %d", ev.Cantidad())
	}

	eventos := ev.Extraer()
	if eventos[0].Nombre != "UsuarioBloqueado" {
		t.Errorf("Expected 'UsuarioBloqueado', got '%s'", eventos[0].Nombre)
	}
}

func TestRegistrarActivacion(t *testing.T) {
	ev := NuevosEventosUsuario()

	ev.RegistrarActivacion("user-123")

	eventos := ev.Extraer()
	if eventos[0].Nombre != "UsuarioActivado" {
		t.Errorf("Expected 'UsuarioActivado', got '%s'", eventos[0].Nombre)
	}
}

func TestRegistrarInactivacion(t *testing.T) {
	ev := NuevosEventosUsuario()

	ev.RegistrarInactivacion("user-123")

	eventos := ev.Extraer()
	if eventos[0].Nombre != "UsuarioInactivado" {
		t.Errorf("Expected 'UsuarioInactivado', got '%s'", eventos[0].Nombre)
	}
}

func TestRegistrarVerificacion(t *testing.T) {
	ev := NuevosEventosUsuario()

	ev.RegistrarVerificacion("user-123")

	eventos := ev.Extraer()
	if eventos[0].Nombre != "CorreoVerificado" {
		t.Errorf("Expected 'CorreoVerificado', got '%s'", eventos[0].Nombre)
	}
}

func TestRegistrarEventoPersonalizado(t *testing.T) {
	ev := NuevosEventosUsuario()

	payload := map[string]string{"accion": "envio_email"}
	ev.RegistrarEvento("CorreoEnviado", payload)

	eventos := ev.Extraer()
	if eventos[0].Nombre != "CorreoEnviado" {
		t.Errorf("Expected 'CorreoEnviado', got '%s'", eventos[0].Nombre)
	}
}

func TestExtraerLimpiaEventos(t *testing.T) {
	ev := NuevosEventosUsuario()

	ev.RegistrarBloqueo("user-123")
	ev.RegistrarActivacion("user-123")

	if ev.Cantidad() != 2 {
		t.Errorf("Expected 2 eventos, got %d", ev.Cantidad())
	}

	eventos := ev.Extraer()
	if len(eventos) != 2 {
		t.Errorf("Expected 2 eventos extraídos, got %d", len(eventos))
	}

	if ev.Cantidad() != 0 {
		t.Errorf("Expected 0 eventos después de Extraer, got %d", ev.Cantidad())
	}
}

func TestMultiplesEventos(t *testing.T) {
	ev := NuevosEventosUsuario()

	ev.RegistrarCreacion("user-123", "test@example.com")
	ev.RegistrarActivacion("user-123")
	ev.RegistrarBloqueo("user-123")
	ev.RegistrarEvento("NotificacionEnviada", map[string]string{"tipo": "sms"})

	if ev.Cantidad() != 4 {
		t.Errorf("Expected 4 eventos, got %d", ev.Cantidad())
	}

	eventos := ev.Extraer()

	expectedNames := []string{
		"UsuarioCreado",
		"UsuarioActivado",
		"UsuarioBloqueado",
		"NotificacionEnviada",
	}

	for i, expected := range expectedNames {
		if eventos[i].Nombre != expected {
			t.Errorf("Evento %d: expected '%s', got '%s'", i, expected, eventos[i].Nombre)
		}
	}
}
