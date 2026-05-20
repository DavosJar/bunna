# 🔍 Informe de Tester — Etapa 5: Seguridad Perimetral

> **Propósito**: Validar los servicios de bloqueo por IP, rate limiting y timeouts contra los 13 escenarios definidos en la especificación.
> **Fecha**: 2026-05-10
> **Tester**: Alexis Jara
> **Developer**: Cesar Ramos

---

## Resultado Global

| Ítem | Resultado |
|------|-----------|
| Tests seguridad | 52 (8 bloqueo_ip + 13 domain + 24 postgres + 7 bcrypt) |
| Tests pasan | ✅ 52/52 |
| Tests rate_limiter | ❌ 0 — no existen |
| Escenarios de spec cubiertos | ⚠️ 5/13 |
| **Veredicto** | 🟡 **No aprueba — faltan tests de rate limiter e integración con login** |

---

## Cobertura vs Especificación

| # | Escenario | Estado | Dependencia |
|---|-----------|--------|-------------|
| 1 | IP no bloqueada → permitir | ✅ | — |
| 2 | IP bloqueada por umbral (20 intentos) | ✅ | — |
| 3 | IP con intentos, login exitoso → contador NO se resetea | ❌ | Requiere integrar con ServicioLogin |
| 4 | IP bloqueada impide login | ❌ | Requiere integrar con ServicioLogin |
| 5 | Bloqueo de IP expirado → permitir | ✅ | — |
| 6 | Limpieza de registros antiguos | ❌ | Interfaz definida, sin impl ni test |
| 7 | Rate limit: dentro del límite | ❌ | Sin test |
| 8 | Rate limit: límite excedido | ❌ | Sin test |
| 9 | Ventana deslizante | ❌ | Sin test |
| 10 | Reset después de ventana | ❌ | Sin test |
| 11 | Timeout absoluto de sesión | ✅ | Cubierto en Etapa 3 |
| 12 | Timeout de inactividad | ✅ | Cubierto en Etapa 4 |
| 13 | Timeout de refresh token | ✅ | Cubierto en Etapa 3 |

---

## Lo implementado

| Componente | Archivo | Estado |
|------------|---------|--------|
| `IntentoPorIP` (entidad) | `seguridad/domain/intento_ip.go` | ✅ |
| `IntentoIPRepositorio` (interfaz) | `seguridad/domain/intento_ip_repositorio.go` | ✅ |
| `ServicioBloqueoIP.Verificar()` | `bloqueo_ip/servicio_bloqueo_ip.go` | ✅ |
| `ServicioBloqueoIP.RegistrarIntentoFallido()` | `bloqueo_ip/servicio_bloqueo_ip.go` | ✅ |
| `ServicioRateLimit.Verificar()` | `rate_limiter/servicio_rate_limiter.go` | ⚠️ Sin tests |
| Tests de BloqueoIP | `bloqueo_ip/servicio_bloqueo_ip_test.go` | ✅ 8 tests |
| Tests de RateLimit | ❌ No existe archivo de test | ❌ |

---

## Lo que falta

### 1. Tests para RateLimiter (crítico)

`rate_limiter/servicio_rate_limiter.go` no tiene archivo de test. Cubre escenarios #7, #8, #9, #10.

### 2. Integración con ServicioLogin

Actualmente `ServicioLogin` no invoca ni `ServicioBloqueoIP` ni `ServicioRateLimit`. La integración debería ser:

```
Login actual:
  1. Validar comando
  2. Transacción: buscar usuario → credenciales → verificar password → crear sesión

Login con seguridad:
  1. Rate limiting (antes de todo)
  2. IP blocking (antes de todo)
  3. Transacción: igual que antes
  4. Si falla: registrar intento fallido por IP
```

### 3. Repositorio concreto para IntentoIP

La interfaz `IntentoIPRepositorio` está definida pero no hay implementación Postgres. Según la spec es parte de Etapa 6 (Infraestructura).

---

## Decisión

**No aprueba.** El servicio de rate limiting no tiene tests y falta la integración con el flujo de login. Se necesita:

1. Crear `rate_limiter/servicio_rate_limiter_test.go` (~6 tests)
2. Integrar `ServicioBloqueoIP` y `ServicioRateLimit` en `ServicioLogin`
3. Tests de integración de login con IP bloqueada (escenarios #3, #4)

Estos cambios tocan `servicio_login.go`, no solo tests.

---

*Fin del informe — Etapa 5: 🟡 Pendiente.*
