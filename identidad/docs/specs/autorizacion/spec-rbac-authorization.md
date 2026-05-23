---
title: Especificación del Módulo RBAC — Roles, Permisos y Autorización
version: 1.0
date_created: 2026-05-22
owner: Equipo Identidad
tags: rbac, roles, permisos, autorizacion, iam, jwt-claims
---

# Especificación del Módulo RBAC — Roles, Permisos y Autorización

> **Propósito**: Definir el sistema de control de acceso basado en roles (RBAC) con permisos atómicos, servicio de autorización y enriquecimiento de claims JWT. Opera sobre el modelo multi-tenant definido en `spec-tenant-management.md`.
>
> **No incluye**: Implementación de handlers HTTP, middleware de autorización, ni código Go de ningún tipo.
>
> **Formato**: Especificación de dominio que define QUÉ se construye, no CÓMO.

---

## 1. Propósito y Alcance

### 1.1 Propósito
Implementar un sistema de control de acceso basado en roles (RBAC) que permita asignar permisos atómicos a usuarios a través de roles. Este módulo extiende el sistema multi-tenant existente añadiendo autorización granular.

### 1.2 Incluye
- Modelo de datos de roles y permisos (3 tablas nuevas)
- Permisos atómicos como constantes de dominio (no modificables en runtime)
- 4 roles de sistema inmutables (sys_admin, administrador, agronomo, caficultor)
- Servicio de autorización con consciencia de tenant
- Enriquecimiento de claims JWT
- Casos de uso de gestión de usuarios (CRUD, contraseñas, asignación de roles)
- Estrategia de seed de permisos y roles

### 1.3 No incluye
- Modelo de tenant (ver `spec-tenant-management.md`)
- Implementación de handlers HTTP (presentación)
- Middleware de autorización (se define su estructura de datos, no su implementación)
- Auditoría de cambios de permisos/roles
- Roles personalizados (creados en runtime por administradores)
- Políticas ABAC o PBAC

---

## 2. Definiciones

| Término | Definición |
|---------|------------|
| **Permiso atómico** | Capacidad indivisible que representa una acción concreta. Formato `modulo:recurso:verbo`. Son constantes de dominio, no modificables en runtime. |
| **Rol de sistema** | Rol predefinido con `es_sistema = true`. Inmutable, no eliminable, no modificable. |
| **SYS_ADMIN** | Rol de sistema global. No atado a ningún tenant. Tiene todos los permisos. |
| **Administrador** | Rol de sistema que actúa como admin del tenant. Permisos completos DENTRO de su tenant. |
| **Seed** | Proceso que siembra permisos y roles de sistema en BD. Idempotente, se ejecuta en cada inicio. |
| **Claims JWT enriquecidos** | Los tokens JWT incluyen información de tenants, roles y permisos del usuario autenticado. |
| **Ámbito de rol** | Un rol puede tener ámbito `global` (aplica en todo el sistema, solo SYS_ADMIN) o `tenant` (aplica dentro de un tenant específico). |
| **`usuario_roles` vs `usuario_tenant_roles`** | `usuario_roles` asigna roles globales (solo para SYS_ADMIN). `usuario_tenant_roles` asigna roles dentro de un tenant específico para el resto de roles. |

---

## 3. Modelo de Datos

### 3.1 Tablas nuevas (módulo RBAC)

```sql
-- Tabla: permiso
CREATE TABLE permisos (
    id VARCHAR(36) PRIMARY KEY,
    codigo VARCHAR(100) NOT NULL UNIQUE,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT NOT NULL DEFAULT '',
    modulo VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- Tabla: rol
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL UNIQUE,
    descripcion TEXT NOT NULL DEFAULT '',
    es_sistema BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Tabla: rol_permiso (relación N:N)
CREATE TABLE rol_permisos (
    rol_id VARCHAR(36) NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permiso_id VARCHAR(36) NOT NULL REFERENCES permisos(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (rol_id, permiso_id)
);
```

### 3.2 Tablas existentes (definidas en spec-tenant-management.md)

Las siguientes tablas ya existen y son prerrequisito de este módulo:

| Tabla | Propósito |
|-------|-----------|
| `tenants` | Organizaciones/empresas (multi-tenant) |
| `usuario_tenants` | Relación N:N usuarios ↔ tenants |
| `usuario_tenant_roles` | Roles asignados a un usuario dentro de un tenant |
| `usuario_roles` | Roles globales asignados a un usuario (solo SYS_ADMIN) |

### 3.3 Índices

```sql
CREATE INDEX idx_rol_permisos_rol_id ON rol_permisos(rol_id);
CREATE INDEX idx_rol_permisos_permiso_id ON rol_permisos(permiso_id);
CREATE INDEX idx_permisos_modulo ON permisos(modulo);
CREATE INDEX idx_permisos_codigo ON permisos(codigo);
```

---

## 4. Permisos del Sistema

Son 8 permisos atómicos, clasificados por módulo. Son **CONSTANTES de dominio**, no modificables en tiempo de ejecución. Se propagan a la BD via seed.

### Módulo: identidad

| Código | Nombre | Descripción |
|--------|--------|-------------|
| `identidad:usuario:crear` | Crear Usuario | Crear nuevos usuarios con asignación opcional de rol |
| `identidad:usuario:modificar` | Modificar Usuario | Modificar datos personales de cualquier usuario |
| `identidad:usuario:eliminar` | Eliminar Usuario | Marcar un usuario como pendiente de eliminación |
| `identidad:usuario:consultar` | Consultar Usuario | Listar y ver detalles de cualquier usuario |
| `identidad:usuario:resetear_password` | Resetear Contraseña | Resetear la contraseña de otro usuario |
| `identidad:rol:asignar` | Asignar Rol | Asignar un rol a un usuario |
| `identidad:rol:revocar` | Revocar Rol | Revocar un rol de un usuario |
| `identidad:permiso:consultar` | Consultar Permisos | Listar permisos de un rol y roles de un usuario |

---

## 5. Roles del Sistema (4, inmutables)

| Rol | Ámbito | `es_sistema` | Descripción |
|-----|--------|:------------:|-------------|
| **sys_admin** | Global | ✅ | Super admin global. NO atado a un tenant. Todos los permisos. |
| **administrador** | Tenant | ✅ | Admin del tenant. Permisos completos dentro de su tenant. |
| **agronomo** | Tenant | ✅ | Permisos intermedios: crear y modificar usuarios, consultar. |
| **caficultor** | Tenant | ✅ | Solo consulta de usuarios. |

### Matriz de permisos por rol

| Permiso | sys_admin | administrador | agronomo | caficultor |
|---------|:---------:|:-------------:|:--------:|:----------:|
| `identidad:usuario:crear` | ✅ | ✅ | ✅ | ❌ |
| `identidad:usuario:modificar` | ✅ | ✅ | ✅ | ❌ |
| `identidad:usuario:eliminar` | ✅ | ✅ | ❌ | ❌ |
| `identidad:usuario:consultar` | ✅ | ✅ | ✅ | ✅ |
| `identidad:usuario:resetear_password` | ✅ | ✅ | ❌ | ❌ |
| `identidad:rol:asignar` | ✅ | ✅ | ❌ | ❌ |
| `identidad:rol:revocar` | ✅ | ✅ | ❌ | ❌ |
| `identidad:permiso:consultar` | ✅ | ✅ | ✅ | ❌ |

### Reglas de roles de sistema

```
Rol con es_sistema = true:
  ├── No se puede modificar nombre, descripción ni es_sistema
  ├── No se puede eliminar
  ├── No se puede modificar su asignación de permisos
  └── Se puede asignar/revocar a usuarios normalmente
```

---

## 6. Estrategia de Seed

El seed de permisos y roles se ejecuta automáticamente en cada inicio de la aplicación. Es **idempotente**:

```
1. Auto-migrate de tablas nuevas (permisos, roles, rol_permisos)

2. Para cada permiso en el catálogo:
   └── Buscar por código en BD
       ├── Si no existe → INSERT
       ├── Si existe pero nombre/descripción cambiaron → UPDATE
       └── Si existe y está igual → no-op

3. Para cada rol de sistema:
   └── Buscar por nombre en BD
       ├── Si no existe → INSERT
       └── Si existe → actualizar descripción si cambió

4. Para cada rol, re-sincronizar permisos:
   ├── Limpiar todos los rol_permisos de ese rol
   └── Insertar los permisos actuales según la matriz
```

---

## 7. Servicio de Autorización

### 7.1 Consulta fundamental
El servicio debe responder:
> **"¿Tiene el usuario X el permiso Y en el contexto del tenant Z?"**

### 7.2 Flujo de verificación

```
1. Recibir usuarioID, tenantID (opcional), codigoPermiso
2. Obtener roles del usuario:
   a. Buscar roles globales del usuario (usuario_roles)
   b. Buscar roles del usuario en el tenant (usuario_tenant_roles)
3. Si el usuario tiene rol "sys_admin" → PERMITIR (sin importar tenant)
4. Si tenantID está vacío y no es SYS_ADMIN → DENEGAR
5. Obtener permisos de todos los roles del usuario
6. Recorrer permisos:
   └── Si algún permiso.codigo == codigoPermiso → PERMITIR
7. Si no se encontró → DENEGAR
```

### 7.3 Carga de claims para JWT

```
1. Recibir usuarioID
2. Verificar si el usuario tiene rol global SYS_ADMIN
   ├── Si es SYS_ADMIN → retornar flag global=true, tenants vacío
   └── Si no → continuar
3. Cargar todos los tenants del usuario (usuario_tenants)
4. Para cada tenant:
   a. Obtener roles del usuario en ese tenant (usuario_tenant_roles)
   b. Obtener permisos de esos roles (roles → rol_permisos → permisos)
5. Retornar estructura con: usuarioID, flag global, lista de tenants con roles y permisos
```

### 7.4 Estrategia de caché (futuro)
En una primera versión, los permisos se consultan de BD en cada verificación. Para alto rendimiento se puede agregar Redis con invalidación por evento, pero queda fuera del alcance de esta spec.

### 7.5 Principio de seguridad
El servicio nunca debe revelar en mensajes de error si un permiso existe o no. El error debe ser genérico: "permiso denegado".

---

## 8. Claims JWT Enriquecidos

Los claims del access token JWT se enriquecen para evitar consultas a BD en cada request:

### Usuario de tenant (NO es SYS_ADMIN)

```json
{
  "sub": "uuid-del-usuario",
  "sid": "uuid-de-la-sesion",
  "iat": 1716300000,
  "exp": 1716300900,
  "typ": "access",
  "global": false,
  "tenants": {
    "tenant-uuid-1": {
      "slug": "nombre-del-tenant",
      "roles": ["administrador"],
      "permisos": ["identidad:usuario:crear", "identidad:usuario:modificar", "..."]
    },
    "tenant-uuid-2": {
      "slug": "otro-tenant",
      "roles": ["caficultor"],
      "permisos": ["identidad:usuario:consultar"]
    }
  }
}
```

### SYS_ADMIN (usuario global)

```json
{
  "sub": "uuid",
  "sid": "uuid",
  "global": true,
  "tenants": {}
}
```

### Reglas de claims

- Los claims se cargan al momento de generar el token (login/refresh)
- Si se asigna/revoca un rol, los cambios solo se reflejan en el siguiente token emitido, no en tiempo real
- El middleware de autenticación extrae estos claims y los pone a disposición de los handlers sin necesidad de consultar BD

---

## 9. Casos de Uso

Todos los casos de uso reciben un `tenantID` (contexto de tenant). Las excepciones son operaciones de SYS_ADMIN que pueden no requerirlo.

### 9.1 Listar Usuarios
- **Requiere permiso**: `identidad:usuario:consultar`
- **Filtros**: nombre, correo, estado, rol
- **Paginación**: página, tamaño de página
- **Ordenación**: por campos permitidos (nombre, correo, fecha_creacion)
- **Respuesta incluye**: datos del usuario + sus roles en el tenant actual
- **Ámbito**: solo usuarios del tenant especificado. SYS_ADMIN ve todos los tenants.

### 9.2 Crear Usuario
- **Requiere permiso**: `identidad:usuario:crear`
- **Datos**: nombre, apellido, correo, teléfono, password
- **Opcional**: asignar un rol al nuevo usuario en el mismo tenant
- **Validaciones**: correo único, password no vacío, nombre no vacío, formato de correo válido
- **Transaccional**: todo en una transacción (usuario + credenciales + rol si aplica + membresía al tenant)

### 9.3 Modificar Usuario
- **Requiere permiso**: `identidad:usuario:modificar`
- **Excepción**: un usuario puede modificar sus propios datos (nombre, apellido, teléfono) sin permiso especial (autogestión)
- **Datos modificables**: nombre, apellido, teléfono

### 9.4 Eliminar Usuario
- **Requiere permiso**: `identidad:usuario:eliminar`
- **Restricción**: un usuario no puede eliminarse a sí mismo
- **Mecanismo**: soft-delete (cambiar estado a PENDIENTE_DE_ELIMINACION)

### 9.5 Cambiar Contraseña
- **NO requiere permiso de autorización** (es autogestión)
- **Requiere**: contraseña actual + contraseña nueva
- **Validación**: la contraseña actual debe coincidir con el hash almacenado

### 9.6 Resetear Contraseña
- **Requiere permiso**: `identidad:usuario:resetear_password`
- **NO requiere**: contraseña actual del usuario objetivo
- **Además**: resetea intentos fallidos a 0
- **Propósito**: para cuando un usuario olvida su contraseña

### 9.7 Asignar Rol a Usuario
- **Requiere permiso**: `identidad:rol:asignar`
- **Validaciones**:
  - El rol debe existir
  - El usuario no debe tener ya ese rol (en el mismo tenant)
  - Si el usuario no es miembro del tenant, primero debe agregársele
  - Si el rol es SYS_ADMIN (global), no acepta tenantID
- **Registro**: en `usuario_tenant_roles` para roles regulares, en `usuario_roles` para SYS_ADMIN

### 9.8 Revocar Rol de Usuario
- **Requiere permiso**: `identidad:rol:revocar`
- **Validación**: el usuario debe tener el rol (en el mismo tenant)
- **Registro**: eliminar de `usuario_tenant_roles` o `usuario_roles`

### 9.9 Listar Roles de un Usuario
- **Requiere permiso**: `identidad:permiso:consultar`
- **Ámbito**: roles del usuario en el tenant especificado
- **Respuesta**: lista de roles con nombre, descripción y si es de sistema

### 9.10 Listar Permisos de un Rol
- **Requiere permiso**: `identidad:permiso:consultar`
- **Respuesta**: lista de permisos con código, nombre, descripción, módulo

---

## 10. Escenarios TDD

### 10.1 Servicio de Autorización

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Usuario tiene permiso en tenant | Usuario con rol admin en tenant T1 | `TienePermiso(uid, T1, crear)` | true |
| 2 | Usuario NO tiene permiso en tenant | Usuario con rol caficultor en T1 | `TienePermiso(uid, T1, eliminar)` | false |
| 3 | SYS_ADMIN tiene todos los permisos | Usuario con rol sys_admin | `TienePermiso(uid, cualquier_tenant, cualquier_permiso)` | true |
| 4 | Usuario sin tenant no tiene permisos | Usuario sin membresía | `TienePermiso(uid, T1, consultar)` | false |
| 5 | Mismo usuario, distinto tenant, distinto rol | Usuario admin en T1, caficultor en T2 | `TienePermiso(uid, T1, crear)` vs `TienePermiso(uid, T2, crear)` | true en T1, false en T2 |

### 10.2 Seed de Permisos y Roles

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 6 | Seed de permisos completo | — | Ejecutar seed | 8 permisos creados |
| 7 | Seed de roles completo | — | Ejecutar seed | 4 roles creados |
| 8 | Seed idempotente | Seed ya ejecutado | Ejecutar seed otra vez | Sin duplicados, sin errores |
| 9 | sys_admin tiene todos los permisos | Seed ejecutado | Consultar permisos de sys_admin | Los 8 permisos |
| 10 | caficultor tiene solo consulta | Seed ejecutado | Consultar permisos de caficultor | Solo 1 permiso |

### 10.3 Gestión de Usuarios

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 11 | Crear usuario sin rol | Permiso crear, datos válidos | Ejecutar en tenant T1 | Usuario creado en T1, sin roles |
| 12 | Crear usuario con rol | Permiso crear, rol especificado | Ejecutar | Usuario + rol asignado en T1 |
| 13 | Crear usuario sin permiso | Sin permiso crear | Ejecutar | Error permiso denegado |
| 14 | Correo duplicado | Correo ya existe | Validación previa | Error sin transacción |
| 15 | Listar usuarios del tenant | Admin en T1 | Listar con tenantID=T1 | Solo usuarios de T1 |
| 16 | Listar usuarios fuera del tenant | Admin en T1, lista usuarios de T2 | Listar con tenantID=T2 | Error (no pertenece a T2) |
| 17 | SYS_ADMIN lista todos | SYS_ADMIN | Listar sin tenantID o con cualquier tenant | Usuarios de todos los tenants |
| 18 | Modificar propio perfil | Usuario modifica sus datos | Ejecutar | Datos actualizados sin permiso |
| 19 | Admin modifica otro usuario | Admin con permiso modificar | Ejecutar | Datos actualizados |
| 20 | Eliminar usuario | Admin con permiso eliminar | Ejecutar | Estado → PENDIENTE_DE_ELIMINACION |
| 21 | Eliminarse a sí mismo | usuarioID == usuarioAEliminar | Validación | Error |

### 10.4 Contraseñas

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 22 | Cambio de password exitoso | Password actual correcto | Ejecutar | Nuevo hash persistido |
| 23 | Password actual incorrecto | Password actual no coincide | Ejecutar | Error |
| 24 | Reset de password exitoso | Admin con permiso resetear | Ejecutar | Nuevo hash, intentos fallidos → 0 |
| 25 | Reset sin permiso | Sin permiso resetear | Ejecutar | Error permiso denegado |

### 10.5 Asignación de Roles

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 26 | Asignar rol en tenant | Admin en T1, rol existe | Ejecutar | usuario_tenant_rol creado en T1 |
| 27 | Asignar rol duplicado | Usuario ya tiene ese rol en T1 | Ejecutar | Error |
| 28 | Asignar rol en tenant donde no es miembro | Usuario no pertenece a T1 | Ejecutar | Error (primero debe agregarlo al tenant) |
| 29 | Revocar rol en tenant | Admin en T1, usuario tiene rol | Ejecutar | usuario_tenant_rol eliminado |
| 30 | Revocar rol que no tiene | Usuario no tiene ese rol | Ejecutar | Error |
| 31 | Asignar SYS_ADMIN (rol global) | SYS_ADMIN ejecuta | Ejecutar | usuario_rol creado (sin tenant) |
| 32 | Asignar SYS_ADMIN por tenant (error) | Se intenta asignar sys_admin con tenantID | Ejecutar | Error: sys_admin es global |

---

## 11. Errores de Dominio

| Error | Significado |
|-------|-------------|
| Permiso denegado | El usuario no tiene el permiso requerido en el contexto especificado |
| Rol no encontrado | El rol solicitado no existe en el sistema |
| Rol ya asignado | El usuario ya tiene ese rol en el contexto especificado |
| Rol no asignado | El usuario no tiene ese rol (no se puede revocar) |
| Rol inmutable | El rol es de sistema y no puede modificarse ni eliminarse |
| Contraseña actual incorrecta | La contraseña actual no coincide con el hash |
| Usuario no pertenece al tenant | No se puede asignar un rol a un usuario que no es miembro del tenant |

---

## 12. Integración con JWT

El `TokenServicio` que genera los tokens JWT debe integrarse con el servicio de autorización para incluir los claims enriquecidos:

```
Flujo de login con claims:
  1. Usuario autentica (login)
  2. TokenServicio.GenerarAccessToken llama internamente a
     AutorizacionServicio.ObtenerClaimsUsuario
  3. Los claims (global, tenants, roles, permisos) se incluyen en el JWT
  4. El access token se retorna con toda la información de autorización embebida
  5. El middleware puede extraer los claims sin consultar BD
```

### Estructura de claims en el middleware

Los claims extraídos del JWT se ponen a disposición de los handlers mediante una estructura que contiene:

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `UsuarioID` | string | ID del usuario autenticado |
| `SesionID` | string | ID de la sesión |
| `Global` | bool | True si es SYS_ADMIN |
| `Tenants` | map[string]TenantClaims | Mapa de tenantID → slug, roles, permisos |

Cada `TenantClaims` contiene:

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `Slug` | string | Slug del tenant |
| `Roles` | []string | Nombres de los roles del usuario en ese tenant |
| `Permisos` | []string | Códigos de los permisos del usuario en ese tenant |

---

## 13. Integración con Registry

El Registry existente (inyección de dependencias) debe extenderse con:

```
Nuevas dependencias a registrar:
  - RolRepositorio (persistencia)
  - PermisoRepositorio (persistencia)
  - UsuarioTenantRolRepositorio (persistencia)
  - TenantRepositorio (persistencia, ya existe desde spec-tenant)
  - AutorizacionServicio (servicio de dominio)
  - ServicioGestionUsuarios (aplicación)
  - ServicioCambiarPassword (aplicación)
  - ServicioResetearPassword (aplicación)
  - ServicioAsignacionRoles (aplicación)

Al iniciar:
  1. Auto-migrate de tablas nuevas (permisos, roles, rol_permisos)
  2. Ejecutar seed de permisos y roles
```

---

## 14. Dependencias

- **`spec-tenant-management.md`** — Define el modelo de tenant del que depende este módulo (tenants, usuario_tenants, usuario_tenant_roles, usuario_roles)
- Módulo de sesiones existente — Para integración JWT (login, refresh, claims)
- Módulo de usuarios existente — Para reutilizar entidad Usuario y Credenciales

---

## 15. Consideraciones de Seguridad

| # | Riesgo | Mitigación |
|---|--------|------------|
| 1 | Escalada de privilegios: usuario con rol "agronomo" intenta asignar roles | El servicio verifica `identidad:rol:asignar` que solo tienen sys_admin y administrador |
| 2 | Cruce de tenants: admin de T1 accede a datos de T2 | Toda operación verifica que el usuario pertenece al tenant. SYS_ADMIN es la única excepción. |
| 3 | Token stale (claims desactualizados) | Aceptable. Los claims reflejan el estado al momento de emisión. Al refrescar se obtienen claims actualizados. |
| 4 | Asignación directa de permisos a usuario sin rol | El modelo no tiene tabla usuario_permiso. Los permisos solo se asignan via roles. |
| 5 | Error informativo revela existencia de permisos | Errores genéricos: "permiso denegado", sin especificar qué permiso falta. |
| 6 | Roles de sistema no editables | Validar `es_sistema` antes de cualquier operación de modificación/eliminación de roles. |
| 7 | Inyección SQL en filtros de listado | Usar columnas permitidas para validar campos de ordenación. |
| 8 | Transaccionalidad en creación de usuario | Usar UnitOfWork para garantizar atomicidad en operaciones multi-tabla. |

---

## 16. Checklist de Validación

- [ ] ¿Existen las tablas `permisos`, `roles`, `rol_permisos`?
- [ ] ¿Los 8 permisos están definidos como constantes de dominio?
- [ ] ¿Los 4 roles de sistema existen y tienen `es_sistema = true`?
- [ ] ¿sys_admin tiene todos los permisos sin restricción de tenant?
- [ ] ¿La matriz de permisos por rol es correcta?
- [ ] ¿El seed es idempotente?
- [ ] ¿El servicio de autorización acepta tenantID?
- [ ] ¿SYS_ADMIN siempre pasa la verificación sin importar tenant?
- [ ] ¿Los claims JWT incluyen `global` y `tenants`?
- [ ] ¿Los errores de permiso son genéricos?
- [ ] ¿Un usuario no puede asignarse permisos a sí mismo?
- [ ] ¿No hay tabla `usuario_permiso` (solo se asigna via roles)?
- [ ] ¿Los casos de uso verifican membresía al tenant antes de operar?
- [ ] ¿Las transacciones son atómicas via UnitOfWork?
- [ ] ¿El rol de sistema SYS_ADMIN no puede asignarse con tenantID?
- [ ] ¿La re-asignación de permisos en seed es completa (limpia y re-inserta)?

---

## 17. Especificaciones Relacionadas

- `spec-tenant-management.md` — Modelo multi-tenant (prerrequisito de este módulo)
- `../../sesiones/login_spec.md` — Login, refresh, logout, JWT (integración de claims)
- `../../registro/spec_registro.md` — Registro con verificación de correo
- `../../adr/architecture-context.md` — Contexto arquitectónico y flujo de capas
