package rbac

import (
	"testing"
)

func TestTodosLosPermisosExisten(t *testing.T) {
	if len(TodosLosPermisos) != 21 {
		t.Errorf("Expected 21 permisos, got %d", len(TodosLosPermisos))
	}
}

func TestPermisosCodigosUnicos(t *testing.T) {
	vistos := make(map[string]bool)
	for _, p := range TodosLosPermisos {
		if vistos[p.Codigo] {
			t.Errorf("Permiso duplicado: %s", p.Codigo)
		}
		vistos[p.Codigo] = true
	}
}

func TestConstantesPermisos(t *testing.T) {
	esperados := []string{
		PermisoUsuarioCrear,
		PermisoUsuarioModificar,
		PermisoUsuarioEliminar,
		PermisoUsuarioConsultar,
		PermisoUsuarioResetearPassword,
		PermisoUsuarioExpulsar,
		PermisoCredencialesConsultar,
		PermisoCredencialesDesbloquear,
		PermisoRolAsignar,
		PermisoRolRevocar,
		PermisoRolCrear,
		PermisoRolModificar,
		PermisoRolEliminar,
		PermisoRolPermisoAsignar,
		PermisoRolPermisoRevocar,
		PermisoPermisoConsultar,
		PermisoSesionConsultar,
		PermisoSesionForzarCierre,
		PermisoTenantConfigurar,
		PermisoIPBloqueadaConsultar,
		PermisoIPDesbloquear,
	}
	for _, codigo := range esperados {
		if codigo == "" {
			t.Errorf("Permiso constante vacío")
		}
	}
}

func TestTodosLosPermisosConModulo(t *testing.T) {
	for _, p := range TodosLosPermisos {
		if p.Modulo == "" {
			t.Errorf("Permiso sin modulo: %s", p.Codigo)
		}
		if p.Nombre == "" {
			t.Errorf("Permiso sin nombre: %s", p.Codigo)
		}
	}
}
