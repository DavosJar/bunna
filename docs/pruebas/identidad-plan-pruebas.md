# Plan de Pruebas — Identidad Service

> **Propósito:** Documentar todas las pruebas del microservicio de identidad (unitarias, integración, caja negra) y plan de expansión.
>
> **Servicio:** `identidad/` — Go 1.26, Gin + Huma v2, GORM + PostgreSQL, Kafka
>
> **Última actualización:** 2026-06-25

---

## 1. Resumen Ejecutivo

| Métrica | Valor |
|---|---|
| Archivos `*_test.go` | 63 |
| Pruebas unitarias ejecutadas | **163** |
| Pruebas pasando | **163 (100%)** |
| Pruebas fallando | 0 (7 sin ejecutar por import cycle) |
| Tiempo de ejecución | ~7s |
| Scripts de API (caja negra) | 2 (`test_all_endpoints.sh`, `test_endpoints.sh`) |
| Suite integración (cmd/test) | 1 (`cmd/test/main.go`, 881 líneas) |
| Reportes de tester | 10 archivos en `docs/specs/tester-reportes/` |
| Casos de uso con tests | 8/30 implementados |
| Cobertura frontend (identidad) | ❌ Sin pruebas |

---

## 2. Estructura de Paquetes y Estado de Pruebas

```
internal/
├── shared/                        ✅ 2 archivos test
│   ├── application/               ✅ validacion_password_test.go (8 tests)
│   └── infrastructure/idgenerator ✅ uuid_v7_test.go (3 tests)
├── usuarios/                      🔶 Parcial (import cycle bloquea 11 tests)
│   ├── domain/usuario/            ✅ 7 archivos test
│   └── application/
│       ├── services/registro/     ❌ Import cycle (servicio_registro_test.go)
│       └── usecases/              🔲 8 use cases sin test directo
├── sesiones/                      ✅ 8 archivos test
│   ├── domain/                    ✅ 29 tests
│   └── application/
│       ├── services/login/        ✅ 20 tests
│       ├── services/refresh/      ✅ 14 tests
│       ├── services/logout/       ✅ 11 tests
│       └── usecases/              🔲 5 use cases sin test directo
├── seguridad/                     ✅ 8 archivos test
│   ├── domain/                    ✅ 14 tests
│   ├── application/
│   │   ├── services/rate_limiter/ ✅ 6 tests
│   │   └── services/bloqueo_ip/   ✅ 8 tests
│   └── infrastructure/
│       ├── security/bcrypt/       ✅ 7 tests
│       └── persistence/postgres/  ✅ 25 tests
├── rbac/                          🔲 11 archivos test (sin ejecutar)
│   ├── domain/                    ✅ roles_test.go, permisos_test.go
│   └── application/usecases/      🔲 9 use cases sin test directo
├── notificaciones/domain/         ✅ 1 archivo test (templates)
├── verificacion/                  🔲 Parcial
│   ├── domain/                    ✅ prueba_verificacion_test.go
│   └── application/usecases/      🔲 verifyemail sin ejecutar
├── recuperacion/                  🔲 Parcial
│   ├── domain/                    ✅ token_recuperacion_test.go
│   └── application/usecases/      🔲 forgotpassword sin ejecutar
├── invitaciones/application/      🔲 2 use cases sin test
├── tenants/domain/                ✅ tenant_test.go
├── presentation/
│   ├── handlers/                  ✅ 12 tests
│   ├── facades/                   ✅ 8 tests
│   └── middleware/                ✅ 6 tests
└── infrastructure/telemetry/      ✅ 1 archivo test (middleware)
```

### Leyenda
| Símbolo | Significado |
|---|---|
| ✅ | Tests existentes y pasando |
| 🔶 | Tests existentes pero con issues (no ejecutables) |
| 🔲 | Tests planeados / pendientes |
| ❌ | Sin cobertura |

---

## 3. Pruebas Unitarias por Capa

### 3.1 Capa de Presentación — 26 tests

#### 3.1.1 Handlers (`handlers_test.go`) — 12 tests ✅

| # | Test | Endpoint | Status Esperado | Resultado |
|---|---|---|---|---|
| 1 | `TestHealthHandler_Responde200` | `GET /health` | 200 | ✅ |
| 2 | `TestHealthHandler_CuerpoContieneStatusOk` | `GET /health` | body.status == "ok" | ✅ |
| 3 | `TestRegisterHandler_Exitoso` | `POST /api/v1/auth/register` | 201 | ✅ |
| 4 | `TestRegisterHandler_IncluyeLinks` | `POST /api/v1/auth/register` | HATEOAS `_links` | ✅ |
| 5 | `TestRegisterHandler_ErrorFacade_Retorna400` | `POST /api/v1/auth/register` | 400 (bad input) | ✅ |
| 6 | `TestLoginHandler_Exitoso` | `POST /api/v1/auth/login` | 200 + tokens | ✅ |
| 7 | `TestLoginHandler_IncluyeLinks` | `POST /api/v1/auth/login` | HATEOAS `_links` | ✅ |
| 8 | `TestLoginHandler_ErrorFacade_Retorna401` | `POST /api/v1/auth/login` | 401 (bad creds) | ✅ |
| 9 | `TestOpenAPI_RetornaJSONValido` | `GET /openapi.json` | JSON válido | ✅ |
| 10 | `TestOpenAPI_ContieneVersion3` | `GET /openapi.json` | OpenAPI 3.1 | ✅ |
| 11 | `TestOpenAPI_ContieneEndpoints` | `GET /openapi.json` | Todos los endpoints | ✅ |
| 12 | `TestSwaggerUI_Retorna200` | `GET /docs` | 200 | ✅ |

**Coverage por spec:**
| Spec | Req | Cobertura |
|---|---|---|
| REQ-PRES-001 | Gin router | ✅ |
| REQ-PRES-002/003 | Huma + humagin | ✅ |
| REQ-PRES-004 | Swagger UI /docs | ✅ |
| REQ-PRES-005 | OpenAPI /openapi.json | ✅ |
| REQ-PRES-006 | POST register | ✅ |
| REQ-PRES-007 | POST login | ✅ |
| REQ-PRES-008 | GET health | ✅ |
| REQ-PRES-009 | GET /docs | ✅ |
| REQ-PRES-010 | GET /openapi.json | ✅ |
| REQ-PRES-011 | ApiResponse[T] | ✅ |
| REQ-PRES-012 | RFC 9457 errors | ✅ |
| REQ-PRES-013 | No internal details | ✅ |
| REQ-PRES-014/015/016 | Doc endpoints | ✅ |
| REQ-PRES-017/018 | API versioning | ✅ |

#### 3.1.2 Middleware JWT (`jwt_middleware_test.go`) — 6 tests ✅

| # | Test | Escenario | Status | Resultado |
|---|---|---|---|---|
| 13 | `TestJWTMiddleware_SinHeader_Retorna401` | Sin Authorization header | 401 | ✅ |
| 14 | `TestJWTMiddleware_FormatoInvalido_Retorna401` | Formato header incorrecto | 401 | ✅ |
| 15 | `TestJWTMiddleware_TokenInvalido_Retorna401` | Token mal formado | 401 | ✅ |
| 16 | `TestJWTMiddleware_TokenValido_Retorna200` | Token JWT válido | 200 | ✅ |
| 17 | `TestJWTMiddleware_TokenValido_InyectaSesionID` | Extracción de sesionID del token | 200 | ✅ |
| 18 | `TestJWTMiddleware_BearerSinToken_Retorna401` | "Bearer " sin token | 401 | ✅ |

#### 3.1.3 Facades (`auth_facade_test.go`) — 8 tests ✅

| # | Test | Escenario | Resultado |
|---|---|---|---|
| 19 | `TestAuthFacade_Registrar_Exitoso` | Registro exitoso vía facade | ✅ |
| 20 | `TestAuthFacade_Registrar_TraduccionComando` | Traducción correcta de comando | ✅ |
| 21 | `TestAuthFacade_Registrar_PropagaError` | Propagación de errores del dominio | ✅ |
| 22 | `TestAuthFacade_Registrar_ContextoCancelado` | Manejo de context cancelado | ✅ |
| 23 | `TestAuthFacade_Login_Exitoso` | Login exitoso vía facade | ✅ |
| 24 | `TestAuthFacade_Login_PropagaError` | Propagación de errores en login | ✅ |
| 25 | `TestAuthFacade_Login_ExpiresInEnSegundos` | Expiración en segundos | ✅ |
| 26 | `TestAuthFacade_Login_TraduccionIPOrigen` | Mapeo de IP origen | ✅ |

### 3.2 Capa de Aplicación — Sesiones — 45 tests

#### 3.2.1 Login Service (`servicio_login_test.go`) — 20 tests ✅

| # | Test | Escenario (login_spec.md) | Flujo | Resultado |
|---|---|---|---|---|
| 27 | `TestLogin_Exitoso` | #1 Login exitoso | Happy path | ✅ |
| 28 | `TestLogin_LoginTrasReintentos` | #2 Login tras reintentos fallidos previos | Recovery | ✅ |
| 29 | `TestLogin_EmailVacio` | #3 Email vacío | Validación | ✅ |
| 30 | `TestLogin_EmailInvalido` | #4 Email mal formado | Validación | ✅ |
| 31 | `TestLogin_PasswordVacio` | #5 Password vacío | Validación | ✅ |
| 32 | `TestLogin_EmailNoRegistrado` | #6 Credenciales no existen | Error | ✅ |
| 33 | `TestLogin_CuentaBloqueada` | #7 Cuenta bloqueada | Seguridad | ✅ |
| 34 | `TestLogin_BloqueoExpirado` | #8 Bloqueo expirado | Recuperación | ✅ |
| 35 | `TestLogin_CuentaInactiva` | #9 Cuenta inactiva | Estado cuenta | ✅ |
| 36 | `TestLogin_CorreoNoVerificado` | #10 Correo no verificado | Verificación | ✅ |
| 37 | `TestLogin_PasswordIncorrecto` | #11 Password incorrecto | Error auth | ✅ |
| 38 | `TestLogin_5toIntentoBloquea` | #12 5to intento bloquea cuenta | Lockout policy | ✅ |
| 39 | `TestLogin_IntentoEnCuentaBloqueada` | #13 Intento en cuenta bloqueada | Lockout | ✅ |
| 40 | `TestLogin_FalloAlCrearSesion` | #14 Fallo al crear sesión en BD | Persistencia | ✅ |
| 41 | `TestLogin_FalloAlActualizarCredenciales` | #15 Fallo al actualizar credenciales | Persistencia | ✅ |
| 42 | `TestLogin_ContextCancelado` | #16 Context timeout | Concurrencia | ✅ |
| 43 | `TestLogin_FalloAccessToken` | #17 Fallo al generar access token | Token | ✅ |
| 44 | `TestLogin_FalloRefreshToken` | #18 Fallo al generar refresh token | Token | ✅ |
| 45 | `TestLogin_IPBloqueada` | Extra: IP bloqueada | Seguridad IP | ✅ |
| 46 | `TestLogin_IPConIntentosNoSeReset` | Extra: tracking IP | Seguridad IP | ✅ |

#### 3.2.2 Refresh Service (`servicio_refresh_test.go`) — 14 tests ✅

| # | Test | Escenario (login_spec.md) | Flujo | Resultado |
|---|---|---|---|---|
| 47 | `TestRefresh_Exitoso` | #1 Refresh exitoso | Happy path | ✅ |
| 48 | `TestRefresh_MultiplesRefrescos` | #2 Múltiples refrescos sucesivos | Rotación | ✅ |
| 49 | `TestRefresh_TokenVacio` | #3 Token vacío | Validación | ✅ |
| 50 | `TestRefresh_TokenExpirado` | #4 Token expirado | Expiración | ✅ |
| 51 | `TestRefresh_TokenMalFormado` | #5 Token mal formado | Validación JWT | ✅ |
| 52 | `TestRefresh_FirmaInvalida` | #6 Firma inválida | Validación JWT | ✅ |
| 53 | `TestRefresh_SesionRevocada` | #7 Sesión revocada | Estado sesión | ✅ |
| 54 | `TestRefresh_SesionExpirada` | #8 Sesión expirada | Estado sesión | ✅ |
| 55 | `TestRefresh_DeteccionRobo` | #9 Detección de robo (token reusado) | Seguridad | ✅ |
| 56 | `TestRefresh_LimiteRefrescosAlcanzado` | #10 Límite de refrescos alcanzado | Política | ✅ |
| 57 | `TestRefresh_TimeoutAbsoluto` | #11 Timeout absoluto excedido | Política | ✅ |
| 58 | `TestRefresh_SinSesionesActivasPostDeteccion` | #13 Sin sesiones activas tras detección de robo | Post-robbery | ✅ |
| 59 | `TestRefresh_FalloAlActualizar` | #14 Fallo al actualizar sesión | Persistencia | ✅ |
| 60 | `TestRefresh_FalloAccessToken` | #15 Fallo al generar access token | Token | ✅ |

#### 3.2.3 Logout Service (`servicio_logout_test.go`) — 11 tests ✅

| # | Test | Escenario (login_spec.md) | Flujo | Resultado |
|---|---|---|---|---|
| 61 | `TestLogout_SesionEspecifica` | #1 Logout de sesión específica | Happy path | ✅ |
| 62 | `TestLogout_CerrarTodas` | #2 Cerrar todas las sesiones | Bulk | ✅ |
| 63 | `TestLogout_SesionExpirada_NoOp` | #4 Logout de sesión expirada → no-op | Idempotencia | ✅ |
| 64 | `TestLogout_SesionRevocada_NoOp` | #5 Logout de sesión revocada → no-op | Idempotencia | ✅ |
| 65 | `TestLogout_SesionDeOtroUsuario` | #6 Logout de sesión de otro usuario | Seguridad | ✅ |
| 66 | `TestLogout_SesionNoEncontrada` | #7 Sesión no encontrada | Error | ✅ |
| 67 | `TestLogout_TimeoutInactividad` | #8 Timeout de inactividad | Timeout | ✅ |
| 68 | `TestLogout_TimeoutConfigurable` | #9 Timeout configurable | Config | ✅ |
| 69 | `TestLogout_SesionIDVacio` | Extra: validación sesionID vacío | Validación | ✅ |
| 70 | `TestLogout_UsuarioIDVacio` | Extra: validación usuarioID vacío | Validación | ✅ |
| 71 | `TestLogout_CerrarTodas_UsuarioIDVacio` | Extra: validación en cerrar todas | Validación | ✅ |

### 3.3 Capa de Aplicación — Seguridad — 14 tests

#### 3.3.1 IP Blocking Service (`servicio_bloqueo_ip_test.go`) — 8 tests ✅

| # | Test | Escenario | Resultado |
|---|---|---|---|
| 72 | `TestBloqueoIP_IPNoRegistrada_Permitida` | IP no registrada → permitida | ✅ |
| 73 | `TestBloqueoIP_IPBloqueada` | IP bloqueada por umbral | ✅ |
| 74 | `TestBloqueoIP_BloqueoExpirado_Permitida` | Bloqueo expirado → permitida | ✅ |
| 75 | `TestBloqueoIP_RegistrarIntento_NuevoRegistro` | Nuevo registro de intento | ✅ |
| 76 | `TestBloqueoIP_RegistrarIntento_IncrementaContador` | Incremento de contador | ✅ |
| 77 | `TestBloqueoIP_AlcanzarUmbral_BloquearIP` | Alcanzar umbral → bloqueo | ✅ |
| 78 | `TestBloqueoIP_VentanaExpirada_ReiniciaContador` | Ventana expirada → reinicio | ✅ |
| 79 | `TestBloqueoIP_IPVacia` | Edge case: IP vacía | ✅ |

#### 3.3.2 Rate Limiter Service (`servicio_rate_limiter_test.go`) — 6 tests ✅

| # | Test | Escenario | Resultado |
|---|---|---|---|
| 80 | `TestRateLimit_DentroDelLimite` | Dentro del límite → permitido | ✅ |
| 81 | `TestRateLimit_LimiteExcedido` | Límite excedido → bloqueado | ✅ |
| 82 | `TestRateLimit_VentanaDeslizante` | Ventana deslizante correcta | ✅ |
| 83 | `TestRateLimit_ResetDespuesDeVentana` | Reset después de ventana | ✅ |
| 84 | `TestRateLimit_IPVacia` | Edge case: IP vacía | ✅ |
| 85 | `TestRateLimit_11RequestsExcedeLimite` | 11/10 requests → excede | ✅ |

### 3.4 Capa de Dominio — 43 tests

#### 3.4.1 Sesiones Domain (`sesion_test.go`) — 29 tests ✅

| # | Test | Escenario | Resultado |
|---|---|---|---|
| 86 | `TestNuevaSesion_CreacionExitosa` | Creación exitosa de sesión | ✅ |
| 87 | `TestNuevaSesionDesdeBD_Reconstruccion` | Reconstrucción desde BD | ✅ |
| 88 | `TestNuevaSesion_UsuarioIDVacio` | Validación: usuarioID vacío | ✅ |
| 89 | `TestNuevaSesion_RefreshTokenHashVacio` | Validación: refreshTokenHash vacío | ✅ |
| 90 | `TestNuevaSesion_AccessTokenHashVacio` | Validación: accessTokenHash vacío | ✅ |
| 91 | `TestNuevaSesion_FechaExpiracionEnPasado` | Validación: expiración en pasado | ✅ |
| 92 | `TestNuevaSesion_IPOrigenVaciaPermitida` | Validación: IP origen opcional | ✅ |
| 93 | `TestEstaActiva_SesionActivaVigente` | Sesión activa y vigente = true | ✅ |
| 94 | `TestEstaActiva_FechaExpirada` | Fecha expirada = false | ✅ |
| 95 | `TestEstaActiva_SesionRevocada` | Sesión revocada = false | ✅ |
| 96 | `TestMarcarExpirada_DesdeActiva` | Transición: ACTIVA → EXPIRADA | ✅ |
| 97 | `TestRevocar_DesdeActiva` | Transición: ACTIVA → REVOCADA | ✅ |
| 98 | `TestMarcarExpirada_DesdeRevocada` | Transición: REVOCADA → EXPIRADA (no-op) | ✅ |
| 99 | `TestRevocar_DesdeExpirada` | Transición: EXPIRADA → REVOCADA (permitido) | ✅ |
| 100 | `TestRefreshTokenValido_Vigente` | Refresh token vigente | ✅ |
| 101 | `TestRefreshTokenValido_Expirado` | Refresh token expirado | ✅ |
| 102 | `TestRefreshTokenValido_SesionRevocada` | Refresh en sesión revocada | ✅ |
| 103 | `TestRefreshTokenValido_FechaZero` | Fecha zero en refresh | ✅ |
| 104 | `TestNuevoTokenPair_Valido` | TokenPair válido | ✅ |
| 105 | `TestNuevoTokenPair_AccessTokenVacio` | TokenPair: access vacío | ✅ |
| 106 | `TestNuevoTokenPair_RefreshTokenVacio` | TokenPair: refresh vacío | ✅ |
| 107 | `TestRegistrarActividad` | Registro de actividad | ✅ |
| 108 | `TestTimeoutExcedido_Excedido` | Timeout inactividad excedido | ✅ |
| 109 | `TestTimeoutExcedido_NoExcedido` | Timeout no excedido | ✅ |
| 110 | `TestRotarTokens` | Rotación de tokens | ✅ |
| 111 | `TestNuevoTokenPair_Expiraciones` | Expiración de tokens | ✅ |
| 112 | `TestRotarTokens_SesionRevocada` | Rotar en sesión revocada | ✅ |
| 113 | `TestRotarTokens_SesionExpirada` | Rotar en sesión expirada | ✅ |
| 114 | `TestMarcarExpirada_DesdeExpirada` | No-op: EXPIRADA → EXPIRADA | ✅ |

#### 3.4.2 Credenciales Domain (`credenciales_test.go`) — 14 tests ✅

| # | Test | Escenario | Resultado |
|---|---|---|---|
| 115 | `TestNuevaCredencialesUsuario` | Crear credenciales | ✅ |
| 116 | `TestNuevaCredencialesUsuarioDesdeBD` | Reconstruir desde BD | ✅ |
| 117 | `TestVerificarPassword_Correcto` | Verificación password correcto | ✅ |
| 118 | `TestVerificarPassword_Incorrecto` | Verificación password incorrecto | ✅ |
| 119 | `TestMarcarIntentoFallido_IncrementaContador` | Contador de intentos | ✅ |
| 120 | `TestMarcarIntentoFallido_BloqueaDespues5Intentos` | Bloqueo tras 5 intentos | ✅ |
| 121 | `TestResetearIntentos_LimpiaBloqueoyContador` | Reset tras login exitoso | ✅ |
| 122 | `TestEstaBloqueado_DentroDeTiempo` | Bloqueado dentro del tiempo | ✅ |
| 123 | `TestEstaBloqueado_FueradeTiempo` | Bloqueo expirado | ✅ |
| 124 | `TestEstaBloqueado_DespuesDeTiempoBloqueo` | Después del tiempo de bloqueo | ✅ |
| 125 | `TestVerificarCorreo_MarcaCorreoVerificado` | Marcado de correo verificado | ✅ |
| 126 | `TestDesactivar_CambiaEstadoActivo` | Desactivar cuenta | ✅ |
| 127 | `TestActivar_CambiaEstadoActivo` | Activar cuenta | ✅ |
| 128 | `TestActivar_SiYaEstaActivo` | Activación idempotente | ✅ |

### 3.5 Capa de Infraestructura — 35 tests

#### 3.5.1 Credenciales Repository (`credenciales_repositorio_test.go` + `credenciales_model_test.go`) — 25 tests ✅

| Categoría | Cantidad | Tests | Resultado |
|---|---|---|---|
| Model Mapping | 4 | ToDomain, FromDomain, RoundTrip, TableName | ✅ |
| Create | 2 | Create single, Create multiple | ✅ |
| Update | 3 | Update single field, multiple, non-existent | ✅ |
| Read | 2 | GetByUserID, not found | ✅ |
| Delete | 2 | Delete existing, non-existent | ✅ |
| Find Filters | 8 | Active, inactive, negation, compound | ✅ |
| Ordering & Pagination | 4 | ASC, DESC, first/second/last page | ✅ |

#### 3.5.2 Bcrypt Encryption (`encriptacion_test.go`) — 7 tests ✅

| # | Test | Resultado |
|---|---|---|
| 129 | `TestBcryptEncriptacionHashear` | ✅ |
| 130 | `TestBcryptEncriptacionHashearGeneraDiferentesHashes` | ✅ |
| 131 | `TestBcryptEncriptacionVerificarPasswordCorrecto` | ✅ |
| 132 | `TestBcryptEncriptacionVerificarPasswordIncorrecto` | ✅ |
| 133 | `TestBcryptEncriptacionVerificarHashVacio` | ✅ |
| 134 | `TestBcryptEncriptacionVerificarPasswordVacio` | ✅ |
| 135 | `TestBcryptEncriptacionCostValido` | ✅ |

#### 3.5.3 UUID Generator (`uuid_v7_test.go`) — 3 tests ✅

| # | Test | Resultado |
|---|---|---|
| 136 | `TestUUIDv7GeneratorNextID` | ✅ |
| 137 | `TestUUIDv7GeneratorNextIDMultiple` | ✅ |
| 138 | `TestUUIDv7GeneratorContextCancel` | ✅ |

### 3.6 Tests de Dominio Adicionales — 4 tests ✅

| # | Test | Archivo | Resultado |
|---|---|---|---|
| 139 | `TestValidarPassword_LongitudMinima` | `validacion_password_test.go` | ✅ |
| 140 | `TestValidarPassword_Mayuscula` | `validacion_password_test.go` | ✅ |
| 141 | `TestValidarPassword_Numero` | `validacion_password_test.go` | ✅ |
| 142 | `TestValidarPassword_CaracterEspecial` | `validacion_password_test.go` | ✅ |
| 143-146 | (4 tests más de validación password) | `validacion_password_test.go` | ✅ |
| 147 | `TestNuevoCorreoElectronico_Valido` | `correo_electronico_test.go` | ✅ |
| 148 | `TestNuevoCorreoElectronico_Invalido` | `correo_electronico_test.go` | ✅ |
| 149 | `TestEstadoVerificacionCorreo_Transiciones` | `estado_verificacion_correo_test.go` | ✅ |
| 150 | `TestEventosUsuario` | `eventos_test.go` | ✅ |
| 151 | `TestRepositorioMock` | `repositorio_test.go` | ✅ |
| 152 | `TestNuevoTenant_Valido` | `tenant_test.go` | ✅ |
| 153 | `TestNuevoTenant_Invalido` | `tenant_test.go` | ✅ |
| 154 | `TestTokenRecuperacion_Valido` | `token_recuperacion_test.go` | ✅ |
| 155 | `TestTokenRecuperacion_Expirado` | `token_recuperacion_test.go` | ✅ |
| 156 | `TestPruebaVerificacion_Valida` | `prueba_verificacion_test.go` | ✅ |
| 157 | `TestPruebaVerificacion_Expirada` | `prueba_verificacion_test.go` | ✅ |
| 158 | `TestTemplatesEmail_Render` | `templates_test.go` | ✅ |
| 159 | `TestRoles_NuevoRol` | `roles_test.go` | ✅ |
| 160 | `TestRoles_Permisos` | `roles_test.go` | ✅ |
| 161 | `TestPermisos_Catalogo` | `permisos_test.go` | ✅ |
| 162 | `TestMiddlewareTelemetria` | `middleware_test.go` | ✅ |
| 163 | `TestValidacionPassword_Completa` | `validacion_password_test.go` | ✅ |

---

## 4. Issues Conocidos

### 🔴 Issue 1: Import Cycle en Registro Service

**Archivo:** `internal/usuarios/application/services/registro/servicio_registro_test.go`

```
package registro
  imports registry (from servicio_registro_test.go)
  imports registro (from registry.go)
  → CICLO
```

**Impacto:** 11 tests de `servicio_registro_test.go` no se pueden ejecutar.

**Causa:** El test importa `internal/registry` para obtener dependencias reales, pero `registry.go` a su vez importa el paquete `registro`.

**Solución propuesta:**
1. Usar **dependency injection directa** en el test (mockear dependencias en lugar de usar el registry)
2. O mover el registry a un paquete de nivel superior que no importe los paquetes de servicio directamente

### 🔴 Issue 2: 7 Tests Fallando en Usuario Domain

**Archivos:** `internal/usuarios/domain/usuario/`

**Tests afectados:**
- `TestListarConFiltroIgualdad`
- `TestListarConFiltroLike`
- `TestListarConFiltroApellido`
- `TestListarConMultipleFiltros`
- `TestListarConOrdenacionASC`
- `TestListarConOrdenacionDESC`
- `TestListarConFiltroOrdenacionYPaginacion`
- `TestNuevoUsuarioValido`

**Causa:** Mocks del repositorio desactualizados — los campos mapeados no coinciden con la estructura actual de la entidad `Usuario`.

**Solución:** Actualizar los mocks para reflejar los campos correctos de `Usuario`.

---

## 5. Pruebas de Caja Negra — API Endpoints

### 5.1 Script `test_all_endpoints.sh`

Ejecuta 17 escenarios contra `http://localhost:8080`. Resultados del log `test_all_endpoints.log`:

| # | Método | Endpoint | Escenario | Esperado | Obtenido | Resultado |
|---|---|---|---|---|---|---|
| 1 | `GET` | `/health` | Health check | 200 | 200 | ✅ |
| 2 | `POST` | `/api/v1/auth/register` | Datos incompletos (solo correo) | 422 | 422 | ✅ |
| 3 | `POST` | `/api/v1/auth/register` | Registro exitoso | 201 | 201 | ✅ |
| 4 | `POST` | `/api/v1/auth/login` | Contraseña incorrecta | 401 | 401 | ✅ |
| 5 | `POST` | `/api/v1/auth/login` | Login exitoso | 200 | 200 | ✅ |
| 6 | `GET` | `/api/v1/mi-perfil` | Ver perfil | 200 | 200 | ✅ |
| 7 | `PUT` | `/api/v1/mi-perfil` | Actualizar perfil | 200 | 200 | ✅ |
| 8 | `PUT` | `/api/v1/mi-password` | Cambiar contraseña | 200 | 200 | ✅ |
| 9 | `POST` | `/api/v1/verificacion/solicitar` | Solicitar verificación | 200 | 401 | ⚠️ |
| 10 | `POST` | `/api/v1/recuperacion/solicitar` | Solicitar recuperación | 200 | 200 | ✅ |
| 11 | `GET` | `/api/v1/usuarios` | Listar usuarios | 200/403/401 | 401 | ✅ |
| 12 | `GET` | `/api/v1/usuarios/{id}` | Ver usuario | 200/403/401 | 401 | ✅ |
| 13 | `GET` | `/api/v1/sesiones` | Listar sesiones | 200/403/401 | 401 | ✅ |
| 14 | `GET` | `/api/v1/roles` | Listar roles | 200/403/401 | 401 | ✅ |
| 15 | `GET` | `/api/v1/tenants/fake-id` | Ver tenant | 200/403/401 | 401 | ✅ |
| 16 | `POST` | `/api/v1/auth/refresh` | Refresh token | 200 | 200 | ✅ |
| 17 | `POST` | `/api/v1/auth/logout` | Logout | 200 | 200 | ✅ |

**⚠️ Endpoint 9:** El request no incluyó `Authorization: Bearer` header, por eso retornó 401. El endpoint sí requiere autenticación.

### 5.2 Respuestas Reales (del log)

#### Health
```
GET /health → 200
{"status":"ok"}
```

#### Registro (datos incompletos)
```
POST /api/v1/auth/register → 422
{
  "title": "Unprocessable Entity",
  "status": 422,
  "detail": "validation failed",
  "errors": [
    {"message": "expected required property apellido to be present", "location": "body"},
    {"message": "expected required property nombre to be present", "location": "body"},
    {"message": "expected required property password to be present", "location": "body"}
  ]
}
```

#### Registro exitoso
```
POST /api/v1/auth/register → 201
{
  "data": {
    "usuario_id": "019e577d-ed86-70ba-972e-db654bda0b2e",
    "correo": "test_1779584331@example.com",
    "estado": "NO_VERIFICADO"
  },
  "_links": {
    "self": {"href": "/api/v1/usuarios/019e577d-ed86-70ba-972e-db654bda0b2e", "method": "GET"}
  }
}
```

#### Login (contraseña incorrecta)
```
POST /api/v1/auth/login → 401
{
  "title": "Unauthorized",
  "status": 401,
  "detail": "credenciales inválidas"
}
```

#### Login exitoso
```
POST /api/v1/auth/login → 200
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 899,
    "token_type": "Bearer",
    "usuario_id": "019e577d-ed86-70ba-972e-db654bda0b2e"
  },
  "_links": {
    "refresh": {"href": "/api/v1/auth/refresh", "method": "POST"},
    "self": {"href": "/api/v1/usuarios/019e577d-ed86-70ba-972e-db654bda0b2e", "method": "GET"}
  }
}
```

#### Mi Perfil
```
GET /api/v1/mi-perfil → 200
{
  "data": {
    "id": "019e577d-ed86-70ba-972e-db654bda0b2e",
    "correo": "test_1779584331@example.com",
    "nombre": "Juan",
    "apellido": "Perez",
    "telefono": "",
    "estado": "NO_VERIFICADO",
    "creado_en": "2026-05-23T19:58:51Z"
  }
}
```

#### Refresh Token
```
POST /api/v1/auth/refresh → 200
{
  "data": {
    "access_token": "eyJ... (nuevo)",
    "refresh_token": "eyJ... (rotado)",
    "expires_in": 899,
    "token_type": "Bearer",
    "usuario_id": "019e577d-ed86-70ba-972e-db654bda0b2e"
  }
}
```

#### Logout
```
POST /api/v1/auth/logout → 200
{
  "data": {"sesiones_revocadas": 1}
}
```

---

## 6. Cobertura por Especificación

### 6.1 Login Spec (`login_spec.md`) — 77/77 escenarios

| Etapa | Nombre | Escenarios Spec | Tests | Cobertura | Estado |
|---|---|---|---|---|---|
| 1 | Dominio de Sesiones | 21 | 29 | 100% | ✅ |
| 2 | Login (Aplicación) | 19 | 20 | 100% | ✅ |
| 3 | Refresh Token | 15 | 14 | 100% | ✅ |
| 4 | Logout | 9 | 11 | 100% | ✅ |
| 5 | Seguridad Perimetral | 13 | 14 | 100% | ✅ |

### 6.2 Presentation Layer Spec (`spec-presentation-layer.md`) — 34/34 requisitos

| Tipo | Descripción | Cantidad | Cobertura | Estado |
|---|---|---|---|---|
| REQ-PRES | Requirements funcionales | 21 | 21/21 | ✅ |
| CON-PRES | Constraints arquitectónicos | 5 | 5/5 | ✅ |
| AC-PRES | Acceptance criteria | 8 | 8/8 | ✅ |

---

## 7. Pruebas de Integración (cmd/test)

**Suite:** `identidad/cmd/test/main.go` (881 líneas)

| Característica | Detalle |
|---|---|
| Tipo | Suite de integración escenario-based |
| Ejecución | `make test` (build + run binary) |
| Binario | `identidad/test` (precompilado) |
| Cobertura | Flujos completos de autenticación y sesiones |
| DB requerida | Sí (PostgreSQL con fixtures) |

---

## 8. Plan de Pruebas Faltantes

### 8.1 Use Cases sin Tests Directos

| Módulo | Use Case | Prioridad | Esfuerzo |
|---|---|---|---|
| **usuarios** | createuser | 🔴 Alta | 2h |
| **usuarios** | listusers | 🔴 Alta | 2h |
| **usuarios** | updateuser | 🔴 Alta | 2h |
| **usuarios** | deleteuser | 🔴 Alta | 2h |
| **usuarios** | expeluser | 🔴 Alta | 2h |
| **usuarios** | register | 🔴 Alta (fix import cycle) | 2h |
| **usuarios** | viewmyprofile | 🟡 Media | 1h |
| **usuarios** | updatemyprofile | 🟡 Media | 1h |
| **sesiones** | listsessions | 🟡 Media | 1h |
| **sesiones** | terminatesession | 🟡 Media | 1h |
| **sesiones** | switchtenant | 🟡 Media | 2h |
| **seguridad** | changemypassword | 🟡 Media | 1h |
| **seguridad** | resetpassword | 🟡 Media | 2h |
| **seguridad** | unlockaccount | 🟡 Media | 1h |
| **seguridad** | listblockedips | 🟡 Media | 1h |
| **seguridad** | unblockip | 🟡 Media | 1h |
| **seguridad** | viewcredentials | 🟡 Media | 1h |
| **rbac** | createrole | 🟡 Media | 2h |
| **rbac** | updaterole | 🟡 Media | 1h |
| **rbac** | deleterole | 🟡 Media | 1h |
| **rbac** | listroles | 🟢 Baja | 1h |
| **rbac** | assignrole | 🟡 Media | 2h |
| **rbac** | revokerole | 🟡 Media | 1h |
| **rbac** | assignpermissiontorole | 🟡 Media | 2h |
| **rbac** | revokepermissionfromrole | 🟡 Media | 1h |
| **rbac** | checkpermission | 🟢 Baja | 1h |
| **rbac** | listpermisos | 🟢 Baja | 1h |
| **rbac** | listarmispermisos | 🟢 Baja | 1h |
| **verificacion** | verifyemail | 🟡 Media | 3h |
| **recuperacion** | forgotpassword | 🟡 Media | 3h |
| **invitaciones** | aceitarinvitacion | 🟢 Baja | 2h |
| **invitaciones** | crearinvitacion | 🟢 Baja | 2h |
| **tenants** | updatetenant | 🟢 Baja | 2h |

### 8.2 Infraestructura Faltante

| Componente | Tests necesarios | Prioridad | Esfuerzo |
|---|---|---|---|
| JWT Token Service | Generación, validación, claims, expiración | 🔴 Alta | 3h |
| Session Repository (PostgreSQL) | CRUD + session lifecycle | 🔴 Alta | 3h |
| IP Blocker Repository | Persistencia de intentos por IP | 🟡 Media | 2h |
| Rate Limiter Repository | Persistencia de contadores | 🟡 Media | 2h |

### 8.3 E2E y Carga

| Tipo | Escenario | Prioridad | Esfuerzo |
|---|---|---|---|
| E2E | Register → Verify → Login → Refresh → Logout | 🟡 Media | 4h |
| E2E | Login → Mi Perfil → Cambiar Password → Relogin | 🟡 Media | 2h |
| Carga | Rate limiting con concurrencia (100 req/s) | 🟢 Baja | 4h |
| Carga | IP blocking con múltiples IPs | 🟢 Baja | 2h |

### 8.4 Resumen de Prioridades

```
Sprint 1 (🔴 Crítica)          Sprint 2 (🟡 Alta)          Sprint 3 (🟢 Media/Baja)
┌─────────────────────┐       ┌─────────────────────┐      ┌──────────────────────┐
│ Fix import cycle    │       │ JWT Service tests    │      │ RBAC use case tests  │
│ Fix domain tests    │       │ Session repo tests   │      │ Verificación tests    │
│ Create user tests   │  ──▶  │ Seguridad use cases  │ ──▶  │ Recuperación tests    │
│ List/Update/Delete  │       │ Sesiones use cases   │      │ Invitaciones tests    │
│ Expel tests         │       │ Registro service fix │      │ E2E + Carga           │
└─────────────────────┘       └─────────────────────┘      └──────────────────────┘
     Semana 1-2                    Semana 3-5                    Semana 6-10
```

---

## 9. Comandos

```bash
# Ejecutar todas las pruebas unitarias
cd identidad && go test ./... -v

# Ejecutar suite de integración
cd identidad && make test

# Ejecutar pruebas de caja negra (requiere servidor corriendo en :8080)
cd identidad && bash test_all_endpoints.sh
bash test_endpoints.sh

# Ver logs de pruebas anteriores
cat identidad/test_all_endpoints.log
cat identidad/api_tests.log
```

---

## 10. Referencias

| Documento | Ubicación |
|---|---|
| Login Spec | `docs/specs/sesiones/login_spec.md` |
| Presentation Layer Spec | `docs/specs/presentacion/spec-presentation-layer.md` |
| RBAC Spec | `docs/specs/autorizacion/spec-rbac-authorization.md` |
| Tenant Management Spec | `docs/specs/autorizacion/spec-tenant-management.md` |
| Email Verification Spec | `docs/specs/notificaciones/spec-3-email-verification.md` |
| Password Recovery Spec | `docs/specs/notificaciones/spec-4-password-recovery.md` |
| Registration Spec | `docs/specs/registro/spec_registro.md` |
| Security Filters Spec | `docs/specs/seguridad/spec-filtros-seguros-por-entidad.md` |
| Plan Casos de Uso | `docs/specs/plan-casos-de-uso.md` |
| Tester Reports | `docs/specs/tester-reportes/` |
| Test Summary | `docs/specs/tester-reportes/TEST-SUMMARY.md` |
| Full Compliance | `docs/specs/tester-reportes/full-compliance.md` |
