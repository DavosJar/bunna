---
title: "Reporte de Implementación — Gestión de Inquilinos (Multi-Tenant)"
version: 1.0
date: 2026-05-22
owner: Equipo Identidad
status: EN_PROGRESO
tags: reporte, tenant, implementacion
---

# Reporte de Implementación: Gestión de Inquilinos (Multi-Tenant)

> **Propósito**: Evaluar el estado actual del código contra lo especificado en `spec-tenant-management.md`.

## 1. Resumen Ejecutivo

| Dimensión | Resultado |
|-----------|-----------|
| **Estado general** | EN_PROGRESO |
| **Tablas BD** | Parcial — faltan migraciones AutoMigrate |
| **Modelos GORM** | Existen (parcial: faltan índices) |
| **Repositorios** | Existen (domain interfaces + implementación PostgreSQL) |
| **Servicios/Use Cases** | Existen (ServicioTenant con 11 métodos) |
| **Build** | Compila |

## 2. Estado por Componente

### 2.1 Modelo de Datos (SQL)
| Tabla | Estado | Archivo |
|-------|--------|---------|
| `tenants` | Existe (modelo GORM) | `internal/tenants/infrastructure/persistence/postgres/tenant_model.go` |
| `usuario_tenants` | Existe (modelo GORM) | `internal/tenants/infrastructure/persistence/postgres/tenant_model.go` |
| `usuario_tenant_roles` | Existe (modelo GORM) | `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` |
| `usuario_roles` | Existe (modelo GORM) | `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` |
| `roles` | Existe (modelo GORM) | `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` |
| `permisos` | Existe (modelo GORM) | `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` |
| `rol_permisos` | Existe (modelo GORM) | `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` |

**Índices SQL especificados**: NO creados. La spec define 7 índices explícitos (idx_tenants_slug, idx_tenants_activo, idx_usuario_tenants_usuario_id, idx_usuario_tenants_tenant_id, idx_usuario_tenant_roles_usuario_id, idx_usuario_tenant_roles_tenant_id, idx_usuario_tenant_roles_rol_id). GORM AutoMigrate los crearía si se agregan las anotaciones correspondientes. Actualmente:
- `TenantModel.Slug` tiene `uniqueIndex` (índice único para slug, cubre `idx_tenants_slug`)
- No hay índice para `tenants.activo`
- No hay índices explícitos en `MembresiaModel`
- No hay índices explícitos en `UsuarioTenantRolModel`

### 2.2 Capa de Dominio

| Componente | Estado | Archivo |
|------------|--------|---------|
| Entidad Tenant | ✅ Encontrado | `internal/tenants/domain/tenant/tenant.go` |
| Value Object Membresia | ✅ Encontrado | `internal/tenants/domain/tenant/membresia.go` |
| Errores de dominio | ✅ Encontrado (12 errores) | `internal/tenants/domain/tenant/errors.go` |
| Repositorio interface (TenantRepositorio) | ✅ Existe | `internal/tenants/domain/tenant/repositorio.go` |
| Repositorio interface (MembresiaRepositorio) | ✅ Existe | `internal/tenants/domain/tenant/repositorio.go` |
| Tests de dominio | ✅ Existen (14 tests) | `internal/tenants/domain/tenant/tenant_test.go` |

**Métodos de TenantRepositorio interface:**
- `Crear(ctx, *Tenant) (*Tenant, error)` ✅
- `ObtenerPorID(ctx, id) (*Tenant, error)` ✅
- `ObtenerPorSlug(ctx, slug) (*Tenant, error)` ✅
- `Actualizar(ctx, *Tenant) (*Tenant, error)` ✅
- `Listar(ctx) ([]*Tenant, error)` ✅
- `ListarPorUsuario(ctx, usuarioID) ([]*Tenant, error)` ✅

**Métodos de MembresiaRepositorio interface:**
- `Crear(ctx, *Membresia) error` ✅
- `Eliminar(ctx, usuarioID, tenantID) error` ✅
- `ExisteMiembro(ctx, usuarioID, tenantID) (bool, error)` ✅
- `ListarUsuariosPorTenant(ctx, tenantID) ([]string, error)` ✅

### 2.3 Capa de Aplicación

| Caso de Uso | Estado | Método en ServicioTenant |
|-------------|--------|--------------------------|
| CU-01: Crear tenant | ✅ Implementado | `CrearTenant` |
| CU-02: Listar todos los tenants | ✅ Implementado | `ListarTodos` |
| CU-03: Listar mis tenants | ✅ Implementado | `ListarMisTenants` |
| CU-04: Obtener tenant por ID | ✅ Implementado | `ObtenerPorID` |
| CU-05: Obtener tenant por slug | ✅ Implementado | `ObtenerPorSlug` |
| CU-06: Activar tenant | ✅ Implementado | `ActivarTenant` |
| CU-07: Desactivar tenant | ✅ Implementado | `DesactivarTenant` |
| CU-08: Agregar usuario a tenant | ✅ Implementado | `AgregarUsuario` |
| CU-09: Remover usuario de tenant | ✅ Implementado | `RemoverUsuario` |
| CU-10: Listar tenants de un usuario | ❌ No implementado | `ListarTenantsDeUsuario` (no existe) |
| CU-11: Listar usuarios de un tenant | ✅ Implementado | `ListarUsuariosDeTenant` |

**Comandos de aplicación:**
- `ComandoCrearTenant` ✅
- `ComandoActivarTenant` ✅
- `ComandoDesactivarTenant` ✅
- `ComandoAgregarUsuario` ✅
- `ComandoRemoverUsuario` ✅

**DTOs de respuesta:**
- `DtoTenant` ✅ (incluye ID, Nombre, Slug, Activo, FechaCreacion)

**Tests de aplicación:** ❌ NO EXISTEN. No hay archivos de test en `internal/tenants/application/services/gestionar_tenant/`.

### 2.4 Capa de Infraestructura

| Componente | Estado | Archivo |
|------------|--------|---------|
| TenantModel (GORM) | ✅ Existe | `internal/tenants/infrastructure/persistence/postgres/tenant_model.go` |
| MembresiaModel (GORM) | ✅ Existe | `internal/tenants/infrastructure/persistence/postgres/tenant_model.go` |
| UsuarioRolModel (GORM) | ✅ Existe | `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` |
| UsuarioTenantRolModel (GORM) | ✅ Existe | `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` |
| TenantRepositorio (PostgreSQL) | ✅ Existe | `internal/tenants/infrastructure/persistence/postgres/tenant_repositorio.go` |
| MembresiaRepositorio (PostgreSQL) | ✅ Existe | `internal/tenants/infrastructure/persistence/postgres/tenant_repositorio.go` |
| AutorizacionServicio | ✅ Existe (usa TenantRepositorio) | `internal/rbac/application/autorizacion_servicio.go` |
| RBAC repositorios (PostgreSQL) | ✅ Existen | `internal/rbac/infrastructure/persistence/postgres/rbac_repositorios.go` |

**Migraciones (AutoMigrate) en `internal/config/database.go`:**
- `usuarios_postgres.UsuarioModel` ✅
- `seguridad_postgres.CredencialesModel` ✅
- `sesiones_postgres.SesionModel` ✅
- `seguridad_postgres.IntentoIPModel` ✅
- `seguridad_postgres.RateLimitIPModel` ✅
- **`tenants_postgres.TenantModel`** ❌ NO INCLUIDO
- **`tenants_postgres.MembresiaModel`** ❌ NO INCLUIDO
- **`rbac_postgres.RolModel`** ❌ NO INCLUIDO
- **`rbac_postgres.PermisoModel`** ❌ NO INCLUIDO
- **`rbac_postgres.RolPermisoModel`** ❌ NO INCLUIDO
- **`rbac_postgres.UsuarioRolModel`** ❌ NO INCLUIDO
- **`rbac_postgres.UsuarioTenantRolModel`** ❌ NO INCLUIDO

### 2.5 Integración

| Componente | Estado | Detalle |
|------------|--------|---------|
| Registro en Registry | ❌ NO | `internal/registry/registry.go` no incluye ningún repositorio ni servicio de tenants ni de RBAC |
| Claims JWT (Global + Tenants) | ✅ Parcial | La estructura `claimsJWT` en `jwt_token_servicio.go` ya incluye `Global` y `Tenants`. La interfaz `TokenServicio.GenerarAccessToken` ya acepta `*rbac.UsuarioClaims`. |
| Carga de claims en login | ❌ NO | `servicio_login.go` pasa `nil` como claims a `GenerarAccessToken` (línea 120). No invoca `AutorizacionServicio.ObtenerClaimsUsuario`. |
| Middleware X-Tenant-ID | ❌ NO | No existe middleware que extraiga valide el header `X-Tenant-ID`. |
| Handlers HTTP de tenants | ❌ NO | No existen handlers para endpoints de tenant management. |
| Router (rutas de tenant) | ❌ NO | `router.go` no registra rutas de tenant. |
| Facade de tenants | ❌ NO | No existe una `TenantFacade` en `internal/presentation/facades/`. |
| Variable TENANT_CONTEXT_REQUIRED | ❌ NO | No definida en `internal/config/env.go` ni en `.env.example`. |
| Seed data (roles/permisos/tenants) | ❌ NO | No hay seed para roles, permisos, o tenants. |

## 3. Checklist de Validación

| # | Ítem | Estado | Evidencia |
|---|------|--------|-----------|
| 1 | ¿La tabla `tenants` existe con los campos especificados (id, nombre, slug, activo, created_at, updated_at)? | ✅ | `TenantModel` en `tenant_model.go` líneas 9-16 |
| 2 | ¿La tabla `usuario_tenants` tiene PK compuesta (usuario_id, tenant_id)? | ✅ | `MembresiaModel` con dos `primaryKey` en `tenant_model.go` líneas 46-47 |
| 3 | ¿La tabla `usuario_tenant_roles` tiene PK compuesta (usuario_id, tenant_id, rol_id)? | ✅ | `UsuarioTenantRolModel` con tres `primaryKey` en `rbac_models.go` líneas 48-53 |
| 4 | ¿La tabla `usuario_roles` existe para roles globales y solo se usa para SYS_ADMIN? | ✅ | `UsuarioRolModel` en `rbac_models.go` líneas 39-45 |
| 5 | ¿Un usuario puede pertenecer a múltiples tenants? | ✅ | Modelo de datos lo permite (PK compuesta no unique por usuario) |
| 6 | ¿Un usuario puede tener diferentes roles en diferentes tenants? | ✅ | `usuario_tenant_roles` asocia usuario+tenant+rol |
| 7 | ¿SYS_ADMIN es un rol global sin necesidad de membresía en tenant? | ✅ | `AutorizacionServicio.TienePermiso` verifica primero si es SYS_ADMIN (sin tenant) |
| 8 | ¿Desactivar un tenant impide nuevas operaciones en ese contexto? | ⚠️ Parcial | `ServicioTenant.AgregarUsuario` valida `!t.EstaActivo()`, pero no hay middleware global que lo valide |
| 9 | ¿El header `X-Tenant-ID` se propaga correctamente desde HTTP hasta la capa de dominio? | ❌ | No hay middleware de X-Tenant-ID |
| 10 | ¿Hay validación de que el usuario pertenece al tenant antes de operar (o es SYS_ADMIN)? | ❌ | No hay middleware de autorización multi-tenant |
| 11 | ¿Los índices están creados para las consultas principales (slug, usuario_id, tenant_id)? | ⚠️ Parcial | Solo `Slug` tiene `uniqueIndex` en GORM. Faltan índices explícitos para `activo`, `usuario_id`, `tenant_id` |
| 12 | ¿La eliminación de un usuario de un tenant también elimina sus roles en cascada? | ⚠️ Parcial | El modelo no tiene `constraint:OnDelete:CASCADE` explícito en GORM. La spec espera que la FK lo maneje |
| 13 | ¿El JWT incluye los claims `global` y `tenants` después del login? | ❌ | `servicio_login.go` pasa `nil` como claims. No se invoca `ObtenerClaimsUsuario` |
| 14 | ¿La variable `TENANT_CONTEXT_REQUIRED` se respeta en el middleware? | ❌ | No existe la variable ni el middleware |

## 4. Detalle de Brechas

### Brecha 1 (CRÍTICA): Migraciones BD no incluyen tablas de tenants/RBAC
**Archivo**: `internal/config/database.go`
**Impacto**: Las tablas `tenants`, `usuario_tenants`, `usuario_tenant_roles`, `usuario_roles`, `roles`, `permisos`, `rol_permisos` no se crean automáticamente al iniciar la aplicación.
**Acción requerida**: Agregar `AutoMigrate` para todos los modelos GORM de tenants y RBAC en `RunMigrations`.

### Brecha 2 (CRÍTICA): Registry no inyecta dependencias de tenants/RBAC
**Archivo**: `internal/registry/registry.go`
**Impacto**: `ServicioTenant`, `AutorizacionServicio`, y todos los repositorios de tenants/RBAC no están disponibles para las capas superiores.
**Acción requerida**: Agregar creación e inyección de:
- `postgres.NewTenantRepositorio(db)`
- `postgres.NewMembresiaRepositorio(db)`
- `postgres.NewRolRepositorio(db)`
- `postgres.NewPermisoRepositorio(db)`
- `postgres.NewRolPermisoRepositorio(db)`
- `postgres.NewUsuarioRolRepositorio(db)`
- `postgres.NewUsuarioTenantRolRepositorio(db)`
- `gestionar_tenant.NuevoServicioTenant(...)`
- `rbac_app.NuevoAutorizacionServicio(...)`

### Brecha 3 (ALTA): Login no carga claims multi-tenant
**Archivo**: `internal/sesiones/application/services/login/servicio_login.go` (línea 120)
**Impacto**: El JWT generado al hacer login no incluye información de tenants, roles ni permisos. Esto rompe el flujo de autorización multi-tenant porque el frontend/cliente no recibe el contexto de tenants del usuario.
**Acción requerida**: Inyectar `AutorizacionServicio` en `ServicioLogin` y llamar `ObtenerClaimsUsuario(usuarioID)` antes de generar el access token.

### Brecha 4 (ALTA): Sin middleware X-Tenant-ID
**Archivo**: No existe (`internal/presentation/middleware/` debería contenerlo)
**Impacto**: No hay forma de extraer y validar el header `X-Tenant-ID` de las peticiones HTTP. Toda la funcionalidad multi-tenant no es operable desde la API.
**Acción requerida**: Crear middleware `TenantMiddleware` que:
- Extraiga `X-Tenant-ID` del header
- Valide que el tenant existe y está activo
- Valide que el usuario pertenece al tenant (o es SYS_ADMIN)
- Inyecte tenantID en el contexto Gin
- Valide `TENANT_CONTEXT_REQUIRED`

### Brecha 5 (ALTA): Sin handlers HTTP de tenant management
**Archivo**: No existe (`internal/presentation/handlers/` debería contener `tenant_handler.go`)
**Impacto**: No hay endpoints REST para los 11 casos de uso de gestión de tenants.
**Acción requerida**: Crear `TenantHandler` y facade correspondiente que exponga:
- `POST /api/v1/tenants` (CU-01)
- `GET /api/v1/tenants` (CU-02)
- `GET /api/v1/me/tenants` (CU-03)
- `GET /api/v1/tenants/{id}` (CU-04)
- `GET /api/v1/tenants/slug/{slug}` (CU-05)
- `PUT /api/v1/tenants/{id}/activate` (CU-06)
- `PUT /api/v1/tenants/{id}/deactivate` (CU-07)
- `POST /api/v1/tenants/{id}/members` (CU-08)
- `DELETE /api/v1/tenants/{id}/members/{usuarioID}` (CU-09)
- `GET /api/v1/usuarios/{id}/tenants` (CU-10)
- `GET /api/v1/tenants/{id}/members` (CU-11)

### Brecha 6 (MEDIA): CU-10 no implementado
**Archivo**: `internal/tenants/application/services/gestionar_tenant/servicio_tenant.go`
**Impacto**: Falta el método `ListarTenantsDeUsuario(ctx, usuarioID)` para listar los tenants de un usuario específico (no el autenticado). El repositorio `ListarPorUsuario` existe, solo falta exponerlo en el servicio.
**Acción requerida**: Agregar método `ListarTenantsDeUsuario` en `ServicioTenant`.

### Brecha 7 (MEDIA): Índices de BD incompletos
**Archivo**: Modelos GORM
**Impacto**: Faltan índices `idx_tenants_activo`, `idx_usuario_tenants_usuario_id`, `idx_usuario_tenants_tenant_id`, `idx_usuario_tenant_roles_*`.
**Acción requerida**: Agregar anotaciones `gorm:"index"` en los campos correspondientes.

### Brecha 8 (MEDIA): Restricciones FK con ON DELETE CASCADE no explícitas
**Archivo**: Modelos GORM
**Impacto**: Las FKs se crean por convención de GORM pero no tienen `ON DELETE CASCADE` explícito.
**Acción requerida**: Agregar `constraint:OnDelete:CASCADE` en las relaciones.

### Brecha 9 (MEDIA): Sin tests de integración ni de aplicación
**Archivo**: No existen
**Impacto**: No hay tests para los 25 escenarios TDD especificados en la spec (secciones 8 y 8.1). Solo existen tests unitarios de dominio (14 tests).
**Acción requerida**: Implementar tests para:
- ServicioTenant (casos de uso CU-01 al CU-11)
- Repositorios PostgreSQL (integración)
- Handlers HTTP (e2e)
- Login con claims multi-tenant

### Brecha 10 (BAJA): Sin seed data
**Archivos**: `cmd/seed/`
**Impacto**: No hay datos iniciales de roles, permisos ni tenants de prueba.
**Acción requerida**: Crear seed que inserte roles de sistema, permisos, y opcionalmente un tenant de ejemplo.

### Brecha 11 (BAJA): TENANT_CONTEXT_REQUIRED no configurable
**Archivo**: `internal/config/env.go`
**Impacto**: La variable de entorno especificada en la sección 10 de la spec no existe en la configuración.
**Acción requerida**: Agregar `TenantContextRequired bool` al struct `Config` y su carga desde `TENANT_CONTEXT_REQUIRED`.

## 5. Recomendaciones

### Prioridad 1 (BLOQUEANTE — Imprescindible para operación básica)
1. **Agregar migraciones AutoMigrate** para todos los modelos de tenants y RBAC en `database.go`.
2. **Completar el Registry** con todas las dependencias de tenants y RBAC.
3. **Integrar login con claims multi-tenant**: inyectar `AutorizacionServicio` en `ServicioLogin`.

### Prioridad 2 (ALTA — Necesario para exponer la funcionalidad)
4. **Crear middleware X-Tenant-ID** con validación de tenant activo y membresía.
5. **Crear handlers HTTP** para los 11 casos de uso de gestión de tenants.
6. **Crear facade de tenants** en la capa de presentación.
7. **Registrar rutas de tenant** en el router.

### Prioridad 3 (MEDIA — Calidad y completitud)
8. **Implementar CU-10** (listar tenants de un usuario específico).
9. **Agregar índices faltantes** en los modelos GORM.
10. **Agregar constraints ON DELETE CASCADE** en las FKs.
11. **Implementar tests** por capas: aplicación (ServicioTenant), infraestructura (repositorios), handlers.
12. **Configurar TENANT_CONTEXT_REQUIRED** en env/config.

### Prioridad 4 (BAJA — Mejora continua)
13. **Crear seed data** para roles, permisos y tenants de prueba.
14. **Implementar escenarios TDD** de integración (login multi-tenant, operaciones con X-Tenant-ID).

## 6. Mapa de Archivos

### Archivos existentes (implementados)

| Archivo | Propósito |
|---------|-----------|
| `internal/tenants/domain/tenant/tenant.go` | Entidad Tenant |
| `internal/tenants/domain/tenant/membresia.go` | Value Object Membresia |
| `internal/tenants/domain/tenant/repositorio.go` | Interfaces repositorio |
| `internal/tenants/domain/tenant/errors.go` | Errores de dominio |
| `internal/tenants/domain/tenant/tenant_test.go` | Tests de dominio |
| `internal/tenants/application/services/gestionar_tenant/servicio_tenant.go` | Servicio de aplicación (casos de uso) |
| `internal/tenants/application/services/gestionar_tenant/comando.go` | Comandos CQRS |
| `internal/tenants/application/services/gestionar_tenant/respuesta.go` | DTOs de respuesta |
| `internal/tenants/infrastructure/persistence/postgres/tenant_model.go` | Modelos GORM (tenants, usuario_tenants) |
| `internal/tenants/infrastructure/persistence/postgres/tenant_repositorio.go` | Implementación PostgreSQL |
| `internal/rbac/domain/claims.go` | Claims JWT (TenantClaims, UsuarioClaims) |
| `internal/rbac/domain/repositorios.go` | Interfaces RBAC |
| `internal/rbac/domain/roles.go` | Roles del sistema |
| `internal/rbac/domain/permisos.go` | Permisos del sistema |
| `internal/rbac/application/autorizacion_servicio.go` | Servicio de autorización (TienePermiso, ObtenerClaimsUsuario) |
| `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` | Modelos GORM RBAC (roles, permisos, usuario_roles, usuario_tenant_roles) |
| `internal/rbac/infrastructure/persistence/postgres/rbac_repositorios.go` | Implementación PostgreSQL RBAC |
| `internal/sesiones/domain/token_servicio.go` | Interfaz TokenServicio con claims |
| `internal/sesiones/infrastructure/security/jwt/jwt_token_servicio.go` | JWT con soporte de Global y Tenants |

### Archivos faltantes (por crear)

| Archivo | Propósito |
|---------|-----------|
| `internal/presentation/handlers/tenant_handler.go` | Handlers HTTP para endpoints de tenant |
| `internal/presentation/facades/tenant_facade.go` | Interfaz facade de tenants |
| `internal/presentation/facades/tenant_facade_impl.go` | Implementación facade de tenants |
| `internal/presentation/middleware/tenant_middleware.go` | Middleware X-Tenant-ID |
| `internal/presentation/dto/tenant_dto.go` | DTOs de entrada/salida para tenant |
| `internal/tenants/application/services/gestionar_tenant/servicio_tenant_test.go` | Tests de aplicación |
| `internal/tenants/infrastructure/persistence/postgres/tenant_repositorio_test.go` | Tests de infraestructura |
| `cmd/seed/seed_tenant.go` | Seed de datos |

### Archivos por modificar

| Archivo | Cambio |
|---------|--------|
| `internal/config/database.go` | Agregar AutoMigrate para modelos de tenants y RBAC |
| `internal/config/env.go` | Agregar `TenantContextRequired` y carga desde `TENANT_CONTEXT_REQUIRED` |
| `internal/registry/registry.go` | Agregar todas las dependencias de tenants y RBAC |
| `internal/sesiones/application/services/login/servicio_login.go` | Inyectar AutorizacionServicio y cargar claims |
| `internal/presentation/router/router.go` | Registrar rutas de tenant, agregar middlewares |
| `main.go` | Conectar facade de tenants |
| `internal/tenants/infrastructure/persistence/postgres/tenant_model.go` | Agregar índices y constraints faltantes |
| `internal/rbac/infrastructure/persistence/postgres/rbac_models.go` | Agregar índices y constraints faltantes |
| `internal/tenants/application/services/gestionar_tenant/servicio_tenant.go` | Agregar método ListarTenantsDeUsuario (CU-10) |
| `.env.example` | Agregar `TENANT_CONTEXT_REQUIRED` |
