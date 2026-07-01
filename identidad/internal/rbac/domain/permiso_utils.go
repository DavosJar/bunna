package rbac

import (
	"context"
	"log"
	"sync"
)

var (
	permisosSistemaMu  sync.RWMutex
	permisosSistemaMap map[string]map[string]bool // rol.Nombre → codigoPermiso → true
	permisosSistemaOk  bool
)

// initPermisosSistema carga los permisos estáticos de identidad definidos en roles.go.
// Se llama la primera vez que se consulta el mapa.
func initPermisosSistema() {
	permisosSistemaMu.Lock()
	defer permisosSistemaMu.Unlock()
	if permisosSistemaOk {
		return
	}
	permisosSistemaMap = make(map[string]map[string]bool, len(RolesDeSistema))
	for _, r := range RolesDeSistema {
		perms := make(map[string]bool, len(r.Permisos))
		for _, p := range r.Permisos {
			perms[p] = true
		}
		permisosSistemaMap[r.Nombre] = perms
	}
	permisosSistemaOk = true
}

// RegistrarPermisosDeModulo agrega permisos de un módulo externo (ej. fincas)
// al mapa en memoria de un rol de sistema. Se llama desde el consumer de Kafka
// cuando llegan permisos de otro servicio.
//
// Después de esta llamada, TienePermisoEnRol retorna true para estos códigos
// sin necesidad de consultar la BD.
func RegistrarPermisosDeModulo(rolNombre string, codigos []string) {
	initPermisosSistema() // asegura que el mapa base existe

	permisosSistemaMu.Lock()
	defer permisosSistemaMu.Unlock()

	perms, ok := permisosSistemaMap[rolNombre]
	if !ok {
		perms = make(map[string]bool)
		permisosSistemaMap[rolNombre] = perms
	}
	for _, c := range codigos {
		perms[c] = true
	}
	log.Printf("[RBAC] %d permisos registrados en memoria para rol '%s': %v", len(codigos), rolNombre, codigos)
}

// TienePermisoEnRol verifica si un rol tiene un permiso específico.
// Para roles de sistema consulta el mapa en memoria (O(1)).
// Para roles personalizados consulta el repositorio.
func TienePermisoEnRol(ctx context.Context, permisoRepo PermisoRepositorio, rol *RolDB, codigoPermiso string, tenantID string) bool {
	if rol.EsSistema {
		initPermisosSistema()

		permisosSistemaMu.RLock()
		perms, ok := permisosSistemaMap[rol.Nombre]
		encontrado := ok && perms[codigoPermiso]
		permisosSistemaMu.RUnlock()

		if encontrado {
			return true
		}
		// Roles de sistema con mapa vacío (agrónomo, caficultor) → consultar BD
	}

	// Roles personalizados o roles de sistema sin permisos estáticos → BD
	permisos, err := permisoRepo.ListarPorRol(ctx, rol.ID, tenantID)
	if err != nil {
		return false
	}
	for _, p := range permisos {
		if p.Codigo == codigoPermiso {
			return true
		}
	}
	return false
}
