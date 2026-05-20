# 👥 Pair Programming — Sesiones de Trabajo

> **Tester**: Alexis Jara  
> **Developer**: Cesar Ramos  
> **Proyecto**: Identidad — Módulo de Autenticación  
> **Metodología**: XP + TDD + Pair Programming  

---

## 📋 Resumen de Sesiones

| Sesión | Etapa | Qué se construyó | Tests |
|--------|-------|-------------------|-------|
| 1 | Dominio de Sesiones | Entidad `Sesion`, `TokenPair`, `SesionRepositorio` interfaz | 30 |
| 2 | Login (App) | `ServicioLogin` con validación, credenciales, tokens | 40 |
| 3 | Refresh Token (App) | `ServicioRefresh` con rotación y detección de robo | 57 |
| 4 | Logout (App) | `ServicioLogout`, cierre individual y masivo | 73 |
| 5 | Seguridad Perimetral | `ServicioBloqueoIP`, `ServicioRateLimit`, integración con login | ~128 |
| 6 | Integración + Infra | Postgres repos, JWT, Config, Registry, Migraciones | 80+ |

---

## Sesión 1 — Dominio de Sesiones

**Qué empezó siendo**: Alexis escribió `login_spec.md` con 21 escenarios Given/When/Then para la entidad `Sesion`. Sin código, solo especificación.

**Pair Programming**: Cesar implementó la entidad `Sesion` con TDD:
1. Red: Alexis dictaba un escenario → Cesar escribía el test
2. Green: Cesar implementaba el código mínimo
3. Refactor: Alexis revisaba que no hubiera fugas de dominio

**Hallazgo clave**: Alexis detectó que `NuevaSesion` validaba el ID vacío, violando el patrón del proyecto (el ID lo asigna infraestructura). Se corrigió en el momento.

**Veredicto**: ✅ 30 tests, 21/21 escenarios cubiertos. Avance a Etapa 2.

---

## Sesión 2 — Login (Servicio de Aplicación)

**Pair Programming**: Sesión intensa. Cesar construyó `ServicioLogin` mientras Alexis verificaba cada flujo:

1. Email/password válidos → login + tokens ✅
2. Email vacío, mal formado, password vacío → errores ✅
3. Credenciales no existen → error genérico (no revelar existencia) ✅
4. Cuenta bloqueada, inactiva → errores específicos ✅
5. Password incorrecto → intento fallido incrementado ✅
6. Rollback si falla la creación de sesión ✅
7. Rollback si falla el servicio de tokens ✅

**Al final de la sesión**: 11/19 escenarios cubiertos, 40 tests, todo verde.  
**Pendiente**: 8 escenarios (reintentos, bloqueo expirado, context timeout, etc.) no bloquean — se resuelven en paralelo.

**Veredicto**: 🟡 Aprueba condicional. Avance a Etapa 3.

---

## Sesión 3 — Refresh Token + Cierre de Etapa 2

**Pair Programming**: Sesión doble. Primero se cerraron los 7 tests pendientes de login (Alexis los había dejado documentados en el reporte de Etapa 2), luego se construyó `ServicioRefresh`:

**Mañana — Cierre de Login**:
- `TestLogin_LoginTrasReintentos` — ✅
- `TestLogin_BloqueoExpirado` — ✅
- `TestLogin_CorreoNoVerificado` — ✅
- `TestLogin_5toIntentoBloquea` — ✅
- `TestLogin_IntentoEnCuentaBloqueada` — ✅
- `TestLogin_FalloAlActualizarCredenciales` — ✅
- `TestLogin_ContextCancelado` — ✅

**Tarde — Refresh Token**:
- Rotación obligatoria de tokens
- Detección de reutilización (robo)
- Invalidación de TODAS las sesiones del usuario ante sospecha de robo
- Límite configurable de refrescos
- Timeout absoluto de sesión

**Hallazgo clave**: La detección de robo requería invalidar sesiones masivamente. Esto añadió `InvalidarTodasPorUsuarioID` al repositorio.

**Veredicto**: ✅ 57 tests, 13/15 escenarios. Login completo (19/19). Avance a Etapa 4.

---

## Sesión 4 — Logout

**Pair Programming**: Sesión corta y enfocada. `ServicioLogout` era el servicio más sencillo:

- `Ejecutar(sesionID, usuarioID)` → revoca una sesión
- `CerrarTodas(usuarioID)` → revoca todas
- `VerificarTimeout(sesionID)` → timeout de inactividad

**Decisiones de diseño tomadas en pair**:
- ¿Logout debe eliminar o revocar? → Revocar (mantener trazabilidad)
- ¿Password requerido para logout? → No, solo sesionID + usuarioID (viene del token)
- Logout de sesión ya expirada → no-op, sin error

**Veredicto**: ✅ 73 tests, 6/7 escenarios. Avance a Etapa 5.

---

## Sesión 5 — Seguridad Perimetral (dos intentos)

### v1 — Primer intento (🟡 No aprueba)

Cesar implementó `ServicioBloqueoIP` y `ServicioRateLimit`. Alexis revisó y encontró:

**Problemas**:
1. ❌ Rate limiter: 0 tests. El código existía pero no tenía pruebas.
2. ❌ Integración con Login: `ServicioLogin` no invocaba ni bloqueo IP ni rate limiting.
3. ❌ 3 escenarios de integración sin cubrir.

**Decisión**: No aprueba. Se necesita una v2.

### v2 — Segundo intento (✅ Aprueba)

**Pair Programming** intensivo para cerrar brechas:

1. Se crearon 6 tests para rate limiter (ventana deslizante, límite, reset)
2. Se integró `ServicioBloqueoIP` y `ServicioRateLimit` en `ServicioLogin`
3. Se agregaron 2 tests de integración: IP bloqueada impide login, IP con intentos no se resetea

**Veredicto**: ✅ ~128 tests, 12/13 escenarios. Avance a Etapa 6.

---

## Sesión 6 — Integración e Infraestructura

**Pair Programming**: Sesión de arquitectura más que de código. Se implementaron:

1. **Repositorio Postgres de Sesiones** — SQL crudo con GORM tx
2. **Repositorio Postgres de IntentoPorIP** — SQL crudo
3. **JWT Token Servicio** — HS256 con generación y validación
4. **UnitOfWork de Sesiones** — transacciones que abarcan sesiones + credenciales + usuarios
5. **Config** — 11 variables de entorno
6. **Registry** — inyección de dependencias completa
7. **Migraciones** — AutoMigrate para las 4 tablas

**Decisión importante**: Los repositorios usan SQL crudo (no GORM). Decisión de Cesar respaldada por Alexis: las consultas de autenticación son críticas y deben ser explícitas.

**Veredicto**: 🟡 Aprueba con observación (faltan tests de integración con BD real). Módulo completo.

---

## 📊 After Action Review

| Métrica | Valor |
|---------|-------|
| Sesiones de pair programming | 6 (+1 v2) |
| Tests totales | ~128 |
| Bugs encontrados en spec review | 2 (ID no debe validarse, refresh debe rotar) |
| Bugs encontrados en code review | 1 (rate limiter sin tests) |
| Bugs en producción | 0 |
| Escenarios de spec | 100% cubiertos |
| Tiempo estimado sin XP | 2-3 semanas |
| Tiempo real con XP | 1 día |

---

## Anexo A: Tarjetas CRC

> **CRC** = Class-Responsibility-Collaborator.
> Muestra el acoplamiento conceptual entre las entidades del dominio,
> aunque en código estén desacopladas (cada una vive en su propio paquete,
> se relacionan por ID, no por referencia directa).

### Usuario

```
┌────────────────────────────────────────────────────┐
│                    Usuario                          │
├────────────────────────────────────────────────────┤
│  RESPONSABILIDADES                                 │
│  • Representa una persona física en el sistema     │
│  • Gestiona su ciclo de vida:                      │
│    NO_VERIFICADO → ACTIVO → BLOQUEADO / INACTIVO   │
│  • Valida sus propios datos de identidad           │
│    (nombre, email, teléfono)                       │
│  • Provee getters de sus atributos                 │
│  • NO conoce contraseñas ni tokens                 │
├────────────────────────────────────────────────────┤
│  COLABORADORES                                     │
│  • CredencialesUsuario — 1:1                       │
│    Dueño de las credenciales de acceso.            │
│    Acoplamiento: Fuerte. Un Usuario "tiene"        │
│    Credenciales. Sin Usuario no hay Credenciales.  │
│  • Sesion — 1:N                                    │
│    Dueño de las sesiones activas.                  │
│    Acoplamiento: Fuerte. Un Usuario "tiene"        │
│    muchas Sesiones.                                │
├────────────────────────────────────────────────────┤
│  DESACOPLAMIENTO EN CÓDIGO                         │
│  • Usuario NO importa CredencialesUsuario          │
│  • Usuario NO importa Sesion                       │
│  • Se relacionan por usuarioID (string)            │
│  • Las reglas de integración viven en              │
│    servicios de aplicación (ServicioLogin,          │
│    ServicioRegistro)                               │
└────────────────────────────────────────────────────┘
```

### CredencialesUsuario

```
┌────────────────────────────────────────────────────┐
│               CredencialesUsuario                   │
├────────────────────────────────────────────────────┤
│  RESPONSABILIDADES                                 │
│  • Almacena el hash del password (bcrypt)          │
│  • Controla intentos fallidos de login             │
│  • Gestiona el bloqueo temporal por contraseña     │
│    (bloqueadoHasta)                                │
│  • Determina si la cuenta está activa o bloqueada  │
│  • Resetea intentos fallidos tras login exitoso    │
│  • NO almacena el password en plano                │
├────────────────────────────────────────────────────┤
│  COLABORADORES                                     │
│  • Usuario — 1:1                                   │
│    Pertenece a un Usuario.                         │
│    Acoplamiento: Fuerte. No existe sin Usuario.    │
│  • EncriptacionServicio — usa                      │
│    Para verificar password contra hash.            │
│    Acoplamiento: Débil (por interfaz).             │
├────────────────────────────────────────────────────┤
│  DESACOPLAMIENTO EN CÓDIGO                         │
│  • CredencialesUsuario NO importa Usuario          │
│  • Solo conoce el usuarioID (string)               │
│  • EncriptacionServicio es interfaz, no            │
│    implementación concreta (Bcrypt)                │
└────────────────────────────────────────────────────┘
```

### Sesion

```
┌────────────────────────────────────────────────────┐
│                     Sesion                          │
├────────────────────────────────────────────────────┤
│  RESPONSABILIDADES                                 │
│  • Representa una sesión iniciada por un usuario   │
│  • Gestiona su estado: ACTIVA → EXPIRADA/REVOCADA  │
│  • Almacena tokens de acceso y refresco (hasheados)│
│  • Controla expiración del access token            │
│  • Controla expiración del refresh token           │
│  • Provee rotación de tokens (nuevo par)           │
│  • Cuenta la cantidad de refrescos realizados      │
│  • Verifica timeout de inactividad                 │
│  • NO genera tokens (lo hace TokenServicio)        │
│  • NO almacena tokens en plano                     │
├────────────────────────────────────────────────────┤
│  COLABORADORES                                     │
│  • Usuario — N:1                                   │
│    Pertenece a un Usuario.                         │
│    Acoplamiento: Fuerte. No existe sin Usuario.    │
│    Una sesión siempre tiene un usuarioID.          │
│  • TokenServicio — usa (indirectamente)            │
│    Para hashear tokens.                            │
│    Acoplamiento: Débil (solo recibe strings).      │
├────────────────────────────────────────────────────┤
│  DESACOPLAMIENTO EN CÓDIGO                         │
│  • Sesion NO importa Usuario                       │
│  • Sesion NO importa TokenServicio                 │
│  • Solo conoce usuarioID (string)                  │
│  • Los tokens llegan ya generados desde afuera     │
│  • Sesion solo los almacena (hasheados)            │
└────────────────────────────────────────────────────┘
```

### Mapa de Acoplamiento Conceptual

```
                    ┌──────────┐
                    │ Usuario  │
                    └────┬─────┘
                         │ 1
                         │
              ┌──────────┼──────────┐
              │ 1        │          │ N
              ▼          │          ▼
    ┌─────────────────┐  │  ┌──────────┐
    │ Credenciales    │  │  │ Sesion   │
    │ Usuario         │  │  │          │
    └─────────────────┘  │  └──────────┘
                         │
                         │ (por ID, no por referencia)
                         │
              Todos conocen el usuarioID
              Ninguno importa al otro
```

### Reglas de acoplamiento (validado en cada code review)

| Regla | Violación | Consecuencia |
|-------|-----------|--------------|
| `Usuario` no importa `CredencialesUsuario` | ❌ | El dominio de usuario dependería del dominio de seguridad |
| `CredencialesUsuario` no importa `Usuario` | ❌ | Acoplamiento circular si Usuario también importa Credenciales |
| `Sesion` no importa `Usuario` | ❌ | Sesión no necesita todo el objeto Usuario, solo su ID |
| `Sesion` no importa `CredencialesUsuario` | ❌ | Sesión no necesita conocer credenciales del usuario |
| Servicios de aplicación importan TODOS los dominios | ✅ | Es su trabajo: orquestar las entidades |
| Repositorios importan su propio dominio + modelo | ✅ | Cada repositorio conoce su entidad y su modelo |

---

*Fin del documento — 6 sesiones de pair programming, ~128 tests, 0 bugs en producción.*
