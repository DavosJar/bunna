---
title: "Spec 4 — Recuperación de Contraseña vía Email"
version: 1.0
date_created: 2026-05-22
owner: Equipo Identidad
tags: password, recovery, email, reset, forgot-password
---

## 1. Propósito y Alcance

Permitir a un usuario recuperar el acceso a su cuenta cuando ha olvidado su contraseña, mediante un flujo seguro de restablecimiento vía correo electrónico.

**Incluye:**
- Solicitud de restablecimiento de contraseña (envía email con token)
- Validación del token de restablecimiento
- Cambio de contraseña usando el token
- Expiración y seguridad de tokens
- Rate limiting por IP y por usuario

**No incluye:**
- La infraestructura de envío de email (spec 3)
- Verificación de correo (spec 3)
- Cambio de contraseña estando autenticado (es autogestión, no recuperación)

## 2. Definiciones

| Término | Definición |
|---------|------------|
| **Token de restablecimiento** | Código único (UUID) generado al solicitar recuperación. Se almacena hasheado (SHA-256) en una tabla dedicada. El valor en plano solo va en el email. |
| **Rate limiting** | Límite de intentos de solicitud de recuperación desde una misma IP o para un mismo usuario en una ventana de tiempo. |

## 3. Flujo Completo

```
USUARIO                              SISTEMA
   │                                     │
   │  1. POST /auth/recover              │
   │  { email: "user@dom.com" }          │
   │─────────────────────────────────────→│
   │                                     │  2. Validar rate limiting por IP
   │                                     │  3. Buscar usuario por email (sin revelar existencia)
   │                                     │  4. Generar token UUID
   │                                     │  5. Hashear token (SHA-256)
   │                                     │  6. Persistir: email + hash + expiración + usado=false
   │                                     │  7. EmailServicio.EnviarTemplate(
   │                                     │       dest, RECUPERACION_CONTRASENA,
   │                                     │       {nombre, token, expiracion_horas})
   │  8. "Si el email existe, recibirás   │
   │     un enlace de recuperación"       │
   │←─────────────────────────────────────│
   │                                     │
   │  9. Usuario abre email → enlace     │
   │     GET /auth/reset?token=xxx       │
   │─────────────────────────────────────→│
   │                                     │  10. Validar token (no expirado, no usado)
   │                                     │  11. Mostrar formulario de nueva contraseña
   │  12. Formulario HTML o API           │
   │  POST /auth/reset/confirm            │
   │  { token: "xxx", password: "..." }   │
   │─────────────────────────────────────→│
   │                                     │  13. Validar token otra vez
   │                                     │  14. Validar fortaleza de nueva contraseña
   │                                     │  15. Hashear nueva contraseña (bcrypt)
   │                                     │  16. Actualizar credenciales
   │                                     │  17. Marcar token como usado
   │                                     │  18. Invalidar TODAS las sesiones del usuario
   │  19. "Contraseña actualizada"        │
   │←─────────────────────────────────────│
```

## 4. Modelo de Datos

```sql
-- Tabla: tokens_recuperacion
CREATE TABLE tokens_recuperacion (
    id VARCHAR(36) PRIMARY KEY,
    usuario_id VARCHAR(36) NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    expira_en TIMESTAMP NOT NULL,
    usado BOOLEAN NOT NULL DEFAULT false,
    creado_en TIMESTAMP NOT NULL,
    usado_en TIMESTAMP
);

CREATE INDEX idx_tokens_recuperacion_hash ON tokens_recuperacion(token_hash);
CREATE INDEX idx_tokens_recuperacion_usuario ON tokens_recuperacion(usuario_id);
```

## 5. Reglas de Negocio

| # | Regla |
|---|-------|
| RN-REC-01 | El token de recuperación se almacena como hash SHA-256 en BD, nunca en plano |
| RN-REC-02 | El token en plano solo se comunica al usuario vía email, nunca en respuesta HTTP |
| RN-REC-03 | El token tiene expiración configurable (default: 1 hora) |
| RN-REC-04 | Cada token es de un solo uso (usado = true después de utilizado) |
| RN-REC-05 | Al cambiar la contraseña con token, se invalidan TODAS las sesiones activas del usuario |
| RN-REC-06 | La solicitud de recuperación NO debe revelar si el email existe o no |
| RN-REC-07 | La fortaleza de la nueva contraseña se valida igual que en el registro |
| RN-REC-08 | El rate limiting de solicitudes es por IP (ej: 3 solicitudes/hora desde misma IP) |
| RN-REC-09 | El rate limiting de solicitudes es por usuario (ej: 1 solicitud/15 min por usuario) |
| RN-REC-10 | Si el token ha expirado, se puede solicitar uno nuevo (el viejo queda inutilizable) |

## 6. Servicio de Recuperación

### 6.1 Operación 1: Solicitar recuperación

```
POST /auth/recover
Request:  { "email": "user@dom.com" }
Response: 200 { "message": "Si el email existe, recibirás un enlace de recuperación" }
```

Flujo:
1. Validar rate limiting por IP (3 intentos/hora)
2. Validar formato de email
3. Buscar usuario por email
4. Si NO existe → retornar éxito genérico (no revelar existencia)
5. Si existe:
   a. Validar rate limiting por usuario (1 solicitud/15 min)
   b. Generar token UUID
   c. Hashear token (SHA-256)
   d. Persistir en `tokens_recuperacion`
   e. Enviar email con template `RECUPERACION_CONTRASENA`
6. Retornar éxito genérico

**Importante**: El paso 3 y 4 son silenciosos. Siempre se retorna el mismo mensaje de éxito, independientemente de si el email existe o no. Esto previene ataques de enumeración de usuarios.

### 6.2 Operación 2: Validar token (GET)

```
GET /auth/reset?token=xxx
Response 200: válido → muestra formulario
Response 410: expirado
Response 404: inválido
```

Flujo:
1. Hashear token recibido
2. Buscar en `tokens_recuperacion` por hash
3. Si no existe → error 404 "enlace inválido"
4. Si `usado = true` → error 410 "enlace ya utilizado"
5. Si `expira_en < ahora` → error 410 "enlace expirado"
6. Si todo OK → retornar éxito (el frontend muestra formulario de nueva contraseña)

### 6.3 Operación 3: Confirmar restablecimiento (POST)

```
POST /auth/reset/confirm
Request:  { "token": "xxx", "password": "NuevaPass123!" }
Response: 200 { "message": "Contraseña actualizada exitosamente" }
```

Flujo:
1. Validar token (misma validación que GET)
2. Validar fortaleza de nueva contraseña
3. Hashear nueva contraseña (bcrypt)
4. Transacción:
   a. Actualizar credenciales del usuario con nuevo hash
   b. Resetear intentos fallidos a 0
   c. Marcar token como usado (usado = true, usado_en = ahora)
   d. Invalidar TODAS las sesiones activas del usuario
5. Retornar éxito

## 7. Seguridad

| # | Riesgo | Mitigación |
|---|--------|------------|
| 1 | **Enumeración de usuarios**: atacante prueba emails para ver cuáles existen | Mensaje de éxito genérico siempre ("Si el email existe..."). No revelar existencia. |
| 2 | **Ataque de fuerza bruta al token**: adivinar el token de recuperación | Token UUID v4 (128 bits de entropía). Rate limiting por IP. |
| 3 | **Reutilización de token**: interceptar token y usarlo después | Token de un solo uso (usado = true). Expiración corta (1h). |
| 4 | **Secuestro de sesión post-recuperación**: sesión antigua sigue activa | Invalidar TODAS las sesiones del usuario al cambiar la contraseña. |
| 5 | **Rate limiting eludido**: múltiples solicitudes desde distintas IPs | Rate limiting combinado: por IP + por usuario (backup). |
| 6 | **Token en logs**: el token en plano podría quedar en logs de servidor | El token viaja en query param (GET) y body (POST), nunca en headers. Loggear solo operaciones, no tokens. |

## 8. Variables de Entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `RECUPERACION_TOKEN_EXPIRACION` | `1h` | Duración del token de recuperación |
| `RECUPERACION_RATE_LIMIT_IP` | `3` | Máximo solicitudes/hora desde misma IP |
| `RECUPERACION_RATE_LIMIT_USUARIO` | `1` | Máximo solicitudes/15 min por usuario |
| `RECUPERACION_RATE_LIMIT_VENTANA` | `15m` | Ventana de rate limiting por usuario |

## 9. Escenarios TDD

### 9.1 Solicitar recuperación

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 1 | Solicitud exitosa | Email existe, rate limit OK | SolicitarRecuperacion(email) | Token generado, hash persistido, email enviado |
| 2 | Solicitud con email inexistente | Email no registrado | SolicitarRecuperacion(email) | Mensaje genérico "Si el email existe...". Sin token, sin email. |
| 3 | Solicitud excede rate limit por IP | 3 solicitudes desde misma IP en 1h | SolicitarRecuperacion(email) | Error "demasiadas solicitudes" |
| 4 | Solicitud excede rate limit por usuario | 1 solicitud para ese usuario en 15 min | SolicitarRecuperacion(email) | Error "demasiadas solicitudes" |
| 5 | Email vacío | email = "" | validación | Error de validación |
| 6 | Email mal formado | email = "invalido" | validación | Error de validación |
| 7 | EmailServicio falla | Email existe, SMTP falla | SolicitarRecuperacion(email) | Token NO persistido (rollback), error retornado |
| 8 | Múltiples tokens para mismo usuario | 2 solicitudes separadas en el tiempo | SolicitarRecuperacion(email) | Ambos tokens válidos hasta que uno se use |

### 9.2 Validar token (GET)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 9 | Token válido | Token no expirado, no usado | ValidarToken(token) | Válido, retorna usuarioID |
| 10 | Token inválido (no existe) | Hash no encontrado en BD | ValidarToken(token) | Error "enlace inválido" |
| 11 | Token expirado | expira_en < ahora | ValidarToken(token) | Error "enlace expirado" |
| 12 | Token ya usado | usado = true | ValidarToken(token) | Error "enlace ya utilizado" |
| 13 | Token mal formado | String que no es UUID | ValidarToken(token) | Error "enlace inválido" |

### 9.3 Confirmar restablecimiento (POST)

| # | Caso | Given | When | Then |
|---|------|-------|------|------|
| 14 | Restablecimiento exitoso | Token válido, password fuerte | ConfirmarRestablecimiento(token, pass) | Nuevo hash, token usado, sesiones invalidadas |
| 15 | Token inválido | Token no existe | ConfirmarRestablecimiento(token, pass) | Error, contraseña no cambiada |
| 16 | Token expirado | Token expirado | ConfirmarRestablecimiento(token, pass) | Error, contraseña no cambiada |
| 17 | Token ya usado | Token ya consumido | ConfirmarRestablecimiento(token, pass) | Error, contraseña no cambiada |
| 18 | Password débil | Password no cumple políticas | ConfirmarRestablecimiento(token, pass) | Error de validación, token NO se marca usado |
| 19 | Sesiones invalidadas post-recuperación | Usuario tenía 3 sesiones activas | ConfirmarRestablecimiento(token, pass) | Las 3 sesiones → REVOCADAS |
| 20 | Intento de reuso del mismo token | Token ya usado en paso anterior | ConfirmarRestablecimiento(token, pass) | Error "enlace ya utilizado" |

## 10. Ubicación en el proyecto

```
internal/recuperacion/
├── domain/
│   ├── token_recuperacion.go        # Entidad TokenRecuperacion
│   ├── token_recuperacion_test.go
│   ├── repositorio.go               # Interfaz TokenRecuperacionRepositorio
│   └── errores.go                   # Errores de dominio
├── application/
│   └── services/
│       └── recuperacion/
│           ├── servicio_recuperacion.go    # Solicitar, Validar, Confirmar
│           ├── comando.go                  # Comandos
│           ├── respuesta.go                # Respuestas
│           └── servicio_recuperacion_test.go
└── infrastructure/
    └── persistence/
        └── postgres/
            ├── token_recuperacion_model.go    # GORM model
            └── token_recuperacion_repositorio.go  # Implementación
```

## 11. Dependencias

- **spec-3-email-verification.md**: Proporciona `EmailServicio` con template `RECUPERACION_CONTRASENA`
- **spec_registro.md**: Las políticas de fortaleza de contraseña deben ser las mismas
- **login_spec.md**: La invalidación de sesiones usa el repositorio de sesiones existente

## 12. Checklist de Validación

- [ ] ¿El token de recuperación se almacena como hash SHA-256 en tabla dedicada?
- [ ] ¿El token en plano nunca se expone en respuestas HTTP?
- [ ] ¿Los tokens son de un solo uso (usado = true)?
- [ ] ¿La expiración del token es configurable?
- [ ] ¿La solicitud no revela si el email existe o no?
- [ ] ¿Hay rate limiting por IP y por usuario?
- [ ] ¿Al confirmar el restablecimiento se invalidan todas las sesiones activas?
- [ ] ¿La nueva contraseña pasa las mismas validaciones que en registro?
- [ ] ¿Los tokens expirados o usados retornan error informativo pero seguro?
- [ ] ¿El endpoint de confirmar cambia la contraseña en transacción atómica?
- [ ] ¿Los intentos fallidos se resetean al restablecer la contraseña?
