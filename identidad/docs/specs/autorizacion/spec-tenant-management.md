---
title: Especificación del Módulo Multi-Tenant — Gestión de Inquilinos
version: 1.0
date_created: 2026-05-22
owner: Equipo Identidad
tags: multitenant, tenant, iam, organizaciones
---

# Especificación del Módulo Multi-Tenant — Gestión de Inquilinos

> **Propósito**: Definir el modelo multi-tenant del sistema de identidad. Un tenant (inquilino) representa una organización/empresa cliente que puede tener N fincas (gestionadas por otro módulo).
>
> **No incluye**: Definición de roles y permisos (va en spec de RBAC), implementación de handlers HTTP, middleware de autorización.

---

## 1. Propósito y Alcance

### 1.1 Propósito
Establecer el modelo de datos, reglas de negocio y casos de uso para la gestión de inquilinos (tenants) en el sistema de identidad. Este módulo permite que múltiples organizaciones compartan la misma instancia del sistema, manteniendo el aislamiento de datos y la asignación granular de roles por tenant.

### 1.2 Incluye
- Entidad Tenant y su ciclo de vida (crear, activar, desactivar, consultar)
- Membresía: relación usuario ↔ tenant (`usuario_tenant`)
- Roles en contexto de tenant (`usuario_tenant_rol`)
- Roles globales (`usuario_rol`) — solo para SYS_ADMIN
- Contexto de tenant en operaciones (header `X-Tenant-ID`)
- Integración con el módulo de sesiones existente

### 1.3 No incluye
- Definición de roles y permisos (cubierto en la Especificación del Módulo de Autorización IAM — RBAC)
- Implementación de handlers HTTP (cubierto en la Especificación de la Capa de Presentación)
- Middleware de autorización (cubierto en la Especificación del Módulo de Autorización IAM — RBAC)
- Políticas de acceso basadas en recursos (PBAC) o atributos (ABAC)

---

## 2. Definiciones

| Término | Definición |
|---------|------------|
| **Tenant** | Organización/empresa cliente del sistema. Cada tenant agrupa usuarios, roles y recursos (fincas). |
| **Membresía** | Relación que indica que un usuario pertenece a un tenant. No implica ningún rol. Se registra en `usuario_tenants`. |
| **Rol en contexto de tenant** | Asignación de un rol a un usuario DENTRO de un tenant específico. Se registra en `usuario_tenant_roles`. |
| **Rol global** | Asignación de un rol a un usuario SIN contexto de tenant. Solo aplica para SYS_ADMIN. Se registra en `usuario_rol`. |
| **SYS_ADMIN** | Super administrador global. No está atado a ningún tenant. Tiene todos los permisos del sistema. |
| **X-Tenant-ID** | Header HTTP que indica en qué contexto de tenant se ejecuta una operación. |
| **Slug** | Identificador único legible para un tenant (ej: `finca-la-esperanza`). Se usa como alternativa al UUID en operaciones de lectura. |

---

## 3. Modelo de Datos

### 3.1 Esquema PostgreSQL

```sql
-- Tabla: tenants
CREATE TABLE tenants (
    id VARCHAR(36) PRIMARY KEY,
    nombre VARCHAR(200) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    activo BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Tabla: usuario_tenant (membresía)
CREATE TABLE usuario_tenants (
    usuario_id VARCHAR(36) NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    tenant_id VARCHAR(36) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (usuario_id, tenant_id)
);

-- Tabla: usuario_tenant_rol (roles dentro de un tenant)
CREATE TABLE usuario_tenant_roles (
    usuario_id VARCHAR(36) NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    tenant_id VARCHAR(36) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    rol_id VARCHAR(36) NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (usuario_id, tenant_id, rol_id)
);

-- Tabla: usuario_rol (roles globales — SOLO para SYS_ADMIN)
CREATE TABLE usuario_roles (
    usuario_id VARCHAR(36) NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    rol_id VARCHAR(36) NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (usuario_id, rol_id)
);

-- Índices
CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_activo ON tenants(activo);
CREATE INDEX idx_usuario_tenants_usuario_id ON usuario_tenants(usuario_id);
CREATE INDEX idx_usuario_tenants_tenant_id ON usuario_tenants(tenant_id);
CREATE INDEX idx_usuario_tenant_roles_usuario_id ON usuario_tenant_roles(usuario_id);
CREATE INDEX idx_usuario_tenant_roles_tenant_id ON usuario_tenant_roles(tenant_id);
CREATE INDEX idx_usuario_tenant_roles_rol_id ON usuario_tenant_roles(rol_id);
```

### 3.2 Restricciones del modelo

| # | Restricción | Descripción |
|---|-------------|-------------|
| 1 | PK compuesta en `usuario_tenants` | `(usuario_id, tenant_id)` — un usuario no puede tener dos membresías al mismo tenant |
| 2 | PK compuesta en `usuario_tenant_roles` | `(usuario_id, tenant_id, rol_id)` — un usuario no puede tener el mismo rol dos veces en el mismo tenant |
| 3 | PK compuesta en `usuario_roles` | `(usuario_id, rol_id)` — un usuario no puede tener el mismo rol global dos veces |
| 4 | `slug` único | No pueden existir dos tenants con el mismo slug |
| 5 | `ON DELETE CASCADE` | Al eliminar un tenant se eliminan todas sus membresías y asignaciones de roles |
| 6 | `activo` por defecto `true` | Todo tenant nuevo se crea activo |

---

## 4. Reglas de Negocio del Modelo Multi-Tenant

| # | Regla | Descripción |
|---|-------|-------------|
| RN-01 | Membresía múltiple | Un usuario puede pertenecer a múltiples tenants. La membresía se registra en `usuario_tenants` y no implica ningún rol. |
| RN-02 | Roles diferenciados por tenant | Un usuario puede tener diferentes roles en diferentes tenants. Ej: "administrador" en tenant A y "caficultor" en tenant B. |
| RN-03 | SYS_ADMIN es global | El rol SYS_ADMIN es un rol global. No se asigna dentro de un tenant. Se asigna vía `usuario_rol` (sin tenant_id). |
| RN-04 | Roles regulares atados a tenant | Los roles regulares (administrador, agronomo, caficultor) SIEMPRE se asignan en contexto de un tenant. Se registran en `usuario_tenant_roles`. |
| RN-05 | Sin duplicación de roles | Un usuario no puede tener el mismo rol dos veces en el mismo tenant. La PK compuesta de `usuario_tenant_roles` lo impide a nivel BD. |
| RN-06 | Desactivación conservadora | Desactivar un tenant debe impedir nuevas operaciones en ese contexto, pero los datos se conservan. |
| RN-07 | Slug inmutable | El slug de un tenant no se puede modificar después de la creación. Es un identificador único y permanente. |
| RN-08 | Eliminación en cascada | Al eliminar un usuario de un tenant (remover membresía), también se eliminan todos sus roles en ese tenant (ON DELETE CASCADE en ambas tablas). |
| RN-09 | SYS_ADMIN omite membresía | SYS_ADMIN puede operar en cualquier tenant sin necesidad de tener membresía explícita en `usuario_tenants`. |

---

## 5. Flujo del Contexto de Tenant

```
Petición HTTP:
  1. Cliente envía header X-Tenant-ID: <uuid>
  2. Middleware extrae el tenant_id del header
  3. Middleware valida que:
     a. El tenant existe en BD
     b. El tenant está activo
     c. El usuario autenticado pertenece al tenant (o es SYS_ADMIN)
  4. Si todo OK → tenant_id queda disponible en el contexto de la petición
  5. Si no:
     - Tenant no existe → error 400 (tenant inválido)
     - Tenant inactivo → error 403 (tenant desactivado)
     - Usuario no miembro → error 403 (acceso denegado)
     - Header ausente y TENANT_CONTEXT_REQUIRED=true → error 400 (tenant requerido)
```

### 5.1 Propagación del contexto

```
Capas:
  HTTP (Gin)          Capa de Presentación    Capa de Aplicación    Capa de Dominio
  ┌──────────┐         ┌──────────────┐        ┌─────────────┐       ┌───────────┐
  │X-Tenant-ID│ ─────→ │ Extraer de   │ ─────→ │ Recibir en  │ ───→ │ Usar para │
  │  header   │         │ contexto Gin │        │ comandos    │       │ autorizar  │
  └──────────┘         └──────────────┘        └─────────────┘       └───────────┘
```

- La capa de presentación extrae el `X-Tenant-ID` del header y lo inyecta en los comandos de aplicación
- La capa de aplicación recibe el `tenantID` como parte del comando y lo pasa al servicio de dominio cuando sea necesario
- La capa de dominio usa el `tenantID` para autorización y consultas

---

## 6. Casos de Uso

| # | Caso de Uso | Descripción | Actor Principal |
|---|-------------|-------------|-----------------|
| CU-01 | **Crear tenant** | Crear una nueva organización con nombre y slug únicos | SYS_ADMIN |
| CU-02 | **Listar todos los tenants** | Listar todos los tenants del sistema | SYS_ADMIN |
| CU-03 | **Listar mis tenants** | Listar los tenants a los que pertenece el usuario autenticado | Usuario autenticado |
| CU-04 | **Obtener tenant por ID** | Ver detalle de un tenant específico | SYS_ADMIN + miembros del tenant |
| CU-05 | **Obtener tenant por slug** | Ver detalle de un tenant usando su slug | SYS_ADMIN + miembros del tenant |
| CU-06 | **Activar tenant** | Cambiar estado de inactivo a activo | SYS_ADMIN |
| CU-07 | **Desactivar tenant** | Cambiar estado de activo a inactivo | SYS_ADMIN |
| CU-08 | **Agregar usuario a tenant** | Establecer membresía de un usuario en un tenant | SYS_ADMIN + administrador del tenant |
| CU-09 | **Remover usuario de tenant** | Eliminar membresía de un usuario en un tenant (en cascada: también elimina roles en ese tenant) | SYS_ADMIN + administrador del tenant |
| CU-10 | **Listar tenants de un usuario** | Obtener todos los tenants a los que un usuario específico pertenece | SYS_ADMIN + administrador del tenant |
| CU-11 | **Listar usuarios de un tenant** | Obtener todos los usuarios miembros de un tenant | SYS_ADMIN + miembros del tenant |

### 6.1 Detalle de casos de uso

#### CU-01: Crear tenant
| Campo | Valor |
|-------|-------|
| **Descripción** | Crea un nuevo tenant con nombre y slug únicos. El tenant se crea siempre activo. |
| **Actor** | SYS_ADMIN |
| **Entrada** | `nombre` (string), `slug` (string) |
| **Validaciones** | `nombre` obligatorio (máx 200 chars), `slug` obligatorio (máx 100 chars, alfanumérico con guiones), slug único en BD |
| **Salida exitosa** | Tenant creado con `activo = true` |
| **Errores** | Slug duplicado, nombre vacío, slug inválido, sin permiso SYS_ADMIN |

#### CU-06: Activar tenant y CU-07: Desactivar tenant
| Campo | Valor |
|-------|-------|
| **Descripción** | Cambia el estado `activo` de un tenant. Al desactivar, el tenant no acepta nuevas operaciones. |
| **Actor** | SYS_ADMIN |
| **Entrada** | `tenant_id` (string) |
| **Validaciones** | Tenant debe existir, no puede activar si ya está activo, no puede desactivar si ya está inactivo |
| **Salida exitosa** | Estado actualizado |
| **Errores** | Tenant no encontrado, sin permiso SYS_ADMIN, estado ya es el solicitado |

#### CU-08: Agregar usuario a tenant
| Campo | Valor |
|-------|-------|
| **Descripción** | Establece membresía de un usuario en un tenant. No asigna ningún rol, solo membresía. |
| **Actor** | SYS_ADMIN + administrador del tenant |
| **Entrada** | `usuario_id` (string), `tenant_id` (string) |
| **Validaciones** | Usuario debe existir, tenant debe existir y estar activo, usuario no debe ser ya miembro |
| **Salida exitosa** | Membresía creada en `usuario_tenants` |
| **Errores** | Usuario no encontrado, tenant no encontrado, usuario ya es miembro, sin permiso |

#### CU-09: Remover usuario de tenant
| Campo | Valor |
|-------|-------|
| **Descripción** | Elimina la membresía de un usuario en un tenant. Por ON DELETE CASCADE, también elimina todos los roles del usuario en ese tenant. |
| **Actor** | SYS_ADMIN + administrador del tenant |
| **Entrada** | `usuario_id` (string), `tenant_id` (string) |
| **Validaciones** | Usuario debe ser miembro del tenant |
| **Salida exitosa** | Membresía y roles eliminados |
| **Errores** | Usuario no es miembro del tenant, sin permiso |

---

## 7. Reglas de Autorización

| # | Acción | ¿Quién puede ejecutarla? |
|---|--------|--------------------------|
| 1 | Crear tenant | Solo SYS_ADMIN |
| 2 | Listar todos los tenants | Solo SYS_ADMIN |
| 3 | Listar mis tenants | Cualquier usuario autenticado (solo los suyos) |
| 4 | Obtener detalle de tenant | SYS_ADMIN + miembros del tenant |
| 5 | Activar tenant | Solo SYS_ADMIN |
| 6 | Desactivar tenant | Solo SYS_ADMIN |
| 7 | Agregar usuario a tenant | SYS_ADMIN + administrador del tenant |
| 8 | Remover usuario de tenant | SYS_ADMIN + administrador del tenant |
| 9 | Listar tenants de un usuario | SYS_ADMIN + administrador del tenant objetivo |
| 10 | Listar usuarios de un tenant | SYS_ADMIN + miembros del tenant |

### 7.1 Matriz de decisión

```
                        ┌──────────┬──────────────┬──────────────┬──────────────┐
                        │SYS_ADMIN │ administrador│  agronomo    │ caficultor   │
                        │ (global) │ (en tenant)  │ (en tenant)  │ (en tenant)  │
┌───────────────────────┼──────────┼──────────────┼──────────────┼──────────────┤
│ Crear tenant          │    ✅    │      ❌      │      ❌      │      ❌      │
│ Listar todos          │    ✅    │      ❌      │      ❌      │      ❌      │
│ Listar mis tenants    │    ✅    │      ✅      │      ✅      │      ✅      │
│ Obtener detalle       │    ✅    │      ✅      │      ✅      │      ✅      │
│ Activar/Desactivar    │    ✅    │      ❌      │      ❌      │      ❌      │
│ Agregar usuario       │    ✅    │      ✅      │      ❌      │      ❌      │
│ Remover usuario       │    ✅    │      ✅      │      ❌      │      ❌      │
│ Listar usuarios       │    ✅    │      ✅      │      ✅      │      ❌      │
└───────────────────────┴──────────┴──────────────┴──────────────┴──────────────┘
```

---

## 8. Escenarios TDD

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Crear tenant exitoso | SYS_ADMIN autenticado, datos válidos (nombre, slug único) | Ejecutar crear tenant | Tenant creado con slug único, activo=true, se retorna el tenant |
| 2 | Crear tenant sin permiso | Usuario NO es SYS_ADMIN | Ejecutar crear tenant | Error de permiso |
| 3 | Slug duplicado | SYS_ADMIN, slug ya existe en BD | Ejecutar crear tenant con el mismo slug | Error de validación (slug duplicado) |
| 4 | Listar todos los tenants | SYS_ADMIN autenticado, existen 3 tenants | Ejecutar listar todos | Retorna los 3 tenants |
| 5 | Listar todos sin permiso | Usuario NO es SYS_ADMIN | Ejecutar listar todos | Error de permiso |
| 6 | Listar mis tenants | Usuario pertenece a 2 tenants | Ejecutar listar mis tenants | Retorna solo sus 2 tenants |
| 7 | Listar mis tenants (sin membresías) | Usuario no pertenece a ningún tenant | Ejecutar listar mis tenants | Retorna lista vacía |
| 8 | Obtener tenant por ID | SYS_ADMIN o miembro del tenant | Ejecutar obtener por ID | Retorna detalle del tenant |
| 9 | Obtener tenant por ID (no miembro) | Usuario NO es miembro del tenant ni SYS_ADMIN | Ejecutar obtener por ID | Error de permiso |
| 10 | Activar tenant | SYS_ADMIN, tenant inactivo | Ejecutar activar tenant | tenant.activo = true |
| 11 | Desactivar tenant | SYS_ADMIN, tenant activo | Ejecutar desactivar tenant | tenant.activo = false |
| 12 | Desactivar tenant sin permiso | Usuario NO es SYS_ADMIN | Ejecutar desactivar tenant | Error de permiso |
| 13 | Agregar usuario a tenant | Administrador del tenant, usuario no es miembro | Ejecutar agregar usuario | Membresía creada en usuario_tenants |
| 14 | Agregar usuario ya miembro | Usuario ya pertenece al tenant | Ejecutar agregar de nuevo | Error (membresía duplicada) |
| 15 | Agregar usuario a tenant inactivo | Administrador, tenant inactivo | Ejecutar agregar usuario | Error (tenant inactivo) |
| 16 | Remover usuario de tenant | Administrador del tenant, usuario es miembro | Ejecutar remover usuario | Membresía + roles eliminados |
| 17 | Remover usuario que no es miembro | Administrador del tenant, usuario NO es miembro | Ejecutar remover usuario | Error (usuario no es miembro) |
| 18 | Remover último administrador | Tenant tiene solo 1 administrador | Ejecutar remover ese administrador | Error (no se puede remover al último administrador) |
| 19 | Listar usuarios de un tenant | Miembro del tenant, tenant tiene 3 usuarios | Ejecutar listar usuarios | Retorna los 3 usuarios del tenant |
| 20 | Listar usuarios sin ser miembro | Usuario NO es miembro del tenant | Ejecutar listar usuarios | Error de permiso |

### 8.1 Escenarios de integración

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 21 | Login con tenants | Usuario pertenece a 2 tenants, hace login | Login exitoso | JWT incluye tenants con roles y permisos |
| 22 | Operación con X-Tenant-ID válido | Usuario miembro del tenant, header presente | Operación en tenant | Éxito |
| 23 | Operación con X-Tenant-ID de otro tenant | Usuario NO miembro del tenant indicado | Operación en tenant | Error 403 |
| 24 | Operación con tenant desactivado | Tenant desactivado, usuario era miembro | Operación en tenant | Error 403 (tenant inactivo) |
| 25 | SYS_ADMIN opera sin membresía | SYS_ADMIN, no es miembro del tenant | Operación en tenant | Éxito (SYS_ADMIN omite membresía) |

---

## 9. Integración con el Módulo de Sesiones

### 9.1 Flujo actual de login
El login actual genera un JWT con `sub` (usuarioID) y `sid` (sesionID).

### 9.2 Cambios requeridos para multi-tenancy

| Aspecto | Cambio |
|---------|--------|
| Carga de datos | Al hacer login, se deben cargar los tenants del usuario y sus roles/permisos por tenant |
| Claims JWT | El JWT debe incluir `global` (bool) y `tenants` (map[tenantID] → slug, roles, permisos) |
| Flujo de login | No cambia: autenticación primero (login), autorización después (RBAC) |
| Refresh token | Al refrescar el token, se deben recargar los claims de autorización |

### 9.3 Diagrama de secuencia (login multi-tenant)

```
Usuario              Servicio Login         Auth Servicio        Tenant Repo
  │                       │                       │                  │
  │──login(creds)────────→│                       │                  │
  │                       │──autenticar()─────────→│                  │
  │                       │←──token(usuarioID)────│                  │
  │                       │                       │                  │
  │                       │──obtenerClaims(uid)──→│                  │
  │                       │                       │──tenants(uid)───→│
  │                       │                       │←──lista tenants─│
  │                       │                       │──roles_permisos─→│
  │                       │                       │←──por tenant────│
  │                       │←──UsuarioConRoles────│                  │
  │                       │                       │                  │
  │                       │──generarJWT(claims)──→│                  │
  │←──token con tenants──│                       │                  │
```

---

## 10. Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `TENANT_CONTEXT_REQUIRED` | `true` | Si es `true`, todas las operaciones requieren el header `X-Tenant-ID` (excepto rutas públicas y operaciones de SYS_ADMIN que no requieren tenant). |

---

## 11. Checklist de Validación

- [ ] ¿La tabla `tenants` existe con los campos especificados (id, nombre, slug, activo, created_at, updated_at)?
- [ ] ¿La tabla `usuario_tenants` tiene PK compuesta (usuario_id, tenant_id)?
- [ ] ¿La tabla `usuario_tenant_roles` tiene PK compuesta (usuario_id, tenant_id, rol_id)?
- [ ] ¿La tabla `usuario_roles` existe para roles globales y solo se usa para SYS_ADMIN?
- [ ] ¿Un usuario puede pertenecer a múltiples tenants?
- [ ] ¿Un usuario puede tener diferentes roles en diferentes tenants?
- [ ] ¿SYS_ADMIN es un rol global sin necesidad de membresía en tenant?
- [ ] ¿Desactivar un tenant impide nuevas operaciones en ese contexto?
- [ ] ¿El header `X-Tenant-ID` se propaga correctamente desde HTTP hasta la capa de dominio?
- [ ] ¿Hay validación de que el usuario pertenece al tenant antes de operar (o es SYS_ADMIN)?
- [ ] ¿Los índices están creados para las consultas principales (slug, usuario_id, tenant_id)?
- [ ] ¿La eliminación de un usuario de un tenant también elimina sus roles en cascada?
- [ ] ¿El JWT incluye los claims `global` y `tenants` después del login?
- [ ] ¿La variable `TENANT_CONTEXT_REQUIRED` se respeta en el middleware?

---

## 12. Referencias

- [Especificación del Módulo de Autorización IAM — RBAC](./spec-iam-rbac.md)
- [Especificación de Sesiones](../sesiones/login_spec.md)
- [Especificación de la Capa de Presentación](../presentacion/spec-presentation-layer.md)
