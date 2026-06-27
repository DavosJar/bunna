# Plan General de Pruebas — Bunna

> **Propósito:** Documentar todas las pruebas existentes en el monorepo (unitarias, integración, caja negra), sus resultados, y establecer un plan de pruebas futuras.
>
> **Última actualización:** 2026-06-25

---

## 1. Resumen Ejecutivo

| Servicio | Lenguaje | Framework de pruebas | Archivos de prueba | Pruebas totales | Estado |
|---|---|---|---|---|---|
| **identidad** | Go 1.26 | `testing` + `testify` | ~63 `*_test.go` | 163+ | ✅ 163/163 PASS |
| **fincas** | Go 1.26 | `testing` + `testify` | 9 `*_test.go` | ~50+ | ✅ Ver detalle |
| **frontend** | JavaScript (React) | Ninguno | 0 | 0 | ❌ Sin cobertura |
| **hardware-monitor-agent** | Rust (2024) | `#[test]` | 1 `*_test.rs` | 1 | ✅ |
| **servicio-monitoreo** | Rust (2021) | `#[cfg(test)]` | 1 inline module | 2 | ✅ |
| **image-service** | Go 1.26 | `testing` | 0 | 0 | ❌ Sin cobertura |
| **YOLOv11** | Python | Ninguno | 0 | 0 | ❌ Sin cobertura |

### Totales del monorepo

| Métrica | Valor |
|---|---|
| Archivos de prueba | ~75 |
| Pruebas unitarias | ~215+ |
| Pruebas de integración (cmd/test) | 2 suites |
| Pruebas de caja negra (shell) | 2 scripts, ~20 escenarios |
| Documentos de testing | 10 reportes + 1 plan de uso |
| Cobertura frontend | 0% |

---

## 2. Pruebas Unitarias — Identidad Service

**Total: 163+ tests | 163 PASS | ~7s de ejecución**

### 2.1 Capa de Presentación

#### Handlers (`internal/presentation/handlers`) — 12 tests

| Test | Resultado | Ref. Spec | Propósito |
|---|---|---|---|
| `TestHealthHandler_Responde200` | ✅ | REQ-PRES-008 | Health endpoint retorna 200 |
| `TestHealthHandler_CuerpoContieneStatusOk` | ✅ | REQ-PRES-008 | Cuerpo health check estructura |
| `TestRegisterHandler_Exitoso` | ✅ | REQ-PRES-006, AC-PRES-003 | Register retorna 201 |
| `TestRegisterHandler_IncluyeLinks` | ✅ | CON-PRES-005 | HATEOAS links en respuesta |
| `TestRegisterHandler_ErrorFacade_Retorna400` | ✅ | AC-PRES-004 | Validación de entrada inválida |
| `TestLoginHandler_Exitoso` | ✅ | REQ-PRES-007, AC-PRES-005 | Login retorna tokens |
| `TestLoginHandler_IncluyeLinks` | ✅ | CON-PRES-005 | Login HATEOAS links |
| `TestLoginHandler_ErrorFacade_Retorna401` | ✅ | AC-PRES-006 | Credenciales incorrectas |
| `TestOpenAPI_RetornaJSONValido` | ✅ | REQ-PRES-005 | OpenAPI spec JSON válido |
| `TestOpenAPI_ContieneVersion3` | ✅ | REQ-PRES-005 | OpenAPI 3.1 format |
| `TestOpenAPI_ContieneEndpoints` | ✅ | REQ-PRES-014 | Todos los endpoints documentados |
| `TestSwaggerUI_Retorna200` | ✅ | REQ-PRES-004 | Swagger UI accesible |

#### Middleware JWT (`internal/presentation/middleware`) — 6 tests

| Test | Resultado | Propósito |
|---|---|---|
| `TestJWTMiddleware_SinHeader_Retorna401` | ✅ | Sin token → 401 |
| `TestJWTMiddleware_FormatoInvalido_Retorna401` | ✅ | Formato header inválido |
| `TestJWTMiddleware_TokenInvalido_Retorna401` | ✅ | Token mal formado |
| `TestJWTMiddleware_TokenValido_Retorna200` | ✅ | Token válido pasa |
| `TestJWTMiddleware_TokenValido_InyectaSesionID` | ✅ | Extrae sesionID del token |
| `TestJWTMiddleware_BearerSinToken_Retorna401` | ✅ | Validación esquema Bearer |

#### Facades (`internal/presentation/facades`) — 8 tests

| Test | Resultado | Ref. Spec |
|---|---|---|
| `TestAuthFacade_Registrar_Exitoso` | ✅ | CON-PRES-001 |
| `TestAuthFacade_Registrar_TraduccionComando` | ✅ | CON-PRES-001 |
| `TestAuthFacade_Registrar_PropagaError` | ✅ | CON-PRES-001 |
| `TestAuthFacade_Registrar_ContextoCancelado` | ✅ | CON-PRES-001 |
| `TestAuthFacade_Login_Exitoso` | ✅ | CON-PRES-001 |
| `TestAuthFacade_Login_PropagaError` | ✅ | CON-PRES-001 |
| `TestAuthFacade_Login_ExpiresInEnSegundos` | ✅ | REQ-PRES-011 |
| `TestAuthFacade_Login_TraduccionIPOrigen` | ✅ | CON-PRES-001 |

### 2.2 Capa de Aplicación — Sesiones

#### Login Service — 20 tests

| Test | Escenario Spec | Resultado |
|---|---|---|
| `TestLogin_Exitoso` | #1 Login exitoso | ✅ |
| `TestLogin_LoginTrasReintentos` | #2 Login tras reintentos | ✅ |
| `TestLogin_EmailVacio` | #3 Email vacío | ✅ |
| `TestLogin_EmailInvalido` | #4 Email mal formado | ✅ |
| `TestLogin_PasswordVacio` | #5 Password vacío | ✅ |
| `TestLogin_EmailNoRegistrado` | #6 Credenciales no existen | ✅ |
| `TestLogin_CuentaBloqueada` | #7 Cuenta bloqueada | ✅ |
| `TestLogin_BloqueoExpirado` | #8 Bloqueo expirado | ✅ |
| `TestLogin_CuentaInactiva` | #9 Cuenta inactiva | ✅ |
| `TestLogin_CorreoNoVerificado` | #10 Correo no verificado | ✅ |
| `TestLogin_PasswordIncorrecto` | #11 Password incorrecto | ✅ |
| `TestLogin_5toIntentoBloquea` | #12 5to intento bloquea | ✅ |
| `TestLogin_IntentoEnCuentaBloqueada` | #13 Intento en cuenta bloqueada | ✅ |
| `TestLogin_FalloAlCrearSesion` | #14 Fallo al crear sesión | ✅ |
| `TestLogin_FalloAlActualizarCredenciales` | #15 Fallo al actualizar credenciales | ✅ |
| `TestLogin_ContextCancelado` | #16 Context timeout | ✅ |
| `TestLogin_FalloAccessToken` | #17 Fallo token access | ✅ |
| `TestLogin_FalloRefreshToken` | #18 Fallo token refresh | ✅ |
| `TestLogin_IPBloqueada` | Seguridad IP bloqueada | ✅ |
| `TestLogin_IPConIntentosNoSeReset` | Seguridad tracking IP | ✅ |

#### Refresh Service — 14 tests

| Test | Escenario Spec | Resultado |
|---|---|---|
| `TestRefresh_Exitoso` | #1 Refresh exitoso | ✅ |
| `TestRefresh_MultiplesRefrescos` | #2 Múltiples refrescos | ✅ |
| `TestRefresh_TokenVacio` | #3 Token vacío | ✅ |
| `TestRefresh_TokenExpirado` | #4 Token expirado | ✅ |
| `TestRefresh_TokenMalFormado` | #5 Token mal formado | ✅ |
| `TestRefresh_FirmaInvalida` | #6 Firma inválida | ✅ |
| `TestRefresh_SesionRevocada` | #7 Sesión revocada | ✅ |
| `TestRefresh_SesionExpirada` | #8 Sesión expirada | ✅ |
| `TestRefresh_DeteccionRobo` | #9 Detección de robo (token reuse) | ✅ |
| `TestRefresh_LimiteRefrescosAlcanzado` | #10 Límite alcanzado | ✅ |
| `TestRefresh_TimeoutAbsoluto` | #11 Timeout absoluto | ✅ |
| `TestRefresh_SinSesionesActivasPostDeteccion` | #13 Sin sesiones tras robo | ✅ |
| `TestRefresh_FalloAlActualizar` | #14 Fallo al actualizar sesión | ✅ |
| `TestRefresh_FalloAccessToken` | #15 Fallo al generar token | ✅ |

#### Logout Service — 11 tests

| Test | Escenario Spec | Resultado |
|---|---|---|
| `TestLogout_SesionEspecifica` | #1 Logout sesión específica | ✅ |
| `TestLogout_CerrarTodas` | #2 Cerrar todas las sesiones | ✅ |
| `TestLogout_SesionExpirada_NoOp` | #4 Logout expirado → no-op | ✅ |
| `TestLogout_SesionRevocada_NoOp` | #5 Logout revocado → no-op | ✅ |
| `TestLogout_SesionDeOtroUsuario` | #6 Logout de otro usuario | ✅ |
| `TestLogout_SesionNoEncontrada` | #7 Sesión no encontrada | ✅ |
| `TestLogout_TimeoutInactividad` | #8 Timeout inactividad | ✅ |
| `TestLogout_TimeoutConfigurable` | #9 Timeout configurable | ✅ |
| `TestLogout_SesionIDVacio` | Validación: sessionID vacío | ✅ |
| `TestLogout_UsuarioIDVacio` | Validación: userID vacío | ✅ |
| `TestLogout_CerrarTodas_UsuarioIDVacio` | Validación: userID vacío (cerrar todas) | ✅ |

### 2.3 Capa de Aplicación — Seguridad

#### IP Blocking Service — 8 tests

| Test | Resultado |
|---|---|
| `TestBloqueoIP_IPNoRegistrada_Permitida` | ✅ |
| `TestBloqueoIP_IPBloqueada` | ✅ |
| `TestBloqueoIP_BloqueoExpirado_Permitida` | ✅ |
| `TestBloqueoIP_RegistrarIntento_NuevoRegistro` | ✅ |
| `TestBloqueoIP_RegistrarIntento_IncrementaContador` | ✅ |
| `TestBloqueoIP_AlcanzarUmbral_BloquearIP` | ✅ |
| `TestBloqueoIP_VentanaExpirada_ReiniciaContador` | ✅ |
| `TestBloqueoIP_IPVacia` | ✅ |

#### Rate Limiter Service — 6 tests

| Test | Resultado |
|---|---|
| `TestRateLimit_DentroDelLimite` | ✅ |
| `TestRateLimit_LimiteExcedido` | ✅ |
| `TestRateLimit_VentanaDeslizante` | ✅ |
| `TestRateLimit_ResetDespuesDeVentana` | ✅ |
| `TestRateLimit_IPVacia` | ✅ |
| `TestRateLimit_11RequestsExcedeLimite` | ✅ |

### 2.4 Capa de Dominio

#### Sesiones Domain — 29 tests (incluye Sesion + TokenPair)

| Test | Escenario Spec | Resultado |
|---|---|---|
| `TestNuevaSesion_CreacionExitosa` | #1 Creación exitosa | ✅ |
| `TestNuevaSesionDesdeBD_Reconstruccion` | #2 Reconstrucción desde BD | ✅ |
| `TestNuevaSesion_UsuarioIDVacio` | #3 usuarioID vacío | ✅ |
| `TestNuevaSesion_RefreshTokenHashVacio` | #4 refreshTokenHash vacío | ✅ |
| `TestNuevaSesion_AccessTokenHashVacio` | #5 accessTokenHash vacío | ✅ |
| `TestNuevaSesion_FechaExpiracionEnPasado` | #6 Fecha expiración en pasado | ✅ |
| `TestNuevaSesion_IPOrigenVaciaPermitida` | #7 IP origen opcional | ✅ |
| `TestEstaActiva_SesionActivaVigente` | #8 Sesión activa vigente | ✅ |
| `TestEstaActiva_FechaExpirada` | #9 Fecha expirada | ✅ |
| `TestEstaActiva_SesionRevocada` | #10 Sesión revocada | ✅ |
| `TestMarcarExpirada_DesdeActiva` | #11 Marcar expirada desde activa | ✅ |
| `TestRevocar_DesdeActiva` | #12 Revocar desde activa | ✅ |
| `TestMarcarExpirada_DesdeRevocada` | #13 Expirada desde revocada → no-op | ✅ |
| `TestRevocar_DesdeExpirada` | #14 Revocar desde expirada → permitido | ✅ |
| `TestRefreshTokenValido_Vigente` | #15 Refresh token vigente | ✅ |
| `TestRefreshTokenValido_Expirado` | #16 Refresh token expirado | ✅ |
| `TestRefreshTokenValido_SesionRevocada` | #17 Token en sesión revocada | ✅ |
| `TestRefreshTokenValido_FechaZero` | #18 Fecha zero | ✅ |
| `TestNuevoTokenPair_Valido` | #19 TokenPair válido | ✅ |
| `TestNuevoTokenPair_AccessTokenVacio` | #20 AccessToken vacío | ✅ |
| `TestNuevoTokenPair_RefreshTokenVacio` | #21 RefreshToken vacío | ✅ |
| `TestRegistrarActividad` | Registro de actividad | ✅ |
| `TestTimeoutExcedido_Excedido` | Timeout inactividad excedido | ✅ |
| `TestTimeoutExcedido_NoExcedido` | Timeout no excedido | ✅ |
| `TestRotarTokens` | Rotación de tokens | ✅ |
| `TestNuevoTokenPair_Expiraciones` | Expiración tokens | ✅ |
| `TestRotarTokens_SesionRevocada` | Rotar sesión revocada | ✅ |
| `TestRotarTokens_SesionExpirada` | Rotar sesión expirada | ✅ |
| `TestMarcarExpirada_DesdeExpirada` | No-op expirada → expirada | ✅ |

#### Seguridad Domain (Credenciales) — 14 tests

| Test | Resultado | Ref. Spec |
|---|---|---|
| `TestNuevaCredencialesUsuario` | ✅ | Etapa 2 login flow |
| `TestNuevaCredencialesUsuarioDesdeBD` | ✅ | Etapa 2 |
| `TestVerificarPassword_Correcto` | ✅ | Etapa 2 |
| `TestVerificarPassword_Incorrecto` | ✅ | Etapa 2 |
| `TestMarcarIntentoFallido_IncrementaContador` | ✅ | Etapa 2 |
| `TestMarcarIntentoFallido_BloqueaDespues5Intentos` | ✅ | Etapa 2 |
| `TestResetearIntentos_LimpiaBloqueoyContador` | ✅ | Etapa 2 |
| `TestEstaBloqueado_DentroDeTiempo` | ✅ | Etapa 2 |
| `TestEstaBloqueado_FueradeTiempo` | ✅ | Etapa 2 |
| `TestEstaBloqueado_DespuesDeTiempoBloqueo` | ✅ | Etapa 2 |
| `TestVerificarCorreo_MarcaCorreoVerificado` | ✅ | Etapa 2 |
| `TestDesactivar_CambiaEstadoActivo` | ✅ | Etapa 2 |
| `TestActivar_CambiaEstadoActivo` | ✅ | Etapa 2 |
| `TestActivar_SiYaEstaActivo` | ✅ | Etapa 2 |

### 2.5 Capa de Infraestructura

#### Credenciales Repository (PostgreSQL) — 25 tests

| Categoría | Cantidad | Tests | Resultado |
|---|---|---|---|
| Model Mapping | 4 | ToDomain, FromDomain, RoundTrip, TableName | ✅ 4/4 |
| Create Operations | 2 | Create single, Create multiple | ✅ 2/2 |
| Update Operations | 3 | Update single field, multiple fields, non-existent | ✅ 3/3 |
| Read Operations | 2 | Get by userID, not found | ✅ 2/2 |
| Delete Operations | 2 | Delete existing, non-existent | ✅ 2/2 |
| Find with Filters | 8 | Active, inactive, negation, compound filters | ✅ 8/8 |
| Ordering & Pagination | 4 | ASC, DESC, first/second/last page | ✅ 4/4 |

#### Bcrypt Encryption — 7 tests

| Test | Resultado |
|---|---|
| `TestBcryptEncriptacionHashear` | ✅ |
| `TestBcryptEncriptacionHashearGeneraDiferentesHashes` | ✅ |
| `TestBcryptEncriptacionVerificarPasswordCorrecto` | ✅ |
| `TestBcryptEncriptacionVerificarPasswordIncorrecto` | ✅ |
| `TestBcryptEncriptacionVerificarHashVacio` | ✅ |
| `TestBcryptEncriptacionVerificarPasswordVacio` | ✅ |
| `TestBcryptEncriptacionCostValido` | ✅ |

#### UUID Generator — 3 tests

| Test | Resultado |
|---|---|
| `TestUUIDv7GeneratorNextID` | ✅ |
| `TestUUIDv7GeneratorNextIDMultiple` | ✅ |
| `TestUUIDv7GeneratorContextCancel` | ✅ |

### 2.6 Capa de Aplicación — Usuarios

| Paquete | Archivo | Tests | Resultado |
|---|---|---|---|
| `createuser` | `usecase_test.go` | — | 🔲 Pendiente |
| `listusers` | `usecase_test.go` | — | 🔲 Pendiente |
| `updateuser` | `usecase_test.go` | — | 🔲 Pendiente |
| `deleteuser` | `usecase_test.go` | — | 🔲 Pendiente |
| `expeluser` | `usecase_test.go` | — | 🔲 Pendiente |
| `viewmyprofile` | `usecase_test.go` | — | 🔲 Pendiente |
| `updatemyprofile` | `usecase_test.go` | — | 🔲 Pendiente |
| `register` | `usecase_test.go` | — | 🔲 Pendiente |

### 2.7 Capa de Aplicación — RBAC

| Paquete | Archivo | Tests | Resultado |
|---|---|---|---|
| `checkpermission` | `usecase_test.go` | — | 🔲 Pendiente |
| `createrole` | `usecase_test.go` | — | 🔲 Pendiente |
| `updaterole` | `usecase_test.go` | — | 🔲 Pendiente |
| `deleterole` | `usecase_test.go` | — | 🔲 Pendiente |
| `listroles` | `usecase_test.go` | — | 🔲 Pendiente |
| `assignrole` | `usecase_test.go` | — | 🔲 Pendiente |
| `revokerole` | `usecase_test.go` | — | 🔲 Pendiente |
| `assignpermissiontorole` | `usecase_test.go` | — | 🔲 Pendiente |
| `revokepermissionfromrole` | `usecase_test.go` | — | 🔲 Pendiente |
| `listpermisos` | `usecase_test.go` | — | 🔲 Pendiente |
| `listarmispermisos` | `usecase_test.go` | — | 🔲 Pendiente |

### 2.8 Capa de Aplicación — Verificación, Recuperación, Invitaciones

| Paquete | Archivo | Tests | Resultado |
|---|---|---|---|
| `verifyemail` | `usecase_test.go` | — | 🔲 Pendiente |
| `forgotpassword` | `usecase_test.go` | — | 🔲 Pendiente |
| `aceptarinvitacion` | `usecase_test.go` | — | 🔲 Pendiente |
| `crearinvitacion` | `usecase_test.go` | — | 🔲 Pendiente |

### 2.9 Issue Conocidos

| Issue | Archivo | Detalle | Impacto |
|---|---|---|---|
| Import cycle | `servicio_registro_test.go` | Ciclo de importación con `registry` | ❌ No se pueden ejecutar 11 tests de registro |
| 7 tests fallando | `internal/usuarios/domain/usuario/` | Mocks desactualizados vs estructura actual | ❌ Filtros y paginación sin cobertura |

---

## 3. Pruebas Unitarias — Fincas Service

**Total: 9 archivos `*_test.go`**

### 3.1 Capa de Dominio

| Archivo | Tests | Resultado |
|---|---|---|
| `finca_test.go` | NuevaFinca, EsPropietario, TieneLotes, Actualizar, CambiarEstado, NewFincaFromPersistence | 🔲 |
| `lote_test.go` | NuevoLote, Actualizar, CambiarEstado, NewLoteFromPersistence | 🔲 |
| `finca_service_test.go` | EliminarFincaConLotes | 🔲 |
| `errores_test.go` | Tipos de error de dominio | 🔲 |

### 3.2 Capa de Presentación (Handlers)

| Archivo | Handler | Tests estimados |
|---|---|---|
| `finca_handler_test.go` | Registrar, Desactivar | ~6 |
| `lote_handler_test.go` | Agregar, Eliminar | ~6 |
| `muestra_handler_test.go` | Tomar, ListarPorLote | ~6 |
| `diagnostico_handler_test.go` | SolicitarManual, Aceptar, Rechazar | ~9 |
| `reporte_handler_test.go` | GenerarPorLote | ~3 |

**Nota:** Para obtener resultados detallados de ejecución, ejecutar:
```bash
cd fincas && go test ./... -v
```

---

## 4. Pruebas Unitarias — Rust Services

### 4.1 hardware-monitor-agent

| Archivo | Test | Resultado |
|---|---|---|
| `src/config_test.rs` | `test_config_from_env` — Verifica que `Config::from_env()` carga correctamente las variables de entorno | ✅ |

### 4.2 servicio-monitoreo

| Archivo | Test | Resultado |
|---|---|---|
| `src/sink/clickhouse.rs` (inline) | `test_escapar_tab_separated` — Escapa caracteres especiales en valores TSV | ✅ |
| `src/sink/clickhouse.rs` (inline) | `test_urlencoding` — Codifica URLs correctamente | ✅ |

---

## 5. Pruebas de Caja Negra — API Endpoints

### 5.1 Script `test_all_endpoints.sh`

Ejecuta 17 escenarios contra `http://localhost:8080` con registro detallado en `test_all_endpoints.log`.

| # | Endpoint | Método | Escenario | Status Esperado | Status Obtenido | Resultado |
|---|---|---|---|---|---|---|
| 1 | `/health` | GET | Health check | 200 | 200 | ✅ |
| 2 | `/api/v1/auth/register` | POST | Datos incompletos (solo correo) | 422 | 422 | ✅ |
| 3 | `/api/v1/auth/register` | POST | Registro exitoso | 201 | 201 | ✅ |
| 4 | `/api/v1/auth/login` | POST | Contraseña incorrecta | 401 | 401 | ✅ |
| 5 | `/api/v1/auth/login` | POST | Login exitoso | 200 | 200 | ✅ |
| 6 | `/api/v1/mi-perfil` | GET | Ver perfil | 200 | 200 | ✅ |
| 7 | `/api/v1/mi-perfil` | PUT | Actualizar perfil | 200 | 200 | ✅ |
| 8 | `/api/v1/mi-password` | PUT | Cambiar contraseña | 200 | 200 | ✅ |
| 9 | `/api/v1/verificacion/solicitar` | POST | Solicitar verificación | 200 | 401 | ⚠️ |
| 10 | `/api/v1/recuperacion/solicitar` | POST | Solicitar recuperación | 200 | 200 | ✅ |
| 11 | `/api/v1/usuarios` | GET | Listar usuarios (sin permisos) | 200/403/401 | 401 | ✅ |
| 12 | `/api/v1/usuarios/{id}` | GET | Ver usuario (sin permisos) | 200/403/401 | 401 | ✅ |
| 13 | `/api/v1/sesiones` | GET | Listar sesiones (sin permisos) | 200/403/401 | 401 | ✅ |
| 14 | `/api/v1/roles` | GET | Listar roles (sin permisos) | 200/403/401 | 401 | ✅ |
| 15 | `/api/v1/tenants/fake-id` | GET | Ver tenant (sin permisos) | 200/403/401 | 401 | ✅ |
| 16 | `/api/v1/auth/refresh` | POST | Refresh token | 200 | 200 | ✅ |
| 17 | `/api/v1/auth/logout` | POST | Logout | 200 | 200 | ✅ |

**⚠️ Endpoint 9:** `/api/v1/verificacion/solicitar` retornó 401 en lugar de 200. Esto se debe a que el script envió el request sin token en la prueba registrada (el log muestra que no se incluyó `Authorization` header).

### 5.2 Script `test_endpoints.sh` (Happy Path)

Ejecuta flujo completo: health → register → login → perfil → sesiones → roles → usuarios.

```
GET    /health                        → 200 ✅
POST   /api/v1/auth/register          → 201 ✅
GET    /api/v1/mi-perfil              → 200 ✅
PUT    /api/v1/mi-perfil              → 200 ✅
GET    /api/v1/sesiones               → 200 ✅
GET    /api/v1/roles                  → 200 ✅
GET    /api/v1/usuarios               → 200 ✅
```

### 5.3 Respuestas Reales (extracto de logs)

#### Registro exitoso
```json
POST /api/v1/auth/register → 201
Body: {
  "data": {
    "usuario_id": "019e577d-ed86-70ba-972e-db654bda0b2e",
    "correo": "test_1779584331@example.com",
    "estado": "NO_VERIFICADO"
  },
  "_links": {
    "self": { "href": "/api/v1/usuarios/019e577d-ed86-70ba-972e-db654bda0b2e", "method": "GET" }
  }
}
```

#### Login exitoso
```json
POST /api/v1/auth/login → 200
Body: {
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 899,
    "token_type": "Bearer",
    "usuario_id": "019e577d-ed86-70ba-972e-db654bda0b2e"
  },
  "_links": {
    "refresh": { "href": "/api/v1/auth/refresh", "method": "POST" },
    "self": { "href": "/api/v1/usuarios/...", "method": "GET" }
  }
}
```

#### Login con credenciales inválidas
```json
POST /api/v1/auth/login → 401
Body: {
  "title": "Unauthorized",
  "status": 401,
  "detail": "credenciales inválidas"
}
```

#### Refresh token
```json
POST /api/v1/auth/refresh → 200
Body: {
  "data": {
    "access_token": "eyJ... (nuevo)",
    "refresh_token": "eyJ... (nuevo, rotado)",
    "expires_in": 899,
    "token_type": "Bearer",
    "usuario_id": "019e577d-ed86-70ba-972e-db654bda0b2e"
  }
}
```

---

## 6. Pruebas de Integración — cmd/test

### 6.1 Fincas — `fincas/cmd/test/main.go`

- **Lenguaje:** Go (main package independiente)
- **Propósito:** Suite de pruebas de integración automatizada contra casos de uso reales
- **Tamaño:** ~1402 líneas
- **Ejecución:** `make run-test` o `go run ./cmd/test`
- **Cobertura:** Casos de uso de fincas (registrar, desactivar, lotes, muestras, diagnósticos, reportes)

### 6.2 Identidad — `identidad/cmd/test/main.go`

- **Lenguaje:** Go (main package independiente)
- **Propósito:** Suite de pruebas de integración escenario-based
- **Tamaño:** ~881 líneas
- **Ejecución:** `make test` (build + run)
- **Binario precompilado:** `identidad/test`
- **Cobertura:** Flujos completos de autenticación (register → verify → login → refresh → logout)

---

## 7. Frontend — Estado Actual

**El frontend NO tiene pruebas.** Esto es una deuda técnica crítica.

| Tipo | Estado | Observación |
|---|---|---|
| Unit tests (Vitest/Jest) | ❌ No existe | Sin framework configurado en `package.json` |
| Component tests (React Testing Library) | ❌ No existe | Sin setup |
| E2E (Playwright/Cypress) | ❌ No existe | Sin setup |
| API mock tests (MSW) | ❌ No existe | Sin setup |
| Lint | ✅ Configurado | ESLint con react-hooks y react-refresh |

**Archivos del frontend sin cobertura:**
- `src/pages/` — 6 páginas (auth, dashboard, landing, admin, perfil)
- `src/components/` — 7 componentes (Layout, Sidebar, Navbar, Toast, etc.)
- `src/services/` — 4 servicios (authApi, identidadApi, yoloApi, toastService)
- `src/context/` — 2 contextos (AuthContext, DiagnosisContext)
- `src/hooks/` — 1 hook (usePermisos)

---

## 8. Resumen de Cobertura vs Especificaciones

### 8.1 Login Spec (`login_spec.md`) — 77/77 escenarios

| Etapa | Nombre | Escenarios | Tests | Estado |
|---|---|---|---|---|
| 1 | Dominio de Sesiones | 21/21 | 29 | ✅ |
| 2 | Login (Aplicación) | 19/19 | 20 | ✅ |
| 3 | Refresh Token | 15/15 | 14 | ✅ |
| 4 | Logout | 9/9 | 11 | ✅ |
| 5 | Seguridad Perimetral | 13/13 | 14 | ✅ |

### 8.2 Presentation Layer Spec — 34/34 requisitos

| Tipo | Cantidad | Cobertura | Estado |
|---|---|---|---|
| REQ-PRES (Requirements) | 21 | 21/21 | ✅ |
| CON-PRES (Constraints) | 5 | 5/5 | ✅ |
| AC-PRES (Acceptance Criteria) | 8 | 8/8 | ✅ |

### 8.3 Fincas Use Cases Spec

| # | Caso de Uso | Tests Unitarios | Tests Integración | Estado |
|---|---|---|---|---|
| 1 | RegistrarFinca | ✅ handler tests | ✅ cmd/test | ✅ |
| 2 | DesactivarFinca | ✅ handler + domain tests | ✅ cmd/test | ✅ |
| 3 | AgregarLote | ✅ handler + domain tests | ✅ cmd/test | ✅ |
| 4 | EliminarLote | ✅ handler + domain tests | ✅ cmd/test | ✅ |
| 5 | TomarMuestra | ✅ handler tests | ✅ cmd/test | ✅ |
| 6 | ListarMuestrasPorLote | ✅ handler tests | ✅ cmd/test | ✅ |
| 7 | SolicitarDiagnosticoManual | ✅ handler tests | ✅ cmd/test | ✅ |
| 8 | RegistrarInferencia | ❌ Sin tests directos | ✅ cmd/test | 🔶 |
| 9 | AceptarDiagnostico | ✅ handler tests | ✅ cmd/test | ✅ |
| 10 | RechazarDiagnostico | ✅ handler tests | ✅ cmd/test | ✅ |
| 11 | GenerarReportePorLote | ✅ handler tests | ✅ cmd/test | ✅ |

---

## 9. Plan de Pruebas Futuras

### 9.1 Prioridad 🔴 CRÍTICA

| # | Tarea | Servicio | Esfuerzo estimado | Dependencias |
|---|---|---|---|---|
| 1 | **Configurar frontend testing** — Agregar Vitest + React Testing Library a `package.json` y crear primer test | frontend | 4h | — |
| 2 | **Test del servicio YOLO** — `POST /api/v1/diagnostico` con multipart form-data | frontend (yoloApi.js) | 2h | #1 |
| 3 | **Fix import cycle en registro** — Refactorizar `servicio_registro_test.go` | identidad | 2h | — |
| 4 | **Fix 7 tests de usuario domain** — Actualizar mocks en `internal/usuarios/domain/usuario/` | identidad | 3h | — |

### 9.2 Prioridad 🟡 ALTA

| # | Tarea | Servicio | Esfuerzo estimado |
|---|---|---|---|
| 5 | **Test JWT Token Service** — Generación, validación, expiración, claims | identidad | 3h |
| 6 | **Test Session Repository** — CRUD operaciones en PostgreSQL | identidad | 3h |
| 7 | **Test image-service** — Procesamiento de imágenes y publicación MQTT | image-service | 4h |
| 8 | **Test use cases de usuarios** — createuser, listusers, updateuser, deleteuser, expeluser | identidad | 6h |
| 9 | **Test use cases de RBAC** — Roles, permisos, asignaciones | identidad | 8h |

### 9.3 Prioridad 🟢 MEDIA

| # | Tarea | Servicio | Esfuerzo estimado |
|---|---|---|---|
| 10 | **Test use cases de sesiones** — listsessions, terminatesession | identidad | 3h |
| 11 | **Test use cases de seguridad** — lockaccount, unlockaccount, listblockedips, unblockip | identidad | 4h |
| 12 | **Test use cases de autogestión** — viewmyprofile, updatemyprofile, changemypassword | identidad | 3h |
| 13 | **Test de verificación de correo** — solicitar, confirmar, reenviar | identidad | 3h |
| 14 | **Test de recuperación de contraseña** — forgotpassword, resetpassword | identidad | 3h |
| 15 | **Test de invitaciones** — crear, aceptar | identidad | 2h |
| 16 | **Test de tenants** — updatetenant | identidad | 2h |
| 17 | **Test E2E: Register → Login → Refresh → Logout** | frontend + identidad | 6h |
| 18 | **Test E2E: YOLO diagnóstico flow** | frontend + image-service + YOLO | 8h |
| 19 | **Test E2E: Monitoreo (hardware → Kafka → ClickHouse)** | monitoreo | 6h |
| 20 | **Configurar CI/CD** — GitHub Actions con test runner para cada servicio | DevOps | 8h |

### 9.4 Prioridad 🔵 BAJA

| # | Tarea | Servicio | Esfuerzo estimado |
|---|---|---|---|
| 21 | **Test de componentes UI** — Layout, Sidebar, Navbar, Toast | frontend | 4h |
| 22 | **Test de páginas** — Login, Register, Dashboard | frontend | 6h |
| 23 | **Test de contextos** — AuthContext, DiagnosisContext | frontend | 3h |
| 24 | **Test de hooks** — usePermisos | frontend | 1h |
| 25 | **Test de validación de contraseña** — frontend | frontend | 1h |
| 26 | **Cobertura de edge cases en handlers** — CORS, errores | identidad + fincas | 3h |
| 27 | **Pruebas de carga** — Rate limiting, concurrencia | identidad | 6h |

### 9.5 Timeline sugerido

```
Sprint 1 (🔴 Crítica)          Sprint 2 (🟡 Alta)          Sprint 3 (🟢 Media)
┌─────────────────────┐       ┌─────────────────────┐      ┌─────────────────────┐
│ #1 Frontend setup   │       │ #5 JWT tests        │      │ #10-16 Use case     │
│ #2 YOLO test        │  ──▶  │ #6 Session repo     │ ──▶  │   tests faltantes   │
│ #3 Fix import cycle │       │ #7 image-service    │      │ #17-19 E2E flows    │
│ #4 Fix domain tests │       │ #8-9 Use cases      │      │ #20 CI/CD           │
└─────────────────────┘       └─────────────────────┘      └─────────────────────┘
     Semana 1-2                    Semana 3-4                   Semana 5-8
```

---

## 10. Herramientas y Comandos

| Servicio | Comando de prueba | Framework | Cobertura |
|---|---|---|---|
| **identidad** | `cd identidad && make test` | `testing` + `testify` | 163 tests |
| **fincas** | `cd fincas && make test` (unit) / `make run-test` (integration) | `testing` + `testify` | ~50 tests |
| **image-service** | `cd image-service && make test` | `testing` | 0 tests |
| **hardware-monitor-agent** | `cd hardware-monitor-agent && cargo test` | `#[test]` | 1 test |
| **servicio-monitoreo** | `cd servicio-monitoreo && cargo test` | `#[cfg(test)]` | 2 tests |
| **frontend** | ❌ No configurado | — | 0 tests |
| **YOLOv11** | ❌ No configurado | — | 0 tests |

---

## 11. Archivos de Prueba por Servicio

### identidad (~63 archivos)

```
internal/shared/infrastructure/idgenerator/uuid_v7_test.go
internal/shared/application/validacion_password_test.go
internal/usuarios/domain/usuario/usuario_test.go
internal/usuarios/domain/usuario/correo_electronico_test.go
internal/usuarios/domain/usuario/estado_verificacion_correo_test.go
internal/usuarios/domain/usuario/eventos_test.go
internal/usuarios/domain/usuario/repositorio_test.go
internal/usuarios/application/services/registro/servicio_registro_test.go
internal/usuarios/application/usecases/createuser/usecase_test.go
internal/usuarios/application/usecases/listusers/usecase_test.go
internal/usuarios/application/usecases/updateuser/usecase_test.go
internal/usuarios/application/usecases/deleteuser/usecase_test.go
internal/usuarios/application/usecases/expeluser/usecase_test.go
internal/usuarios/application/usecases/register/usecase_test.go
internal/usuarios/application/usecases/viewmyprofile/usecase_test.go
internal/usuarios/application/usecases/updatemyprofile/usecase_test.go
internal/sesiones/domain/sesion_test.go
internal/sesiones/application/services/login/servicio_login_test.go
internal/sesiones/application/services/refresh/servicio_refresh_test.go
internal/sesiones/application/services/logout/servicio_logout_test.go
internal/sesiones/application/usecases/login/usecase_test.go
internal/sesiones/application/usecases/refresh/usecase_test.go
internal/sesiones/application/usecases/logout/usecase_test.go
internal/sesiones/application/usecases/listsessions/usecase_test.go
internal/sesiones/application/usecases/terminatesession/usecase_test.go
internal/sesiones/application/usecases/switchtenant/usecase_test.go
internal/seguridad/domain/credenciales_test.go
internal/seguridad/application/services/rate_limiter/servicio_rate_limiter_test.go
internal/seguridad/application/services/bloqueo_ip/servicio_bloqueo_ip_test.go
internal/seguridad/application/usecases/changemypassword/usecase_test.go
internal/seguridad/application/usecases/resetpassword/usecase_test.go
internal/seguridad/application/usecases/unlockaccount/usecase_test.go
internal/seguridad/application/usecases/unblockip/usecase_test.go
internal/seguridad/application/usecases/viewcredentials/usecase_test.go
internal/seguridad/application/usecases/listblockedips/usecase_test.go
internal/seguridad/infrastructure/security/bcrypt/encriptacion_test.go
internal/seguridad/infrastructure/persistence/postgres/credenciales_model_test.go
internal/seguridad/infrastructure/persistence/postgres/credenciales_repositorio_test.go
internal/rbac/domain/roles_test.go
internal/rbac/domain/permisos_test.go
internal/rbac/application/usecases/checkpermission/usecase_test.go
internal/rbac/application/usecases/createrole/usecase_test.go
internal/rbac/application/usecases/updaterole/usecase_test.go
internal/rbac/application/usecases/deleterole/usecase_test.go
internal/rbac/application/usecases/listroles/usecase_test.go
internal/rbac/application/usecases/assignrole/usecase_test.go
internal/rbac/application/usecases/revokerole/usecase_test.go
internal/rbac/application/usecases/assignpermissiontorole/usecase_test.go
internal/rbac/application/usecases/revokepermissionfromrole/usecase_test.go
internal/rbac/application/usecases/listpermisos/usecase_test.go
internal/rbac/application/usecases/listarmispermisos/usecase_test.go
internal/notificaciones/domain/templates_test.go
internal/verificacion/domain/prueba_verificacion_test.go
internal/verificacion/application/usecases/verifyemail/usecase_test.go
internal/recuperacion/domain/token_recuperacion_test.go
internal/recuperacion/application/usecases/forgotpassword/usecase_test.go
internal/invitaciones/application/usecases/aceptarinvitacion/usecase_test.go
internal/invitaciones/application/usecases/crearinvitacion/usecase_test.go
internal/tenants/domain/tenant/tenant_test.go
internal/presentation/handlers/handlers_test.go
internal/presentation/facades/auth_facade_test.go
internal/presentation/middleware/jwt_middleware_test.go
internal/infrastructure/telemetry/middleware/middleware_test.go
```
