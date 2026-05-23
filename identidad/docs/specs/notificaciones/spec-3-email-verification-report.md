---
title: "Reporte de Implementación — Infraestructura de Email y Verificación de Correo"
version: 1.0
date: 2026-05-22
owner: Equipo Identidad
status: EN_PROGRESO
tags: reporte, email, verificacion, implementacion
---

# Reporte de Implementación: Infraestructura de Email y Verificación de Correo

> **Propósito**: Evaluar el estado actual del código contra lo especificado en `spec-3-email-verification.md`.

## 1. Resumen Ejecutivo

| Dimensión | Resultado |
|-----------|-----------|
| **Estado general** | EN_PROGRESO |
| **Módulo notificaciones/** | existe |
| **EmailServicio (interfaz)** | existe |
| **Impl. SMTP** | existe |
| **Templates de email** | existen |
| **Envío asíncrono** | implementado (parcialmente) |
| **Servicio de verificación** | implementado (lógica de negocio completa) |
| **Value Object PruebaVerificacion** | existe |
| **Build** | compila |

**Hallazgo principal**: El dominio y la lógica de aplicación están mayoritariamente implementados, pero la capa de infraestructura (implementación concreta del repositorio de verificación, configuración SMTP vía env vars, handlers HTTP, registro en Registry) está **incompleta**. El flujo punta a punta no está conectado.

## 2. Estado por Componente

### 2.1 Módulo de Notificaciones (Parte A)

| Componente | Estado | Archivos |
|------------|--------|----------|
| Estructura de carpetas (`internal/notificaciones/`) | ✅ Existe | `domain/` e `infrastructure/email/` |
| Interfaz `EmailServicio` | ✅ Existe con `Enviar` y `EnviarTemplate` | `internal/notificaciones/domain/email_servicio.go` |
| Tipos de email (`VERIFICACION_CORREO`, `RECUPERACION_CONTRASENA`) | ✅ Definidos | `internal/notificaciones/domain/tipos_email.go` |
| Templates de texto plano | ✅ Existen con marcadores `{{variable}}` | `internal/notificaciones/domain/templates.go` |
| Implementación SMTP con `net/smtp` | ✅ Existe con STARTTLS (puerto 587) y TLS directo | `internal/notificaciones/infrastructure/email/smtp_servicio.go` |
| Configuración vía variables de entorno | ❌ No implementada | `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`, `SMTP_FROM`, `EMAIL_ASYNC` no están en `internal/config/env.go` |
| Envío asíncrono (post-commit goroutine) | ✅ Implementado en `SMTPServicio.Enviar()` | `smtp_servicio.go` línea 53-59 (goroutine cuando `config.Async == true`) |
| `constructor_email.go` | ❌ No existe | No hay builder de cuerpos de email |
| Tests — templates | ✅ Existen (5 tests, todos pasan) | `internal/notificaciones/domain/templates_test.go` |
| Tests — SMTP (`smtp_servicio_test.go`) | ❌ No existen | No hay tests unitarios ni de integración para `SMTPServicio` |

### 2.2 Servicio de Verificación (Parte B)

| Componente | Estado | Archivos |
|------------|--------|----------|
| Estructura de carpetas (`internal/verificacion/`) | ✅ Existe | `domain/` y `application/services/verificacion/` |
| Value Object `PruebaVerificacion` | ✅ Existe con métodos `Expiro()`, `EstaPendiente()`, `CoincideCon()` | `internal/verificacion/domain/prueba_verificacion.go` |
| Tests `PruebaVerificacion` | ✅ Existen (7 tests, todos pasan) | `internal/verificacion/domain/prueba_verificacion_test.go` |
| Servicio de verificación — Solicitar | ✅ Implementado | `servicio_verificacion.go` líneas 53-97 |
| Servicio de verificación — Confirmar | ✅ Implementado | `servicio_verificacion.go` líneas 100-130 |
| Servicio de verificación — Reenviar | ✅ Implementado | `servicio_verificacion.go` líneas 133-152 |
| Token hasheado (SHA-256) | ✅ Implementado | `prueba_verificacion.go` función `hashearToken()` |
| Rate limiting de reenvíos | ✅ Implementado (configurable) | `servicio_verificacion.go` líneas 144-145 |
| Tests `servicio_verificacion_test.go` | ❌ No existen | No hay tests de integración para el servicio |
| `verificacion/infrastructure/` | ❌ No existe | No hay implementación concreta del repositorio (`verificacion_repositorio.go`) |
| Comandos (DTO entrada) | ✅ Definidos | `comando.go` (3 comandos) |
| Respuestas (DTO salida) | ✅ Definidos | `respuesta.go` (2 respuestas) |

### 2.3 Integración

| Componente | Estado | Detalle |
|------------|--------|---------|
| Registro en Registry | ❌ No implementado | `internal/registry/registry.go` no incluye `EmailServicio` ni `ServicioVerificacion` |
| Variables de entorno SMTP | ❌ No definidas | No aparecen en `internal/config/env.go` ni en `.env.example` |
| Variables de entorno VERIFICACION_* | ❌ No definidas | `VERIFICACION_TOKEN_EXPIRACION`, `VERIFICACION_MAX_REENVIOS`, `VERIFICACION_VENTANA_REENVIOS` no están en `env.go` |
| Dependencias con spec-0 (CorreoElectronico VO) | ✅ Integradas | `CorreoElectronico` VO existe, `EstadoVerificacionCorreo` con máquina de estados completa (4 estados, tests OK) |
| Campo `estado_verificacion_correo` en BD | ✅ Existe | `UsuarioModel` tiene columna `estado_verificacion_correo` y se persiste correctamente |
| Handlers HTTP de verificación | ❌ No existen | No hay endpoints para solicitar/confirmar/reenviar verificación |
| Rutas de verificación en router | ❌ No existen | `router.go` solo registra health, register y login |
| Registro con verificación post-commit | ❌ No implementado | `servicio_registro.go` no invoca `ServicioVerificacion.SolicitarVerificacion` después del commit |
| Interfaz `VerificacionRepositorio` | ✅ Definida | `repositorio.go` con 4 métodos |
| `UsuarioVerificacion` proyección | ✅ Definida | DTO con datos necesarios para el repositorio |
| Evento `CorreoVerificado` | ❌ No emitido | `servicio_verificacion.go` confirma pero no emite eventos de dominio |

## 3. Checklist de Validación

Para cada ítem del checklist de la spec (sección 7):

- [✅] **¿EmailServicio es una interfaz en el dominio de notificaciones?**
  → Sí, en `internal/notificaciones/domain/email_servicio.go` con métodos `Enviar` y `EnviarTemplate`.

- [✅] **¿La implementación SMTP usa `net/smtp` (stdlib)?**
  → Sí, `smtp_servicio.go` importa `net/smtp`. Usa `smtp.PlainAuth`, `smtp.Dial`, `smtp.NewClient`, `client.StartTLS`, `client.Auth`, `client.Mail`, `client.Rcpt`, `client.Data`.

- [✅] **¿Los templates son de texto plano con marcadores `{{variable}}`?**
  → Sí, `templates.go` define templates con `{{nombre}}`, `{{token}}`, `{{expiracion_horas}}`. `RenderizarTemplate` los reemplaza con `strings.ReplaceAll`.

- [✅] **¿El envío asíncrono no bloquea la respuesta HTTP?**
  → Sí, `smtp_servicio.go` línea 53-59: cuando `config.Async == true`, lanza una goroutine y retorna inmediatamente. Sin embargo, esta funcionalidad no está probada ni conectada desde Registry.

- [✅] **¿El token de verificación se almacena como hash SHA-256 en BD?**
  → Sí, `prueba_verificacion.go` usa `sha256.Sum256` y almacena solo el hash (`secretoHash`). El token en plano nunca se persiste.

- [✅] **¿El token en plano nunca se expone en respuestas HTTP?**
  → Sí, `servicio_verificacion.go` retorna solo mensajes de texto en las respuestas. El token solo se pasa al template de email.

- [✅] **¿El servicio de verificación está desacoplado del registro?**
  → Sí, `ServicioVerificacion` es un servicio independiente que depende solo de `VerificacionRepositorio`, `EmailServicio`, y `GeneradorID`. No depende del módulo de registro. Sin embargo, el registro aún no lo invoca.

- [✅] **¿Un token inválido no altera el estado del dominio?**
  → Sí, `servicio_verificacion.go` línea 110-111: si el hash no existe, retorna `ErrEnlaceInvalido` sin modificar estado.

- [✅] **¿La verificación tiene rate limiting de reenvíos configurable?**
  → Sí, `ConfigVerificacion` tiene `MaxReenvios` y `VentanaReenvios` con defaults (5, 24h). El servicio verifica `usuario.ContadorReenvios >= config.MaxReenvios`.

- [✅] **¿El fracaso del email no revierte la operación principal?**
  → Sí, `servicio_verificacion.go` línea 80-92: el email se envía en goroutine y el error solo se loguea.

- [✅] **¿Los tipos de template son extensibles para futuros usos (password recovery)?**
  → Sí, `tipos_email.go` ya define `TipoRecuperacionContrasena` y el template está en `templates.go`. El enum `TipoEmail` es extensible.

## 4. Brechas Detectadas

Lista detallada de lo que falta, ordenado por prioridad:

### Prioridad Alta (bloquea el flujo completo)

| # | Brecha | Componente | Archivo esperado | Impacto |
|---|--------|------------|------------------|---------|
| 1 | ❌ Falta implementación concreta del repositorio de verificación | `verificacion/infrastructure/persistence/postgres/` | `verificacion_repositorio.go` | Sin esto, `ServicioVerificacion` no puede persistir/consultar datos. El servicio compila pero no funciona. |
| 2 | ❌ Falta configuración de variables de entorno SMTP | `internal/config/env.go` | Campos `SMTPHost`, `SMTPPort`, `SMTPUser`, `SMTPPassword`, `SMTPFrom`, `EmailAsync` en `Config` | Sin esto, `SMTPServicio` no se puede construir con valores configurables. |
| 3 | ❌ Falta configuración de variables de entorno de verificación | `internal/config/env.go` | Campos `VerificacionTokenExpiracion`, `VerificacionMaxReenvios`, `VerificacionVentanaReenvios` en `Config` | Sin esto, `ConfigVerificacion` no se puede construir desde env vars. |
| 4 | ❌ Registry no registra EmailServicio ni ServicioVerificacion | `internal/registry/registry.go` | Falta `emailServicio`, `servicioVerificacion` y sus constructores | Sin registro, ningún otro servicio puede usarlos. |
| 5 | ❌ No hay handlers HTTP de verificación | `internal/presentation/handlers/` | `verificacion_handler.go` | No hay endpoints expuestos para la API. |

### Prioridad Media (completan la funcionalidad)

| # | Brecha | Componente | Archivo esperado | Impacto |
|---|--------|------------|------------------|---------|
| 6 | ❌ Falta integración post-commit en registro | `servicio_registro.go` | Llamada a `ServicioVerificacion.SolicitarVerificacion()` después del commit | Los usuarios nuevos no reciben email de verificación al registrarse. |
| 7 | ❌ No hay rutas de verificación en el router | `router.go` | Registro de handlers de verificación | Los endpoints no son accesibles. |
| 8 | ❌ No se emite evento `CorreoVerificado` | `servicio_verificacion.go` | Llamada a `usuario.Eventos().RegistrarVerificacion()` | Los eventos de dominio no se disparan. |

### Prioridad Baja (mejora calidad)

| # | Brecha | Componente | Archivo esperado | Impacto |
|---|--------|------------|------------------|---------|
| 9 | ❌ No existe `smtp_servicio_test.go` | `notificaciones/infrastructure/email/` | `smtp_servicio_test.go` | La implementación SMTP no tiene cobertura de pruebas (ni unitarias ni de integración). |
| 10 | ❌ No existe `servicio_verificacion_test.go` | `verificacion/application/services/verificacion/` | `servicio_verificacion_test.go` | El servicio de verificación no tiene tests de integración. |
| 11 | ❌ No existe `constructor_email.go` | `notificaciones/infrastructure/email/` | `constructor_email.go` | No hay builder de cuerpos de email (separación de concerns). |
| 12 | ❌ `.env.example` no incluye variables SMTP | `.env.example` | Variables SMTP y de verificación | Dificulta la configuración inicial del entorno. |

## 5. Recomendaciones

### Inmediatas (completar infraestructura para desbloquear el flujo)

1. **Implementar repositorio de verificación** — Crear `internal/verificacion/infrastructure/persistence/postgres/verificacion_repositorio.go` que implemente `VerificacionRepositorio` usando GORM. La tabla `usuarios` ya tiene el campo `estado_verificacion_correo`, pero se necesita además persistir el hash del token y el contador de reenvíos. Considerar agregar columnas `secreto_hash_verificacion`, `expiracion_verificacion`, `contador_reenvios`, `ultimo_reenvio` al `UsuarioModel`.

2. **Agregar variables de entorno en `env.go`** — Añadir campos de configuración SMTP y de verificación al struct `Config`, con sus funciones `parsarDuracion`/`getEnv` correspondientes y defaults según la spec.

3. **Conectar Registry** — Registrar `SMTPServicio` (como `EmailServicio`) y `ServicioVerificacion` en `registry.go`. Pasar la config desde `cfg`.

4. **Crear handlers HTTP** — Implementar endpoints:
   - `POST /verify/solicitar` — Solicitar verificación (requiere auth)
   - `GET /verify/confirmar?token=xxx` — Confirmar verificación
   - `POST /verify/reenviar` — Reenviar verificación (requiere auth)

5. **Integrar verificación en registro** — En `servicio_registro.go`, después del commit exitoso (fuera de la transacción), invocar `ServicioVerificacion.SolicitarVerificacion()` con el usuario creado.

### Corto plazo (calidad y robustez)

6. **Agregar tests para SMTP** — Implementar `smtp_servicio_test.go` con mock del servidor SMTP o usando `github.com/emersion/go-smtp` (o similar) para pruebas de integración.

7. **Agregar tests para servicio de verificación** — Implementar `servicio_verificacion_test.go` con mock de `VerificacionRepositorio` y `EmailServicio` para cubrir los 10 escenarios TDD de la spec.

### Medio plazo (extensibilidad)

8. **Emitir eventos de dominio** — `servicio_verificacion.go` debería emitir `CorreoVerificado` usando el sistema de eventos existente en `internal/usuarios/domain/usuario/eventos.go`.

9. **Mejorar envío asíncrono** — Considerar el reemplazo de la goroutine simple por una cola en memoria con reintentos para evitar pérdida de emails en fallos transitorios del SMTP.

---

## Archivos Relevantes

### Implementados (20 archivos)

| Archivo | Líneas | Propósito |
|---------|--------|-----------|
| `internal/notificaciones/domain/email_servicio.go` | 9 | Interfaz EmailServicio |
| `internal/notificaciones/domain/tipos_email.go` | 9 | Enum TipoEmail |
| `internal/notificaciones/domain/templates.go` | 55 | Templates + RenderizarTemplate |
| `internal/notificaciones/domain/errores.go` | 11 | Errores del módulo |
| `internal/notificaciones/domain/templates_test.go` | 78 | Tests de templates |
| `internal/notificaciones/infrastructure/email/smtp_servicio.go` | 142 | Implementación SMTP |
| `internal/verificacion/domain/prueba_verificacion.go` | 68 | Value Object PruebaVerificacion |
| `internal/verificacion/domain/prueba_verificacion_test.go` | 84 | Tests de PruebaVerificacion |
| `internal/verificacion/domain/errores.go` | 12 | Errores de verificación |
| `internal/verificacion/domain/repositorio.go` | 22 | Interfaz VerificacionRepositorio |
| `internal/verificacion/application/services/verificacion/comando.go` | 16 | Comandos (3) |
| `internal/verificacion/application/services/verificacion/respuesta.go` | 11 | Respuestas (2) |
| `internal/verificacion/application/services/verificacion/servicio_verificacion.go` | 152 | Servicio de verificación (Solicitar, Confirmar, Reenviar) |
| `internal/usuarios/domain/usuario/estado_verificacion_correo.go` | 30 | Máquina de estados (4 estados) |
| `internal/usuarios/domain/usuario/estado_verificacion_correo_test.go` | 152 | Tests de máquina de estados |
| `internal/usuarios/domain/usuario/correo_electronico.go` | 71 | VO CorreoElectronico con estado |
| `internal/usuarios/domain/usuario/correo_electronico_test.go` | 156 | Tests de CorreoElectronico |
| `internal/usuarios/domain/usuario/usuario.go` | 166 | Usuario (con VerificarCorreo, SolicitarReenvioVerificacion, MarcarEnlaceExpirado) |
| `internal/usuarios/infrastructure/persistence/postgres/usuario_model.go` | 60 | UsuarioModel con estado_verificacion_correo |
| `internal/usuarios/infrastructure/persistence/postgres/usuario_repositorio.go` | 161 | Repositorio con persistencia de estado_verificacion_correo |

### Faltantes (13 archivos)

| Archivo esperado | Prioridad |
|------------------|-----------|
| `internal/verificacion/infrastructure/persistence/postgres/verificacion_repositorio.go` | Alta |
| `internal/config/env.go` (añadir campos SMTP y VERIFICACION) | Alta |
| `internal/registry/registry.go` (añadir EmailServicio y ServicioVerificacion) | Alta |
| `internal/presentation/handlers/verificacion_handler.go` | Alta |
| `internal/usuarios/application/services/registro/servicio_registro.go` (invocar verificación) | Media |
| `internal/presentation/router/router.go` (añadir rutas) | Media |
| `internal/notificaciones/infrastructure/email/smtp_servicio_test.go` | Baja |
| `internal/verificacion/application/services/verificacion/servicio_verificacion_test.go` | Baja |
| `internal/notificaciones/infrastructure/email/constructor_email.go` | Baja |
| `.env.example` (añadir SMTP*) | Baja |
| Migración BD para hash_token, contador_reenvios (en UsuarioModel) | Alta* |

*Nota: La tabla `usuarios` actualmente tiene `estado_verificacion_correo` pero **no** tiene columnas para almacenar `secreto_hash_verificacion`, `expiracion_verificacion`, `contador_reenvios`, `ultimo_reenvio`. Estas columnas son necesarias para que `VerificacionRepositorio` funcione. Pueden añadirse al `UsuarioModel` existente y GORM las migrará automáticamente.

## Resumen de Validaciones

- `go build ./...` — ✅ Compila sin errores
- `go test ./internal/notificaciones/...` — ✅ 5 tests pasan (templates)
- `go test ./internal/verificacion/domain/...` — ✅ 7 tests pasan (PruebaVerificacion)
- `go test ./internal/usuarios/domain/usuario/...` — ⚠️ Falla 1 test no relacionado (`TestListarSinFiltros` — nil pointer en `repositorio_test.go`)
- `go test ./internal/verificacion/application/...` — ⚠️ No hay test files
- `go test ./internal/notificaciones/infrastructure/...` — ⚠️ No hay test files
