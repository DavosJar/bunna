# Especificaciones del Proyecto — Identidad

Este directorio contiene todas las especificaciones técnicas del proyecto organizadas por dominio/caso de uso.

## Estructura de Carpetas

| Carpeta | Contenido |
|---------|-----------|
| `registro/` | Especificación del caso de uso de registro con verificación de correo |
| `sesiones/` | Especificación completa de sesiones: login, refresh, logout, seguridad perimetral e integración |
| `presentacion/` | Especificación de la capa de presentación: Gin + Huma + OpenAPI |
| `tester-reportes/` | Reportes de pruebas generados por el agente tester durante el desarrollo |

## Mapa de Especificaciones

| Archivo | Descripción | Estado |
|---------|-------------|--------|
| `registro/spec_registro.md` | Registro de usuario con verificación por correo electrónico | ⏳ En elaboración |
| `sesiones/login_spec.md` | Login, refresh token, logout, seguridad perimetral e integración JWT | ✅ Completada |
| `presentacion/spec-presentation-layer.md` | Capa de presentación con Gin, Huma v2 y OpenAPI/Swagger | ✅ Completada |
| `tester-reportes/01-etapa1.md` a `06-etapa6.md` | Reportes de pruebas por etapa | ✅ Completados |
| `tester-reportes/TEST-SUMMARY.md` | Resumen general de pruebas | ✅ Completado |
| `tester-reportes/PAIR-PROGRAMMING.md` | Reporte de pair programming | ✅ Completado |
| `tester-reportes/full-compliance.md` | Reporte de cumplimiento completo | ✅ Completado |

## Documentación Relacionada

- `docs/adr/architecture-context.md` — Contexto arquitectónico y flujo de capas
- `docs/adr/feature-template.md` — Template para nuevas features
