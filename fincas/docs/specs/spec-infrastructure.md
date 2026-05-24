---
title: Especificación de Infraestructura — Microservicio Fincas (Patrón Specification)
version: 1.0
date_created: 2026-05-23
owner: Equipo Catastro
tags: fincas, infraestructura, specification, persistencia, postgres
---

# Especificación de Infraestructura — Microservicio Fincas

> **Propósito**: Definir la arquitectura de la capa de infraestructura del microservicio `fincas`, implementando el patrón Specification para consultas complejas (filtrado dinámico, ordenación y paginación) y el patrón Repository para persistencia. Esta capa es la única que conoce PostgreSQL y GORM.
>
> **Formato**: Especificación técnica con ejemplos de código. Define CÓMO se implementa, no QUÉ se implementa.

---

## 1. Principios Arquitectónicos

### 1.1 Dependencias

La capa de infraestructura:
- **Depende del Dominio** (implementa sus interfaces)
- **No depende de la capa de Aplicación** ni de Presentación
- **No exporta tipos** fuera del paquete `infrastructure` (solo exporta el constructor)
- **El dominio no conoce la infraestructura** (regla de Clean Architecture)

### 1.2 Responsabilidades

| Responsabilidad | Descripción |
|----------------|-------------|
| **Persistencia** | CRUD de entidades contra PostgreSQL |
| **Consultas complejas** | Filtrado dinámico, ordenación y paginación vía Specification |
| **Transacciones** | Unit of Work para operaciones atómicas |
| **Mapeo** | Conversión entre modelos de persistencia y entidades de dominio |
| **Generación de IDs** | UUID v7 (time-ordered) |

### 1.3 Lo que NO hace la infraestructura

- ❌ Lógica de negocio (eso es dominio)
- ❌ Orquestación de casos de uso (eso es aplicación)
- ❌ Autorización (eso es middleware)

---

## 2. Estructura de Carpetas

```
fincas/internal/fincas/infrastructure/
└── persistence/
    └── postgres/
        ├── finca_model.go           ← Modelo GORM para Finca
        ├── lote_model.go            ← Modelo GORM para Lote
        ├── finca_repositorio.go     ← Implementación de FincaRepositorio
        ├── lote_repositorio.go      ← Implementación de LoteRepositorio
        ├── unit_of_work.go          ← Implementación de UnitOfWork
        └── mappers.go              ← Mapeo dominio ↔ modelo (opcional)

fincas/internal/shared/
├── domain/
│   └── specifications.go    ← Tipos genéricos: CriterioFiltro, Paginacion, Ordenacion
└── infrastructure/
    └── idgenerator/
        └── uuid_v7.go       ← Implementación de GeneradorID con UUID v7
```

---

## 3. Shared Domain (Tipos Genéricos)

Se crea un paquete compartido `internal/shared/domain` con los tipos genéricos para consultas. Estos tipos NO pertenecen a un módulo específico porque se reutilizan en todos.

### 3.1 CriterioFiltro

```go
package domain

// CriterioFiltro define una condición de filtro para consultas
type CriterioFiltro struct {
    Campo    string
    Operador string
    Valor    any
}

// Paginacion define los parámetros de paginación y ordenación
type Paginacion struct {
    Pagina       int
    TamanoPagina int
    Ordenaciones []Ordenacion
}

// Ordenacion define un criterio de ordenación
type Ordenacion struct {
    Campo string
    Tipo  TipoOrdenacion
}

// TipoOrdenacion define la dirección del orden
type TipoOrdenacion string

const (
    ASC  TipoOrdenacion = "ASC"
    DESC TipoOrdenacion = "DESC"
)
```

### 3.2 GeneradorID

```go
type GeneradorID interface {
    NextID(ctx context.Context) (string, error)
}
```

---

## 4. Especificación de Entidad (Dominio)

Cada entidad que requiera consultas complejas define su propia `Especificacion` en el dominio. Esto permite controlar qué campos son expuestos para búsqueda.

```go
package domain

import shared "ruta/al/shared/domain"

// EspecificacionFinca envuelve los criterios de búsqueda para fincas
type EspecificacionFinca struct {
    Filtros []shared.CriterioFiltro
}

// EspecificacionLote envuelve los criterios de búsqueda para lotes
type EspecificacionLote struct {
    Filtros []shared.CriterioFiltro
}

// ColumnasPermitidasFincas define los campos seguros para búsqueda
// Solo los campos que no exponen datos sensibles ni permiten enumeración
var ColumnasPermitidasFincas = map[string]bool{
    "nombre":    true,
    "ubicacion": true,
    "estado":    true,
    "usuarioID": true,
    "tenantID":  true,
}

// ColumnasPermitidasLotes define los campos seguros para búsqueda de lotes
var ColumnasPermitidasLotes = map[string]bool{
    "nombre":  true,
    "area":    true,
    "estado":  true,
    "fincaID": true,
}
```

### 4.1 Firma del Repositorio

La interfaz del repositorio en el dominio DEBE incluir una firma estandarizada para consultas:

```go
type FincaRepositorio interface {
    // ... métodos CRUD ...
    
    Listar(
        ctx context.Context,
        especificacion EspecificacionFinca,
        paginacion shared.Paginacion,
    ) ([]Finca, error)
}

type LoteRepositorio interface {
    // ... métodos CRUD ...
    
    Listar(
        ctx context.Context,
        especificacion EspecificacionLote,
        paginacion shared.Paginacion,
    ) ([]Lote, error)
}
```

---

## 5. Implementación Postgres (Ejemplo de Código)

### 5.1 Modelo GORM

```go
package postgres

import (
    "time"
    "ruta/dominio"
)

type FincaModel struct {
    ID          string     `gorm:"type:varchar(36);primaryKey;column:id"`
    Nombre      string     `gorm:"column:nombre"`
    Ubicacion   string     `gorm:"column:ubicacion"`
    Descripcion string     `gorm:"column:descripcion"`
    UsuarioID   string     `gorm:"column:usuario_id;index"`
    TenantID    *string    `gorm:"column:tenant_id;index"`
    Estado      string     `gorm:"column:estado;default:ACTIVA"`
    CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (FincaModel) TableName() string { return "fincas" }

func (m *FincaModel) ToDomain() *dominio.Finca {
    return dominio.NewFincaFromPersistence(
        m.ID, m.Nombre, m.Ubicacion, m.Descripcion, m.UsuarioID,
        m.TenantID, dominio.EstadoFinca(m.Estado), m.CreatedAt, m.UpdatedAt,
    )
}

func FromDomainFinca(f *dominio.Finca) *FincaModel {
    return &FincaModel{
        ID:          f.ID(),
        Nombre:      f.Nombre(),
        Ubicacion:   f.Ubicacion(),
        Descripcion: f.Descripcion(),
        UsuarioID:   f.UsuarioID(),
        TenantID:    f.TenantID(),
        Estado:      string(f.Estado()),
        CreatedAt:   f.CreatedAt(),
        UpdatedAt:   f.UpdatedAt(),
    }
}
```

### 5.2 Repositorio con Specification

```go
package postgres

import (
    "context"
    "errors"

    "ruta/dominio"
    shared "ruta/shared/domain"
    "gorm.io/gorm"
)

type fincaRepositorio struct {
    db *gorm.DB
}

func (r *fincaRepositorio) Listar(
    ctx context.Context,
    especificacion dominio.EspecificacionFinca,
    paginacion shared.Paginacion,
) ([]dominio.Finca, error) {
    query := r.db.WithContext(ctx).Model(&FincaModel{})

    // Mapeo de nombres: campo dominio → columna BD
    mapeoColumnas := map[string]string{
        "nombre":    "nombre",
        "ubicacion": "ubicacion",
        "estado":    "estado",
        "usuarioID": "usuario_id",
        "tenantID":  "tenant_id",
    }

    // Aplicar filtros (solo campos permitidos)
    for _, filtro := range especificacion.Filtros {
        if !dominio.ColumnasPermitidasFincas[filtro.Campo] {
            continue // ← degradación elegante: ignora filtros no permitidos
        }

        columnaDB, ok := mapeoColumnas[filtro.Campo]
        if !ok {
            continue
        }

        switch filtro.Operador {
        case "=":
            query = query.Where(columnaDB+" = ?", filtro.Valor)
        case "!=":
            query = query.Where(columnaDB+" != ?", filtro.Valor)
        case "LIKE":
            query = query.Where(columnaDB+" LIKE ?", filtro.Valor)
        }
    }

    // Aplicar ordenación (solo campos permitidos)
    for _, ord := range paginacion.Ordenaciones {
        if !dominio.ColumnasPermitidasFincas[ord.Campo] {
            continue
        }
        columnaDB, ok := mapeoColumnas[ord.Campo]
        if !ok {
            continue
        }
        orden := "ASC"
        if ord.Tipo == shared.DESC {
            orden = "DESC"
        }
        query = query.Order(columnaDB + " " + orden)
    }

    // Paginación segura
    pagina := paginacion.Pagina
    if pagina < 1 {
        pagina = 1
    }
    tamano := paginacion.TamanoPagina
    if tamano < 1 {
        tamano = 10
    }
    offset := (pagina - 1) * tamano

    var models []FincaModel
    result := query.Offset(offset).Limit(tamano).Find(&models)
    if result.Error != nil {
        return nil, result.Error
    }

    fincas := make([]dominio.Finca, len(models))
    for i, m := range models {
        fincas[i] = *m.ToDomain()
    }
    return fincas, nil
}
```

### 5.3 Unit of Work

```go
package postgres

import (
    "context"
    "ruta/dominio"
    "gorm.io/gorm"
)

type UnitOfWorkPostgres struct {
    db          *gorm.DB
    fincaRepo   dominio.FincaRepositorio
    loteRepo    dominio.LoteRepositorio
    generadorID dominio.GeneradorID
}

func (uw *UnitOfWorkPostgres) Transaccional(
    ctx context.Context,
    fn func(tx *UnitOfWorkPostgres) error,
) error {
    return uw.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
        txUow := &UnitOfWorkPostgres{
            db:          txDB,
            fincaRepo:   uw.fincaRepo,
            loteRepo:    uw.loteRepo,
            generadorID: uw.generadorID,
        }
        return fn(txUow)
    })
}

func (uw *UnitOfWorkPostgres) FincaRepository() dominio.FincaRepositorio {
    return uw.fincaRepo
}

func (uw *UnitOfWorkPostgres) LoteRepository() dominio.LoteRepositorio {
    return uw.loteRepo
}

func (uw *UnitOfWorkPostgres) GeneradorID() dominio.GeneradorID {
    return uw.generadorID
}
```

---

## 6. Flujo de una Consulta Compleja

```
Request: GET /api/v1/fincas?nombre=progreso&estado=ACTIVA&pagina=1&tamano=10&orden=nombre,ASC

1. Handler recibe query params
       ↓
2. Handler construye EspecificacionFinca y Paginacion
       ↓
3. Facade llama al servicio de aplicación
       ↓
4. Servicio llama a FincaRepositorio.Listar(ctx, espec, pag)
       ↓
5. Repositorio traduce a GORM:
   - Ignora campos no permitidos (ColumnasPermitidasFincas)
   - Mapea nombres dominio → columnas BD
   - Aplica filtros parametrizados (SQL injection safe)
   - Aplica ordenación
   - Aplica LIMIT/OFFSET
       ↓
6. Resultado: []FincaModel → ToDomain() → []Finca
       ↓
7. Respuesta atraviesa capas de vuelta
```

---

## 7. Reglas de Seguridad

| Regla | Descripción |
|-------|-------------|
| **Degradación elegante** | Si un filtro usa campo no permitido, se IGNORA (no error) |
| **Sin IDs en ColumnasPermitidas** | `id`, `fincaID` NO deben estar en el whitelist para evitar enumeración |
| **SQL Injection** | Todos los valores se pasan como parámetros (`?`), nunca concatenados |
| **Paginación forzada** | `Pagina` < 1 → 1. `TamanoPagina` < 1 → 10 |
| **Sin datos sensibles** | No incluir en whitelist campos que expongan datos de otros usuarios |

---

## 8. Estrategia de Pruebas

| Nivel | Qué se prueba | Cómo |
|-------|--------------|------|
| **Unitarias (Mappers)** | ToDomain / FromDomain | Crear modelo, convertir a dominio, verificar campos |
| **Integración (Repositorio)** | CRUD + Specification + paginación | PostgreSQL de prueba + GORM real |
| **Filtros** | Cada operador en cada campo permitido | Test parametrizado con tabla de casos |
| **Seguridad** | Filtros con campos no permitidos | Se ignoran, no rompen |
| **Paginación** | Límites, offsets, páginas vacías | Test con conjuntos de datos conocidos |

---

## 9. Dependencias Tecnológicas

| Dependencia | Propósito |
|-------------|-----------|
| GORM + driver PostgreSQL | ORM y persistencia |
| google/uuid | Generación de UUID v4/v7 |
| testify | Asserts y mocks para pruebas |
```
