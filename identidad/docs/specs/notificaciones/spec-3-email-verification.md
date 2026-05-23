---
title: "Spec 3 — Infraestructura de Email y Verificación de Correo"
version: 1.0
date_created: 2026-05-22
owner: Equipo Identidad
tags: email, verificacion, notificaciones, smtp
---

## 1. Propósito y Alcance

**Incluye:**

- Módulo de notificaciones (`notificaciones/`) con interfaz abstracta para envío de emails
- Implementación SMTP con Gmail (configurable)
- Sistema de templates de email (texto plano)
- Envío asíncrono (no bloqueante)
- Servicio de verificación de correo desacoplado del registro
- Estados, máquina de estados y flujo de verificación
- Token de verificación con expiración
- Reenvío de verificación con rate limiting

**No incluye:**

- Recuperación de contraseña (spec 4)
- Integración con el registro (ya existe en spec_registro.md)
- Cola de mensajes persistente (los emails se envían best-effort)
- Logs de emails enviados en BD

## 2. Definiciones

| Término | Definición |
|---------|------------|
| **EmailServicio** | Abstracción para enviar emails. Recibe destinatario, asunto y cuerpo. No sabe nada de verificación ni recuperación. |
| **Template de email** | Texto plano con marcadores de posición reemplazables ({{token}}, {{nombre}}, etc.) |
| **Envío asíncrono** | El email se envía en segundo plano (goroutine/cola en memoria). La respuesta HTTP no espera a que el email se envíe. |
| **Secreto de verificación** | Token único (UUID) generado al solicitar verificación. Se almacena hasheado (SHA-256) en BD. El valor en plano solo va en el email. |
| **PruebaVerificacion** | Value Object que encapsula el hash del secreto y su fecha de expiración. |

## 3. Módulo de Notificaciones (Parte A)

### 3.1 Arquitectura

El módulo `notificaciones/` es un módulo independiente y reutilizable:

- **Domain**: Define la interfaz `EmailServicio` y los tipos de email (enum de templates)
- **Infrastructure**: Implementación SMTP, templates, envío asíncrono
- NO depende de ningún otro módulo del proyecto

### 3.2 Interfaz EmailServicio

La interfaz debe permitir:

- `Enviar(ctx, destinatario, asunto, cuerpo) -> error` — Envío síncrono básico
- Opcionalmente: `EnviarTemplate(ctx, destinatario, tipoTemplate, datos) -> error` — Envío con template

Tipos de email (templates) que debe soportar desde el inicio:

| Tipo | Propósito | Variables |
|------|-----------|-----------|
| `VERIFICACION_CORREO` | Verificar dirección de correo | `{{nombre}}`, `{{token}}`, `{{expiracion_horas}}` |
| `RECUPERACION_CONTRASENA` | Restablecer contraseña (futuro, spec 4) | `{{nombre}}`, `{{token}}`, `{{expiracion_horas}}` |

### 3.3 Implementación SMTP

- Usa `net/smtp` de la stdlib (sin librerías externas)
- Conexión TLS en puerto 587 (STARTTLS)
- Configurable vía variables de entorno
- Timeout configurable para la conexión SMTP
- Log de errores (sin fallar la operación principal)

### 3.4 Envío Asíncrono

- El envío del email se hace post-commit (después de la transacción de BD)
- Si falla el envío, se loguea el error pero NO se revierte la operación
- Estrategia: goroutine simple (`go emailServicio.Enviar(...)`) en primera versión
- Futuro: cola de mensajes (Redis, RabbitMQ) — fuera de alcance

### 3.5 Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `SMTP_HOST` | `smtp.gmail.com` | Host del servidor SMTP |
| `SMTP_PORT` | `587` | Puerto SMTP |
| `SMTP_USER` | (requerido) | Email de la cuenta remitente |
| `SMTP_PASSWORD` | (requerido) | Contraseña de aplicación |
| `SMTP_FROM` | mismo que SMTP_USER | Dirección FROM en los correos |
| `EMAIL_ASYNC` | `true` | Si es true, los emails se envían asíncronamente |

### 3.6 Ubicación en el proyecto

```
internal/notificaciones/
├── domain/
│   ├── email_servicio.go     # Interfaz EmailServicio
│   ├── tipos_email.go        # Enum de tipos de template
│   ├── templates.go          # Templates de texto plano
│   └── errores.go            # Errores del módulo
└── infrastructure/
    └── email/
        ├── smtp_servicio.go       # Implementación SMTP
        ├── smtp_servicio_test.go  # Tests de integración
        └── constructor_email.go   # Builder de cuerpos de email
```

### 3.7 Escenarios TDD — EmailServicio

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Envío exitoso | SMTP configurado correctamente | `EmailServicio.Enviar(ctx, to, subject, body)` | nil error |
| 2 | Envío con credenciales inválidas | SMTP_USER / SMTP_PASSWORD incorrectos | Enviar email | Error de autenticación |
| 3 | Envío con destinatario inválido | Email destino mal formado | Enviar email | Error de validación |
| 4 | Template reemplaza variables | Template con `{{nombre}}` | Renderizar con datos {nombre: "Juan"} | "Hola Juan" |
| 5 | Envío asíncrono no bloquea | Email lento (5s) | Enviar asíncrono | La llamada retorna inmediatamente |
| 6 | Timeout de conexión SMTP | Servidor SMTP no responde | Enviar con timeout corto | Error de timeout |

## 4. Validación de Correo (Parte B)

### 4.1 Concepto

La validación de correo es un servicio independiente. NO está acoplado al registro. Cualquier flujo puede:

1. Solicitar verificación de un correo -> genera token + envía email
2. Confirmar verificación con token -> marca correo como VERIFICADO

Esto permite que el registro use verificación, y también otros flujos futuros (cambio de email, verificación forzada por admin, etc.)

### 4.2 Máquina de Estados

La verificación de correo tiene 4 estados:

```
PENDIENTE_VERIFICACION ---> VERIFICADO (confirmación exitosa)
PENDIENTE_VERIFICACION ---> ENLACE_EXPIRADO (pasó tiempo de expiración)
PENDIENTE_VERIFICACION ---> REENVIO_SOLICITADO (usuario pide reenvío)
ENLACE_EXPIRADO        ---> REENVIO_SOLICITADO (usuario pide reenvío)
REENVIO_SOLICITADO     ---> VERIFICADO (confirmación exitosa)
REENVIO_SOLICITADO     ---> ENLACE_EXPIRADO (expiró nuevo token)
VERIFICADO             ---> (terminal)
```

### 4.3 Value Object: PruebaVerificacion

Encapsula:

- `secretoHash` — hash SHA-256 del token (nunca el token en plano en BD)
- `expiraEn` — fecha de expiración del token

Comportamiento:

- `Expiro(ahora) bool` — true si ya expiró
- `EstaPendiente() bool` — true si tiene un secreto asignado

### 4.4 Servicio de Verificación (desacoplado)

**Operación 1: Solicitar verificación**

1. Recibe: `usuarioID` (o `correo`)
2. Genera token UUID único
3. Hashea el token (SHA-256)
4. Persiste el hash + expiración en el usuario (`PruebaVerificacion`)
5. Envía email con el token en plano usando `EmailServicio` (verificación template)
6. Cambia estado de verificación a `REENVIO_SOLICITADO` (si ya existía una verificación previa)
7. Retorna éxito (sin incluir el token en la respuesta)

**Operación 2: Confirmar verificación**

1. Recibe: token en plano
2. Hashea el token recibido
3. Busca usuario por hash del token
4. Si no existe -> error genérico "enlace inválido" (no altera estado)
5. Si el token expiró -> cambia estado a `ENLACE_EXPIRADO`, error "enlace expirado"
6. Si el token es válido -> cambia estado a `VERIFICADO`
7. Emite evento `CorreoVerificado`

**Operación 3: Reenviar verificación**

1. Recibe: `usuarioID`
2. Verifica que no haya excedido límite de reenvíos en ventana de tiempo
3. Genera nuevo token, actualiza PruebaVerificacion
4. Envía nuevo email
5. Incrementa contador de reenvíos

### 4.5 Reglas de negocio

| # | Regla |
|---|-------|
| RN-VER-01 | El token de verificación se almacena como hash SHA-256 en BD, nunca en plano |
| RN-VER-02 | El token en plano solo se comunica al usuario vía email, nunca en respuesta HTTP |
| RN-VER-03 | La expiración del token es configurable (default: 24h) |
| RN-VER-04 | Un token inválido NO altera el estado del dominio (solo retorna error) |
| RN-VER-05 | El máximo de reenvíos es configurable (default: 5) |
| RN-VER-06 | La ventana de reenvíos es configurable (default: 24h) |
| RN-VER-07 | Un correo ya verificado no puede solicitar reenvío |
| RN-VER-08 | El envío del email es best-effort: si falla, la operación no se revierte |

### 4.6 Flujo: Registro + Verificación

```
Registro (existente en spec_registro.md):
  1. Crear usuario + credenciales (transacción)
  2. Post-commit: SolicitarVerificacion(usuarioID)
     a. Generar token
     b. Persistir hash + expiración
     c. EmailServicio.EnviarTemplate(destino, VERIFICACION_CORREO, {token, nombre})
     d. Si email falla -> log, no bloquear registro

Confirmación (GET /verify?token=xxx):
  1. Recibir token del enlace
  2. ConfirmarVerificacion(token)
     a. Hashear token
     b. Buscar por hash
     c. Validar expiración
     d. Marcar VERIFICADO
     e. Retornar éxito
```

### 4.7 Variables de Entorno (adicionales)

| Variable | Default | Descripción |
|----------|---------|-------------|
| `VERIFICACION_TOKEN_EXPIRACION` | `24h` | Duración del token de verificación |
| `VERIFICACION_MAX_REENVIOS` | `5` | Máximo de reenvíos permitidos |
| `VERIFICACION_VENTANA_REENVIOS` | `24h` | Ventana para contar reenvíos |

### 4.8 Ubicación en el proyecto

```
internal/verificacion/             # Módulo independiente de verificación
├── domain/
│   ├── prueba_verificacion.go     # Value Object
│   ├── prueba_verificacion_test.go
│   ├── errores.go                 # Errores de verificación
│   └── repositorio.go             # Interfaz: buscar por hash de token
├── application/
│   └── services/
│       └── verificacion/
│           ├── servicio_verificacion.go      # Solicitar, Confirmar, Reenviar
│           ├── comando.go                    # Comandos (solicitar, confirmar, reenviar)
│           ├── respuesta.go                  # Respuestas DTO
│           └── servicio_verificacion_test.go
└── infrastructure/
    └── persistence/
        └── postgres/
            └── verificacion_repositorio.go   # Buscar por hash en tabla usuarios

internal/notificaciones/           # Módulo de comunicación (Parte A)
├── domain/
│   ├── email_servicio.go
│   ├── tipos_email.go
│   ├── templates.go
│   └── errores.go
└── infrastructure/
    └── email/
        ├── smtp_servicio.go
        ├── smtp_servicio_test.go
        └── constructor_email.go
```

### 4.9 Escenarios TDD — Verificación

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Solicitar verificación exitosa | Usuario con correo PENDIENTE | SolicitarVerificacion(usuarioID) | Token generado, hash persistido, email enviado, estado -> REENVIO_SOLICITADO |
| 2 | Solicitar verificación por primera vez | Usuario nuevo, nunca solicitó | SolicitarVerificacion(usuarioID) | Token generado, estado -> PENDIENTE_VERIFICACION |
| 3 | Confirmar verificación exitosa | Token válido, no expirado | ConfirmarVerificacion(token) | Estado -> VERIFICADO, prueba limpiada |
| 4 | Confirmar con token inválido | Hash del token no existe en BD | ConfirmarVerificacion(token) | Error genérico, estado NO cambia |
| 5 | Confirmar con token expirado | Token expirado (expiraEn < ahora) | ConfirmarVerificacion(token) | Estado -> ENLACE_EXPIRADO, error |
| 6 | Confirmar cuando ya VERIFICADO | Usuario ya VERIFICADO | ConfirmarVerificacion(token) | Error |
| 7 | Reenviar verificación exitoso | Usuario PENDIENTE, dentro de límite | ReenviarVerificacion(usuarioID) | Nuevo token, email enviado, contador+1 |
| 8 | Reenviar excede límite | 5 reenvíos en 24h | ReenviarVerificacion(usuarioID) | Error "demasiados intentos" |
| 9 | Reenviar cuando ya VERIFICADO | Correo ya verificado | ReenviarVerificacion(usuarioID) | Error |
| 10 | EmailServicio falla en solicitud | SMTP no disponible | SolicitarVerificacion(usuarioID) | Token persistido, email no enviado, SIN error al usuario |

## 5. Integración con Registry

```
Nuevas dependencias:
  - EmailServicio -> smtp_servicio (implementación SMTP)
  - ServicioVerificacion -> con dependencias: EmailServicio, repositorios
  - Al iniciar: registrar en Registry para que otros servicios lo usen
```

## 6. Dependencias

- **spec-0-correo-electronico-vo.md**: Define `CorreoElectronico` VO con el estado de verificación
- **spec_registro.md**: Usa este servicio de verificación post-registro
- **spec-4-password-recovery.md** (futura): Usa el mismo `EmailServicio`

## 7. Checklist de Validación

- [ ] ¿EmailServicio es una interfaz en el dominio de notificaciones?
- [ ] ¿La implementación SMTP usa `net/smtp` (stdlib)?
- [ ] ¿Los templates son de texto plano con marcadores `{{variable}}`?
- [ ] ¿El envío asíncrono no bloquea la respuesta HTTP?
- [ ] ¿El token de verificación se almacena como hash SHA-256 en BD?
- [ ] ¿El token en plano nunca se expone en respuestas HTTP?
- [ ] ¿El servicio de verificación está desacoplado del registro?
- [ ] ¿Un token inválido no altera el estado del dominio?
- [ ] ¿La verificación tiene rate limiting de reenvíos configurable?
- [ ] ¿El fracaso del email no revierte la operación principal?
- [ ] ¿Los tipos de template son extensibles para futuros usos (password recovery)?
