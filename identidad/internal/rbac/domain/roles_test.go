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
	expectedPermisos := len(TodosLosPermisos)
	if len(sysAdmin.Permisos) != expectedPermisos {
		t.Errorf("sys_admin debe tener %d permisos, got %d", expectedPermisos, len(sysAdmin.Permisos))
	}
}

func TestCaficultorPermisosDinamicos(t *testing.T) {
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
	if len(caficultor.Permisos) != 0 {
		t.Errorf("caficultor debe tener 0 permisos (dinámicos), got %d", len(caficultor.Permisos))
	}
}

func TestAgronomoPermisosDinamicos(t *testing.T) {
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
	if len(agronomo.Permisos) != 0 {
		t.Errorf("agronomo debe tener 0 permisos (dinámicos), got %d", len(agronomo.Permisos))
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
