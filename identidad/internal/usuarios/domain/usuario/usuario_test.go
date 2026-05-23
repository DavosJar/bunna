package usuario

import (
	"testing"
)

func TestNuevoUsuarioValido(t *testing.T) {
	u, err := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u == nil {
		t.Fatal("Expected Usuario instance, got nil")
	}
	if u.Nombre() != "Juan" {
		t.Errorf("Expected nombre 'Juan', got '%s'", u.Nombre())
	}
	if u.Apellido() != "Pérez" {
		t.Errorf("Expected apellido 'Pérez', got '%s'", u.Apellido())
	}
	if u.Correo() != "juan@example.com" {
		t.Errorf("Expected correo 'juan@example.com', got '%s'", u.Correo())
	}
	if u.Telefono() != "+34666666666" {
		t.Errorf("Expected telefono '+34666666666', got '%s'", u.Telefono())
	}
	if u.Estado() != NO_VERIFICADO {
		t.Errorf("Expected estado NO_VERIFICADO, got %s", u.Estado())
	}
}

func TestNuevoUsuarioSinID(t *testing.T) {
	u, err := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.ID() != "" {
		t.Errorf("Expected ID '', got '%s'", u.ID())
	}
}

func TestNuevoUsuarioSinNombre(t *testing.T) {
	u, err := NuevoUsuario("", "juan@example.com", "", "Pérez", "+34666666666")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Nombre() != "" {
		t.Errorf("Expected nombre '', got '%s'", u.Nombre())
	}
}

func TestNuevoUsuarioSinCorreo(t *testing.T) {
	_, err := NuevoUsuario("", "", "Juan", "Pérez", "+34666666666")
	if err != ErrCorreoRequerido {
		t.Errorf("Expected ErrCorreoRequerido, got %v", err)
	}
}

func TestNuevoUsuarioSinApellido(t *testing.T) {
	u, err := NuevoUsuario("", "juan@example.com", "Juan", "", "+34666666666")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Apellido() != "" {
		t.Errorf("Expected apellido '', got '%s'", u.Apellido())
	}
}

func TestNuevoUsuarioSinTelefono(t *testing.T) {
	u, err := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Telefono() != "" {
		t.Errorf("Expected telefono '', got '%s'", u.Telefono())
	}
}

func TestNuevoUsuarioGeneraEvento(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	eventos := u.PullEventos()
	if len(eventos) != 1 {
		t.Errorf("Expected 1 evento, got %d", len(eventos))
	}
	if eventos[0].Nombre != "UsuarioCreado" {
		t.Errorf("Expected 'UsuarioCreado', got '%s'", eventos[0].Nombre)
	}
}

func TestCambiarEstadoValido(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.PullEventos()

	err := u.CambiarEstado(ACTIVO)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Estado() != ACTIVO {
		t.Errorf("Expected estado ACTIVO, got %s", u.Estado())
	}
	eventos := u.PullEventos()
	if len(eventos) != 1 {
		t.Errorf("Expected 1 evento, got %d", len(eventos))
	}
	if eventos[0].Nombre != "EstadoUsuarioCambiado" {
		t.Errorf("Expected 'EstadoUsuarioCambiado', got '%s'", eventos[0].Nombre)
	}
}

func TestCambiarEstadoInvalido(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	err := u.CambiarEstado(INACTIVO)
	if err != ErrTransicionNoPermitida {
		t.Errorf("Expected ErrTransicionNoPermitida, got %v", err)
	}
	if u.Estado() != NO_VERIFICADO {
		t.Errorf("Expected estado NO_VERIFICADO, got %s", u.Estado())
	}
}

func TestBloquearUsuario(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.CambiarEstado(ACTIVO)
	u.PullEventos()

	err := u.Bloquear()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Estado() != BLOQUEADO {
		t.Errorf("Expected estado BLOQUEADO, got %s", u.Estado())
	}
	eventos := u.PullEventos()
	found := false
	for _, e := range eventos {
		if e.Nombre == "UsuarioBloqueado" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'UsuarioBloqueado' evento")
	}
}

func TestActivarUsuario(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.CambiarEstado(ACTIVO)
	u.CambiarEstado(INACTIVO)
	u.PullEventos()

	err := u.Activar()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Estado() != ACTIVO {
		t.Errorf("Expected estado ACTIVO, got %s", u.Estado())
	}
	eventos := u.PullEventos()
	found := false
	for _, e := range eventos {
		if e.Nombre == "UsuarioActivado" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'UsuarioActivado' evento")
	}
}

func TestActivarUsuarioDesdeBloqueado(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.CambiarEstado(ACTIVO)
	u.CambiarEstado(BLOQUEADO)
	u.PullEventos()

	err := u.Activar()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Estado() != ACTIVO {
		t.Errorf("Expected estado ACTIVO, got %s", u.Estado())
	}
}

func TestPendienteEliminacionEsTerminal(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.CambiarEstado(ACTIVO)
	u.CambiarEstado(PENDIENTE_DE_ELIMINACION)
	u.PullEventos()

	estadosInvalidos := []EstadoUsuario{NO_VERIFICADO, ACTIVO, INACTIVO, BLOQUEADO}
	for _, destino := range estadosInvalidos {
		err := u.CambiarEstado(destino)
		if err != ErrTransicionNoPermitida {
			t.Errorf("PENDIENTE_DE_ELIMINACION -> %s: expected ErrTransicionNoPermitida, got %v", destino, err)
		}
	}
}

func TestInactivarUsuario(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.CambiarEstado(ACTIVO)
	u.PullEventos()

	err := u.Inactivar()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Estado() != INACTIVO {
		t.Errorf("Expected estado INACTIVO, got %s", u.Estado())
	}
	eventos := u.PullEventos()
	found := false
	for _, e := range eventos {
		if e.Nombre == "UsuarioInactivado" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'UsuarioInactivado' evento")
	}
}

func TestNuevoUsuarioEstadoVerificacionInicial(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	if u.EstadoVerificacionCorreo() != PENDIENTE_VERIFICACION {
		t.Errorf("Expected estado verificacion PENDIENTE_VERIFICACION, got %s", u.EstadoVerificacionCorreo())
	}
}

func TestVerificarCorreoValidoDesdePendiente(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.PullEventos()

	err := u.VerificarCorreo()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.EstadoVerificacionCorreo() != VERIFICADO {
		t.Errorf("Expected estado VERIFICADO, got %s", u.EstadoVerificacionCorreo())
	}
	eventos := u.PullEventos()
	if len(eventos) != 1 || eventos[0].Nombre != "CorreoVerificado" {
		t.Error("Expected 'CorreoVerificado' evento")
	}
}

func TestVerificarCorreoInvalidoDesdeVerificado(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.VerificarCorreo()
	u.PullEventos()

	err := u.VerificarCorreo()
	if err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("Expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
}

func TestSolicitarReenvioValidoDesdePendiente(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.PullEventos()

	err := u.SolicitarReenvioVerificacion()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.EstadoVerificacionCorreo() != REENVIO_SOLICITADO {
		t.Errorf("Expected estado REENVIO_SOLICITADO, got %s", u.EstadoVerificacionCorreo())
	}
}

func TestSolicitarReenvioValidoDesdeEnlaceExpirado(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.MarcarEnlaceExpirado()
	u.PullEventos()

	err := u.SolicitarReenvioVerificacion()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.EstadoVerificacionCorreo() != REENVIO_SOLICITADO {
		t.Errorf("Expected estado REENVIO_SOLICITADO, got %s", u.EstadoVerificacionCorreo())
	}
}

func TestSolicitarReenvioInvalidoDesdeVerificado(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.VerificarCorreo()
	u.PullEventos()

	err := u.SolicitarReenvioVerificacion()
	if err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("Expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
}

func TestMarcarEnlaceExpiradoValidoDesdePendiente(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.PullEventos()

	err := u.MarcarEnlaceExpirado()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.EstadoVerificacionCorreo() != ENLACE_EXPIRADO {
		t.Errorf("Expected estado ENLACE_EXPIRADO, got %s", u.EstadoVerificacionCorreo())
	}
}

func TestMarcarEnlaceExpiradoInvalidoDesdeVerificado(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.VerificarCorreo()
	u.PullEventos()

	err := u.MarcarEnlaceExpirado()
	if err != ErrTransicionVerificacionNoPermitida {
		t.Errorf("Expected ErrTransicionVerificacionNoPermitida, got %v", err)
	}
}

func TestFlujoCompletoVerificacion(t *testing.T) {
	u, _ := NuevoUsuario("", "juan@example.com", "Juan", "Pérez", "+34666666666")
	u.PullEventos()

	err := u.MarcarEnlaceExpirado()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = u.SolicitarReenvioVerificacion()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = u.VerificarCorreo()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if u.EstadoVerificacionCorreo() != VERIFICADO {
		t.Errorf("Expected estado VERIFICADO, got %s", u.EstadoVerificacionCorreo())
	}
}