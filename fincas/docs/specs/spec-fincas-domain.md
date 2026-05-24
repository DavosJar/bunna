---
title: Especificación del Microservicio Fincas — Dominio y API
version: 1.0
date_created: 2026-05-23
owner: Equipo Catastro
tags: fincas, lotes, dominio, api, rbac
---

# Especificación del Microservicio Fincas — Dominio y API

> **Propósito**: Definir el modelo de dominio, la API REST (contratos de integración), la integración RBAC y la persistencia del microservicio `fincas`. Los casos de uso de aplicación se definen en una especificación separada (`spec-fincas-application.md`). Este documento constituye la especificación detallada "AI-Ready" que define QUÉ se construye, siguiendo la arquitectura limpia y patrones establecidos en el servicio `identidad`.
>
> **No incluye**: Casos de uso de aplicación (ver `spec-fincas-application.md`), implementación concreta de handlers HTTP, repositorios GORM, migraciones de base de datos, pruebas unitarias detalladas (más allá de la estrategia), ni configuración de despliegue.
>
> **Formato**: Especificación de dominio + API que define QUÉ construir, no CÓMO implementarlo.

---

## 1. Propósito y Alcance

### 1.1 Propósito

Implementar un microservicio completo para la gestión de Fincas y Lotes del sistema Bunna/CafeScan, que permita a los usuarios (caficultores, agrónomos y administradores) realizar operaciones CRUD completas sobre sus unidades productivas, con control de acceso basado en roles (RBAC) y validación de propiedad de recursos.

### 1.2 Incluye

- Modelo de dominio completo: entidades `Finca` y `Lote` con validaciones de negocio
- API REST con 10 endpoints protegidos por JWT y autorización RBAC
- Middleware de autenticación JWT local (secreto compartido con `identidad`)
- Middleware de autorización con verificación de permisos atómicos
- 8 permisos RBAC nuevos para el módulo `fincas` (catálogo para seed en identidad)
- Esquema de persistencia PostgreSQL (2 tablas + índices)
- Integración con el sistema de autorización multi-tenant existente (futuro: claims enriquecidos)

### 1.3 No incluye

- Implementación de handlers HTTP concretos (capa presentación)
- Implementación de repositorios GORM (capa infraestructura)
- Migraciones de base de datos versionadas
- Comunicación con otros microservicios (futuro)
- Frontend React (capa de presentación web)
- Despliegue, contenedores Docker, CI/CD
- Auditoría de cambios sobre Fincas/Lotes

### 1.4 Audiencia

- Desarrolladores que implementarán el microservicio `fincas`
- Equipo de identidad (para coordinar seed de permisos)
- Arquitectos del sistema Bunna
- Generación de código asistida por IA

---

## 2. Definiciones

| Término | Definición |
|---------|-----------|
| **Finca** | Unidad principal de gestión agrícola. Entidad raíz del agregado. Pertenece a un usuario (caficultor o agronómo). Contiene nombre, ubicación y descripción. |
| **Lote** | Subdivisión espacial de una Finca. Entidad parte del agregado Finca. Contiene nombre, área (hectáreas) y descripción. No puede existir sin una Finca. |
| **usuarioID** | Identificador UUID del usuario propietario de la finca. Referencia lógica al sistema de identidad. Se obtiene del claim `sub` del JWT. |
| **tenantID** | Identificador UUID del tenant (organización) al que pertenece la finca. Opcional en versión inicial. Futuro: obligatorio para multi-tenant. |
| **Permiso atómico** | Capacidad indivisible que representa una acción concreta. Formato `modulo:recurso:verbo`. Constante de dominio definida en el catálogo de identidad. |
| **sys_admin** | Rol de sistema global. No atado a ningún tenant. Tiene todos los permisos sobre todos los módulos, incluyendo `fincas`. |
| **Clean Architecture** | Patrón arquitectónico con capas: Dominio → Aplicación → Infraestructura → Presentación. Dependencias hacia adentro. |
| **Agregado** | Grupo de entidades que se tratan como una unidad coherente de persistencia y consistencia transaccional. `Finca` es el agregado raíz; `Lote` es parte del agregado. |
| **ApiResponse[T]** | Estructura genérica de respuesta HTTP. Contiene `data` (tipo T), `_links` (HATEOAS), `message` (opcional). |
| **Problem Details** | Formato de error estándar RFC 9457. Campos: `type`, `title`, `status`, `detail`, `instance`. |
| **HATEOAS** | Hypermedia As The Engine Of Application State. Enlaces en las respuestas API que indican acciones disponibles sobre el recurso. |

---

## 3. Requisitos, Restricciones y Guías

### 3.1 Requisitos Funcionales

| ID | Descripción |
|----|-------------|
| **REQ-FIN-001** | El sistema debe permitir registrar una nueva Finca con nombre, ubicación, descripción (opcional) y usuarioID (del token JWT). |
| **REQ-FIN-002** | El sistema debe permitir consultar una Finca por su ID, solo si pertenece al usuario autenticado (o si el usuario es sys_admin/administrador). |
| **REQ-FIN-003** | El sistema debe permitir listar todas las Fincas del usuario autenticado, con filtros opcionales. El usuario sys_admin ve todas las fincas del sistema. |
| **REQ-FIN-004** | El sistema debe permitir actualizar los datos de una Finca existente, solo si el usuario es el propietario. |
| **REQ-FIN-005** | El sistema debe permitir eliminar una Finca existente. Si tiene Lotes asociados, requiere confirmación explícita (`confirm=true`). |
| **REQ-FIN-006** | El sistema debe permitir crear un Lote dentro de una Finca existente, siempre que el usuario sea el propietario de la Finca. |
| **REQ-FIN-007** | El sistema debe permitir consultar un Lote por su ID, solo si la Finca padre pertenece al usuario autenticado. |
| **REQ-FIN-008** | El sistema debe permitir listar todos los Lotes de una Finca, solo si la Finca pertenece al usuario. |
| **REQ-FIN-009** | El sistema debe permitir actualizar un Lote existente, solo si el usuario es propietario de la Finca padre. |
| **REQ-FIN-010** | El sistema debe permitir eliminar un Lote existente. En el futuro verificará muestras/diagnósticos asociados. |
| **REQ-FIN-011** | El sistema debe validar todos los datos de entrada antes de persistir (nombre no vacío, área positiva, etc.). |
| **REQ-FIN-012** | El sistema debe persistir todas las operaciones en PostgreSQL. |
| **REQ-FIN-013** | El sistema debe retornar errores descriptivos siguiendo RFC 9457 (Problem Details) cuando una operación falle. |
| **REQ-FIN-014** | El sistema debe exponer las respuestas exitosas en formato ApiResponse[T] con enlaces HATEOAS. |

### 3.2 Requisitos de Arquitectura

| ID | Descripción |
|----|-------------|
| **ARQ-FIN-001** | Backend implementado en Go (golang). |
| **ARQ-FIN-002** | Estructura de carpetas sigue Clean Architecture con capas: `domain/`, `application/`, `infrastructure/`, `presentation/`. |
| **ARQ-FIN-003** | Las dependencias fluyen hacia adentro: Presentación → Aplicación → Dominio. Infraestructura implementa interfaces del Dominio. |
| **ARQ-FIN-004** | La capa de presentación sigue el flujo `Handler → Facade → Mapper → Domain`. Handler nunca importa domain ni mapper directamente. Facade nunca importa gin ni nada HTTP. Mapper no tiene lógica de negocio, solo conversión de structs. |
| **ARQ-FIN-005** | Las respuestas HTTP usan formato genérico `ApiResponse[T]` con HATEOAS construido por el Handler. |
| **ARQ-FIN-006** | Los errores HTTP siguen RFC 9457 (Problem Details). |
| **ARQ-FIN-007** | Se usa Gin como router HTTP. |
| **ARQ-FIN-008** | La validación JWT se realiza localmente con el secreto compartido (`JWT_SECRET`), sin llamar al servicio identidad. |

### 3.3 Restricciones Técnicas

| ID | Descripción |
|----|-------------|
| **CON-FIN-001** | Base de datos: PostgreSQL. |
| **CON-FIN-002** | IDs de Finca y Lote: UUID v4 (almacenados como VARCHAR(36)). |
| **CON-FIN-003** | Fechas: `time.Time` en UTC. |
| **CON-FIN-004** | El campo `area` del Lote usa precisión DECIMAL(10,2) en BD, float64 en Go (hectáreas con 2 decimales). |
| **CON-FIN-005** | No se almacenan contraseñas ni secretos en este dominio. |
| **CON-FIN-006** | No existe código cíclico entre capas. Dominio no importa infraestructura ni presentación. |
| **CON-FIN-007** | `usuario_id` en fincas es una referencia lógica al sistema de identidad. No existe FK física a tabla `usuarios` de otro servicio. |
| **CON-FIN-008** | El JWT se valida localmente en el middleware del servicio fincas usando `golang-jwt/jwt/v5`. |
| **CON-FIN-009** | La versión inicial no incluye claims enriquecidos de tenant/permisos en el JWT. La autorización se basa en verificación de propiedad del recurso. |

### 3.4 Guías de Diseño

| ID | Descripción |
|----|-------------|
| **GUD-FIN-001** | El dominio modela entidades ricas con comportamiento y validaciones, no estructuras anémicas. |
| **GUD-FIN-002** | Los errores de dominio se definen como variables exportadas en un archivo `domain/errores.go`. |
| **GUD-FIN-003** | Las interfaces de repositorio se definen en el dominio, se implementan en infraestructura. |
| **GUD-FIN-004** | Los casos de uso (servicios de aplicación) se implementan como struct con métodos, un archivo por funcionalidad. Se definen en `spec-fincas-application.md`. |
| **GUD-FIN-005** | Los mappers no contienen lógica de negocio: solo convierten entre estructuras de dominio, DTOs de aplicación y DTOs de presentación. |
| **GUD-FIN-006** | Un Lote nunca existe sin una Finca. La creación de un Lote siempre valida que la Finca exista y pertenezca al usuario. |
| **GUD-FIN-007** | La eliminación de una Finca con Lotos requiere confirmación explícita del usuario. |
| **GUD-FIN-008** | Los mensajes de error al usuario son descriptivos pero no revelan detalles internos de implementación. |
| **GUD-FIN-009** | Toda operación de escritura debe ser transaccional (Unit of Work o transacción de base de datos). |

---

## 4. RBAC Integration

### 4.1 Permisos del Módulo Fincas

Se agregan 8 permisos atómicos al catálogo de permisos del sistema identidad:

| Código | Nombre | Descripción |
|--------|--------|-------------|
| `fincas:finca:crear` | Crear Finca | Crear una nueva finca |
| `fincas:finca:modificar` | Modificar Finca | Actualizar datos de una finca existente |
| `fincas:finca:eliminar` | Eliminar Finca | Eliminar una finca |
| `fincas:finca:consultar` | Consultar Finca | Listar y ver detalles de fincas |
| `fincas:lote:crear` | Crear Lote | Crear un nuevo lote dentro de una finca |
| `fincas:lote:modificar` | Modificar Lote | Actualizar datos de un lote existente |
| `fincas:lote:eliminar` | Eliminar Lote | Eliminar un lote |
| `fincas:lote:consultar` | Consultar Lote | Listar y ver detalles de lotes |

### 4.2 Matriz de Permisos por Rol

| Permiso | sys_admin | administrador | agronomo | caficultor |
|---------|:---------:|:-------------:|:--------:|:----------:|
| `fincas:finca:crear` | ✅ | ✅ | ✅ | ❌ |
| `fincas:finca:modificar` | ✅ | ✅ | ✅ | ❌ |
| `fincas:finca:eliminar` | ✅ | ✅ | ❌ | ❌ |
| `fincas:finca:consultar` | ✅ | ✅ | ✅ | ✅ |
| `fincas:lote:crear` | ✅ | ✅ | ✅ | ❌ |
| `fincas:lote:modificar` | ✅ | ✅ | ✅ | ❌ |
| `fincas:lote:eliminar` | ✅ | ✅ | ❌ | ❌ |
| `fincas:lote:consultar` | ✅ | ✅ | ✅ | ✅ |

### 4.3 Coordinación con Identidad

Este spec NO modifica identidad directamente. El equipo de identidad debe:

1. Agregar los 8 permisos del módulo `fincas` al catálogo de permisos (constantes de dominio en identidad).
2. Actualizar la matriz de permisos en el seed de roles de identidad para que los roles incluyan estos permisos.
3. El seed debe ser idempotente: los permisos nuevos se insertan sin duplicar los existentes.

**Tabla de referencia de roles existentes en identidad:**

| Rol | Ámbito | `es_sistema` |
|-----|--------|:------------:|
| `sys_admin` | Global | ✅ |
| `administrador` | Tenant | ✅ |
| `agronomo` | Tenant | ✅ |
| `caficultor` | Tenant | ✅ |

### 4.4 Estrategia de Verificación de Autorización (Versión Inicial)

En la versión inicial, el middleware de autorización en `fincas` NO consulta la base de datos de identidad. En su lugar:

1. **Verificación de propiedad**: Para operaciones sobre recursos existentes (modificar, eliminar, consultar), se verifica que el `usuarioID` del recurso coincida con el `usuarioID` del token JWT.
2. **Verificación granular por permiso**: El handler verifica el permiso requerido según el endpoint. En versión inicial esto se implementa como una función simple que valida propiedad de recurso.
3. **sys_admin**: El JWT no incluye claim `global` en versión inicial. El rol sys_admin se identifica por configuración o se omite inicialmente. Cuando esté disponible el claim enriquecido, sys_admin (`global: true`) será siempre autorizado sin verificar propiedad.

### 4.5 Estrategia de Verificación de Autorización (Versión Futura)

Cuando identidad implemente claims JWT enriquecidos con `global`, `tenants[].roles[]` y `tenants[].permisos[]`, el middleware de autorización en `fincas` deberá:

1. Extraer `global` del JWT: si es `true`, autorizar siempre.
2. Extraer `tenants[].permisos[]`: verificar si el permiso requerido está en la lista de permisos del tenant actual.
3. Si no tiene el permiso → HTTP 403 Forbidden con error genérico (sin revelar qué permiso falta).

---

## 5. Domain Model

### 5.1 Entidad Finca

```
Entidad: Finca (Agregado Raíz)
├── id:          string (UUID v4, inmutable)
├── nombre:      string (requerido, 3-200 caracteres)
├── ubicacion:   string (requerido, max 500 caracteres)
├── descripcion: string (opcional, max 1000 caracteres, default "")
├── usuarioID:   string (ID del propietario, del token JWT)
├── tenantID:    *string (opcional, ID del tenant para multi-tenant futuro)
├── createdAt:   time.Time
└── updatedAt:   time.Time
```

**Comportamiento / Invariantes:**

- `nombre` no puede estar vacío, debe tener entre 3 y 200 caracteres.
- `ubicacion` no puede estar vacío, máximo 500 caracteres.
- `descripcion` es opcional, máximo 1000 caracteres.
- `id` se genera como UUID v4 al crear la entidad. Es inmutable después de la creación.
- `usuarioID` se asigna del token JWT, nunca se recibe del cliente.
- `tenantID` es opcional en versión inicial.
- `createdAt` y `updatedAt` se establecen en el momento de creación/actualización.

**Métodos de dominio:**

| Método | Descripción |
|--------|-------------|
| `NewFinca(nombre, ubicacion, descripcion, usuarioID)` | Constructor con validaciones. Retorna `(Finca, error)`. |
| `Actualizar(nombre, ubicacion, descripcion)` | Actualiza campos permitidos. Retorna `error`. |
| `EsPropietario(usuarioID)` | Retorna `true` si el usuarioID coincide con el propietario. |
| `TieneLotes(cantidad)` | Verificación externa (recibe conteo de lotes, no consulta BD directamente). |

### 5.2 Entidad Lote

```
Entidad: Lote (Parte del Agregado Finca)
├── id:          string (UUID v4, inmutable)
├── fincaID:     string (UUID, FK lógica a Finca)
├── nombre:      string (requerido, 3-200 caracteres)
├── area:        float64 (requerido, > 0, hectáreas con 2 decimales)
├── descripcion: string (opcional, max 1000 caracteres, default "")
├── createdAt:   time.Time
└── updatedAt:   time.Time
```

**Comportamiento / Invariantes:**

- `nombre` no puede estar vacío, debe tener entre 3 y 200 caracteres.
- `area` debe ser mayor a 0 (positiva). Se almacena con precisión de 2 decimales.
- `descripcion` es opcional, máximo 1000 caracteres.
- `fincaID` es obligatorio. Un Lote siempre pertenece a una Finca.
- `id` se genera como UUID v4 al crear la entidad. Es inmutable después de la creación.
- `createdAt` y `updatedAt` se establecen en el momento de creación/actualización.

**Métodos de dominio:**

| Método | Descripción |
|--------|-------------|
| `NewLote(fincaID, nombre, area, descripcion)` | Constructor con validaciones. Retorna `(Lote, error)`. |
| `Actualizar(nombre, area, descripcion)` | Actualiza campos permitidos. Retorna `error`. |

### 5.3 Errores de Dominio

| Variable | Mensaje | Significado |
|----------|---------|-------------|
| `ErrNombreFincaRequerido` | "El nombre de la finca es requerido" | nombre vacío o menor a 3 chars |
| `ErrNombreFincaLargo` | "El nombre de la finca no puede exceder 200 caracteres" | nombre > 200 chars |
| `ErrUbicacionRequerida` | "La ubicación de la finca es requerida" | ubicación vacía |
| `ErrUbicacionLarga` | "La ubicación no puede exceder 500 caracteres" | ubicación > 500 chars |
| `ErrDescripcionLarga` | "La descripción no puede exceder 1000 caracteres" | descripción > 1000 chars |
| `ErrNombreLoteRequerido` | "El nombre del lote es requerido" | nombre vacío o menor a 3 chars |
| `ErrNombreLoteLargo` | "El nombre del lote no puede exceder 200 caracteres" | nombre > 200 chars |
| `ErrAreaRequerida` | "El área del lote es requerida y debe ser mayor a 0" | area <= 0 |
| `ErrFincaNoEncontrada` | "Finca con ID {id} no encontrada" | finca no existe |
| `ErrLoteNoEncontrado` | "Lote con ID {id} no encontrado" | lote no existe |
| `ErrNoPropietario` | "No tienes permisos sobre este recurso" | usuario no es propietario |
| `ErrFincaConLotes` | "La finca tiene {n} lote(s) asociado(s). Confirma la eliminación" | finca con lotes sin confirm |

### 5.4 Interfaces de Repositorio (Definidas en Dominio)

Los repositorios definen exclusivamente operaciones de persistencia (CRUD real). NO contienen lógica de negocio ni representan casos de uso. Cada método tiene una responsabilidad única y bien definida.

```go
// FincaRepositorio define las operaciones de persistencia para Finca.
// NO contiene lógica de negocio. Solo CRUD + consultas específicas.
type FincaRepositorio interface {
    // Crear persiste una nueva finca. Retorna error si ya existe una con el mismo ID.
    Crear(ctx context.Context, finca *Finca) error
    
    // ObtenerPorID busca una finca por su ID. Retorna nil si no existe.
    ObtenerPorID(ctx context.Context, id string) (*Finca, error)
    
    // ListarPorUsuario retorna todas las fincas de un usuario. Array vacío si no hay.
    ListarPorUsuario(ctx context.Context, usuarioID string) ([]Finca, error)
    
    // ListarTodas retorna todas las fincas del sistema (solo para administradores).
    ListarTodas(ctx context.Context) ([]Finca, error)
    
    // Actualizar persiste los cambios de una finca existente.
    Actualizar(ctx context.Context, finca *Finca) error
    
    // Eliminar borra una finca por su ID. Error si no existe.
    Eliminar(ctx context.Context, id string) error
    
    // ContarLotes retorna la cantidad de lotes asociados a una finca.
    ContarLotes(ctx context.Context, fincaID string) (int, error)
}

// LoteRepositorio define las operaciones de persistencia para Lote.
// NO contiene lógica de negocio. Solo CRUD + consultas específicas.
type LoteRepositorio interface {
    // Crear persiste un nuevo lote. Retorna error si ya existe uno con el mismo ID.
    Crear(ctx context.Context, lote *Lote) error
    
    // ObtenerPorID busca un lote por su ID. Retorna nil si no existe.
    ObtenerPorID(ctx context.Context, id string) (*Lote, error)
    
    // ListarPorFinca retorna todos los lotes de una finca. Array vacío si no hay.
    ListarPorFinca(ctx context.Context, fincaID string) ([]Lote, error)
    
    // Actualizar persiste los cambios de un lote existente.
    Actualizar(ctx context.Context, lote *Lote) error
    
    // Eliminar borra un lote por su ID. Error si no existe.
    Eliminar(ctx context.Context, id string) error
}
```

---

## 6. Casos de Uso (Futuro)

Los casos de uso del microservicio `fincas` se definirán en una sub-especificación posterior (`spec-fincas-application.md`). 
Esta especificación se centra exclusivamente en la capa de dominio.

Los casos de uso representan capacidades del sistema, no operaciones CRUD. 
Ejemplos de casos de uso futuros:
- "Registrar una nueva finca" 
- "Transferir una finca a otro agricultor"
- "Generar un nuevo lote dentro de una finca"
- "Dividir un lote existente en dos lotes más pequeños"
- "Cerrar una finca (eliminación lógica con verificación de dependencias)"

---

## 7. API Contracts

### 7.1 Formato General de Respuesta Exitosa

```json
{
  "data": { },
  "_links": {
    "self": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000", "method": "GET" },
    "list": { "href": "/api/v1/fincas", "method": "GET" },
    "update": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000", "method": "PUT" },
    "delete": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000", "method": "DELETE" }
  }
}
```

Para respuestas de lista:

```json
{
  "data": [ ],
  "_links": {
    "self": { "href": "/api/v1/fincas", "method": "GET" },
    "create": { "href": "/api/v1/fincas", "method": "POST" }
  }
}
```

### 7.2 Formato General de Error (RFC 9457)

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "Finca con ID 123e4567-e89b-12d3-a456-426614174000 no encontrada",
  "instance": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000"
}
```

### 7.3 Endpoints

#### POST /api/v1/fincas

> Crear una nueva finca.

- **Permiso requerido**: `fincas:finca:crear`
- **Autenticación**: Requerida (JWT)
- **Body**:

```json
{
  "nombre": "Finca El Progreso",
  "ubicacion": "Vereda La Esperanza, Municipio de Andes, Antioquia",
  "descripcion": "Finca dedicada al cultivo de café orgánico certificado"
}
```

- **Response 201**:

```json
{
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "nombre": "Finca El Progreso",
    "ubicacion": "Vereda La Esperanza, Municipio de Andes, Antioquia",
    "descripcion": "Finca dedicada al cultivo de café orgánico certificado",
    "usuarioID": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "createdAt": "2026-05-23T14:30:00Z",
    "updatedAt": "2026-05-23T14:30:00Z"
  },
  "_links": {
    "self": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000", "method": "GET" },
    "list": { "href": "/api/v1/fincas", "method": "GET" },
    "update": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000", "method": "PUT" },
    "delete": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000", "method": "DELETE" }
  }
}
```

- **Response 400** (validación):

```json
{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "detail": "El nombre de la finca es requerido y debe tener entre 3 y 200 caracteres",
  "instance": "/api/v1/fincas"
}
```

#### GET /api/v1/fincas

> Listar fincas del usuario autenticado.

- **Permiso requerido**: `fincas:finca:consultar`
- **Autenticación**: Requerida (JWT)
- **Query params opcionales**: `?nombre=progreso` (búsqueda parcial, futura implementación)
- **Response 200**:

```json
{
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "nombre": "Finca El Progreso",
      "ubicacion": "Vereda La Esperanza, Municipio de Andes, Antioquia",
      "descripcion": "Finca dedicada al cultivo de café orgánico certificado",
      "usuarioID": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "createdAt": "2026-05-23T14:30:00Z",
      "updatedAt": "2026-05-23T14:30:00Z"
    }
  ],
  "_links": {
    "self": { "href": "/api/v1/fincas", "method": "GET" },
    "create": { "href": "/api/v1/fincas", "method": "POST" }
  }
}
```

#### GET /api/v1/fincas/{id}

> Obtener detalle de una finca.

- **Permiso requerido**: `fincas:finca:consultar`
- **Autenticación**: Requerida (JWT)
- **Response 200**: Igual que el objeto `data` de creación, con enlaces HATEOAS.
- **Response 404**:

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "Finca con ID 123e4567-e89b-12d3-a456-426614174000 no encontrada",
  "instance": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000"
}
```

- **Response 403** (no propietario):

```json
{
  "type": "about:blank",
  "title": "Forbidden",
  "status": 403,
  "detail": "No tienes permisos sobre este recurso",
  "instance": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000"
}
```

#### PUT /api/v1/fincas/{id}

> Actualizar una finca existente.

- **Permiso requerido**: `fincas:finca:modificar`
- **Autenticación**: Requerida (JWT)
- **Body**:

```json
{
  "nombre": "Finca El Progreso (Renovada)",
  "ubicacion": "Vereda La Esperanza, Municipio de Andes, Antioquia",
  "descripcion": "Finca renovada con nuevas variedades de café"
}
```

- **Response 200**: Objeto `data` actualizado con `updatedAt` actualizado, más enlaces HATEOAS.
- **Response 400**: Error de validación.
- **Response 403**: No propietario.
- **Response 404**: Finca no encontrada.

#### DELETE /api/v1/fincas/{id}

> Eliminar una finca.

- **Permiso requerido**: `fincas:finca:eliminar`
- **Autenticación**: Requerida (JWT)
- **Query params**: `?confirm=true` (opcional, requerido si la finca tiene lotes)
- **Response 200** (eliminación exitosa):

```json
{
  "data": null,
  "_links": {
    "list": { "href": "/api/v1/fincas", "method": "GET" },
    "create": { "href": "/api/v1/fincas", "method": "POST" }
  }
}
```

- **Response 409** (finca con lotes, sin confirmación):

```json
{
  "type": "about:blank",
  "title": "Conflict",
  "status": 409,
  "detail": "La finca tiene 3 lote(s) asociado(s). Confirma la eliminación con ?confirm=true",
  "instance": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000"
}
```

#### POST /api/v1/fincas/{fincaId}/lotes

> Crear un lote dentro de una finca.

- **Permiso requerido**: `fincas:lote:crear`
- **Autenticación**: Requerida (JWT)
- **Body**:

```json
{
  "nombre": "Lote A",
  "area": 5.25,
  "descripcion": "Lote principal de café variedad Castillo"
}
```

- **Response 201**:

```json
{
  "data": {
    "id": "223e4567-e89b-12d3-a456-426614174001",
    "fincaID": "123e4567-e89b-12d3-a456-426614174000",
    "nombre": "Lote A",
    "area": 5.25,
    "descripcion": "Lote principal de café variedad Castillo",
    "createdAt": "2026-05-23T15:00:00Z",
    "updatedAt": "2026-05-23T15:00:00Z"
  },
  "_links": {
    "self": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000/lotes/223e4567-e89b-12d3-a456-426614174001", "method": "GET" },
    "list": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000/lotes", "method": "GET" },
    "update": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000/lotes/223e4567-e89b-12d3-a456-426614174001", "method": "PUT" },
    "delete": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000/lotes/223e4567-e89b-12d3-a456-426614174001", "method": "DELETE" },
    "finca": { "href": "/api/v1/fincas/123e4567-e89b-12d3-a456-426614174000", "method": "GET" }
  }
}
```

#### GET /api/v1/fincas/{fincaId}/lotes

> Listar lotes de una finca.

- **Permiso requerido**: `fincas:lote:consultar`
- **Autenticación**: Requerida (JWT)
- **Response 200**: Array de lotes con enlaces HATEOAS.

#### GET /api/v1/fincas/{fincaId}/lotes/{id}

> Obtener detalle de un lote.

- **Permiso requerido**: `fincas:lote:consultar`
- **Autenticación**: Requerida (JWT)
- **Response 200**: Objeto lote con enlaces HATEOAS.
- **Response 404**: Lote no encontrado.

#### PUT /api/v1/fincas/{fincaId}/lotes/{id}

> Actualizar un lote existente.

- **Permiso requerido**: `fincas:lote:modificar`
- **Autenticación**: Requerida (JWT)
- **Body**: Mismos campos que creación (nombre, area, descripcion).
- **Response 200**: Objeto lote actualizado con enlaces HATEOAS.

#### DELETE /api/v1/fincas/{fincaId}/lotes/{id}

> Eliminar un lote.

- **Permiso requerido**: `fincas:lote:eliminar`
- **Autenticación**: Requerida (JWT)
- **Query params**: `?confirm=true` (opcional, para versión futura con muestras)
- **Response 200**: Eliminación exitosa.
- **Response 404**: Lote no encontrado.

### 7.4 Mapa de Códigos HTTP

| Condición | Código HTTP | Uso |
|-----------|------------|-----|
| Operación exitosa (creación) | 201 Created | POST (fincas y lotes) |
| Operación exitosa (lectura, actualización, eliminación) | 200 OK | GET, PUT, DELETE |
| Datos inválidos (validación de dominio) | 400 Bad Request | Campos con error específico |
| No autenticado | 401 Unauthorized | Token ausente, inválido o expirado |
| Sin permisos sobre el recurso | 403 Forbidden | No es propietario o no tiene el permiso |
| Recurso no encontrado | 404 Not Found | Finca o Lote no existe |
| Conflicto (dependencias activas) | 409 Conflict | Finca con lotes sin confirmar |
| Error interno del servidor | 500 Internal Server Error | Error inesperado, mensaje genérico |

### 7.5 Resumen de Endpoints

| Método | Ruta | Permiso Requerido | Códigos de Respuesta |
|--------|------|-------------------|---------------------|
| POST | `/api/v1/fincas` | `fincas:finca:crear` | 201, 400, 401, 500 |
| GET | `/api/v1/fincas` | `fincas:finca:consultar` | 200, 401, 500 |
| GET | `/api/v1/fincas/{id}` | `fincas:finca:consultar` | 200, 401, 403, 404, 500 |
| PUT | `/api/v1/fincas/{id}` | `fincas:finca:modificar` | 200, 400, 401, 403, 404, 500 |
| DELETE | `/api/v1/fincas/{id}` | `fincas:finca:eliminar` | 200, 401, 403, 404, 409, 500 |
| POST | `/api/v1/fincas/{fincaId}/lotes` | `fincas:lote:crear` | 201, 400, 401, 403, 404, 500 |
| GET | `/api/v1/fincas/{fincaId}/lotes` | `fincas:lote:consultar` | 200, 401, 403, 404, 500 |
| GET | `/api/v1/fincas/{fincaId}/lotes/{id}` | `fincas:lote:consultar` | 200, 401, 403, 404, 500 |
| PUT | `/api/v1/fincas/{fincaId}/lotes/{id}` | `fincas:lote:modificar` | 200, 400, 401, 403, 404, 500 |
| DELETE | `/api/v1/fincas/{fincaId}/lotes/{id}` | `fincas:lote:eliminar` | 200, 401, 403, 404, 500 |

---

## 8. Authorization Flow

### 8.1 Middleware JWT (Autenticación)

El middleware se implementa en `internal/presentation/middleware/jwt_middleware.go`.

**Flujo:**

```
1. Extraer header: Authorization: Bearer <token>
   │
2. Validar que el header exista y tenga formato "Bearer <token>"
   │   └── Si no → HTTP 401, error: "Token de autenticación no proporcionado"
   │
3. Parsear y validar el JWT usando JWT_SECRET (golang-jwt/jwt/v5, algoritmo HS256)
   │   ├── Token inválido → HTTP 401, error: "Token inválido o expirado"
   │   ├── Token expirado → HTTP 401, error: "Token expirado"
   │   └── Válido → continuar
   │
4. Extraer claims:
   │   ├── sub (usuarioID) → string
   │   ├── sid (sesionID) → string
   │   └── (futuro) global, tenants → desde claims enriquecidos
   │
5. Inyectar usuarioID en el contexto Gin: c.Set("usuarioID", usuarioID)
   │
6. Llamar a c.Next() para continuar con el siguiente middleware/handler
```

**Estructura del JWT validado (versión inicial):**

```go
type Claims struct {
    jwt.RegisteredClaims
    // Versión inicial: solo sub y sid
    // Versión futura: global bool, tenants map[string]TenantInfo
}
```

**Configuración:**

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `JWT_SECRET` | Secreto compartido con identidad | (sin default, requerida) |
| `JWT_ISSUER` | Issuer del token (validación) | `bunna-identidad` |

### 8.2 Middleware de Autorización (Authz)

El middleware se implementa en `internal/presentation/middleware/authz_middleware.go`.

**Flujo (versión inicial):**

```
1. Recibir el permiso requerido como parámetro del middleware
   │   Ejemplo: router.Use(authzMiddleware.RequerirPermiso("fincas:finca:crear"))
   │
2. Obtener usuarioID del contexto Gin (inyectado por JWT middleware)
   │
3. Por ahora: la verificación es delegada al handler o facade
   │   (en versión inicial, la propiedad del recurso se verifica en el caso de uso)
   │
4. Versión futura:
   │   ├── Extraer claims enriquecidos del contexto (global, tenants[].permisos)
   │   ├── Si global=true → autorizar siempre
   │   ├── Si permiso está en la lista del tenant actual → autorizar
   │   └── Si no → HTTP 403 Forbidden con error genérico
   │
5. Llamar a c.Next() si autorizado, o abortar con 403
```

**Flujo de verificación de propiedad (en caso de uso / facade):**

```
Para operaciones sobre recursos existentes (GET, PUT, DELETE por ID):
  1. Obtener el recurso (Finca o Lote) del repositorio por ID
  2. Comparar recurso.usuarioID con usuarioID del contexto
  3. Si no coinciden → error ErrNoPropietario → HTTP 403 Forbidden
```

### 8.3 Claims JWT de Identidad (Referencia)

El servicio `identidad` genera JWT con los siguientes claims. El microservicio `fincas` debe parsear al menos `sub` y `sid`.

**Claims actuales (identidad):**

```json
{
  "sub": "uuid-usuario",
  "sid": "uuid-sesion",
  "iat": 1716300000,
  "exp": 1716300900,
  "typ": "access"
}
```

**Claims futuros (enriquecidos con RBAC):**

```json
{
  "sub": "uuid-usuario",
  "sid": "uuid-sesion",
  "global": false,
  "tenants": {
    "tenant-uuid-1": {
      "slug": "nombre-del-tenant",
      "roles": ["administrador"],
      "permisos": ["fincas:finca:crear", "fincas:finca:modificar", "..."]
    }
  }
}
```

---

## 9. Persistence

### 9.1 Esquema PostgreSQL

```sql
-- Tabla: fincas
CREATE TABLE fincas (
    id          VARCHAR(36)     PRIMARY KEY,
    nombre      VARCHAR(200)    NOT NULL,
    ubicacion   VARCHAR(500)    NOT NULL,
    descripcion TEXT            NOT NULL DEFAULT '',
    usuario_id  VARCHAR(36)     NOT NULL,
    tenant_id   VARCHAR(36),
    created_at  TIMESTAMP       NOT NULL,
    updated_at  TIMESTAMP       NOT NULL,
    deleted_at  TIMESTAMP
);

-- Tabla: lotes
CREATE TABLE lotes (
    id          VARCHAR(36)     PRIMARY KEY,
    finca_id    VARCHAR(36)     NOT NULL REFERENCES fincas(id) ON DELETE CASCADE,
    nombre      VARCHAR(200)    NOT NULL,
    area        DECIMAL(10,2)   NOT NULL,
    descripcion TEXT            NOT NULL DEFAULT '',
    created_at  TIMESTAMP       NOT NULL,
    updated_at  TIMESTAMP       NOT NULL,
    deleted_at  TIMESTAMP
);

-- Índices
CREATE INDEX idx_fincas_usuario_id ON fincas(usuario_id);
CREATE INDEX idx_fincas_tenant_id  ON fincas(tenant_id);
CREATE INDEX idx_lotes_finca_id    ON lotes(finca_id);
```

### 9.2 Notas sobre el Esquema

- `deleted_at` soporta soft-delete opcional (NULL = activo, NOT NULL = eliminado lógicamente). En versión inicial puede omitirse (eliminación física).
- `usuario_id` en `fincas` es FK lógica. No existe restricción `REFERENCES` hacia una tabla `usuarios` porque el usuario pertenece al servicio `identidad`.
- `finca_id` en `lotes` SÍ tiene FK física con `ON DELETE CASCADE` porque ambas tablas están en el mismo servicio.
- `tenant_id` es nullable. En versión inicial será NULL para todos los registros.
- La columna `area` usa `DECIMAL(10,2)` para precisión exacta de 2 decimales en hectáreas.
- Se recomienda usar migraciones versionadas (archivos SQL) para producción y `AutoMigrate` de GORM para desarrollo.

### 9.3 Modelo de Persistencia vs. Modelo de Dominio

El modelo de persistencia (estructuras GORM) puede diferir del modelo de dominio. La capa de infraestructura es responsable de mapear entre ambos:

| Campo Dominio | Columna BD | Tipo BD | Notas |
|---------------|-----------|---------|-------|
| `Finca.ID` | `id` | VARCHAR(36) | UUID v4 como string |
| `Finca.Nombre` | `nombre` | VARCHAR(200) | - |
| `Finca.Ubicacion` | `ubicacion` | VARCHAR(500) | - |
| `Finca.Descripcion` | `descripcion` | TEXT | Opcional, default "" |
| `Finca.UsuarioID` | `usuario_id` | VARCHAR(36) | FK lógica a identidad |
| `Finca.TenantID` | `tenant_id` | VARCHAR(36) | Nullable, futuro multi-tenant |
| `Finca.CreatedAt` | `created_at` | TIMESTAMP | UTC |
| `Finca.UpdatedAt` | `updated_at` | TIMESTAMP | UTC |
| `Lote.ID` | `id` | VARCHAR(36) | UUID v4 como string |
| `Lote.FincaID` | `finca_id` | VARCHAR(36) | FK a fincas |
| `Lote.Nombre` | `nombre` | VARCHAR(200) | - |
| `Lote.Area` | `area` | DECIMAL(10,2) | Float64 en Go, decimal en BD |
| `Lote.Descripcion` | `descripcion` | TEXT | Opcional, default "" |
| `Lote.CreatedAt` | `created_at` | TIMESTAMP | UTC |
| `Lote.UpdatedAt` | `updated_at` | TIMESTAMP | UTC |

---

## 10. Estructura del Módulo

```
fincas/
├── cmd/
│   └── main.go                          ← Punto de entrada
├── internal/
│   ├── config/
│   │   └── config.go                    ← Variables de entorno (JWT_SECRET, DB, puerto, etc.)
│   ├── fincas/
│   │   ├── domain/
│   │   │   ├── finca.go                 ← Entidad Finca con constructor y métodos
│   │   │   ├── lote.go                  ← Entidad Lote con constructor y métodos (futuro: mover a módulo lotes/)
│   │   │   ├── errores.go              ← Errores de dominio
│   │   │   └── repositories.go          ← Interfaces FincaRepositorio, LoteRepositorio
│   │   ├── application/                 ← Casos de uso (se implementarán en spec-fincas-application.md)
│   │   ├── infrastructure/
│   │   │   ├── finca_repositorio.go     ← Implementación PostgreSQL de FincaRepositorio
│   │   │   └── lote_repositorio.go      ← Implementación PostgreSQL de LoteRepositorio
│   │   └── presentation/
│   │       ├── handler/
│   │       │   └── finca_handler.go     ← Handlers HTTP para fincas y lotes
│   │       ├── facade/
│   │       │   └── finca_facade.go      ← Facade que orquesta operaciones del dominio y aplica mappers
│   │       ├── mapper/
│   │       │   └── finca_mapper.go      ← Mapper dominio ↔ DTO presentación
│   │       └── dto/
│   │           └── finca_dto.go         ← DTOs de presentación (request/response)
│   ├── lotes/
│   │   └── ...                          ← En versión futura, mover lote aquí
│   ├── presentation/
│   │   ├── middleware/
│   │   │   ├── jwt_middleware.go        ← Validación JWT, extrae usuarioID
│   │   │   └── authz_middleware.go      ← Verificación de permisos RBAC
│   │   └── router.go                    ← Configuración de rutas Gin
│   ├── registry/
│   │   └── registry.go                  ← Inyección de dependencias
│   └── shared/
│       ├── api_response.go              ← ApiResponse[T] genérico
│       └── errors.go                    ← Manejo de errores HTTP (RFC 9457)
└── docs/
    └── specs/
        ├── spec-arquitectura-fincas.md  ← Especificación de arquitectura (existente)
        └── spec-fincas-domain.md        ← Esta especificación
```

---

## 11. Acceptance Criteria

### 11.1 Criterios de Aceptación Funcionales

- **AC-FIN-001**: Dado un usuario autenticado sin fincas, Cuando accede al listado (`GET /api/v1/fincas`), Entonces recibe un array vacío con enlace para crear.
- **AC-FIN-002**: Dado un usuario autenticado, Cuando crea una finca con datos válidos, Entonces recibe HTTP 201 con la finca creada y enlaces HATEOAS.
- **AC-FIN-003**: Dado un usuario autenticado, Cuando crea una finca con nombre vacío, Entonces recibe HTTP 400 con error de validación.
- **AC-FIN-004**: Dado un usuario autenticado, Cuando consulta una finca que le pertenece, Entonces recibe HTTP 200 con los datos de la finca.
- **AC-FIN-005**: Dado un usuario autenticado, Cuando consulta una finca de otro usuario, Entonces recibe HTTP 403.
- **AC-FIN-006**: Dado un usuario autenticado, Cuando consulta una finca inexistente, Entonces recibe HTTP 404.
- **AC-FIN-007**: Dado un usuario autenticado propietario de una finca, Cuando actualiza sus datos, Entonces recibe HTTP 200 con los datos actualizados.
- **AC-FIN-008**: Dado un usuario autenticado propietario de una finca sin lotes, Cuando elimina la finca, Entonces recibe HTTP 200 y la finca desaparece del listado.
- **AC-FIN-009**: Dado un usuario autenticado propietario de una finca con lotes, Cuando elimina la finca sin `confirm=true`, Entonces recibe HTTP 409 con advertencia y conteo de lotes.
- **AC-FIN-010**: Dado un usuario autenticado propietario de una finca con lotes, Cuando elimina la finca con `confirm=true`, Entonces recibe HTTP 200 y tanto la finca como sus lotes se eliminan.
- **AC-FIN-011**: Dado un usuario autenticado propietario de una finca, Cuando crea un lote con datos válidos, Entonces recibe HTTP 201 con el lote creado.
- **AC-FIN-012**: Dado un usuario autenticado propietario de una finca, Cuando crea un lote con área negativa, Entonces recibe HTTP 400 con error de validación.
- **AC-FIN-013**: Dado un usuario autenticado, Cuando lista lotes de una finca que le pertenece, Entonces recibe HTTP 200 con el array de lotes.
- **AC-FIN-014**: Dado un request sin token JWT, Cuando accede a cualquier endpoint protegido, Entonces recibe HTTP 401.
- **AC-FIN-015**: Dado un request con token JWT inválido o expirado, Cuando accede a cualquier endpoint protegido, Entonces recibe HTTP 401.

### 11.2 Criterios de Aceptación Técnicos

- **AC-TEC-001**: El seed de permisos en identidad incluye los 8 permisos del módulo `fincas` y la matriz de roles está actualizada.
- **AC-TEC-002**: El middleware JWT extrae correctamente `usuarioID` del claim `sub` y lo inyecta en el contexto Gin.
- **AC-TEC-003**: Las respuestas siguen el formato `ApiResponse[T]` con `data` y `_links`.
- **AC-TEC-004**: Los errores siguen el formato RFC 9457 (Problem Details).
- **AC-TEC-005**: La eliminación de una finca con `ON DELETE CASCADE` elimina sus lotes en BD.
- **AC-TEC-006**: La entidad Finca no permite crear una finca con nombre vacío.
- **AC-TEC-007**: La entidad Lote no permite crear un lote con área <= 0.

---

## 12. Test Strategy

### 12.1 Niveles de Prueba

| Nivel | Qué se prueba | Framework | Dependencias Externas |
|-------|--------------|-----------|----------------------|
| **Unitarias (Dominio)** | Entidades Finca y Lote, validaciones, errores | `testing` + `testify` | Ninguna |
| **Unitarias (Aplicación)** | Casos de uso con repositorios mockeados | `testing` + `testify/mock` | Ninguna (mocks) |
| **Unitarias (Presentación)** | Handlers con facades mockeados | `testing` + `testify/mock` + `httptest` | Ninguna (mocks) |
| **Integración (Persistencia)** | Repositorios PostgreSQL real | `testing` + `testify` + `gorm` | PostgreSQL (test DB) |
| **Integración (API)** | End-to-end: HTTP request → handler → repositorio real | `testing` + `httptest` + `testify` | PostgreSQL (test DB) |

### 12.2 Cobertura Mínima

- **Dominio**: 95% (reglas de negocio críticas)
- **Aplicación**: 90% (orquestación de casos de uso)
- **Infraestructura**: 80% (repositorios contra BD de prueba)
- **Presentación (Backend)**: 85% (handlers, mappers, facades)

### 12.3 Escenarios Clave de Prueba

**Dominio Finca:**
- Crear finca con datos válidos → entidad creada con todos los campos
- Crear finca con nombre vacío → error `ErrNombreFincaRequerido`
- Crear finca con nombre de 250 caracteres → error `ErrNombreFincaLargo`
- Crear finca con ubicación vacía → error `ErrUbicacionRequerida`
- Verificar propiedad con `EsPropietario()` → true/false según corresponda
- Actualizar finca con datos válidos → campos actualizados, `updatedAt` modificado

**Dominio Lote:**
- Crear lote válido asociado a finca → lote creado con todos los campos
- Crear lote con área 0 → error `ErrAreaRequerida`
- Crear lote con área negativa → error `ErrAreaRequerida`
- Crear lote con nombre vacío → error `ErrNombreLoteRequerido`
- Crear lote con nombre de 250 caracteres → error `ErrNombreLoteLargo`

**Aplicación (Servicios):**
- CrearFinca con datos válidos → repositorio.Crear llamado con entidad correcta
- CrearFinca con datos inválidos → error, repositorio.Crear NO llamado
- EliminarFinca con lotes sin confirmación → error `ErrFincaConLotes`, no se elimina
- EliminarFinca con lotes con confirmación → elimina finca y lotes
- ObtenerFinca de otro usuario → error `ErrNoPropietario`
- ListarFincas → solo fincas del usuario autenticado

**Middleware JWT:**
- Token válido → extrae usuarioID correctamente, llama a `c.Next()`
- Token ausente → HTTP 401, no llama a `c.Next()`
- Token inválido (mala firma) → HTTP 401
- Token expirado → HTTP 401

**API (Handlers):**
- POST `/api/v1/fincas` con datos válidos → 201 + `ApiResponse` con HATEOAS
- POST `/api/v1/fincas` con datos inválidos → 400 + Problem Details
- GET `/api/v1/fincas/{id}` sin autenticación → 401
- GET `/api/v1/fincas/{id}` de otro usuario → 403
- DELETE `/api/v1/fincas/{id}` con lotes sin confirmación → 409 + advertencia

---

## 13. Dependencies & External Integrations

### 13.1 External Systems

| ID | Sistema | Propósito | Tipo de Integración |
|----|---------|-----------|---------------------|
| **EXT-001** | Servicio Identidad | Autenticación JWT y RBAC | Secreto JWT compartido (`JWT_SECRET`). Validación local sin llamada HTTP. |
| **EXT-002** | Base de Datos PostgreSQL | Persistencia de Fincas y Lotes | Conexión directa vía GORM. |

### 13.2 Technology Platform Dependencies

| ID | Dependencia | Propósito | Restricción |
|----|-------------|-----------|-------------|
| **PLT-001** | Go 1.22+ | Lenguaje de programación | Obligatorio |
| **PLT-002** | Gin | Router HTTP | Obligatorio (estándar del proyecto) |
| **PLT-003** | GORM + driver PostgreSQL | ORM y persistencia | Obligatorio |
| **PLT-004** | golang-jwt/jwt/v5 | Validación de tokens JWT | Obligatorio |
| **PLT-005** | google/uuid | Generación de UUID v4 | Obligatorio |
| **PLT-006** | testify | Asserts y mocks para pruebas | Obligatorio |

### 13.3 Data Dependencies

| ID | Dependencia | Formato | Acceso |
|----|-------------|---------|--------|
| **DAT-001** | JWT_SECRET compartido | String (simétrico, HS256) | Variable de entorno. Debe coincidir con el secreto del servicio identidad. |
| **DAT-002** | usuarioID del JWT | UUID en claim `sub` | Extraído del token en cada request. |

### 13.4 Infrastructure Dependencies

| ID | Componente | Requisito |
|----|------------|-----------|
| **INF-001** | PostgreSQL 16+ | Base de datos transaccional |
| **INF-002** | Docker + docker-compose | Entorno de desarrollo local (opcional) |

### 13.5 Compliance Dependencies

| ID | Requisito | Impacto |
|----|-----------|---------|
| **COM-001** | Los datos de fincas/lotes pertenecen al usuario propietario. Ningún otro usuario (excepto admin/sys_admin) puede acceder a ellos. | Validación de propiedad en cada operación. |
| **COM-002** | Las contraseñas o secretos nunca se almacenan en este dominio. Toda autenticación delega en identidad. | Sin tabla de credenciales en la BD. |

### 13.6 Nota sobre Dependencias Futuras

- **Integración con image-service**: Cuando se implemente la funcionalidad de cargar imágenes de fincas/lotes, el endpoint expondrá un mecanismo para asociar imágenes a una finca o lote.
- **Integración con muestras/diagnósticos**: Cuando exista el módulo de muestras, la eliminación de un lote deberá verificar si tiene muestras asociadas y requerir confirmación.
- **Multi-tenant**: Cuando el sistema migre a multi-tenant completo, `tenant_id` se volverá obligatorio y la autorización verificará membresía al tenant.

---

## 14. Examples & Edge Cases

### 14.1 Creación de Finca — Success Path

```
POST /api/v1/fincas
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "nombre": "Finca La Pradera",
  "ubicacion": "Km 5 vía a San Jerónimo, Antioquia",
  "descripcion": ""
}

→ 201 Created

{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "nombre": "Finca La Pradera",
    "ubicacion": "Km 5 vía a San Jerónimo, Antioquia",
    "descripcion": "",
    "usuarioID": "user-1234",
    "createdAt": "2026-05-23T14:30:00Z",
    "updatedAt": "2026-05-23T14:30:00Z"
  },
  "_links": {
    "self":   { "href": "/api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890", "method": "GET" },
    "list":   { "href": "/api/v1/fincas", "method": "GET" },
    "update": { "href": "/api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890", "method": "PUT" },
    "delete": { "href": "/api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890", "method": "DELETE" }
  }
}
```

### 14.2 Validación de Nombre Vacío — Error Path

```
POST /api/v1/fincas
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "nombre": "",
  "ubicacion": "Vereda La Esperanza"
}

→ 400 Bad Request

{
  "type": "about:blank",
  "title": "Bad Request",
  "status": 400,
  "detail": "El nombre de la finca es requerido y debe tener entre 3 y 200 caracteres",
  "instance": "/api/v1/fincas"
}
```

### 14.3 Eliminación de Finca con Lotes — Conflict Path

```
DELETE /api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890
Authorization: Bearer <jwt>

→ 409 Conflict

{
  "type": "about:blank",
  "title": "Conflict",
  "status": 409,
  "detail": "La finca tiene 2 lote(s) asociado(s). Confirma la eliminación con ?confirm=true",
  "instance": "/api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### 14.4 Eliminación de Finca con Lotes — Confirmado

```
DELETE /api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890?confirm=true
Authorization: Bearer <jwt>

→ 200 OK

{
  "data": null,
  "_links": {
    "list":   { "href": "/api/v1/fincas", "method": "GET" },
    "create": { "href": "/api/v1/fincas", "method": "POST" }
  }
}
```

### 14.5 Acceso a Finca de Otro Usuario — Forbidden Path

```
GET /api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890
Authorization: Bearer <jwt-de-otro-usuario>

→ 403 Forbidden

{
  "type": "about:blank",
  "title": "Forbidden",
  "status": 403,
  "detail": "No tienes permisos sobre este recurso",
  "instance": "/api/v1/fincas/a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### 14.6 Edge Cases

| Escenario | Comportamiento Esperado |
|-----------|------------------------|
| **Nombre de finca exactamente 3 caracteres** | Válido. Se acepta. |
| **Nombre de finca exactamente 200 caracteres** | Válido. Se acepta. |
| **Nombre de finca 201 caracteres** | Inválido. Error `ErrNombreFincaLargo`. HTTP 400. |
| **Área de lote con más de 2 decimales (ej: 5.256)** | Se redondea o trunca a 2 decimales. Decisión de implementación: truncar. |
| **Área de lote muy grande (ej: 99999999.99)** | Válido si cabe en DECIMAL(10,2) — máximo 99999999.99 hectáreas. |
| **Token JWT sin claim `sub`** | Error de autenticación. HTTP 401. |
| **Token JWT con `sub` vacío** | Error de autenticación. HTTP 401. |
| **Descripción con exactamente 1000 caracteres** | Válido. Se acepta. |
| **Descripción con 1001 caracteres** | Inválido. Error de validación. HTTP 400. |
| **Listar fincas sin crear ninguna** | Array vacío. HTTP 200. |
| **Eliminar finca ya eliminada** | Error `ErrFincaNoEncontrada`. HTTP 404 (o 410 Gone si se implementa soft-delete). |
| **Crear lote en finca que no existe** | Error `ErrFincaNoEncontrada`. HTTP 404. |
| **Crear lote en finca de otro usuario** | Error `ErrNoPropietario`. HTTP 403. |

---

## 15. Validation Criteria

- [ ] ¿Los 8 permisos del módulo `fincas` están definidos en el catálogo de permisos de identidad?
- [ ] ¿La matriz de permisos por rol incluye los nuevos permisos de `fincas`?
- [ ] ¿La entidad Finca valida nombre (3-200 chars), ubicación (1-500 chars) y descripción (0-1000 chars)?
- [ ] ¿La entidad Lote valida nombre (3-200 chars), área (>0) y descripción (0-1000 chars)?
- [ ] ¿Los casos de uso verifican propiedad del recurso antes de modificar/eliminar?
- [ ] ¿La eliminación de finca con lotes requiere confirmación explícita?
- [ ] ¿El middleware JWT extrae `usuarioID` del claim `sub` y lo inyecta en el contexto Gin?
- [ ] ¿El middleware JWT retorna HTTP 401 para tokens ausentes, inválidos o expirados?
- [ ] ¿Todas las respuestas exitosas usan el formato `ApiResponse[T]`?
- [ ] ¿Todos los errores siguen el formato RFC 9457 (Problem Details)?
- [ ] ¿Los enlaces HATEOAS se incluyen en las respuestas exitosas?
- [ ] ¿Las tablas `fincas` y `lotes` existen con el esquema correcto?
- [ ] ¿El índice `idx_fincas_usuario_id` existe para filtrar por usuario?
- [ ] ¿El índice `idx_lotes_finca_id` existe para listar lotes por finca?
- [ ] ¿Las pruebas unitarias del dominio pasan sin dependencias externas?
- [ ] ¿Las pruebas de integración con BD pasan en entorno de prueba?
- [ ] ¿La cobertura mínima del dominio es >= 95%?
- [ ] ¿La cobertura mínima de aplicación es >= 90%?

---

## 16. Related Specifications / Further Reading

- `spec-arquitectura-fincas.md` — Especificación general de arquitectura del microservicio Fincas (existente).
- `../../identidad/docs/specs/autorizacion/spec-rbac-authorization.md` — Especificación del módulo RBAC con permisos atómicos, roles y matriz de permisos (referencia para integración).
- `../../identidad/docs/adr/architecture-context.md` — Contexto arquitectónico del proyecto y patrón de capas (Handler → Facade → Mapper → Domain).
- `../../identidad/docs/specs/presentacion/spec-presentation-layer.md` — Especificación de la capa de presentación con Gin del servicio identidad.
- Casos de uso originales: CU-RF4 (Gestión de Fincas), CU-RF5 (Gestión de Lotes) — Documento de requisitos del sistema CafeScan.
- RFC 9457 — Problem Details for HTTP APIs (https://www.rfc-editor.org/rfc/rfc9457).
