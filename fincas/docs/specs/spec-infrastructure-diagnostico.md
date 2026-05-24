---
title: Especificación de Infraestructura — Diagnósticos, Muestras e Inferencias
version: 1.0
date_created: 2026-05-24
owner: Fincas Team
tags: diagnostico, infraestructura, persistencia, postgres, gorm
---

# Especificación de Infraestructura — Diagnósticos, Muestras e Inferencias

> **Propósito**: Definir la arquitectura de la capa de infraestructura para el módulo de diagnósticos: modelos de persistencia, repositorios para Diagnóstico, Muestra y CandidatoReentrenamiento, y la extensión del Unit of Work para operaciones atómicas multi-entidad. Esta especificación complementa a `spec-infrastructure.md` (que cubre fincas y lotes).
>
> **Formato**: Especificación técnica. Define CÓMO se implementa, no QUÉ se implementa. Sin bloques de código.

---

## 1. Principios Arquitectónicos

### 1.1 Dependencias

La capa de infraestructura del módulo de diagnósticos:
- Depende del dominio de diagnósticos (implementa sus interfaces).
- No depende de la capa de Aplicación ni de Presentación.
- Reutiliza los tipos compartidos de `internal/shared/domain` (CriterioFiltro, Paginacion, Ordenacion, GeneradorID).
- Puede depender del paquete postgres de fincas para compartir la misma instancia de base de datos y el Unit of Work.
- El dominio de diagnósticos no conoce la infraestructura.

### 1.2 Responsabilidades

| Responsabilidad | Descripción |
|----------------|-------------|
| Persistencia | CRUD de Diagnóstico, Muestra y CandidatoReentrenamiento contra PostgreSQL |
| Consultas complejas | Filtrado dinámico, ordenación y paginación vía Specification |
| Transacciones | Unit of Work unificado que abarca repositorios de fincas y de diagnósticos |
| Mapeo | Conversión entre modelos GORM y entidades de dominio |
| Value Objects embebidos | Ubicacion (latitud, longitud) y ResultadoInferencia se mapean como columnas planas en la tabla |

### 1.3 Lo que NO hace

- Lógica de negocio (eso es dominio).
- Orquestación de casos de uso (eso es aplicación).
- Autorización (eso es middleware).

---

## 2. Estructura de Carpetas

Se crea `internal/diagnostico/infrastructure/persistence/postgres/` siguiendo el mismo patrón que `internal/fincas/infrastructure/persistence/postgres/`.

```
internal/diagnostico/infrastructure/
└── persistence/
    └── postgres/
        ├── diagnostico_model.go      ← Modelo GORM para Diagnóstico
        ├── muestra_model.go          ← Modelo GORM para Muestra
        ├── candidato_model.go        ← Modelo GORM para CandidatoReentrenamiento
        ├── diagnostico_repositorio.go ← Implementación de DiagnosticoRepositorio
        ├── muestra_repositorio.go     ← Implementación de MuestraRepositorio
        ├── candidato_repositorio.go   ← Implementación de CandidatoReentrenamientoRepositorio
        └── unit_of_work.go           ← Extensión del UnitOfWork con repositorios de diagnóstico
```

El Unit of Work unificado puede ubicarse en `internal/shared/infrastructure` o puede extenderse el existente de fincas para incluir también los repositorios de diagnóstico, según decisión arquitectónica.

---

## 3. Modelos de Persistencia (GORM)

### 3.1 Modelo Diagnostico

Tabla en BD: `diagnosticos`

| Columna | Tipo GORM | Descripción |
|---------|-----------|-------------|
| id | varchar(36) PK | UUID v7 |
| nombre | varchar(200) | Nombre generado automáticamente, formato "INF-{YYYYMMDD}-{random}" |
| muestras_id | varchar(36) FK | Referencia a la muestra asociada. Índice. |
| tenant_id | varchar(36) | Identificador del tenant. Índice. |
| estado | varchar(20) | PENDIENTE, ACEPTADO o RECHAZADO. Default PENDIENTE. |
| image_url | text | URL de la imagen procesada por YOLO. |
| tiene_clorosis | boolean | Resultado de la inferencia. |
| confianza | decimal(5,4) | Valor entre 0.0000 y 1.0000. |
| procesado_at | timestamp | Momento en que YOLO procesó la imagen. |
| created_at | timestamp | Auto gestionado por GORM. |
| updated_at | timestamp | Auto gestionado por GORM. |

El Value Object ResultadoInferencia se mapea como columnas planas (image_url, tiene_clorosis, confianza, procesado_at) en lugar de como tabla separada o JSON embebido, para facilitar consultas y evitar joins innecesarios.

Métodos de conversión:
- ToDomain: reconstruye la entidad Diagnostico con su ResultadoInferencia y Ubicacion usando los constructores NewDiagnosticoFromStorage y NewResultadoInferencia.
- FromDomain: extrae todos los campos de la entidad Diagnostico y su ResultadoInferencia al modelo plano.

### 3.2 Modelo Muestra

Tabla en BD: `muestras`

| Columna | Tipo GORM | Descripción |
|---------|-----------|-------------|
| id | varchar(36) PK | UUID v7 |
| lote_id | varchar(36) | Referencia al lote. Índice. |
| tenant_id | varchar(36) | Identificador del tenant. Índice. |
| latitud | double precision | Coordenada GPS, rango [-90, 90]. |
| longitud | double precision | Coordenada GPS, rango [-180, 180]. |
| created_at | timestamp | Auto gestionado por GORM. |
| updated_at | timestamp | Auto gestionado por GORM. |

El Value Object Ubicacion se mapea como columnas planas (latitud, longitud) en la misma tabla. Una vez creada, la ubicación no se modifica porque implicaría un cambio de suelo que invalidaría el historial de diagnósticos.

Métodos de conversión:
- ToDomain: reconstruye la entidad Muestra usando NewMusetraFromStorage y NewUbicacion.
- FromDomain: extrae los campos de la entidad Muestra y su Ubicacion al modelo plano.

### 3.3 Modelo CandidatoReentrenamiento

Tabla en BD: `candidatos_reentrenamiento`

| Columna | Tipo GORM | Descripción |
|---------|-----------|-------------|
| id | varchar(36) PK | UUID v7 |
| diagnostico_id | varchar(36) FK | Referencia al diagnóstico rechazado. Índice único. |
| image_url | text | URL de la imagen asociada al diagnóstico. |
| tiene_clorosis | boolean | Resultado de la inferencia original. |
| confianza | decimal(5,4) | Confianza de la inferencia original. |
| motivo | text nullable | Razón del rechazo proporcionada por el usuario. |
| rechazado_por_usuario_id | varchar(36) | ID del usuario que rechazó el diagnóstico. |
| created_at | timestamp | Auto gestionado por GORM. |

No tiene método ToDomain porque esta tabla no tiene una entidad de dominio equivalente. Es un modelo de persistencia puro al servicio del caso de uso RechazarDiagnostico. Su repositorio expone únicamente Crear y ListarPendientes.

---

## 4. Repositorios

### 4.1 DiagnosticoRepositorio

Implementa la interfaz `domain.DiagnosticoRepositorio` definida en `internal/diagnostico/domain/repositories.go`.

Paquete: `internal/diagnostico/infrastructure/persistence/postgres`
Tipo: privado `diagnosticoRepositorio`
Constructor exportado: `NewDiagnosticoRepositorio(db *gorm.DB) domain.DiagnosticoRepositorio`

Métodos a implementar:

| Método | Comportamiento |
|--------|---------------|
| Crear | Convierte Diagnostico a DiagnosticoModel con FromDomain. Ejecuta Create de GORM. |
| ObtenerPorID | Busca por id con First. Si no existe, retorna domain.ErrDiagnosticoNoEncontrado. Convierte con ToDomain. |
| ListarPorFinca | Busca diagnosticos cuyo muestras_id corresponda a muestras del lote de la finca. Requiere join con muestras y lotes, o consulta en dos pasos. |
| Actualizar | Convierte a modelo y ejecuta Updates de GORM con where id. |
| Eliminar | Ejecuta Delete de GORM con where id. |
| Buscar | Implementa Specification pattern: aplica filtros por campo, ordenación y paginación. Mapeo de columnas: nombre, muestraID → muestras_id, tenantID → tenant_id, estado → estado. Validación de columnas permitidas contra domain.ColumnasPermitidasDiagnostico. |

### 4.2 MuestraRepositorio

Implementa la interfaz `domain.MuestraRepositorio` definida en `internal/diagnostico/domain/repositories.go`.

Tipo: privado `muestraRepositorio`
Constructor exportado: `NewMuestraRepositorio(db *gorm.DB) domain.MuestraRepositorio`

Métodos a implementar:

| Método | Comportamiento |
|--------|---------------|
| Crear | Convierte Muestra a MuestraModel con FromDomain. Ejecuta Create de GORM. |
| ObtenerPorID | Busca por id con First. Si no existe, retorna error de dominio. Convierte con ToDomain. |
| ListarPorDiagnostico | Busca muestras cuyo id coincida con el muestras_id del diagnóstico dado. Equivalente a ObtenerPorID porque la relación es 1:1 (una muestra tiene un diagnóstico). |
| Actualizar | Convierte a modelo y ejecuta Updates con where id. |
| Eliminar | Ejecuta Delete con where id. |
| Buscar | Implementa Specification pattern. Mapeo de columnas: nombre → nombre, loteID → lote_id, tenantID → tenant_id. Validación contra domain.ColumnasPermitidasMuestra. |

### 4.3 CandidatoReentrenamientoRepositorio

Nuevo repositorio necesario para el caso de uso RechazarDiagnostico. No existe interfaz en el dominio actual porque se crea como parte de esta especificación.

Ubicación de la interfaz: `internal/diagnostico/domain/repositories.go` (añadir)
Ubicación de la implementación: `internal/diagnostico/infrastructure/persistence/postgres/candidato_repositorio.go`

Interfaz:

| Método | Descripción |
|--------|-------------|
| Crear(ctx, CandidatoReentrenamiento) | Persiste un nuevo candidato a reentrenamiento. |
| ObtenerPorDiagnosticoID(ctx, diagnosticoID) | Retorna el candidato asociado a un diagnóstico específico. |
| ListarPendientes(ctx, limite) | Retorna los N candidatos más antiguos no procesados aún para reentrenamiento. |

Tipo: privado `candidatoRepositorio`
Constructor exportado: `NewCandidatoReentrenamientoRepositorio(db *gorm.DB) domain.CandidatoReentrenamientoRepositorio`

El modelo CandidatoModel es un modelo de persistencia plano, sin entidad de dominio equivalente. El repositorio trabaja directamente con el modelo.

---

## 5. Unit of Work Unificado

### 5.1 Problema

Actualmente existe `UnitOfWorkPostgres` en `internal/fincas/infrastructure/persistence/postgres/` que maneja FincaRepositorio y LoteRepositorio. El caso de uso RechazarDiagnostico necesita un Unit of Work que abarque DiagnosticoRepositorio y CandidatoReentrenamientoRepositorio.

### 5.2 Solución propuesta

Opción recomendada: crear un Unit of Work propio del módulo de diagnósticos en `internal/diagnostico/infrastructure/persistence/postgres/unit_of_work.go` que maneje los repositorios de diagnóstico. Si en el futuro se necesita una transacción que cruce fincas y diagnósticos, se crea un Unit of Work unificado en `internal/shared/infrastructure/`.

El Unit of Work de diagnósticos expone:

| Método | Descripción |
|--------|-------------|
| Transaccional(ctx, fn) | Ejecuta fn dentro de una transacción GORM. Si fn retorna error, rollback. Si fn retorna nil, commit. |
| DiagnosticoRepo() | Retorna instancia de DiagnosticoRepositorio vinculada a la transacción. |
| MuestraRepo() | Retorna instancia de MuestraRepositorio vinculada a la transacción. |
| CandidatoRepo() | Retorna instancia de CandidatoReentrenamientoRepositorio vinculada a la transacción. |
| GeneradorID() | Retorna el generador de IDs compartido. |

Dentro de Transaccional, cada llamada a los métodos de repositorio crea una nueva instancia del repositorio con el *gorm.DB transaccional, de la misma forma que lo hace el UnitOfWorkPostgres existente.

### 5.3 Uso desde casos de uso

El caso de uso RechazarDiagnostico recibe el UnitOfWork por constructor. Dentro de su flujo:
1. Valida permisos y carga el diagnóstico (fuera de la transacción).
2. Ejecuta uw.Transaccional(ctx, func(tx) { tx.DiagnosticoRepo().Actualizar(...); tx.CandidatoRepo().Crear(...) }).
3. Si la transacción es exitosa, publica el evento.
4. Si la transacción falla, retorna error y no publica nada.

---

## 6. Esquema de Base de Datos

### 6.1 Tablas nuevas

Se crean dos tablas nuevas y una existente que se completa:

```sql
-- Tabla: muestras (completa la existente si ya fue creada por migraciones previas)
CREATE TABLE muestras (
    id          VARCHAR(36) PRIMARY KEY,
    lote_id     VARCHAR(36) NOT NULL,
    tenant_id   VARCHAR(36) NOT NULL,
    latitud     DOUBLE PRECISION NOT NULL,
    longitud    DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_muestras_lote_id ON muestras(lote_id);
CREATE INDEX idx_muestras_tenant_id ON muestras(tenant_id);

-- Tabla: diagnosticos
CREATE TABLE diagnosticos (
    id              VARCHAR(36) PRIMARY KEY,
    nombre          VARCHAR(200) NOT NULL,
    muestras_id     VARCHAR(36) NOT NULL REFERENCES muestras(id),
    tenant_id       VARCHAR(36) NOT NULL,
    estado          VARCHAR(20) NOT NULL DEFAULT 'PENDIENTE',
    image_url       TEXT NOT NULL,
    tiene_clorosis  BOOLEAN NOT NULL,
    confianza       DECIMAL(5,4) NOT NULL,
    procesado_at    TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_diagnosticos_muestras_id ON diagnosticos(muestras_id);
CREATE INDEX idx_diagnosticos_tenant_id ON diagnosticos(tenant_id);

-- Tabla: candidatos_reentrenamiento
CREATE TABLE candidatos_reentrenamiento (
    id                          VARCHAR(36) PRIMARY KEY,
    diagnostico_id              VARCHAR(36) NOT NULL UNIQUE REFERENCES diagnosticos(id),
    image_url                   TEXT NOT NULL,
    tiene_clorosis              BOOLEAN NOT NULL,
    confianza                   DECIMAL(5,4) NOT NULL,
    motivo                      TEXT,
    rechazado_por_usuario_id    VARCHAR(36) NOT NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_candidatos_created_at ON candidatos_reentrenamiento(created_at);
```

### 6.2 Consideraciones de migración

- Las migraciones pueden ser automáticas con GORM AutoMigrate o versionadas con archivos SQL.
- El orden de creación debe ser: muestras → diagnosticos → candidatos_reentrenamiento, respetando las FK.
- La tabla `muestras` referencia a `lotes` (de fincas) mediante `lote_id`, pero no se crea una FK física debido al aislamiento entre módulos. La integridad referencial se valida en la capa de aplicación.
- La columna `tenant_id` en diagnosticos y muestras permite filtrar por tenant sin necesidad de joins con fincas/lotes.

---

## 7. Reglas de Specification (Consultas Complejas)

### 7.1 Columnas permitidas para Diagnostico

| Campo dominio | Columna BD | Tipo de filtro |
|---------------|------------|----------------|
| nombre | nombre | =, !=, LIKE |
| muestraID | muestras_id | =, != |
| tenantID | tenant_id | = |
| estado | estado | =, != |

### 7.2 Columnas permitidas para Muestra

| Campo dominio | Columna BD | Tipo de filtro |
|---------------|------------|----------------|
| loteID | lote_id | = |
| tenantID | tenant_id | = |

### 7.3 Reglas de seguridad en consultas

- Si un filtro usa un campo no permitido, se ignora (degradación elegante, no error).
- Todos los valores se pasan como parámetros (?), nunca concatenados, para prevenir SQL injection.
- Paginación forzada: pagina < 1 → 1, tamano < 1 → 10.
- No se incluyen id, muestras_id ni lote_id como columnas de filtro público para evitar enumeración.
- Los filtros por tenantID se agregan automáticamente en el caso de uso, no en el repositorio.

---

## 8. Nuevas Interfaces de Dominio Requeridas

### 8.1 CandidatoReentrenamientoRepositorio

Se añade a `internal/diagnostico/domain/repositories.go`:

```go
type CandidatoReentrenamientoRepositorio interface {
    Crear(ctx context.Context, candidato *CandidatoReentrenamiento) error
    ObtenerPorDiagnosticoID(ctx context.Context, diagnosticoID string) (*CandidatoReentrenamiento, error)
    ListarPendientes(ctx context.Context, limite int) ([]CandidatoReentrenamiento, error)
}
```

### 8.2 Entidad CandidatoReentrenamiento

Nuevo archivo `internal/diagnostico/domain/candidato_reentrenamiento.go`:

Entidad de dominio plana (no es un aggregate raíz, es una entidad dependiente de Diagnostico). Contiene:

| Atributo | Tipo | Notas |
|----------|------|-------|
| id | string | Generado al crear |
| diagnosticoID | string | FK al diagnóstico rechazado |
| imageURL | string | Copia del resultado de inferencia |
| tieneClorosis | bool | Copia del resultado de inferencia |
| confianza | float64 | Copia del resultado de inferencia |
| motivo | string | Opcional, razón del rechazo |
| rechazadoPorUsuarioID | string | Quién rechazó |
| createdAt | time.Time | Asignado al crear |

### 8.3 Errores de dominio adicionales

Se añaden a `internal/diagnostico/domain/errores.go`:

- ErrCandidatoNoEncontrado: el candidato a reentrenamiento no existe.
- ErrCandidatoDuplicado: ya existe un candidato para ese diagnóstico.

---

## 9. Integración con Infraestructura Existente

### 9.1 Dependencia entre paquetes

El paquete `internal/diagnostico/infrastructure/persistence/postgres` importa:
- `gorm.io/gorm` para el ORM.
- `internal/diagnostico/domain` para las interfaces y entidades.
- `internal/shared/domain` para los tipos genéricos (CriterioFiltro, Paginacion).
- Opcional: `google/uuid` si se requiere generación de IDs en los repositorios (aunque idealmente los IDs se generan en la capa de aplicación).

No importa nada de `internal/fincas/` para mantener el aislamiento entre módulos.

### 9.2 Conexión a base de datos

Los repositorios de diagnóstico reciben la misma instancia de `*gorm.DB` que los repositorios de fincas, configurada en el punto de entrada de la aplicación (`cmd/main.go`). Ambos módulos comparten la misma conexión PostgreSQL y las mismas transacciones cuando sea necesario.

### 9.3 Inicialización en main.go

```
db = conectarPostgres()
generadorID = idgenerator.NewUUIDV7Generator()

// Repositorios de fincas
fincaRepo = postgres.NewFincaRepositorio(db)
loteRepo = postgres.NewLoteRepositorio(db)
uowFincas = postgres.NewUnitOfWorkPostgres(db, generadorID)

// Repositorios de diagnóstico
diagnosticoRepo = pgdiagnostico.NewDiagnosticoRepositorio(db)
muestraRepo = pgdiagnostico.NewMuestraRepositorio(db)
candidatoRepo = pgdiagnostico.NewCandidatoReentrenamientoRepositorio(db)
uowDiagnostico = pgdiagnostico.NewUnitOfWorkDiagnostico(db, generadorID)
```

---

## 10. Estrategia de Pruebas

| Nivel | Qué se prueba | Cómo |
|-------|--------------|------|
| Unitarias (Modelos) | ToDomain / FromDomain para DiagnosticoModel y MuestraModel | Crear modelo, convertir a dominio, verificar todos los campos incluyendo VOs anidados |
| Integración (DiagnosticoRepositorio) | CRUD + Buscar con filtros | PostgreSQL de prueba + GORM real |
| Integración (MuestraRepositorio) | CRUD + Buscar por loteID | PostgreSQL de prueba + GORM real |
| Integración (CandidatoRepositorio) | Crear, ObtenerPorDiagnosticoID, ListarPendientes | PostgreSQL de prueba + GORM real |
| Integración (UnitOfWork) | Transacción exitosa y con rollback | Ambos repos dentro de una transacción, forzar error y verificar que no persista nada |
| Filtros (Specification) | Cada operador y campo permitido en diagnosticos y muestras | Test parametrizado con tabla de casos |
| Seguridad | Filtros con campos no permitidos | Se ignoran, no rompen |
| VOs embebidos | Ubicacion y ResultadoInferencia como columnas planas | Crear y recuperar, verificar valores correctos |

---

## 11. Dependencias Tecnológicas

| Dependencia | Propósito |
|-------------|-----------|
| GORM + driver PostgreSQL | ORM y persistencia (compartido con fincas) |
| google/uuid | Generación de UUID v7 |
| testify | Asserts y mocks para pruebas |
| internal/shared/domain | Tipos genéricos CriterioFiltro, Paginacion, GeneradorID |

No se introducen dependencias nuevas. Todo se resuelve con las mismas librerías que ya usa el módulo de fincas.

---

## 12. Relación con Otras Especificaciones

- `spec-infrastructure.md` — Infraestructura del módulo de fincas (modelos Finca y Lote, repositorios, Unit of Work, Specification pattern base).
- `spec-fincas-domain.md` — Entidades de dominio Finca y Lote.
- `spec-application-casos-de-uso.md` — Casos de uso que orquestan estos repositorios (TomarMuestra, ListarMuestrasPorLote, SolicitarDiagnosticoManual, RegistrarInferencia, AceptarDiagnostico, RechazarDiagnostico, GenerarReportePorLote).
- `internal/diagnostico/domain/` — Interfaces de repositorio, entidades Diagnostico, Muestra, Ubicacion, ResultadoInferencia, errores y especificaciones que esta capa implementa.
