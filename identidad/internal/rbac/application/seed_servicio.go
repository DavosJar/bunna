package application

import (
	"context"
	"log"
	"strings"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

// SeedServicio siembra permisos y roles de sistema en BD (idempotente)
type SeedServicio struct {
	rolRepo        rbac.RolRepositorio
	permisoRepo    rbac.PermisoRepositorio
	rolPermisoRepo rbac.RolPermisoRepositorio
	idGenerator    shareddomain.GeneradorID
}

func NuevoSeedServicio(
	rolRepo rbac.RolRepositorio,
	permisoRepo rbac.PermisoRepositorio,
	rolPermisoRepo rbac.RolPermisoRepositorio,
	idGenerator shareddomain.GeneradorID,
) *SeedServicio {
	return &SeedServicio{
		rolRepo:        rolRepo,
		permisoRepo:    permisoRepo,
		rolPermisoRepo: rolPermisoRepo,
		idGenerator:    idGenerator,
	}
}

// Ejecutar siembra permisos y roles de sistema de forma idempotente
func (s *SeedServicio) Ejecutar(ctx context.Context) error {
	log.Println("[Seed] Iniciando seed de permisos y roles...")

	// 1. Sembrar permisos
	permisoIDs := make(map[string]string) // codigo → id
	for _, info := range rbac.TodosLosPermisos {
		existente, err := s.permisoRepo.ObtenerPorCodigo(ctx, info.Codigo)
		if err != nil {
			return err
		}
		if existente == nil {
			nuevoID, err := s.idGenerator.NextID(ctx)
			if err != nil {
				return err
			}
			permiso := &rbac.PermisoDB{
				ID:          nuevoID,
				Codigo:      info.Codigo,
				Nombre:      info.Nombre,
				Descripcion: info.Descripcion,
				Modulo:      info.Modulo,
			}
			if err := s.permisoRepo.Crear(ctx, permiso); err != nil {
				return err
			}
			permisoIDs[info.Codigo] = nuevoID
			log.Printf("[Seed] Permiso creado: %s", info.Codigo)
		} else {
			permisoIDs[info.Codigo] = existente.ID
			if existente.Nombre != info.Nombre || existente.Descripcion != info.Descripcion {
				if err := s.permisoRepo.ActualizarNombreDescripcion(ctx, existente.ID, info.Nombre, info.Descripcion); err != nil {
					return err
				}
				log.Printf("[Seed] Permiso actualizado: %s", info.Codigo)
			}
		}
	}

	// 2. Sembrar roles y sus permisos
	for _, rolInfo := range rbac.RolesDeSistema {
		existente, err := s.rolRepo.ObtenerPorNombre(ctx, rolInfo.Nombre)
		var rolID string

		if err == rbac.ErrRolNoEncontrado {
			nuevoID, err := s.idGenerator.NextID(ctx)
			if err != nil {
				return err
			}
			rol := &rbac.RolDB{
				ID:          nuevoID,
				Nombre:      rolInfo.Nombre,
				Descripcion: rolInfo.Descripcion,
				EsSistema:   rolInfo.EsSistema,
			}
			if err := s.rolRepo.Crear(ctx, rol); err != nil {
				return err
			}
			rolID = nuevoID
			log.Printf("[Seed] Rol creado: %s", rolInfo.Nombre)
		} else if err != nil {
			return err
		} else {
			rolID = existente.ID
			if existente.Descripcion != rolInfo.Descripcion {
				if err := s.rolRepo.ActualizarDescripcion(ctx, rolID, rolInfo.Descripcion); err != nil {
					return err
				}
			}
		}

		// 3. Limpiar permisos que ya no corresponden al rol
		permisosActuales, err := s.permisoRepo.ListarPorRol(ctx, rolID, rbac.TenantIDSistema)
		if err != nil {
			return err
		}

		deseados := make(map[string]bool)
		for _, c := range rolInfo.Permisos {
			deseados[c] = true
		}

		for _, p := range permisosActuales {
			if !deseados[p.Codigo] {
				// Solo limpiamos permisos que pertenezcan al módulo de identidad.
				// Los permisos de otros módulos (ej. fincas) se manejan vía Kafka y no deben ser eliminados por este seed.
				if strings.HasPrefix(p.Codigo, "identidad:") {
					if err := s.rolPermisoRepo.EliminarPermiso(ctx, rolID, p.ID, rbac.TenantIDSistema); err != nil {
						return err
					}
					log.Printf("[Seed] Permiso eliminado del rol %s: %s", rolInfo.Nombre, p.Codigo)
				}
			}
		}

		// 4. Re-sincronizar permisos del rol (upsert idempotente)
		for _, codigoPermiso := range rolInfo.Permisos {
			permisoID := permisoIDs[codigoPermiso]
			if err := s.rolPermisoRepo.AsignarPermiso(ctx, rolID, permisoID, rbac.TenantIDSistema, ""); err != nil {
				return err
			}
		}
		log.Printf("[Seed] Permisos sincronizados para rol: %s", rolInfo.Nombre)
	}

	log.Println("[Seed] Seed completado exitosamente")
	return nil
}
