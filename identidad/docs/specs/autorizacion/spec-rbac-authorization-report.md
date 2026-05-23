---
title: "Reporte de Implementación — RBAC: Roles, Permisos y Autorización"
version: 1.0
date: 2026-05-22
owner: Equipo Identidad
status: EN_PROGRESO
tags: reporte, rbac, roles, permisos, implementacion
---

# Reporte de Implementación: RBAC — Roles, Permisos y Autorización

> **Propósito**: Evaluar el estado actual del código contra lo especificado en `spec-rbac-authorization.md`.

## 1. Resumen Ejecutivo

| Dimensión | Resultado |
|-----------|-----------|
| **Estado general** | EN_PROGRESO |
| **Tablas BD nuevas** | creadas (modelos GORM implementados, pero SIN AutoMigrate) |
| **Permisos como constantes** | definidos (8 constantes) |
| **Roles de sistema** | definidos (4 roles con matriz de permisos) |
| **Seed de permisos/roles** | implementado (idempotente, upsert + re-sincronización) |
| **Servicio de autorización** | implementado (TienePermiso + ObtenerClaimsUsuario) |
| **Claims JWT enriquecidos** | implementado en TokenServicio, pero NO integrado en login/refresh |
| **Casos de uso de gestión** | NO implementados |
| **AutoMigrate en database.go** | NO implementado para modelos RBAC |
| **Registry (DI)** | NO implementado (faltan repos, servicios y seed) |
| **Build** | compila (verificado) |

## 2. Estado por Componente

### 2.1 Modelo de Datos (SQL)
- Tabla `permisos`: **existe** — `PermisoModel` en `internal/rbac/infrastructure/persistence/postgres/rbac_models.go`
- Tabla `roles`: **existe** — `RolModel` en `internal/rbac/infrastructure/persistence/postgres/rbac_models.go`
- Tabla `rol_permisos`: **existe** — `RolPermisoModel` en `internal/rbac/infrastructure/persistence/postgres/rbac_models.go`
- Tabla `usuario_roles`: **existe** — `UsuarioRolModel` en `internal/rbac/infrastructure/persistence/postgres/rbac_models.go`
- Tabla `usuario_tenant_roles`: **existe** — `UsuarioTenantRolModel` en `internal/rbac/infrastructure/persistence/postgres/rbac_models.go`
- Índices: **no creados** — no hay `CREATE INDEX` explícito en migraciones (GORM `AutoMigrate` los crea automáticamente para `uniqueIndex`, pero no los índices adicionales como `idx_rol_permisos_rol_id`)

### 2.2 Permisos del Sistema
- 8 permisos definidos como constantes de dominio: **sí**
  - `identidad:usuario:crear` (`PermisoUsuarioCrear`)
  - `identidad:usuario:modificar` (`PermisoUsuarioModificar`)
  - `identidad:usuario:eliminar` (`PermisoUsuarioEliminar`)
  - `identidad:usuario:consultar` (`PermisoUsuarioConsultar`)
  - `identidad:usuario:resetear_password` (`PermisoUsuarioResetearPassword`)
  - `identidad:rol:asignar` (`PermisoRolAsignar`)
  - `identidad:rol:revocar` (`PermisoRolRevocar`)
  - `identidad:permiso:consultar` (`PermisoPermisoConsultar`)
- Formato `modulo:recurso:verbo`: **cumple**
- Ubicación: `internal/rbac/domain/permisos.go`

### 2.3 Roles del Sistema
- 4 roles de sistema inmutables: **sí**
  - `sys_admin` — todos los permisos, ámbito global
  - `administrador` — todos los permisos, ámbito tenant
  - `agronomo` — crear, modificar, consultar usuarios + consultar permisos
  - `caficultor` — solo consultar usuarios
- Matriz de permisos por rol: **implementada** — `RolesDeSistema` en `internal/rbac/domain/roles.go`
- sys_admin con todos los permisos: **sí** (8 permisos)
- Ubicación: `internal/rbac/domain/roles.go`

### 2.4 Seed
- Estrategia de seed: **implementada**
- Idempotente: **sí** — upsert de permisos (busca por código, crea o actualiza), upsert de roles, re-sincronización de `rol_permisos` (limpia y re-inserta)
- Ubicación: `internal/rbac/application/seed_servicio.go`
- **NO se ejecuta automáticamente** — no está invocado en `main.go` ni en `cmd/main.go`, ni registrado en el Registry

### 2.5 Servicio de Autorización
- Consulta "¿Tiene permiso?" con tenant context: **implementado** — `TienePermiso(ctx, usuarioID, tenantID, codigoPermiso)`
- SYS_ADMIN bypass: **implementado** — si el usuario tiene rol `sys_admin`, retorna `true` sin verificar tenant
- Errores genéricos: **sí** — `ErrPermisoDenegado` es el error de dominio
- ObtenerClaimsUsuario para JWT: **implementado** — carga global flag, tenants, roles y permisos
- Ubicación: `internal/rbac/application/autorizacion_servicio.go`
- **NOTA**: No incluye validación de membresía al tenant (no inyecta `MembresiaRepositorio`)

### 2.6 Claims JWT
- Claim `global`: **incluido** en `claimsJWT` struct
- Claim `tenants` con roles/permisos: **incluido** en `claimsJWT` struct
- TokenClaims expone Global y Tenants: **sí**
- Integración con login/refresh: **NO implementada** — tanto `servicio_login.go` (línea 120) como `servicio_refresh.go` (línea 130) pasan `nil` como claims a `GenerarAccessToken`
- Middleware JWT extrae claims pero **NO expone** Global y Tenants en el contexto Gin (solo usuarioID y sesionID)

### 2.7 Casos de Uso (9.1 al 9.10)
- **Listar usuarios con permisos**: NO implementado
- **Crear usuario con rol**: NO implementado (existe `ServicioRegistro` pero no asigna roles en tenant)
- **Modificar/Eliminar usuario**: NO implementado
- **Cambiar/Resetear contraseña**: NO implementado
- **Asignar/Revocar rol**: NO implementado
- **Listar roles/permisos**: NO implementado

### 2.8 Integración
- AutoMigrate de tablas nuevas: **NO** — `config/database.go` no incluye `AutoMigrate` para `PermisoModel`, `RolModel`, `RolPermisoModel`, `UsuarioRolModel`, `UsuarioTenantRolModel`
- Registro en Registry: **NO** — `internal/registry/registry.go` no registra repositorios RBAC, `AutorizacionServicio`, `SeedServicio`, ni servicios de aplicación de gestión
- UnitOfWork extendido: **NO** — las interfaces `UnitOfWork` en sesiones y usuarios no incluyen accessors para repositorios RBAC
- Seed al iniciar: **NO** — no se ejecuta `SeedServicio.Ejecutar()` en `main.go` ni `cmd/main.go`

## 3. Checklist de Validación

| # | Ítem | Estado | Evidencia |
|---|------|--------|-----------|
| 1 | ¿Existen las tablas `permisos`, `roles`, `rol_permisos`? | ❌ Sin AutoMigrate | Modelos GORM existen en `rbac_models.go` pero no se migran en `database.go` |
| 2 | ¿Los 8 permisos están definidos como constantes de dominio? | ✅ | `internal/rbac/domain/permisos.go` — 8 constantes |
| 3 | ¿Los 4 roles de sistema existen y tienen `es_sistema = true`? | ✅ | `internal/rbac/domain/roles.go` — `RolesDeSistema` con `EsSistema: true` |
| 4 | ¿sys_admin tiene todos los permisos sin restricción de tenant? | ✅ | `roles.go` sys_admin tiene 8 permisos; `autorizacion_servicio.go` tiene SYS_ADMIN bypass |
| 5 | ¿La matriz de permisos por rol es correcta? | ✅ | `roles.go` coincide exactamente con la spec |
| 6 | ¿El seed es idempotente? | ✅ | `seed_servicio.go` — upsert + re-sincronización de permisos |
| 7 | ¿El servicio de autorización acepta tenantID? | ✅ | `TienePermiso(ctx, usuarioID, tenantID, codigoPermiso)` |
| 8 | ¿SYS_ADMIN siempre pasa la verificación sin importar tenant? | ✅ | Bypass en `TienePermiso` líneas 38-44 |
| 9 | ¿Los claims JWT incluyen `global` y `tenants`? | ✅ | `claimsJWT` struct y `TokenClaims` en dominio de sesiones |
| 10 | ¿Los errores de permiso son genéricos? | ✅ | `ErrPermisoDenegado = errors.New("permiso denegado")` |
| 11 | ¿Un usuario no puede asignarse permisos a sí mismo? | ❌ Caso de uso no implementado | Sin `ServicioAsignacionRoles` |
| 12 | ¿No hay tabla `usuario_permiso`? | ✅ | No existe modelo para `usuario_permiso` |
| 13 | ¿Los casos de uso verifican membresía al tenant antes de operar? | ❌ Casos de uso no implementados | — |
| 14 | ¿Las transacciones son atómicas via UnitOfWork? | ❌ Casos de uso no implementados | Los UnitOfWork existentes no incluyen repos RBAC |
| 15 | ¿El rol de sistema SYS_ADMIN no puede asignarse con tenantID? | ❌ Caso de uso no implementado | — |
| 16 | ¿La re-asignación de permisos en seed es completa (limpia y re-inserta)? | ✅ | `seed_servicio.go` — `LimpiarPermisosDeRol` + `AsignarPermiso` para cada permiso |

## 4. Brechas Detectadas

Listadas por prioridad de implementación:

### Prioridad ALTA (bloqueante para funcionalidad básica)

1. **AutoMigrate de modelos RBAC** — `config/database.go` no migra `PermisoModel`, `RolModel`, `RolPermisoModel`, `UsuarioRolModel`, `UsuarioTenantRolModel`. Sin esto, las tablas no se crean automáticamente.

2. **Registro en Registry** — `internal/registry/registry.go` no registra ninguna dependencia RBAC:
   - `RolRepositorio`, `PermisoRepositorio`, `RolPermisoRepositorio`, `UsuarioRolRepositorio`, `UsuarioTenantRolRepositorio`
   - `AutorizacionServicio`
   - `SeedServicio`
   - `TenantRepositorio` (existe en módulo tenants pero no en registry)
   - `MembresiaRepositorio` (existe en módulo tenants pero no en registry)

3. **Seed no ejecutado al iniciar** — `main.go` y `cmd/main.go` no invocan `SeedServicio.Ejecutar()`. Los permisos y roles nunca se siembran en BD.

4. **Login/Refresh sin claims JWT** — `servicio_login.go` y `servicio_refresh.go` pasan `nil` como claims a `GenerarAccessToken`. El `AutorizacionServicio` no está inyectado, por lo que los tokens JWT no incluyen información de autorización.

### Prioridad MEDIA (casos de uso core)

5. **Casos de uso de gestión no implementados** — No existen servicios de aplicación para:
   - `ServicioGestionUsuarios` (listar, crear, modificar, eliminar)
   - `ServicioCambiarPassword`
   - `ServicioResetearPassword`
   - `ServicioAsignacionRoles` (asignar, revocar, listar roles/permisos)

6. **UnitOfWork no extendido** — Las interfaces `UnitOfWork` existentes no exponen repositorios RBAC (usuario_roles, usuario_tenant_roles), necesarios para operaciones transaccionales de gestión de usuarios.

7. **Middleware JWT no expone Global/Tenants** — `jwt_middleware.go` solo inyecta `usuarioID` y `sesionID` en el contexto Gin, pero no `Global` ni `Tenants`. Los claims de autorización están presentes en el `TokenClaims` validado pero no se propagan a los handlers.

### Prioridad BAJA (mejoras)

8. **Sin tests para seed y autorización** — No existen archivos `*_test.go` en `internal/rbac/application/`. Solo hay tests de dominio (permisos, roles).

9. **Sin handlers HTTP para RBAC** — No hay handlers para endpoints de gestión de usuarios, roles o permisos. Solo existen `health_handler.go`, `login_handler.go` y `register_handler.go`.

10. **Sin router para nuevos endpoints** — `router.go` solo registra login, register y health.

11. **MembresiaRepositorio no disponible en AutorizacionServicio** — El servicio de autorización no puede verificar pertenencia a tenant.

## 5. Recomendaciones

### Inmediatas (para habilitar funcionalidad base)

1. **Agregar AutoMigrate en `database.go`** para los 5 modelos RBAC:
   ```go
   db.AutoMigrate(&rbac_postgres.PermisoModel{})
   db.AutoMigrate(&rbac_postgres.RolModel{})
   db.AutoMigrate(&rbac_postgres.RolPermisoModel{})
   db.AutoMigrate(&rbac_postgres.UsuarioRolModel{})
   db.AutoMigrate(&rbac_postgres.UsuarioTenantRolModel{})
   ```

2. **Extender Registry** con todas las dependencias RBAC:
   - Repositorios: `RolRepositorio`, `PermisoRepositorio`, `RolPermisoRepositorio`, `UsuarioRolRepositorio`, `UsuarioTenantRolRepositorio`
   - Servicios: `AutorizacionServicio`, `SeedServicio`
   - `TenantRepositorio` y `MembresiaRepositorio` desde el módulo tenants

3. **Ejecutar seed al inicio** en `main.go` tras migraciones:
   ```go
   seedSvc := rbac_application.NuevoSeedServicio(...)
   if err := seedSvc.Ejecutar(ctx); err != nil { ... }
   ```

4. **Integrar AutorizacionServicio en login y refresh**:
   - Inyectar en `ServicioLogin` y `ServicioRefresh`
   - En login: llamar `ObtenerClaimsUsuario(usuarioID)` antes de `GenerarAccessToken`
   - En refresh: mismo approach

### Corto plazo (casos de uso de gestión)

5. **Implementar casos de uso en orden**:
   - `ServicioGestionUsuarios` (listar, crear, modificar, eliminar)
   - `ServicioCambiarPassword` y `ServicioResetearPassword`
   - `ServicioAsignacionRoles` (asignar, revocar, listar roles/permisos)

6. **Extender UnitOfWork** con repositorios RBAC necesarios para operaciones transaccionales.

7. **Actualizar middleware JWT** para exponer `global` y `tenants` en el contexto Gin.

### Mediano plazo

8. **Agregar tests unitarios** para los servicios de aplicación (seed, autorización).

9. **Implementar handlers HTTP** para los nuevos casos de uso.

10. **Actualizar router** con los nuevos endpoints.

---

## Archivos relevantes por componente

### Implementados (capa dominio y aplicación RBAC)

| Archivo | Propósito |
|---------|-----------|
| `internal/rbac/domain/permisos.go` | 8 constantes de permiso + catálogo `TodosLosPermisos` |
| `internal/rbac/domain/permisos_test.go` | Tests de constantes y catálogo |
| `internal/rbac/domain/roles.go` | 4 roles de sistema + matriz de permisos |
| `internal/rbac/domain/roles_test.go` | Tests de roles y matriz |
| `internal/rbac/domain/claims.go` | `TenantClaims` y `UsuarioClaims` para JWT |
| `internal/rbac/domain/errors.go` | Todos los errores de dominio |
| `internal/rbac/domain/repositorios.go` | Interfaces de repositorios + `RolDB`/`PermisoDB` |
| `internal/rbac/application/autorizacion_servicio.go` | `TienePermiso` + `ObtenerClaimsUsuario` |
| `internal/rbac/application/seed_servicio.go` | Seed idempotente |
| `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` | 5 modelos GORM |
| `internal/rbac/infrastructure/persistence/postgres/rbac_repositorios.go` | 4 implementaciones de repositorios |

### Integración JWT

| Archivo | Estado |
|---------|--------|
| `internal/sesiones/domain/token_servicio.go` | ✅ Interfaz actualizada con `*rbac.UsuarioClaims` |
| `internal/sesiones/domain/tokens.go` | ✅ `TokenClaims` incluye `Global` y `Tenants` |
| `internal/sesiones/infrastructure/security/jwt/jwt_token_servicio.go` | ✅ Claims JWT enriquecidos |

### Por implementar

| Archivo | Brecha |
|---------|--------|
| `internal/config/database.go` | ❌ Faltan AutoMigrate para modelos RBAC |
| `internal/registry/registry.go` | ❌ Faltan dependencias RBAC |
| `main.go` / `cmd/main.go` | ❌ Falta ejecución de seed |
| `internal/sesiones/application/services/login/servicio_login.go` | ❌ `AutorizacionServicio` no inyectado; `nil` en claims |
| `internal/sesiones/application/services/refresh/servicio_refresh.go` | ❌ `AutorizacionServicio` no inyectado; `nil` en claims |
| `internal/presentation/middleware/jwt_middleware.go` | ❌ Global/Tenants no expuestos en contexto Gin |
| (no existe) | ❌ `ServicioGestionUsuarios` |
| (no existe) | ❌ `ServicioCambiarPassword` |
| (no existe) | ❌ `ServicioResetearPassword` |
| (no existe) | ❌ `ServicioAsignacionRoles` |
| `internal/usuarios/domain/usuario/unit_of_work.go` | ❌ Sin repos RBAC |
| `internal/sesiones/domain/unit_of_work.go` | ❌ Sin repos RBAC |
