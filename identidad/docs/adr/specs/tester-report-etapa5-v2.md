# 🔍 Informe de Tester — Etapa 5: Seguridad Perimetral (v2)

> **Propósito**: Validar los servicios de bloqueo por IP, rate limiting e integración con login.
> **Fecha**: 2026-05-10
> **Tester**: Alexis Jara
> **Developer**: Cesar Ramos

---

## Resultado Global

| Ítem | v1 | v2 |
|------|----|----|
| Tests rate_limiter | ❌ 0 | ✅ 6 |
| Tests login | 18 | ✅ 20 (+2 integración IP) |
| Tests totales (seguridad + sesiones) | ~110 | ✅ ~128 |
| Escenarios cubiertos | ⚠️ 5/13 | ✅ 12/13 |
| **Veredicto** | 🟡 No aprueba | ✅ **APRUEBA — Avanzar a Etapa 6** |

---

## Cobertura vs Especificación

| # | Escenario | v1 | v2 |
|---|-----------|----|----|
| 1 | IP no bloqueada → permitir | ✅ | ✅ |
| 2 | IP bloqueada por umbral | ✅ | ✅ |
| 3 | IP con intentos, login exitoso → contador NO se resetea | ❌ | ✅ `TestLogin_IPConIntentosNoSeReset` |
| 4 | IP bloqueada impide login | ❌ | ✅ `TestLogin_IPBloqueada` |
| 5 | Bloqueo de IP expirado → permitir | ✅ | ✅ |
| 6 | Limpieza de registros antiguos | ❌ | ❌ *(Etapa 6: infraestructura)* |
| 7 | Rate limit: dentro del límite | ❌ | ✅ `TestRateLimit_DentroDelLimite` |
| 8 | Rate limit: límite excedido | ❌ | ✅ `TestRateLimit_LimiteExcedido` |
| 9 | Ventana deslizante | ❌ | ✅ `TestRateLimit_VentanaDeslizante` |
| 10 | Reset después de ventana | ❌ | ✅ `TestRateLimit_ResetDespuesDeVentana` |
| 11 | Timeout absoluto de sesión | ✅ | ✅ (Etapa 3) |
| 12 | Timeout de inactividad | ✅ | ✅ (Etapa 4) |
| 13 | Timeout de refresh token | ✅ | ✅ (Etapa 3) |

---

## Tests agregados en v2

### rate_limiter (6 nuevos)

| Test | Escenario |
|------|-----------|
| `TestRateLimit_DentroDelLimite` | #7 — 5 requests, límite 10 → permitir |
| `TestRateLimit_LimiteExcedido` | #8 — 11 requests, límite 10 → error |
| `TestRateLimit_VentanaDeslizante` | #9 — requests en t=0 y t=1, evalúa correctamente |
| `TestRateLimit_ResetDespuesDeVentana` | #10 — esperar ventana → permitir de nuevo |
| `TestRateLimit_IPVacia` | validación IP requerida |
| `TestRateLimit_11RequestsExcedeLimite` | #8 — caso borde de límite |

### login (2 nuevos con integración IP)

| Test | Escenario |
|------|-----------|
| `TestLogin_IPBloqueada` | #4 — IP bloqueada → login rechazado antes de validar credenciales |
| `TestLogin_IPConIntentosNoSeReset` | #3 — login exitoso desde IP con intentos previos, contador IP no se resetea |

---

## Decisión

**Se aprueba Etapa 5 y se avanza a Etapa 6 (Infraestructura).** El único escenario pendiente (#6 limpieza de registros) corresponde a la implementación concreta del repositorio en Postgres/Redis.

---

*Fin del informe — Etapa 5: ✅ Lista para avanzar a Etapa 6.*
