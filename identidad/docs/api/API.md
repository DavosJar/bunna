# API del Servicio Identidad

API REST para autenticación, administración de usuarios, control de acceso basado en roles (RBAC),
gestión de sesiones, verificación de correo electrónico y recuperación de contraseñas.

---

## Índice

- [Base URL](#base-url)
- [Autenticación](#autenticación)
- [Formato de respuestas](#formato-de-respuestas)
  - [Éxito](#éxito)
  - [Error](#error)
  - [Códigos de error comunes](#códigos-de-error-comunes)
- [Endpoints](#endpoints)
  - [Sistema](#sistema)
  - [Autenticación](#autenticación-1)
  - [Mi Perfil (autogestión)](#mi-perfil-autogestión)
  - [Usuarios (admin)](#usuarios-admin)
  - [Seguridad (admin)](#seguridad-admin)
  - [Sesiones (admin)](#sesiones-admin)
  - [Roles y Permisos (RBAC, admin)](#roles-y-permisos-rbac-admin)
  - [Tenants (admin)](#tenants-admin)
  - [Verificación de Correo](#verificación-de-correo)
  - [Recuperación de Contraseña](#recuperación-de-contraseña)

---

## Base URL

```
http://localhost:8080
```

---

## Autenticación

La API utiliza tokens JWT (JSON Web Tokens) para autenticación. El flujo es:

1. **Obtener token**: enviar `POST /api/v1/auth/login` con credenciales.
2. **Usar token**: incluir el `access_token` en el header `Authorization` de las
   requests protegidas.
3. **Renovar token**: cuando el `access_token` expire, usar `POST /api/v1/auth/refresh`
   con el `refresh_token`.
4. **Cerrar sesión**: enviar `POST /api/v1/auth/logout` para revocar la sesión actual.

```
Authorization: Bearer <access_token>
```

Los endpoints marcados con **🔒** requieren un token JWT válido. Sin token,
responden con `401 Unauthorized`.

Los siguientes endpoints son **públicos** (no requieren autenticación):

| Endpoint | Método |
|---|---|
| `/health` | GET |
| `/api/v1/auth/register` | POST |
| `/api/v1/auth/login` | POST |
| `/api/v1/auth/refresh` | POST |
| `/api/v1/verificacion/confirmar` | POST |
| `/api/v1/recuperacion/solicitar` | POST |
| `/api/v1/recuperacion/validar` | POST |
| `/api/v1/recuperacion/confirmar` | POST |

---

## Formato de respuestas

### Éxito

Todas las respuestas exitosas (excepto `/health`) usan un envoltorio uniforme:

```json
{
  "data": { ... },
  "_links": {
    "self": { "href": "/api/v1/usuarios/{id}", "method": "GET" }
  }
}
```

| Campo | Tipo | Descripción |
|---|---|---|
| `data` | `object` | Payload de la respuesta (varía por endpoint) |
| `_links` | `object` | (opcional) Enlaces HATEOAS a recursos relacionados |

**Excepción**: el endpoint `GET /health` retorna directamente `{ "status": "ok" }`
sin envoltorio.

### Error

Los errores siguen el formato [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457)
(Problem Details):

```json
{
  "title": "Bad Request",
  "status": 400,
  "detail": "descripción del error específico",
  "errors": [
    {
      "message": "detalle adicional del campo",
      "field": "correo"
    }
  ]
}
```

| Campo | Tipo | Descripción |
|---|---|---|
| `title` | `string` | Título del error |
| `status` | `number` | Código de estado HTTP |
| `detail` | `string` | Descripción específica del error |
| `errors` | `array` | (opcional) Errores detallados por campo |

### Códigos de error comunes

| Código | Significado | Causas comunes |
|---|---|---|
| `400 Bad Request` | Datos inválidos en la solicitud | Campos requeridos faltantes, formato incorrecto |
| `401 Unauthorized` | No autenticado | Token ausente, inválido o expirado; credenciales incorrectas |
| `403 Forbidden` | Sin permisos | El usuario no tiene el rol/permiso necesario |
| `404 Not Found` | Recurso no encontrado | ID de usuario, rol, sesión o tenant inválido |
| `409 Conflict` | Conflicto de estado | Rol ya asignado, correo ya verificado, slug duplicado |
| `422 Unprocessable Entity` | Regla de negocio violada | Cuenta bloqueada, IP bloqueada, rate limit excedido, enlace expirado |
| `500 Internal Server Error` | Error del servidor | Contactar al equipo de plataforma |

---

## Endpoints

### Sistema

#### `GET /health` — Health check

Verifica que el servicio está activo y respondiendo. Sin autenticación.

**Response** `200 OK`

```json
{
  "status": "ok"
}
```

---

### Autenticación

#### `POST /api/v1/auth/register` — Registro de usuario

Crea un nuevo usuario con sus credenciales. Sin autenticación.

**Request body**

```json
{
  "nombre": "Juan",
  "apellido": "Pérez",
  "correo": "juan@correo.com",
  "password": "secreto123",
  "telefono": "0999999999"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `nombre` | `string` | Sí | Nombre del usuario |
| `apellido` | `string` | Sí | Apellido del usuario |
| `correo` | `string` | Sí | Correo electrónico válido |
| `password` | `string` | Sí | Contraseña (mínimo 8 caracteres) |
| `telefono` | `string` | No | Teléfono de contacto |

**Response** `201 Created`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "correo": "juan@correo.com",
    "estado": "NO_VERIFICADO"
  },
  "_links": {
    "self": {
      "href": "/api/v1/usuarios/01926b1e-dead-beef-cafe-000000000001",
      "method": "GET"
    }
  }
}
```

---

#### `POST /api/v1/auth/login` — Inicio de sesión

Autentica al usuario y devuelve tokens JWT de acceso y refresco. Sin autenticación.

**Request body**

```json
{
  "correo": "juan@correo.com",
  "password": "secreto123"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `correo` | `string` | Sí | Correo electrónico del usuario |
| `password` | `string` | Sí | Contraseña del usuario |

**Response** `200 OK`

```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "token_type": "Bearer",
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001"
  },
  "_links": {
    "self": {
      "href": "/api/v1/usuarios/01926b1e-dead-beef-cafe-000000000001",
      "method": "GET"
    },
    "refresh": {
      "href": "/api/v1/auth/refresh",
      "method": "POST"
    }
  }
}
```

---

#### `POST /api/v1/auth/refresh` — Renovar sesión

Renueva el `access_token` usando el `refresh_token`. Aplica rotación de tokens
(el refresh token anterior se invalida). Sin autenticación.

**Request body**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response** `200 OK`

```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "token_type": "Bearer",
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001"
  }
}
```

---

#### `POST /api/v1/auth/logout` — Cerrar sesión 🔒

Cierra la sesión actual del usuario autenticado.

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "sesiones_revocadas": 1
  }
}
```

---

#### `POST /api/v1/auth/logout/all` — Cerrar todas las sesiones 🔒

Cierra todas las sesiones activas del usuario autenticado.

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "sesiones_revocadas": 3
  }
}
```

---

### Mi Perfil (autogestión)

#### `GET /api/v1/mi-perfil` — Ver mi perfil 🔒

Obtiene los datos del perfil del usuario autenticado.

**Response** `200 OK`

```json
{
  "data": {
    "id": "01926b1e-dead-beef-cafe-000000000001",
    "correo": "juan@correo.com",
    "nombre": "Juan",
    "apellido": "Pérez",
    "telefono": "0999999999",
    "estado": "ACTIVO",
    "creado_en": "2026-05-23T12:00:00Z"
  }
}
```

---

#### `PUT /api/v1/mi-perfil` — Modificar mi perfil 🔒

Actualiza los datos del perfil del usuario autenticado.

**Request body**

```json
{
  "nombre": "Juan Carlos",
  "apellido": "Pérez García"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `nombre` | `string` | Sí | Nuevo nombre |
| `apellido` | `string` | Sí | Nuevo apellido |

**Response** `200 OK`

```json
{
  "data": {
    "id": "01926b1e-dead-beef-cafe-000000000001",
    "correo": "juan@correo.com",
    "nombre": "Juan Carlos",
    "apellido": "Pérez García",
    "modificado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `PUT /api/v1/mi-password` — Cambiar mi contraseña 🔒

Cambia la contraseña del usuario autenticado.

**Request body**

```json
{
  "password_actual": "vieja123",
  "nueva_password": "nueva123"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `password_actual` | `string` | Sí | Contraseña actual |
| `nueva_password` | `string` | Sí | Nueva contraseña |

**Response** `200 OK`

```json
{
  "data": {
    "modificado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

### Usuarios (admin)

#### `POST /api/v1/usuarios` — Crear usuario 🔒

Crea un nuevo usuario en el sistema (requiere permisos de administración).

**Request body**

```json
{
  "correo": "juan@correo.com",
  "nombre": "Juan",
  "apellido": "Pérez",
  "password": "secreto123"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `correo` | `string` | Sí | Correo electrónico del usuario |
| `nombre` | `string` | Sí | Nombre del usuario |
| `apellido` | `string` | Sí | Apellido del usuario |
| `password` | `string` | Sí | Contraseña (mínimo 8 caracteres) |

**Response** `201 Created`

```json
{
  "data": {
    "id": "01926b1e-dead-beef-cafe-000000000001",
    "correo": "juan@correo.com",
    "nombre": "Juan",
    "apellido": "Pérez",
    "activo": true,
    "creado_en": "2026-05-23T12:00:00Z"
  }
}
```

---

#### `GET /api/v1/usuarios` — Listar usuarios 🔒

Lista usuarios del sistema con filtros y paginación.

**Query parameters**

| Parámetro | Tipo | Default | Descripción |
|---|---|---|---|
| `pagina` | `int` | `1` | Número de página (1-based) |
| `tamano` | `int` | `20` | Elementos por página (máx. 100) |
| `correo` | `string` | — | Filtrar por correo exacto |
| `estado` | `string` | — | Filtrar por estado (ej: `ACTIVO`, `INACTIVO`) |

**Response** `200 OK`

```json
{
  "data": {
    "usuarios": [
      {
        "id": "01926b1e-dead-beef-cafe-000000000001",
        "correo": "juan@correo.com",
        "nombre": "Juan",
        "apellido": "Pérez",
        "estado": "ACTIVO",
        "creado_en": "2026-05-23T12:00:00Z"
      }
    ],
    "total": 1,
    "pagina": 1,
    "tamano": 20
  }
}
```

---

#### `PUT /api/v1/usuarios/{usuarioID}` — Modificar usuario 🔒

Modifica los datos de un usuario existente.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario a modificar |

**Request body**

```json
{
  "nombre": "Juan Actualizado",
  "apellido": "Pérez Actualizado"
}
```

**Response** `200 OK`

```json
{
  "data": {
    "id": "01926b1e-dead-beef-cafe-000000000001",
    "correo": "juan@correo.com",
    "nombre": "Juan Actualizado",
    "apellido": "Pérez Actualizado",
    "modificado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `DELETE /api/v1/usuarios/{usuarioID}` — Dar de baja usuario 🔒

Desactiva un usuario del sistema (baja lógica).

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario a dar de baja |

**Request body** (opcional)

```json
{
  "motivo": "Cierre de cuenta"
}
```

**Response** `200 OK`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "estado": "INACTIVO",
    "baja_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `POST /api/v1/usuarios/{usuarioID}/expulsar` — Expulsar usuario 🔒

Expulsa a un usuario del sistema, desactivándolo e invalidando todas sus sesiones.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario a expulsar |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "estado": "EXPULSADO",
    "sesiones_revocadas": 5,
    "expulsado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

### Seguridad (admin)

#### `POST /api/v1/usuarios/{usuarioID}/reset-password` — Resetear contraseña 🔒

Resetea la contraseña de un usuario (requiere permisos administrativos).

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario |

**Request body**

```json
{
  "nueva_password": "nueva123"
}
```

**Response** `200 OK`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "modificado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `POST /api/v1/usuarios/{usuarioID}/unlock` — Desbloquear cuenta 🔒

Desbloquea la cuenta de un usuario bloqueada por intentos fallidos.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario a desbloquear |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "desbloqueado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `GET /api/v1/ips-bloqueadas` — Listar IPs bloqueadas 🔒

Lista las direcciones IP bloqueadas temporalmente por exceso de intentos.

**Query parameters**

| Parámetro | Tipo | Default | Descripción |
|---|---|---|---|
| `pagina` | `int` | `1` | Número de página (1-based) |
| `tamano` | `int` | `20` | Elementos por página (máx. 100) |

**Response** `200 OK`

```json
{
  "data": {
    "ips": [
      {
        "ip": "192.168.1.100",
        "intentos": 10,
        "bloqueado_hasta": "2026-05-23T15:00:00Z"
      }
    ],
    "total": 1,
    "pagina": 1
  }
}
```

---

#### `DELETE /api/v1/ips-bloqueadas/{ip}` — Desbloquear IP 🔒

Elimina el bloqueo de una dirección IP.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `ip` | `string` | Dirección IP a desbloquear |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "ip": "192.168.1.100",
    "desbloqueado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `GET /api/v1/credenciales/{usuarioID}` — Consultar credenciales 🔒

Obtiene el estado de las credenciales de un usuario (bloqueo, intentos, verificación).

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "activo": true,
    "correo_verificado": true,
    "intentos_fallidos": 0,
    "bloqueado_hasta": ""
  }
}
```

---

### Sesiones (admin)

#### `GET /api/v1/sesiones` — Listar sesiones 🔒

Lista las sesiones activas del sistema con paginación.

**Query parameters**

| Parámetro | Tipo | Default | Descripción |
|---|---|---|---|
| `pagina` | `int` | `1` | Número de página (1-based) |
| `tamano` | `int` | `20` | Elementos por página (máx. 100) |

**Response** `200 OK`

```json
{
  "data": {
    "sesiones": [
      {
        "id": "ses-001",
        "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
        "ip_origen": "192.168.1.100",
        "estado": "ACTIVA",
        "ultima_actividad": "2026-05-23T13:00:00Z"
      }
    ],
    "total": 1,
    "pagina": 1
  }
}
```

---

#### `DELETE /api/v1/sesiones/{sesionID}` — Forzar cierre de sesión 🔒

Fuerza el cierre de una sesión específica (requiere permisos administrativos).

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `sesionID` | `string` | ID de la sesión a cerrar |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "sesion_id": "ses-001",
    "estado": "REVOCADA",
    "revocado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

### Roles y Permisos (RBAC, admin)

#### `GET /api/v1/roles` — Listar roles 🔒

Lista los roles del sistema con paginación.

**Query parameters**

| Parámetro | Tipo | Default | Descripción |
|---|---|---|---|
| `pagina` | `int` | `1` | Número de página (1-based) |
| `tamano` | `int` | `20` | Elementos por página (máx. 100) |

**Response** `200 OK`

```json
{
  "data": {
    "roles": [
      {
        "id": "rol-001",
        "nombre": "Administrador",
        "descripcion": "Acceso total al sistema",
        "es_sistema": true,
        "permisos": ["usuario:crear", "usuario:leer", "usuario:modificar", "usuario:eliminar"]
      }
    ],
    "total": 1,
    "pagina": 1
  }
}
```

---

#### `POST /api/v1/roles` — Crear rol 🔒

Crea un nuevo rol en el sistema con permisos opcionales.

**Request body**

```json
{
  "nombre": "Editor",
  "descripcion": "Puede editar contenidos",
  "permisos": ["usuario:leer"]
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `nombre` | `string` | Sí | Nombre del rol |
| `descripcion` | `string` | Sí | Descripción del rol |
| `permisos` | `string[]` | No | Códigos de permisos iniciales |

**Response** `201 Created`

```json
{
  "data": {
    "id": "rol-002",
    "nombre": "Editor",
    "descripcion": "Puede editar contenidos",
    "es_sistema": false,
    "creado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `PUT /api/v1/roles/{rolID}` — Modificar rol 🔒

Actualiza el nombre y descripción de un rol.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `rolID` | `string` | ID del rol a modificar |

**Request body**

```json
{
  "nombre": "Editor Senior",
  "descripcion": "Puede editar y publicar contenidos"
}
```

**Response** `200 OK`

```json
{
  "data": {
    "id": "rol-002",
    "nombre": "Editor Senior",
    "descripcion": "Puede editar y publicar contenidos",
    "modificado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `DELETE /api/v1/roles/{rolID}` — Eliminar rol 🔒

Elimina un rol del sistema (no se pueden eliminar roles de sistema).

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `rolID` | `string` | ID del rol a eliminar |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "rol_id": "rol-002",
    "eliminado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `POST /api/v1/usuarios/{usuarioID}/roles` — Asignar rol a usuario 🔒

Asigna un rol a un usuario, opcionalmente en un tenant específico.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario |

**Request body**

```json
{
  "rol_id": "rol-001",
  "tenant_id": ""
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `rol_id` | `string` | Sí | ID del rol a asignar |
| `tenant_id` | `string` | No | ID del tenant (vacío = global) |

**Response** `201 Created`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "rol_id": "rol-001",
    "tenant_id": "",
    "asignado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `DELETE /api/v1/usuarios/{usuarioID}/roles/{rolID}` — Revocar rol de usuario 🔒

Revoca un rol asignado a un usuario.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `usuarioID` | `string` | ID del usuario |
| `rolID` | `string` | ID del rol a revocar |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "rol_id": "rol-001",
    "tenant_id": "",
    "revocado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `POST /api/v1/roles/{rolID}/permisos` — Asignar permiso a rol 🔒

Asigna un permiso a un rol específico.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `rolID` | `string` | ID del rol |

**Request body**

```json
{
  "permiso_codigo": "usuario:crear"
}
```

**Response** `201 Created`

```json
{
  "data": {
    "rol_id": "rol-001",
    "permiso_codigo": "usuario:crear",
    "asignado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

#### `DELETE /api/v1/roles/{rolID}/permisos/{codigo}` — Revocar permiso de rol 🔒

Revoca un permiso previamente asignado a un rol.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `rolID` | `string` | ID del rol |
| `codigo` | `string` | Código del permiso a revocar |

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "rol_id": "rol-001",
    "permiso_codigo": "usuario:crear",
    "revocado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

### Tenants (admin)

#### `PUT /api/v1/tenants/{tenantID}` — Configurar tenant 🔒

Actualiza la configuración de un tenant.

**Parámetros de ruta**

| Parámetro | Tipo | Descripción |
|---|---|---|
| `tenantID` | `string` | ID del tenant a configurar |

**Request body**

```json
{
  "nombre": "Mi Empresa",
  "slug": "mi-empresa"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `nombre` | `string` | Sí | Nuevo nombre del tenant |
| `slug` | `string` | Sí | Nuevo slug del tenant |

**Response** `200 OK`

```json
{
  "data": {
    "tenant_id": "tenant-001",
    "nombre": "Mi Empresa",
    "slug": "mi-empresa",
    "modificado_en": "2026-05-23T14:00:00Z"
  }
}
```

---

### Verificación de Correo

#### `POST /api/v1/verificacion/solicitar` — Solicitar verificación de correo 🔒

Envía un enlace de verificación al correo del usuario autenticado.

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "mensaje": "Se ha enviado un enlace de verificación al correo registrado"
  }
}
```

---

#### `POST /api/v1/verificacion/confirmar` — Confirmar verificación de correo

Confirma la verificación del correo electrónico usando el token recibido.
Sin autenticación (el token de verificación es el medio de autenticación).

**Request body**

```json
{
  "token": "abc123..."
}
```

**Response** `200 OK`

```json
{
  "data": {
    "mensaje": "Correo verificado exitosamente"
  }
}
```

---

#### `POST /api/v1/verificacion/reenviar` — Reenviar verificación de correo 🔒

Reenvía el enlace de verificación al correo del usuario autenticado.

**Request body**: vacío

**Response** `200 OK`

```json
{
  "data": {
    "mensaje": "Se ha reenviado el enlace de verificación"
  }
}
```

---

### Recuperación de Contraseña

#### `POST /api/v1/recuperacion/solicitar` — Solicitar recuperación de contraseña

Envía un enlace de recuperación al correo electrónico proporcionado.
Sin autenticación.

**Request body**

```json
{
  "correo": "juan@correo.com"
}
```

**Response** `200 OK`

```json
{
  "data": {
    "mensaje": "Si el correo está registrado, recibirás un enlace de recuperación"
}
```

> ⚠️ Por seguridad, la respuesta es la misma exista o no el correo en el sistema.
> Esto evita que un atacante pueda enumerar usuarios registrados.

---

#### `POST /api/v1/recuperacion/validar` — Validar token de recuperación

Valida si un token de recuperación es válido y devuelve el ID del usuario asociado.
Sin autenticación.

**Request body**

```json
{
  "token": "abc123..."
}
```

**Response** `200 OK`

```json
{
  "data": {
    "usuario_id": "01926b1e-dead-beef-cafe-000000000001",
    "valido": true
  }
}
```

---

#### `POST /api/v1/recuperacion/confirmar` — Confirmar restablecimiento de contraseña

Restablece la contraseña usando el token de recuperación. Sin autenticación.

**Request body**

```json
{
  "token": "abc123...",
  "nueva_password": "nueva123"
}
```

| Campo | Tipo | Requerido | Descripción |
|---|---|---|---|
| `token` | `string` | Sí | Token de recuperación |
| `nueva_password` | `string` | Sí | Nueva contraseña |

**Response** `200 OK`

```json
{
  "data": {
    "mensaje": "Contraseña restablecida exitosamente"
  }
}
```
