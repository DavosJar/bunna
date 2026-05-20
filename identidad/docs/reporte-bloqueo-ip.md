# Reporte de Prueba: Bloqueo por IP ante reintentos fallidos

**Fecha:** 2026-05-13  
**Tester:** Alexis Jara  
**Developer:** Cesar Ramos  
**Servicio:** identidad (bunna)  
**Endpoint:** `POST /api/v1/auth/login`

---

## Configuración actual del sistema

| Variable | Valor | Default en |
|---|---|---|
| `BLOQUEO_IP_MAX_INTENTOS` | 20 | `internal/config/env.go:66` |
| `BLOQUEO_IP_VENTANA` | 15m | `internal/config/env.go:70` |
| `BLOQUEO_IP_DURACION` | 30m | `internal/config/env.go:74` |
| `RATE_LIMIT_MAX_REQUESTS` | 10 | `internal/config/env.go:78` |
| `RATE_LIMIT_VENTANA` | 1m | `internal/config/env.go:82` |

### Mecanismos de protección

| Mecanismo | ¿Qué cuenta? | Umbral | Ventana | Penalidad | Tabla |
|---|---|---|---|---|---|
| **Rate Limiter** | Toda request (exitosa o no) | 10 requests | 1 minuto | Bloqueo de 1 minuto | `rate_limit_ip` |
| **Bloqueo por IP** | Solo password incorrecto | 20 fallos | 15 minutos | Bloqueo de 30 min (TTL fijo) | `intentos_por_ip` |

### Flujo de código

`internal/sesiones/application/services/login/servicio_login.go:53-176`

```
1. Validar email + password
2. Rate limiter Verificar()               ← bloquea TODA request si excede 10/min
3. Bloqueo IP Verificar()                 ← bloquea TODA request si IP está bloqueada
4. [Transacción DB]
   a. Buscar usuario por email
   b. Obtener credenciales
   c. Verificar password
   d. Si ok: crear sesión + tokens
   e. Si falla: marca passwordIncorrecto=true
5. Si passwordIncorrecto:
   → Bloqueo IP RegistrarIntentoFallido() ← acumula contador
```

---

## Prueba realizada

### Escenario

Desde una misma IP, 20 intentos con password incorrecto (a través de dos ventanas del rate limiter), seguido de 1 intento con password correcto.

**IP utilizada:** `10.0.0.222`  
**Seed:** `admin@bunna.com` / `admin1234` (usuario ACTIVO y VERIFICADO)

### Fase 1 — 12 intentos fallidos rápidos

```
  #01  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.792s  credenciales inválidas
  #02  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.611s  credenciales inválidas
  #03  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.578s  credenciales inválidas
  #04  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.439s  credenciales inválidas
  #05  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.655s  credenciales inválidas
  #06  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.594s  credenciales inválidas
  #07  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.449s  credenciales inválidas
  #08  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.417s  credenciales inválidas
  #09  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.521s  credenciales inválidas
  #10  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.603s  credenciales inválidas
────── rate limiter se activa ──────
  #11  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.004s  demasiados intentos, intente más tarde
  #12  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.001s  demasiados intentos, intente más tarde
```

**Contador bloqueo IP post-Fase 1:** 10 fallos

### Pausa: 65 segundos

Espera para que la ventana del rate limiter (1 minuto) expire.

### Fase 3 — 12 intentos fallidos más

```
  #13  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.647s  credenciales inválidas
  #14  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.544s  credenciales inválidas
  #15  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.540s  credenciales inválidas
  #16  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.723s  credenciales inválidas
  #17  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.591s  credenciales inválidas
  #18  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.450s  credenciales inválidas
  #19  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.444s  credenciales inválidas
  #20  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.477s  credenciales inválidas
  #21  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.514s  credenciales inválidas
  #22  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.503s  credenciales inválidas
────── en el intento #22 se activa el BLOQUEO IP (20° fallo) ──────
  #23  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.002s  demasiados intentos, intente más tarde
  #24  POST /api/v1/auth/login  X-Real-IP:10.0.0.222  wrongpass  HTTP 401  0.002s  demasiados intentos, intente más tarde
```

**Contador bloqueo IP post-Fase 3:** 20 fallos → **IP BLOQUEADA por 30 minutos**

### Pausa: 65 segundos

Espera para que la ventana del rate limiter expire. **El bloqueo IP continúa activo** (TTL fijo de 30 minutos, no se ve afectado por el reseteo del rate limiter).

### Fase 5 — Login con credenciales CORRECTAS

```
        POST /api/v1/auth/login  X-Real-IP:10.0.0.222  admin1234 (correcto)  HTTP 401  0.041s  IP bloqueada temporalmente por exceso de intentos
```

### Resultado

| Intento | ¿Debe pasar? | ¿Pasó? | ¿Qué pasó? |
|---|---|---|---|
| Fallos 1-10 | Pasan (validación de password) | ✅ | HTTP 401 ~530ms "credenciales inválidas" |
| Fallo 11 | Rate limiter bloquea | ✅ | HTTP 401 ~2ms "demasiados intentos" |
| Fallos 13-22 | Pasan (validación de password) | ✅ | HTTP 401 ~530ms "credenciales inválidas" |
| Fallo 22 | **Bloqueo IP se activa** (contador 20≥20) | ✅ | Contador 20, bloquea por 30 min |
| Login correcto | **Bloqueo IP bloquea** | ✅ | HTTP 401 ~41ms "IP bloqueada temporalmente por exceso de intentos" |

---

## Comportamiento verificado

### 1. TTL fijo — El bloqueo NO se libera con login correcto

El intento #22 activó el bloqueo IP. El login correcto (Fase 5) fue rechazado. Esto es correcto: si el login correcto liberara el bloqueo, el atacante sabría que acertó la contraseña (oráculo). TTL fijo de 30 minutos no filtra información.

### 2. Bloqueo por IP persiste entre ventanas de rate limiter

El bloqueo IP usa su propia tabla (`intentos_por_ip`) y su propia ventana (15 min). El rate limiter usa `rate_limit_ip` con ventana de 1 minuto. Son independientes. El contador de bloqueo IP acumuló fallos a través de dos ventanas del rate limiter sin reiniciarse.

### 3. Firmas de tiempo distinguen cada mecanismo

| Mecanismo | Tiempo | Explicación |
|---|---|---|
| Login normal | ~400-800ms | Consulta BD: usuario, credenciales, bcrypt verify |
| Rate limiter bloqueando | ~1-4ms | Solo lee `rate_limit_ip`, sin bcrypt |
| Bloqueo IP bloqueando | ~41ms | Solo lee `intentos_por_ip`, sin bcrypt |

---

## Archivos involucrados

| Archivo | Propósito |
|---|---|
| `internal/sesiones/application/services/login/servicio_login.go` | Orquestación del login: rate limiter → bloqueo IP → transacción |
| `internal/seguridad/application/services/rate_limiter/servicio_rate_limiter.go` | Rate limiter preventivo (10 requests/min) |
| `internal/seguridad/application/services/bloqueo_ip/servicio_bloqueo_ip.go` | Bloqueo por IP (20 fallos → 30 min TTL) |
| `internal/seguridad/domain/intento_ip.go` | Entidad compartida entre rate limiter y bloqueo IP |
| `internal/config/env.go` | Configuración vía variables de entorno |
| `internal/registry/registry.go` | Wiring de dependencias (líneas 87-103) |
