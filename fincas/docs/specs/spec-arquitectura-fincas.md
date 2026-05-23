---
title: Especificación General de Arquitectura — Microservicio Fincas (CU-RF4 / CU-RF5)
version: 1.0
date_created: 2026-05-20
owner: Equipo Catastro
tags: arquitectura, fincas, lotes, clean-architecture, microservicio
---

# Especificación General de Arquitectura — Microservicio Fincas

> **Nota**: Este documento define la arquitectura general del microservicio `fincas`. A partir de esta especificación general surgirán sub-especificaciones detalladas por capa y por caso de uso. Este documento no contiene código fuente ni estructuras de implementación concretas.

---

## 1. Propósito y Alcance

Definir la arquitectura general, los límites del dominio, el aislamiento por capas y las guías estructurales del microservicio `fincas`, el cual gestiona los casos de uso **CU-RF4 (Gestión de Fincas)** y **CU-RF5 (Gestión de Lotes)** del sistema CafeScan.

**Alcance:**
- Modelado del dominio: `Finca` y `Lote` como entidades raíz.
- Casos de uso completos: Crear, Editar, Listar y Eliminar para ambos recursos.
- Arquitectura limpia por capas: Dominio, Aplicación, Infraestructura, Presentación.
- Estrategia de pruebas: unitarias, integración y end-to-end.
- Aislamiento total del microservicio: opera como monolito modular sin comunicación con otros servicios hasta que se defina la integración.
- Frontend en React + Vite dentro del mismo módulo.
- Persistencia en PostgreSQL.

**Fuera de alcance (para sub-especificaciones posteriores):**
- Implementación concreta de handlers HTTP, routers, middleware.
- Implementación concreta de repositorios GORM.
- Contratos detallados de API REST.
- Estrategia de despliegue y contenedores.
- Comunicación con otros microservicios (identidad, image-service, etc.).
- UI/UX del frontend (componentes, páginas, estados).

---

## 2. Definiciones

| Término | Definición |
|---------|-----------|
| **Finca** | Unidad principal de gestión agrícola. Pertenece a un Agricultor. Contiene nombre, ubicación y descripción. |
| **Lote** | Subdivisión espacial de una Finca. Pertenece a una Finca. Contiene nombre, área (hectáreas) y descripción. |
| **Agricultor** | Actor con rol de negocio que posee fincas. Se identifica por su `usuarioID` del sistema de identidad. |
| **Administrador** | Actor con permisos globales sobre todas las fincas del sistema. |
| **Clean Architecture** | Patrón de arquitectura que separa el software en capas con dependencias hacia adentro: Dominio → Aplicación → Infraestructura → Presentación. |
| **Monolito Modular** | Arquitectura donde el código de un microservicio está aislado en su propio módulo, pero no se comunica por red con otros servicios — toda la comunicación es interna al proceso. |
| **Agregado** | Grupo de entidades que se tratan como una unidad coherente de persistencia y consistencia transaccional. |
| **Value Object** | Objeto inmutable sin identidad propia, definido por sus atributos. |
| **Domain Event** | Evento que registra un hecho significativo ocurrido en el dominio. |
| **Unit of Work** | Patrón que agrupa múltiples operaciones de persistencia en una sola transacción atómica. |

---

## 3. Requisitos, Restricciones y Guías

### 3.1 Requisitos Funcionales (derivados de casos de uso)

- **RF-FIN-001**: El sistema debe permitir al Agricultor registrar una nueva Finca con nombre, ubicación y descripción.
- **RF-FIN-002**: El sistema debe permitir al Agricultor editar los datos de una Finca existente.
- **RF-FIN-003**: El sistema debe permitir al Agricultor listar todas sus Fincas registradas.
- **RF-FIN-004**: El sistema debe permitir al Agricultor eliminar una Finca existente, advirtiendo si tiene Lotes asociados (eliminación en cascada previa confirmación).
- **RF-FIN-005**: El sistema debe permitir al Agricultor, con una Finca seleccionada, registrar un nuevo Lote con nombre, área y descripción.
- **RF-FIN-006**: El sistema debe permitir al Agricultor editar los datos de un Lote existente.
- **RF-FIN-007**: El sistema debe permitir al Agricultor listar todos los Lotes de una Finca.
- **RF-FIN-008**: El sistema debe permitir al Agricultor eliminar un Lote, advirtiendo si tiene muestras o diagnósticos asociados (eliminación en cascada previa confirmación).
- **RF-FIN-009**: El sistema debe validar que los datos ingresados sean completos y válidos antes de persistir.
- **RF-FIN-010**: El sistema debe persistir todas las operaciones en PostgreSQL.
- **RF-FIN-011**: El sistema debe mostrar mensajes de error descriptivos cuando una operación falle.
- **RF-FIN-012**: El sistema debe mostrar una pantalla vacía con opción de crear cuando no existan Fincas ni Lotes.

### 3.2 Requisitos de Arquitectura

- **ARQ-FIN-001**: El backend se implementa en Go (golang).
- **ARQ-FIN-002**: La estructura de carpetas sigue Clean Architecture con las capas: `domain/`, `application/`, `infrastructure/`, `presentation/`.
- **ARQ-FIN-003**: Las dependencias fluyen hacia adentro: Presentación → Aplicación → Dominio. Infraestructura implementa interfaces definidas en Dominio.
- **ARQ-FIN-004**: El frontend se implementa en React con Vite como bundler, dentro del mismo módulo.
- **ARQ-FIN-005**: No existe comunicación con otros microservicios en esta fase. El módulo opera como monolito modular autónomo.
- **ARQ-FIN-006**: La autenticación y autorización se validan localmente mediante JWT (el token llega desde el frontend, se valida con la misma clave que el servicio identidad).
- **ARQ-FIN-007**: El microservicio expone API REST sobre HTTP.

### 3.3 Restricciones Técnicas

- **CON-FIN-001**: La base de datos es PostgreSQL, manejada vía GORM (migraciones automáticas o manuales).
- **CON-FIN-002**: Los IDs de Finca y Lote son UUID v4.
- **CON-FIN-003**: Toda fecha se maneja como `time.Time` en UTC.
- **CON-FIN-004**: Las transacciones que modifican múltiples agregados usan Unit of Work.
- **CON-FIN-005**: Las contraseñas o secretos nunca se almacenan en este dominio (eso pertenece al servicio identidad).
- **CON-FIN-006**: El campo `area` del Lote usa precisión decimal fija (hectáreas con 2 decimales).
- **CON-FIN-007**: No se permite código cíclico entre capas. Una capa de dominio no puede importar nada de infraestructura ni presentación.

### 3.4 Guías de Diseño

- **GUD-FIN-001**: El dominio modela entidades ricas con comportamiento, no simples "anémicos" con getters/setters.
- **GUD-FIN-002**: Los errores de dominio se definen como variables exportadas usando `errors.New()` en el paquete `domain/errores.go`.
- **GUD-FIN-003**: Los repositorios se definen como interfaces en el dominio, se implementan en infraestructura.
- **GUD-FIN-004**: Los casos de uso se implementan como servicios de aplicación, cada uno en su propio archivo.
- **GUD-FIN-005**: La capa de presentación sigue el flujo `Handler → Facade → Mapper → Domain` (según el patrón establecido en `identidad`).
- **GUD-FIN-006**: Las respuestas HTTP usan un formato genérico `ApiResponse[T]` con HATEOAS. Los errores siguen RFC 9457 (Problem Details).
- **GUD-FIN-007**: Un `Lote` nunca existe sin una `Finca`. La creación de un Lote requiere una Finca existente y válida.
- **GUD-FIN-008**: La eliminación de una Finca con Lotes requiere confirmación explícita del usuario (eliminación en cascada).
- **GUD-FIN-009**: Los mensajes de error al usuario son descriptivos pero no revelan detalles internos de implementación.
- **GUD-FIN-010**: Los mappers (dominio ↔ DTO) no contienen lógica de negocio, solo conversión de estructura.

---

## 4. Arquitectura General

### 4.1 Vista General del Microservicio

El microservicio `fincas` sigue el patrón de **Clean Architecture** (Arquitectura Limpia) con 4 capas:

```
┌──────────────────────────────────────────────────────────┐
│                    PRESENTACIÓN                          │
│  (Handlers HTTP, DTOs, Mappers, Facades, Router, Gin)   │
│                    ↓ (depende de)                        │
├──────────────────────────────────────────────────────────┤
│                    APLICACIÓN                            │
│  (Casos de Uso / Servicios de Aplicación, DTOs internos) │
│                    ↓ (depende de)                        │
├──────────────────────────────────────────────────────────┤
│                    DOMINIO                               │
│  (Entidades, Value Objects, Interfaces de Repositorio,   │
│   Errores de Dominio, Eventos de Dominio)                │
│                    ↑ (implementa)                        │
├──────────────────────────────────────────────────────────┤
│                    INFRAESTRUCTURA                       │
│  (Repositorios PostgreSQL/GORM, DB Connection, Migraciones)│
└──────────────────────────────────────────────────────────┘
```

**Reglas de dependencia**:
- Dominio: No conoce ninguna otra capa.
- Aplicación: Conoce solo al Dominio.
- Infraestructura: Conoce al Dominio (implementa sus interfaces).
- Presentación: Conoce a Aplicación y Dominio (vía DTOs y Mappers, nunca directamente).

### 4.2 Límites del Dominio

```
┌───────────────────────────────────────────┐
│           BOUNDED CONTEXT:                 │
│         GESTIÓN DE FINCAS                  │
│                                           │
│  ┌─────────────┐     ┌─────────────┐     │
│  │   Finca     │─────│   Lote      │     │
│  │             │1   N│             │     │
│  │ - id        │     │ - id        │     │
│  │ - nombre    │     │ - fincaID   │     │
│  │ - ubicacion │     │ - nombre    │     │
│  │ - descripcion│    │ - area      │     │
│  │ - usuarioID │     │ - descripcion│    │
│  │ - createdAt │     │ - createdAt │     │
│  │ - updatedAt │     │ - updatedAt │     │
│  └─────────────┘     └─────────────┘     │
│                                           │
│  Agregados:                               │
│  - Finca es el agregado raíz              │
│  - Lote pertenece al agregado Finca       │
└───────────────────────────────────────────┘
```

- `Finca` es el agregado raíz. Un `Lote` no puede existir sin una `Finca`.
- `usuarioID` en Finca es una referencia al usuario del sistema de identidad, no una entidad embebida.
- No existen relaciones directas con entidades de otros dominios (muestras, diagnósticos) — esas referencias se manejarán cuando se defina la integración.

### 4.3 Capa de Dominio

Contiene el modelo de negocio puro. No tiene dependencias externas (ni frameworks, ni bases de datos, ni librerías HTTP).

**Componentes:**
- **Entidades**: `Finca`, `Lote` — con comportamiento y validaciones de negocio.
- **Value Objects**: `Ubicacion`, `Area`, `NombreFinca` (si aplica).
- **Interfaces de Repositorio**: `FincaRepositorio`, `LoteRepositorio` — definen contratos de persistencia.
- **Errores de Dominio**: Errores específicos del dominio (`ErrFincaNoEncontrada`, `ErrLoteNoEncontrado`, `ErrDatosInvalidos`, etc.).
- **Eventos de Dominio**: (opcional) `FincaCreada`, `FincaEliminada`, `LoteCreado`, `LoteEliminado`.

### 4.4 Capa de Aplicación

Orquesta los casos de uso. Coordina el flujo entre el dominio y la infraestructura.

**Componentes:**
- **Servicios de Aplicación**: Un archivo por caso de uso o grupo relacionado.
  - `CrearFincaService` — Recibe datos de entrada, valida reglas de negocio vía dominio, persiste vía repositorio.
  - `EditarFincaService` — Recibe ID y datos actualizados, valida existencia, actualiza, persiste.
  - `ListarFincasService` — Recibe criterios de filtro (usuarioID), retorna lista.
  - `EliminarFincaService` — Recibe ID, valida existencia, verifica lotes asociados, elimina en cascada.
  - `CrearLoteService`, `EditarLoteService`, `ListarLotesService`, `EliminarLoteService` — análogos para Lote.
- **DTOs de Aplicación**: Estructuras de entrada/salida para cada servicio.
- **Interfaces de servicios externos**: (futuro) interfaces para integración con identidad u otros servicios.

### 4.5 Capa de Infraestructura

Implementa las interfaces definidas en el dominio. Contiene todo lo relacionado con tecnología externa.

**Componentes:**
- **Persistencia**: Repositorios PostgreSQL usando GORM.
  - `FincaRepositorioPostgres` — implementa `FincaRepositorio`.
  - `LoteRepositorioPostgres` — implementa `LoteRepositorio`.
- **Migraciones**: Esquemas de base de datos y migraciones automáticas o versionadas.
- **Conexión**: Configuración y gestión de la conexión a PostgreSQL.
- **Transacciones**: Implementación de Unit of Work para operaciones atómicas.

### 4.6 Capa de Presentación (Backend)

Expone la funcionalidad del microservicio vía API REST. No contiene lógica de negocio.

**Componentes:**
- **Handlers HTTP**: Funciones que reciben requests, delegan en Facades, retornan respuestas.
- **Facades**: Agrupan casos de uso relacionados y aplican Mappers para convertir dominio → DTOs de presentación.
- **Mappers**: Convierten entre entidades de dominio y DTOs de presentación.
- **DTOs de Presentación**: Estructuras planas serializables a JSON.
- **Router**: Definición de rutas (Gin).
- **Middleware**: Autenticación JWT, logging, CORS, recovery.
- **ApiResponse**: Estructura genérica para respuestas estandarizadas.
- **Manejo de errores**: Traducción de errores de dominio/ aplicacion a códigos y mensajes HTTP.

### 4.7 Capa de Presentación (Frontend)

Aplicación React con Vite para interactuar con la API REST.

**Componentes:**
- **Páginas**: `FincasPage`, `FincaDetallePage`, `LotesPage`, `LoteFormPage`.
- **Componentes**: Formularios de creación/edición, listados, tarjetas de finca, modales de confirmación.
- **Servicios**: Cliente HTTP para consumir la API REST del backend.
- **Estado**: Context API o similar para estado de sesión y datos locales.
- **Enrutamiento**: React Router para navegación entre vistas.

---

## 5. Flujo de Datos por Capa

```
Request HTTP
    │
    ▼
[Router Gin] ────→ [Middleware JWT (extrae usuarioID)]
    │
    ▼
[Handler] ────→ Recibe request, extrae parámetros, llama a Facade
    │
    ▼
[Facade] ────→ Agrupa casos de uso, aplica Mapper (DTO ← Dominio)
    │
    ▼
[Servicio Aplicación] ────→ Orquesta lógica, llama a Repositorio (vía interfaz)
    │
    ▼
[Repositorio (Impl.)] ────→ GORM → PostgreSQL
    │
    ▼ (respuesta)
[Servicio Aplicación] ────→ Retorna entidad de dominio o error
    │
    ▼
[Facade] ────→ Mapper: Dominio → DTO
    │
    ▼
[Handler] ────→ Construye ApiResponse[T] con HATEOAS, retorna JSON
```

---

## 6. Contratos de Interfaz (Alto Nivel)

### 6.1 API REST — Endpoints

El microservicio expone los siguientes endpoints de alto nivel. Los contratos detallados (request/response bodies, códigos HTTP, headers) se definen en una sub-especificación de presentación.

| Método | Ruta | Recurso | Operación |
|--------|------|---------|-----------|
| POST | `/api/v1/fincas` | Finca | Crear nueva finca |
| GET | `/api/v1/fincas` | Finca | Listar fincas del usuario autenticado |
| GET | `/api/v1/fincas/{id}` | Finca | Obtener detalle de una finca |
| PUT | `/api/v1/fincas/{id}` | Finca | Actualizar finca existente |
| DELETE | `/api/v1/fincas/{id}` | Finca | Eliminar finca (con confirmación si tiene lotes) |
| POST | `/api/v1/fincas/{fincaId}/lotes` | Lote | Crear nuevo lote en una finca |
| GET | `/api/v1/fincas/{fincaId}/lotes` | Lote | Listar lotes de una finca |
| GET | `/api/v1/fincas/{fincaId}/lotes/{id}` | Lote | Obtener detalle de un lote |
| PUT | `/api/v1/fincas/{fincaId}/lotes/{id}` | Lote | Actualizar lote existente |
| DELETE | `/api/v1/fincas/{fincaId}/lotes/{id}` | Lote | Eliminar lote (con confirmación si tiene dependencias) |

### 6.2 Formato de Respuesta

Todas las respuestas exitosas siguen el formato `ApiResponse[T]`:

```
{
  "data": { ... },         // Recurso solicitado (objeto o array)
  "_links": {
    "self": { "href": "...", "method": "GET" },
    "create": { "href": "...", "method": "POST" },
    "update": { "href": "...", "method": "PUT" },
    "delete": { "href": "...", "method": "DELETE" }
  }
}
```

Los errores siguen el formato RFC 9457 (Problem Details):

```
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "Finca con ID xyz no encontrada",
  "instance": "/api/v1/fincas/xyz"
}
```

### 6.3 Autenticación

- El microservicio recibe un JWT en el header `Authorization: Bearer <token>`.
- El middleware JWT valida el token usando el mismo secreto que el servicio identidad.
- Extrae del token: `usuarioID` (subject) y `rol` (claim).
- Todas las operaciones de Finca/Lote filtran por `usuarioID` excepto para el rol `Administrador`.
- Los endpoints de listado retornan solo los recursos del usuario autenticado (o todos si es Administrador).

---

## 7. Manejo de Flujos Alternativos y Errores

### 7.1 Flujos Alternativos (de casos de uso)

| ID | Condición | Comportamiento del Sistema | Capa Responsable |
|----|-----------|---------------------------|------------------|
| FA-2a | Usuario sin Fincas registradas | Retorna lista vacía con indicación de "sin datos". La presentación muestra pantalla vacía con botón "Crear primera finca" | Presentación (Frontend) |
| FA-2a (Lotes) | Finca sin Lotes registrados | Retorna lista vacía. La presentación muestra pantalla vacía con botón "Crear primer lote" | Presentación (Frontend) |
| FA-5a | Datos inválidos o incompletos | Servicio de aplicación retorna error de validación. Handler traduce a HTTP 400 con detalles de campos erróneos | Aplicación + Presentación |
| FA-6a | Error de conexión con BD | Repositorio retorna error. Servicio de aplicación lo propaga. Handler traduce a HTTP 500 con mensaje genérico | Infraestructura → Aplicación → Presentación |
| FA-7a | Eliminar Finca con Lotes asociados | Servicio de aplicación verifica lotes existentes. API requiere confirmación explícita (query param `?confirm=true` o两步). Si no hay confirmación, retorna HTTP 409 con advertencia | Aplicación + Presentación |
| FA-7a (Lotes) | Eliminar Lote con muestras/diagnósticos asociados | Misma lógica: verifica dependencias, requiere confirmación, retorna advertencia si no se confirma | Aplicación + Presentación |

### 7.2 Mapa de Errores

| Condición | Código HTTP | Mensaje |
|-----------|------------|---------|
| Datos inválidos (validación) | 400 (Bad Request) | Campos específicos con error |
| No autenticado | 401 (Unauthorized) | Token ausente o inválido |
| Sin permisos (otro usuario) | 403 (Forbidden) | No tiene permisos para este recurso |
| Recurso no encontrado | 404 (Not Found) | Finca/Lote con ID X no encontrado |
| Conflicto (dependencias) | 409 (Conflict) | La Finca tiene N lotes asociados. Confirme la eliminación |
| Error interno del servidor | 500 (Internal Server Error) | No se pudo completar la operación. Intente más tarde |

---

## 8. Estrategia de Pruebas

### 8.1 Niveles de Prueba

| Nivel | Qué se prueba | Framework | Dependencias Externas |
|-------|--------------|-----------|----------------------|
| **Unitarias (Dominio)** | Entidades, Value Objects, reglas de negocio | `testing` + `testify` | Ninguna |
| **Unitarias (Aplicación)** | Servicios de aplicación con repositorios mockeados | `testing` + `testify/mock` | Ninguna (mocks) |
| **Unitarias (Presentación)** | Handlers con facades mockeados | `testing` + `testify/mock` + `httptest` | Ninguna (mocks) |
| **Integración (Persistencia)** | Repositorios contra PostgreSQL real | `testing` + `testify` + `gorm` | PostgreSQL (test DB) |
| **Integración (API)** | End-to-end: HTTP request → handler → facade → repositorio real | `testing` + `httptest` + `testify` | PostgreSQL (test DB) |
| **Frontend (Componentes)** | Componentes React con datos mockeados | Vitest + Testing Library | Ninguna (mocks) |
| **Frontend (Integración)** | Flujo completo de página | Vitest + MSW (Mock Service Worker) | API mockeada |

### 8.2 Cobertura Mínima

- **Dominio**: 95% (reglas de negocio críticas)
- **Aplicación**: 90% (orquestación de casos de uso)
- **Infraestructura**: 80% (repositorios contra BD de prueba)
- **Presentación (Backend)**: 85% (handlers, mappers, facades)
- **Frontend**: 80% (componentes, páginas, servicios)

### 8.3 Organización de Pruebas

Las pruebas se organizan junto al código que prueban, siguiendo la convención de Go (`_test.go`):

- `internal/fincas/domain/finca_test.go` — Pruebas de dominio de Finca
- `internal/lotes/domain/lote_test.go` — Pruebas de dominio de Lote
- `internal/fincas/application/crear_finca_test.go` — Pruebas de caso de uso
- `internal/fincas/infrastructure/persistence/finca_repositorio_test.go` — Pruebas de integración con BD
- `internal/fincas/presentation/handler/finca_handler_test.go` — Pruebas de handler
- `web/src/__tests__/` — Pruebas de frontend

### 8.4 Escenarios Clave de Prueba

**Dominio Finca:**
- Crear finca con datos válidos → entidad creada con todos los campos
- Crear finca con nombre vacío → error de dominio
- Crear finca con nombre demasiado largo → error de dominio
- Eliminar finca sin lotes → éxito
- Eliminar finca con lotes → error o advertencia (según diseño)
- Actualizar campos de finca → cambios reflejados correctamente
- Finca pertenece a usuario específico → validación de propiedad

**Dominio Lote:**
- Crear lote válido asociado a finca → lote creado
- Crear lote con área negativa → error de dominio
- Crear lote con área cero → error de dominio
- Crear lote sin finca existente → error
- Listar lotes de finca → solo lotes de esa finca

**Aplicación (Servicios):**
- CrearFinca: datos válidos → repositorio.Create llamado con entidad correcta
- CrearFinca: datos inválidos → error de validación, repositorio.Create NO llamado
- EliminarFinca: finca con lotes sin confirmación → error, no se elimina
- EliminarFinca: finca con lotes con confirmación → elimina finca y lotes
- ListarFincas: usuario autenticado → solo sus fincas
- ListarFincas: administrador → todas las fincas

**Presentación:**
- POST /api/v1/fincas con datos válidos → 201 + ApiResponse
- POST /api/v1/fincas con datos inválidos → 400 + Problem Details
- GET /api/v1/fincas/{id} sin autenticación → 401
- DELETE /api/v1/fincas/{id} con lotes sin confirmación → 409 + advertencia
- DELETE /api/v1/fincas/{id} con lotes con confirmación → 200 + eliminación

**Frontend:**
- Renderizar lista vacía → mostrar mensaje "No hay fincas" + botón crear
- Renderizar lista con datos → mostrar tarjetas de finca
- Formulario de creación → validar campos requeridos
- Confirmación de eliminación → modal de confirmación
- Error de servidor → mostrar mensaje de error amigable

---

## 9. Estrategia de Persistencia

### 9.1 Esquema General de Base de Datos

La base de datos PostgreSQL contendrá dos tablas principales:

- `fincas`: Almacena las fincas registradas.
  - Columnas: id (UUID PK), usuario_id (UUID, FK referencial a usuarios), nombre, ubicacion, descripcion, created_at, updated_at, deleted_at (soft delete opcional).
  
- `lotes`: Almacena los lotes asociados a una finca.
  - Columnas: id (UUID PK), finca_id (UUID, FK a fincas), nombre, area (decimal), descripcion, created_at, updated_at, deleted_at (soft delete opcional).

- `usuario_id` en `fincas` es una referencia lógica al usuario del sistema de identidad. No existe una FK física hacia la tabla `usuarios` de otro servicio (por aislamiento de microservicio). La validación de existencia del usuario se delega al frontend y al token JWT.

### 9.2 Migraciones

Las migraciones de base de datos se manejan de una de las siguientes formas (decisión en sub-especificación):
- Migraciones automáticas de GORM (`AutoMigrate`) para desarrollo.
- Migraciones versionadas (archivos SQL) para producción.
- Ambas: automáticas en desarrollo, versionadas en producción.

---

## 10. Organización del Módulo

### 10.1 Estructura Raíz

```
fincas/                          ← CARPETA RAIZ DEL MÓDULO
├── docs/                        ← Documentación y especificaciones
│   ├── specs/                   ← Sub-especificaciones detalladas
│   │   ├── spec-arquitectura-fincas.md   ← Este documento
│   │   ├── spec-domain.md               ← Sub-especificación de dominio (futuro)
│   │   ├── spec-application.md           ← Sub-especificación de aplicación (futuro)
│   │   ├── spec-infrastructure.md        ← Sub-especificación de infraestructura (futuro)
│   │   ├── spec-presentation-backend.md  ← Sub-especificación de presentación backend (futuro)
│   │   └── spec-presentation-frontend.md ← Sub-especificación de frontend (futuro)
│   └── adr/                     ← Architecture Decision Records (futuro)
├── cmd/                         ← Puntos de entrada de la aplicación
├── internal/                    ← Código interno (no exportable como librería)
│   ├── config/                  ← Configuración (variables de entorno)
│   ├── fincas/                  ← Módulo Finca
│   │   ├── domain/              ← Entidades, VO, interfaces, errores
│   │   ├── application/         ← Casos de uso / servicios
│   │   ├── infrastructure/      ← Repositorios PostgreSQL
│   │   └── presentation/        ← Handlers, DTOs, Mappers, Facades
│   ├── lotes/                   ← Módulo Lote (misma estructura que fincas)
│   └── shared/                  ← Código compartido entre módulos
├── web/                         ← Frontend React + Vite
└── [archivos de configuración]  ← go.mod, Makefile, Dockerfile, .env, etc.
```

> **Nota**: Esta estructura completa se presenta como referencia. Las sub-especificaciones detallarán la estructura exacta de cada subdirectorio. Por ahora solo se define la carpeta raíz (`fincas/`) y el subdirectorio de documentación (`docs/`).

---

## 11. Justificación y Contexto

### 11.1 ¿Por qué un microservicio separado para Fincas?

- **Límite del dominio claro**: Fincas y Lotes son conceptos del dominio catastral/agrícola, distintos del dominio de identidad (usuarios, sesiones) y del dominio de procesamiento de imágenes.
- **Aislamiento futuro**: Cuando se requiera escalar, este dominio puede independizarse sin afectar a otros servicios.
- **Responsabilidad única**: El servicio maneja exclusivamente la gestión de unidades productivas (fincas y lotes), sin mezclar lógica de autenticación, roles o procesamiento de imágenes.

### 11.2 ¿Por qué monolito modular inicialmente?

- **Simplicidad operativa**: Un solo binario, un solo despliegue. Facilita el desarrollo inicial, las pruebas y el debugging.
- **Sin overhead de red**: Las llamadas entre módulos son locales (funciones Go), no requieren serialización/deserialización ni latencia de red.
- **Madurez del dominio**: Permite validar la lógica de negocio y los casos de uso antes de añadir la complejidad de comunicación entre servicios.
- **Migración gradual**: Cuando se requiera, la separación física es natural porque las capas ya están aisladas por paquetes.

### 11.3 ¿Por qué la misma estructura de identidad?

- **Consistencia del proyecto**: Mantener el mismo patrón de Clean Architecture en todos los microservicios facilita la comprensión, el onboarding y el mantenimiento.
- **Reuso de patrones**: Mappers, Facades, ApiResponse, Unit of Work son patrones ya probados en el servicio identidad.
- **Herramientas compartidas**: GORM, Gin, Huma, JWT, testify son dependencias ya validadas en el proyecto.

### 11.4 Relación con otros servicios

En esta fase (monolito modular), `fincas` no se comunica con otros servicios. Las referencias externas son:
- `usuarioID`: Se obtiene del token JWT que llega en cada request. No se consulta al servicio identidad para validar existencia.
- `Finca` como referencia para Muestras/Diagnósticos: Se definirá en una fase posterior cuando se construyan esos dominios.

---

## 12. Dependencias Tecnológicas

### 12.1 Backend (Go)

| Dependencia | Propósito | Tipo |
|-------------|-----------|------|
| Go 1.26+ | Lenguaje de programación | Platforma |
| Gin | Router HTTP | Framework |
| Huma v2 + humagin | Documentación OpenAPI automática | Framework |
| GORM + driver PostgreSQL | ORM y persistencia | Librería |
| golang-jwt/jwt/v5 | Validación de tokens JWT | Librería |
| google/uuid | Generación de UUID v4 | Librería |
| testify | Asserts y mocks para pruebas | Testing |
| golang.org/x/crypto | Utilidades criptográficas (bcrypt si aplica) | Librería |

### 12.2 Frontend (Web)

| Dependencia | Propósito |
|-------------|-----------|
| React 18+ | Librería de UI |
| Vite | Bundler y dev server |
| React Router | Enrutamiento del lado del cliente |
| fetch / axios | Cliente HTTP para API calls |
| Vitest + Testing Library | Pruebas unitarias y de componentes |

### 12.3 Infraestructura

| Componente | Propósito |
|------------|-----------|
| PostgreSQL 18+ | Base de datos relacional |
| Docker + docker-compose | Entorno de desarrollo local |

---

## 13. Criterios de Aceptación

- **AC-FIN-001**: Dado un usuario autenticado sin fincas, Cuando accede al listado, Entonces ve una pantalla vacía con opción de crear su primera finca.
- **AC-FIN-002**: Dado un usuario autenticado, Cuando crea una finca con datos válidos, Entonces la finca se persiste y se muestra en el listado con mensaje de confirmación.
- **AC-FIN-003**: Dado un usuario autenticado, Cuando crea una finca con datos inválidos, Entonces recibe errores de validación en los campos correspondientes.
- **AC-FIN-004**: Dado un usuario autenticado con fincas existentes, Cuando edita una finca, Entonces los cambios se persisten y se reflejan en el listado.
- **AC-FIN-005**: Dado un usuario autenticado, Cuando elimina una finca sin lotes, Entonces la finca se elimina permanentemente.
- **AC-FIN-006**: Dado un usuario autenticado, Cuando elimina una finca con lotes sin confirmar, Entonces recibe una advertencia y la operación no se ejecuta.
- **AC-FIN-007**: Dado un usuario autenticado, Cuando elimina una finca con lotes confirmando la acción, Entonces la finca y sus lotes se eliminan en cascada.
- **AC-FIN-008**: Dado un usuario autenticado con una finca seleccionada, Cuando crea un lote con datos válidos, Entonces el lote se asocia a la finca y se muestra en el detalle.
- **AC-FIN-009**: Dado un usuario autenticado, Cuando lista lotes de una finca sin lotes, Entonces ve pantalla vacía con opción de crear el primer lote.
- **AC-FIN-010**: Dado un error de conexión con la base de datos, Cuando se ejecuta cualquier operación, Entonces el usuario recibe un mensaje de error indicando que no se pudo completar la operación.
- **AC-FIN-011**: Dado un request sin token JWT, Cuando se accede a cualquier endpoint protegido, Entonces el sistema retorna HTTP 401.
- **AC-FIN-012**: Dado un token JWT de otro usuario, Cuando se intenta acceder a una finca de otro agricultor, Entonces el sistema retorna HTTP 403.
- **AC-FIN-013**: Dado un usuario administrador, Cuando accede al módulo de fincas, Entonces puede ver y gestionar todas las fincas del sistema.
- **AC-FIN-014**: Dado que los tests se ejecutan, Cuando se miden las coberturas, Entonces se cumplen los umbrales mínimos definidos en esta especificación.

---

## 14. Criterios de Validación del Microservicio

- [ ] El servidor inicia correctamente y expone los endpoints REST en el puerto configurado.
- [ ] Los endpoints de salud (`/health`) y documentación (`/docs`) están disponibles.
- [ ] La especificación OpenAPI se genera automáticamente y es accesible.
- [ ] El CRUD completo de Fincas funciona contra PostgreSQL.
- [ ] El CRUD completo de Lotes funciona contra PostgreSQL.
- [ ] Las validaciones de datos rechazan entradas inválidas con mensajes descriptivos.
- [ ] La autenticación JWT protege todos los endpoints.
- [ ] Los flujos alternativos (datos inválidos, recursos no encontrados, conflictos, errores de BD) se manejan correctamente.
- [ ] Las pruebas unitarias del dominio pasan sin dependencias externas.
- [ ] Las pruebas de integración con BD pasan en entorno de prueba.
- [ ] Las pruebas de frontend pasan con datos mockeados.
- [ ] La cobertura de pruebas cumple los umbrales definidos.
- [ ] El frontend React se renderiza correctamente y se comunica con la API.

---

## 15. Especificaciones Relacionadas / Lectura Adicional

- `./spec-domain.md` — Sub-especificación detallada de la capa de dominio (Futuro).
- `./spec-application.md` — Sub-especificación detallada de la capa de aplicación (Futuro).
- `./spec-infrastructure.md` — Sub-especificación detallada de la capa de infraestructura (Futuro).
- `./spec-presentation-backend.md` — Sub-especificación detallada de la capa de presentación backend (Futuro).
- `./spec-presentation-frontend.md` — Sub-especificación detallada del frontend React+Vite (Futuro).
- `../../identidad/docs/adr/architecture-context.md` — Patrón de capas de presentación del proyecto (Handler → Facade → Mapper → Domain).
- `../../identidad/docs/specs/presentacion/spec-presentation-layer.md` — Especificación de capa de presentación con Gin + Huma v2 del servicio identidad.
- Casos de uso originales: CU-RF4 (Gestión de Fincas), CU-RF5 (Gestión de Lotes) — Documento de requisitos del sistema CafeScan.
