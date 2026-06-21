# Especificaciones del Proyecto — Identidad

Este directorio contiene todas las especificaciones técnicas del proyecto organizadas por dominio/caso de uso.

## Estructura de Carpetas

| Carpeta | Contenido |
|---------|-----------|
| `registro/` | Especificación del caso de uso de registro con verificación de correo |
| `sesiones/` | Especificación completa de sesiones: login, refresh, logout, seguridad perimetral e integración |
| `presentacion/` | Especificación de la capa de presentación: Gin + Huma + OpenAPI |
| `usuarios/` | Especificaciones del módulo de usuarios: entidad, VO y refactors del dominio |
| `notificaciones/` | Especificaciones del sistema de notificaciones: email, verificación de correo y recuperación de contraseña |
| `autorizacion/` | Especificaciones del módulo de autorización IAM (RBAC + multi-tenant) |
| `tester-reportes/` | Reportes de pruebas generados por el agente tester durante el desarrollo |
| `telemetria/` | 6 mini-specs del sistema de telemetría asíncrona (API, NEGOCIO, BD) |

## Mapa de Especificaciones

| Archivo | Descripción | Estado |
|---------|-------------|--------|
| `registro/spec_registro.md` | Registro de usuario con verificación por correo electrónico | ⏳ En elaboración |
| `sesiones/login_spec.md` | Login, refresh token, logout, seguridad perimetral e integración JWT | ✅ Completada |
| `presentacion/spec-presentation-layer.md` | Capa de presentación con Gin, Huma v2 y OpenAPI/Swagger | ✅ Completada |
| `usuarios/spec-0-correo-electronico-vo.md` | Spec 0 — Refactor: CorreoElectronico como Value Object, eliminación de VERIFICACION_FALLIDA | 📋 En revisión |
| `autorizacion/spec-tenant-management.md` | Modelo multi-tenant: tenants, membresía, contexto de tenant | ✅ Completada |
| `notificaciones/spec-3-email-verification.md` | Spec 3 — Infraestructura de email (SMTP, templates) y verificación de correo desacoplada | 📋 En revisión |
| `notificaciones/spec-4-password-recovery.md` | Spec 4 — Recuperación de contraseña vía email con token de un solo uso | 📋 En revisión |
| `autorizacion/spec-rbac-authorization.md` | RBAC: roles, permisos atómicos, servicio de autorización, claims JWT, gestión de usuarios | ✅ Completada |
| `telemetria/01-arquitectura-y-organizacion.md` | Fundación: paquetes, interfaces, payload unificado, feature flag | 📋 Nueva |
| `telemetria/02-buffer-async-kafka.md` | Buffer ring asíncrono, productor Kafka, backpressure, métricas Prometheus | 📋 Nueva |
| `telemetria/03-middleware-api.md` | Middleware Gin para captura de logs API con trace_id | 📋 Nueva |
| `telemetria/04-decorador-negocio.md` | Decoradores de casos de uso para auditoría de negocio vía Facades | 📋 Nueva |
| `telemetria/05-plugin-gorm.md` | Plugin GORM para captura de métricas de base de datos | 📋 Nueva |
| `telemetria/06-cableado-final.md` | Cableado en Registry, feature flag, graceful shutdown, pruebas | 📋 Nueva |
| `autorizacion/spec-iam-rbac.deprecated.md` | Versión anterior del módulo IAM (reemplazada por las 2 specs de autorización) | ❌ Deprecated |
| `tester-reportes/01-etapa1.md` a `06-etapa6.md` | Reportes de pruebas por etapa | ✅ Completados |
| `tester-reportes/TEST-SUMMARY.md` | Resumen general de pruebas | ✅ Completado |
| `tester-reportes/PAIR-PROGRAMMING.md` | Reporte de pair programming | ✅ Completado |
| `tester-reportes/full-compliance.md` | Reporte de cumplimiento completo | ✅ Completado |

## Documentación Relacionada

- `docs/adr/architecture-context.md` — Contexto arquitectónico y flujo de capas
- `docs/adr/feature-template.md` — Template para nuevas features
