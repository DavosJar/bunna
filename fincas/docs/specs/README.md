# Especificaciones del Microservicio Fincas

## Estructura de Carpetas

```
specs/
├── README.md                              ← Este archivo
├── spec-arquitectura.md                   ← Arquitectura general del microservicio
├── fincas/
│   └── spec-domain.md                     ← Dominio (Finca, Lote, API, RBAC, persistencia)
├── application/
│   └── spec-application.md                ← Aplicación (casos de uso / servicios)
├── infraestructura/
│   └── spec-infrastructure.md             ← Infraestructura (repositorios, BD, configuración)
├── presentacion/
│   └── spec-presentation-backend.md       ← Presentación backend (handlers, middleware, router)
└── frontend/
    └── spec-presentation-frontend.md      ← Frontend React + Vite

adr/
├── architecture-context.md                ← Contexto arquitectónico y reglas invariantes
└── feature-template.md                    ← Template para creación de nuevas features
```

## Índice de Especificaciones

| Archivo | Estado | Versión | Descripción |
|---------|--------|:-------:|-------------|
| `spec-arquitectura.md` | ✅ Completado | 1.0 | Arquitectura general, capas, flujo de datos, estrategia de pruebas |
| `fincas/spec-domain.md` | ✅ Completado | 1.0 | Modelo de dominio, contratos API REST, integración RBAC, esquema PostgreSQL |
| `application/spec-application.md` | ✅ Completado | 1.0 | Casos de uso, servicios de aplicación, DTOs, transacciones |
| `infraestructura/spec-infrastructure.md` | ✅ Completado | 1.0 | Repositorios GORM, conexión BD, migraciones, configuración |
| `presentacion/spec-presentation-backend.md` | ✅ Completado | 1.0 | Handlers Gin, facades, mappers, middleware JWT+authz, router |
| `frontend/spec-presentation-frontend.md` | ✅ Completado | 1.0 | Páginas React, componentes, servicios HTTP, enrutamiento |
| `adr/architecture-context.md` | ✅ Completado | 1.0 | Reglas invariantes del flujo Handler→Facade→Mapper→Domain |
| `adr/feature-template.md` | ✅ Completado | 1.0 | Checklist y template para creación de nuevas features |

## Convenciones

- Todos los specs siguen el formato **"AI-Ready"**: nivel de detalle suficiente para generación de código automatizada.
- El contenido se organiza por capa de Clean Architecture dentro de cada spec.
- Los identificadores siguen el formato `TIPO-FIN-NNN` (ej: `REQ-FIN-001`, `AC-FIN-001`, `ERR-FIN-001`).

## Especificaciones Relacionadas (Identidad)

- `../../identidad/docs/adr/architecture-context.md` — Patrón de capas de presentación
- `../../identidad/docs/specs/presentacion/spec-presentation-layer.md` — Presentación con Gin + Huma v2
- `../../identidad/docs/specs/autorizacion/spec-rbac-authorization.md` — Sistema RBAC de identidad
