package rbac

import (
	"testing"
)

func TestRolesDeSistemaSon4(t *testing.T) {
	if len(RolesDeSistema) != 4 {
		t.Errorf("Expected 4 roles, got %d", len(RolesDeSistema))
	}
}

func TestTodosLosRolesSonDeSistema(t *testing.T) {
	for _, r := range RolesDeSistema {
		if !r.EsSistema {
			t.Errorf("Rol %s debe tener es_sistema=true", r.Nombre)
		}
	}
}

func TestSysAdminTieneTodosLosPermisos(t *testing.T) {
	var sysAdmin *RolInfo
	for i, r := range RolesDeSistema {
		if r.Nombre == RolSysAdmin {
			sysAdmin = &RolesDeSistema[i]
			break
		}
	}
	if sysAdmin == nil {
		t.Fatal("Rol sys_admin no encontrado")
	}
	if len(sysAdmin.Permisos) != 8 {
		t.Errorf("sys_admin debe tener 8 permisos, got %d", len(sysAdmin.Permisos))
	}
}

func TestCaficultorSoloPuedeConsultar(t *testing.T) {
	var caficultor *RolInfo
	for i, r := range RolesDeSistema {
		if r.Nombre == RolCaficultor {
			caficultor = &RolesDeSistema[i]
			break
		}
	}
	if caficultor == nil {
		t.Fatal("Rol caficultor no encontrado")
	}
	if len(caficultor.Permisos) != 1 {
		t.Errorf("caficultor debe tener 1 permiso, got %d", len(caficultor.Permisos))
	}
	if caficultor.Permisos[0] != PermisoUsuarioConsultar {
		t.Errorf("caficultor debe tener solo %s", PermisoUsuarioConsultar)
	}
}

func TestAgronomoNoTienePermisosDeRol(t *testing.T) {
	var agronomo *RolInfo
	for i, r := range RolesDeSistema {
		if r.Nombre == RolAgronomo {
			agronomo = &RolesDeSistema[i]
			break
		}
	}
	if agronomo == nil {
		t.Fatal("Rol agronomo no encontrado")
	}
	for _, p := range agronomo.Permisos {
		if p == PermisoRolAsignar || p == PermisoRolRevocar {
			t.Errorf("agronomo no debe tener permiso %s", p)
		}
	}
}

func TestRolesNombresConstantes(t *testing.T) {
	nombres := map[string]bool{
		RolSysAdmin:      false,
		RolAdministrador: false,
		RolAgronomo:      false,
		RolCaficultor:    false,
	}
	for _, r := range RolesDeSistema {
		if _, ok := nombres[r.Nombre]; !ok {
			t.Errorf("Rol inesperado: %s", r.Nombre)
		}
		nombres[r.Nombre] = true
	}
	for nombre, encontrado := range nombres {
		if !encontrado {
			t.Errorf("Rol no encontrado: %s", nombre)
		}
	}
}