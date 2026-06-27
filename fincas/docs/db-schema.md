# Base de Datos de Fincas — Esquema Completo

## Stack

- **Motor:** PostgreSQL 16+
- **ORM:** GORM (AutoMigrate, sin migraciones SQL manuales)
- **Base:** `bunna_fincas` (usuario: `fincas_user`)

---

## 1. Tabla `fincas` — Fincas/parcelas

Creada por `FincaModel` → `fincas_model.go`

| Columna | Tipo | Nulable | Default | Descripción |
|---|---|---|---|---|
| `id` | `varchar(36)` | NO | — | UUID primario |
| `nombre` | `text` | SÍ | — | Nombre de la finca |
| `ubicacion` | `text` | SÍ | — | JSON de coordenadas |
| `descripcion` | `text` | SÍ | — | Descripción opcional |
| `usuario_id` | `text` | SÍ | — | ID del usuario creador |
| `tenant_id` | `text` | SÍ | — | **Tenant (nulable)** |
| `estado` | `text` | SÍ | `'ACTIVA'` | `ACTIVA` / `PENDIENTE_ELIMINACION` |
| `created_at` | `timestamptz` | SÍ | — | Fecha de creación |
| `updated_at` | `timestamptz` | SÍ | — | Última actualización |

**Índices:** `id` (PK), `usuario_id`, `tenant_id`

---

## 2. Tabla `lotes` — Lotes de una finca

Creada por `LoteModel` → `lote_model.go`

| Columna | Tipo | Nulable | Default | Descripción |
|---|---|---|---|---|
| `id` | `varchar(36)` | NO | — | UUID primario |
| `finca_id` | `text` | SÍ | — | FK lógica a `fincas.id` |
| `tenant_id` | `text` | SÍ | — | **Tenant** |
| `nombre` | `text` | SÍ | — | Nombre del lote |
| `area` | `numeric` | SÍ | — | Área en m² |
| `descripcion` | `text` | SÍ | — | Descripción opcional |
| `estado` | `text` | SÍ | `'ACTIVO'` | `ACTIVO` / `ELIMINADO` |
| `created_at` | `timestamptz` | SÍ | — | Fecha de creación |
| `updated_at` | `timestamptz` | SÍ | — | Última actualización |

**Índices:** `id` (PK), `finca_id`, `tenant_id`

---

## 3. Tabla `nodos` — Nodos IoT / Cámaras

Creada por `NodoModel` → `nodo_model.go`

| Columna | Tipo | Nulable | Default | Descripción |
|---|---|---|---|---|
| `id` | `varchar(36)` | NO | — | UUID primario |
| `finca_id` | `varchar(36)` | SÍ | — | FK lógica a `fincas.id` |
| `lote_id` | `varchar(36)` | SÍ | — | FK lógica a `lotes.id` (opcional) |
| `tenant_id` | `varchar(36)` | SÍ | — | **Tenant** |
| `nombre` | `varchar(100)` | SÍ | — | Nombre descriptivo |
| `node_key` | `varchar(100)` | SÍ | — | **Clave única del dispositivo** (hardware) |
| `nodo_id` | `varchar(50)` | SÍ | — | Alias o ID secundario |
| `estado` | `varchar(20)` | SÍ | `'ACTIVO'` | `ACTIVO` / `INACTIVO` / `MANTENIMIENTO` |
| `ubicacion` | `text` | SÍ | — | JSON de ubicación |
| `ultimo_heartbeat` | `timestamptz` | SÍ | — | Último latido del dispositivo |
| `firmware_version` | `varchar(20)` | SÍ | — | Versión de firmware |
| `created_at` | `timestamptz` | SÍ | — | Fecha de creación |
| `updated_at` | `timestamptz` | SÍ | — | Última actualización |
| `creado_en` | `timestamptz` | SÍ | — | (duplicado de created_at) |
| `actualizado_en` | `timestamptz` | SÍ | — | (duplicado de updated_at) |

**Índices:** `id` (PK), `finca_id`, `lote_id`, `tenant_id`, `node_key` (ÚNICO), `nodo_id` (ÚNICO)

> **Nota:** Las columnas `creado_en` / `actualizado_en` son duplicados de `created_at` / `updated_at` (artefacto de cambios en el modelo GORM). Se prefieren `created_at` / `updated_at`.

---

## 4. Tabla `muestras` — Muestras de diagnóstico

Creada por `MuestraModel` → `muestra_model.go`

| Columna | Tipo | Nulable | Default | Descripción |
|---|---|---|---|---|
| `id` | `varchar(36)` | NO | — | UUID primario |
| `lote_id` | `varchar(36)` | SÍ | — | FK lógica a `lotes.id` |
| `tenant_id` | `varchar(36)` | SÍ | — | **Tenant** |
| `latitud` | `double precision` | SÍ | — | Coordenada GPS |
| `longitud` | `double precision` | SÍ | — | Coordenada GPS |
| `created_at` | `timestamptz` | SÍ | — | Fecha de creación |
| `updated_at` | `timestamptz` | SÍ | — | Última actualización |

**Índices:** `id` (PK), `lote_id`, `tenant_id`

---

## 5. Tabla `diagnosticos` — Diagnósticos de clorosis

Creada por `DiagnosticoModel` → `diagnostico_model.go`

| Columna | Tipo | Nulable | Default | Descripción |
|---|---|---|---|---|
| `id` | `varchar(36)` | NO | — | UUID primario |
| `nombre` | `varchar(200)` | SÍ | — | Nombre del diagnóstico |
| `muestras_id` | `varchar(36)` | SÍ | — | FK lógica a `muestras.id` |
| `tenant_id` | `varchar(36)` | SÍ | — | **Tenant** |
| `estado` | `varchar(20)` | SÍ | `'PENDIENTE'` | `PENDIENTE` / `ACEPTADO` / `RECHAZADO` |
| `image_url` | `text` | SÍ | — | URL de la imagen inferida |
| `tiene_clorosis` | `boolean` | SÍ | — | Resultado de la inferencia |
| `confianza` | `numeric(5,4)` | SÍ | — | Nivel de confianza (0.0000–9.9999) |
| `procesado_at` | `timestamptz` | SÍ | — | Cuándo se procesó la inferencia |
| `created_at` | `timestamptz` | SÍ | — | Fecha de creación |
| `updated_at` | `timestamptz` | SÍ | — | Última actualización |

**Índices:** `id` (PK), `muestras_id`, `tenant_id`

---

## 6. Tabla `candidatos_reentrenamiento` — Candidatos a reentrenar modelo

Creada por `CandidatoModel` → `candidato_model.go`

| Columna | Tipo | Nulable | Default | Descripción |
|---|---|---|---|---|
| `id` | `varchar(36)` | NO | — | UUID primario |
| `diagnostico_id` | `varchar(36)` | SÍ | — | **FK única** a `diagnosticos.id` (1:1) |
| `image_url` | `text` | SÍ | — | URL de la imagen |
| `tiene_clorosis` | `boolean` | SÍ | — | Resultado |
| `confianza` | `numeric(5,4)` | SÍ | — | Confianza |
| `motivo` | `text` | SÍ | — | Motivo del rechazo (nulable) |
| `rechazado_por_usuario_id` | `varchar(36)` | SÍ | — | Usuario que rechazó |
| `created_at` | `timestamptz` | SÍ | — | Fecha de creación |

**Índices:** `id` (PK), `diagnostico_id` (ÚNICO)

---

## Relaciones (lógicas, sin FK explícitas)

```
fincas (1) ───< lotes (N)
fincas (1) ───< nodos (N)
lotes  (1) ───< muestras (N)
lotes  (1) ───< nodos (N)          [lote_id nullable]
muestras (1) ──< diagnosticos (N)   [vía muestras_id]
diagnosticos (1) ── (0..1) candidatos_reentrenamiento  [1:1 por uniqueIndex]
```

No hay `FOREIGN KEY` explícitas — la integridad referencial se mantiene en la capa de dominio/ aplicación.

---

## Multi-tenancy

Todas las tablas tienen `tenant_id`:

| Tabla | tenant_id | Notas |
|---|---|---|
| `fincas` | Nulable (`text`) | Una finca puede no tener tenant asignado |
| `lotes` | Nulable (`text`) | Hereda tenant de la finca |
| `nodos` | Nulable (`varchar(36)`) | Hereda tenant de la finca |
| `muestras` | Nulable (`varchar(36)`) | Hereda tenant del lote |
| `diagnosticos` | Nulable (`varchar(36)`) | Hereda tenant de la muestra |
| `candidatos_reentrenamiento` | — | No tiene tenant (dato interno del modelo) |

El tenant se asigna desde el `AuthContext.TenantID`, que se extrae del JWT (`tenant_id` claim).

---

## Estados (máquinas de estado)

### Finca (`fincas.estado`)
- `ACTIVA` → `PENDIENTE_ELIMINACION`

### Lote (`lotes.estado`)
- `ACTIVO` → `ELIMINADO`

### Nodo (`nodos.estado`)
- `ACTIVO` → `INACTIVO` | `MANTENIMIENTO`

### Diagnóstico (`diagnosticos.estado`)
- `PENDIENTE` → `ACEPTADO` | `RECHAZADO`

---

## Flujo de autenticación JWT

```
POST /api/v1/auth/login (identidad)
  ↓
JWT con: sub (usuario_id), sid (sesion_id),
         tenant_id, rol (sys_admin/admin/agronomo/caficultor)
  ↓
Frontend → Authorization: Bearer <JWT>
  ↓
Middleware fincas valida: firma HMAC + expiración + tipo="access"
  ↓
Extrae: usuarioID, sesionID, tenantID, rol
  ↓
AuthContext { UsuarioID, TenantID, Permisos: nil }
  ↓
Use cases: TienePermiso("fincas:finca:crear")
           → si Permisos=nil → true (desarrollo)
           → si hay lista → busca match
```

---

## API Endpoints

### Públicos (sin JWT)
| Método | Ruta | Uso |
|---|---|---|
| `GET` | `/api/v1/nodos/validar?nodeKey=:key` | YOLO valida cámara |
| `POST` | `/api/v1/diagnosticos/inferencia` | YOLO envía resultado |

### Protegidos (con JWT)
| Método | Ruta | Descripción |
|---|---|---|
| `POST` | `/fincas` | Crear finca |
| `POST` | `/fincas/:id/desactivar` | Desactivar finca |
| `POST` | `/fincas/:id/lotes` | Agregar lote |
| `DELETE` | `/fincas/:id/lotes/:loteID` | Eliminar lote |
| `GET` | `/nodos` | Listar nodos |
| `GET` | `/nodos/:id` | Obtener nodo |
| `POST` | `/nodos` | Registrar nodo |
| `PUT` | `/nodos/:id` | Editar nodo |
| `POST` | `/nodos/:id/desactivar` | Desactivar nodo |

### Endpoints de diagnóstico
| Método | Ruta | Descripción |
|---|---|---|
| `POST` | `/muestras/:loteID` | Tomar muestra |
| `GET` | `/muestras/:loteID` | Listar muestras de un lote |
| `POST` | `/diagnosticos/:diagnosticoID/aceptar` | Aceptar diagnóstico |
| `POST` | `/diagnosticos/:diagnosticoID/rechazar` | Rechazar diagnóstico |
| `POST` | `/diagnosticos/:muestraID/solicitar` | Solicitar diagnóstico manual |

---

## Permisos del módulo fincas

| Código | Descripción |
|---|---|
| `fincas:finca:crear` | Crear finca |
| `fincas:finca:desactivar` | Desactivar finca |
| `fincas:lote:crear` | Crear lote |
| `fincas:lote:eliminar` | Eliminar lote |
| `fincas:muestra:crear` | Tomar muestra |
| `fincas:muestra:consultar` | Ver muestras |
| `fincas:diagnostico:solicitar` | Solicitar diagnóstico |
| `fincas:diagnostico:aceptar` | Aceptar diagnóstico |
| `fincas:diagnostico:rechazar` | Rechazar diagnóstico |
| `fincas:reporte:generar` | Generar reporte |
| `fincas:nodo:crear` | Crear nodo IoT |
| `fincas:nodo:consultar` | Consultar nodos |
| `fincas:nodo:editar` | Editar nodo |
| `fincas:nodo:desactivar` | Desactivar nodo |
