# Arquitectura Flutter — App móvil CafeScan/bunna

Diseño aprobado como base de implementación. Alcance: cliente móvil de los backends
existentes `identidad` (:8080) y `fincas` (:8082) + YOLO API externa. **Sin modo
offline**: la app siempre lee del backend en vivo; el único dato persistido
localmente es el par de tokens JWT en Flutter Secure Storage. Caché de datos de
negocio = solo memoria (providers Riverpod) durante la sesión activa.

Base URLs (emulador Android, configurables por `--dart-define`):
- Identidad: `http://10.0.2.2:8080`
- Fincas: `http://10.0.2.2:8082`
- YOLO: `https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com` (misma que usa la web)

## 0. Decisiones transversales

| Decisión | Elección | Justificación |
|---|---|---|
| Estado | Riverpod 3 con codegen (`@riverpod`) | AsyncValue integra carga/error/dato; autoDispose por defecto = "caché en memoria de sesión" exacto al requisito |
| Errores | Excepciones tipadas (`sealed class AppException`) + `AsyncValue.guard`, **no** Either/fpdart | Menos dependencias, integración natural con Riverpod, más simple de implementar correctamente |
| DTO vs Entidad | DTOs (Freezed+json_serializable) espejo exacto del wire en `data/`; entidades limpias camelCase en `domain/`; mappers `toDomain()` | Confina el caos de naming de fincas a la capa data |
| Clientes HTTP | 3 instancias Dio (identidad, fincas, yolo) vía providers | Base URLs, mappers de error y timeouts distintos por servicio |
| Refresh | `QueuedInterceptor` por cliente + `TokenRefreshCoordinator` compartido (single-flight) | Réplica del patrón `isRefreshing`+cola del frontend, generalizado a 2 clientes concurrentes |
| Tokens | `TokenStore`: copia en memoria (lectura síncrona para interceptores) respaldada por secure storage | Rotación atómica; interceptores no hacen I/O de disco por request |
| Casos de uso | Una clase por operación (`call()`), inyectadas por provider | Pedido explícito; espeja la estructura de los backends Go |
| Home móvil | `/fincas` para todos los roles (la web manda admin a `/admin`, pero el módulo admin está diferido en móvil) | Divergencia deliberada, documentada |
| 403 | **Nunca** desloguea (los permisos de fincas dependen de la sync Kafka y pueden estar vacíos con sesión válida) | Riesgo #6 del análisis |
| Bootstrap post-login | `mi-perfil` → `mis-tenants` → `mis-permisos` **en secuencia**, no en paralelo | Rate limit de identidad: 10 req/1s por IP |

## 1. Estructura de carpetas

```
mobile/lib/
├── main.dart                          # ProviderScope + carga inicial de TokenStore
├── app/
│   ├── app.dart                       # MaterialApp.router (tema, locale es)
│   ├── router/
│   │   ├── app_router.dart            # GoRouter: rutas + redirect por estado de sesión
│   │   ├── routes.dart                # constantes de paths/nombres
│   │   └── shell_scaffold.dart        # StatefulShellRoute: tabs Fincas | Diagnóstico | Perfil
│   └── theme/
│       └── app_theme.dart
├── core/
│   ├── config/
│   │   └── app_env.dart               # base URLs por --dart-define con defaults 10.0.2.2
│   ├── error/
│   │   ├── app_exception.dart         # sealed AppException + subtipos
│   │   ├── error_mapper.dart          # interfaz ErrorMapper + mapeo común de red/timeouts
│   │   ├── identidad_error_mapper.dart# RFC 9457 → AppException
│   │   └── fincas_error_mapper.dart   # {error, detalle} → AppException
│   ├── network/
│   │   ├── api_envelope.dart          # ApiEnvelope<T> — lee solo `data`, ignora links/_links
│   │   ├── dio_providers.dart         # identidadDioProvider, fincasDioProvider, yoloDioProvider
│   │   └── interceptors/
│   │       ├── auth_interceptor.dart      # adjunta Bearer; omite rutas públicas
│   │       ├── refresh_interceptor.dart   # QueuedInterceptor 401 → coordinator → retry
│   │       └── error_interceptor.dart     # DioException → AppException (usa el mapper del servicio)
│   ├── session/
│   │   ├── auth_tokens.dart           # {accessToken, refreshToken, expiresAt}
│   │   ├── token_store.dart           # memoria + flutter_secure_storage; save/clear/loadFromDisk
│   │   ├── token_refresh_coordinator.dart # single-flight de refresh, compartido entre clientes
│   │   ├── jwt_claims.dart            # decode base64 SIN verificar firma: sub, sid, tenant_id, rol, exp
│   │   └── session_events.dart        # StreamController<SessionEvent> (expired) → escucha authController
│   └── widgets/                       # LoadingView, ErrorRetryView, EmptyView, etc.
├── features/
│   ├── auth/
│   │   ├── data/
│   │   │   ├── dtos/                  # login_response_dto, perfil_dto, mis_tenants_dto, permiso_dto…
│   │   │   ├── auth_api.dart          # datasource remoto (identidadDio)
│   │   │   └── auth_repository_impl.dart
│   │   ├── domain/
│   │   │   ├── entities/              # auth_session, perfil, tenant_con_rol, permiso, mis_tenants
│   │   │   ├── auth_repository.dart   # contrato (§2)
│   │   │   └── usecases/              # login, restore_session, logout, switch_tenant,
│   │   │                              #   cargar_contexto_usuario (perfil+tenants+permisos secuencial)
│   │   └── presentation/
│   │       ├── auth_controller.dart   # AsyncNotifier<AuthState> keepAlive — estado global de sesión
│   │       ├── permisos.dart          # puedeProvider(codigo): espejo de usePermisos web (sys_admin ⇒ true)
│   │       └── pages/                 # splash_page, login_page
│   ├── fincas/                        # misma estructura data/domain/presentation
│   ├── lotes/
│   ├── muestras/
│   ├── diagnosticos/                  # YOLO (cámara) + vincular resultado + aceptar/rechazar
│   ├── reportes/                      # fase 5 — carpeta reservada
│   └── perfil/                        # fase ≥2: editar perfil, cambiar password
mobile/test/                           # espeja lib/: mappers, coordinator, repos (mocktail), widget login
```

Reglas de dependencia: `features/*` puede importar `core/`; `core/` jamás importa
`features/`; dentro de una feature, `presentation → domain ← data` (data implementa
contratos de domain). Entre features solo se importan contratos/entidades de
`domain` (p. ej. `diagnosticos` usa la entidad `Muestra`).

## 2. Contratos de repositorios (domain)

Convención: los repositorios devuelven `Future<T>` y **lanzan `AppException`**;
nunca exponen DTOs ni DioException. Cada método documenta su endpoint real.

```dart
// features/auth/domain/auth_repository.dart
abstract interface class AuthRepository {
  /// POST /api/v1/identidad/auth/login — persiste tokens en TokenStore.
  Future<AuthSession> login({required String correo, required String password});

  /// Arranque: lee tokens de secure storage; si el access expiró intenta un
  /// refresh vía coordinator. Devuelve null si no hay sesión recuperable.
  Future<AuthSession?> restoreSession();

  /// POST /auth/logout | /auth/logout/all (best-effort) y limpia TokenStore SIEMPRE.
  Future<void> logout({bool todasLasSesiones = false});

  /// POST /auth/switch-tenant — reemplaza AMBOS tokens (rotación) en TokenStore.
  Future<AuthSession> switchTenant(String tenantId);

  /// GET /mi-perfil
  Future<Perfil> getMiPerfil();

  /// GET /tenants/mis-tenants
  Future<MisTenants> getMisTenants();

  /// GET /mis-permisos
  Future<List<Permiso>> getMisPermisos();

  // Diferido a fases posteriores: register, recuperación y verificación de correo.
}
```

Entidades de auth: `AuthSession {usuarioId, tenantId, rol, sesionId}` (de la
respuesta de login + claims del JWT); `MisTenants {List<TenantConRol> tenants,
String propioId}`; `TenantConRol {id, nombre, slug, rol, esPropio}`;
`Permiso {codigo, nombre, descripcion, modulo}`; `Perfil {id, correo, nombre,
apellido, telefono, estado, creadoEn}`.

```dart
// features/fincas/domain/fincas_repository.dart
abstract interface class FincasRepository {
  /// GET /api/v1/fincas/fincas (el backend ya filtra por tenant del token)
  Future<List<Finca>> listar();

  /// POST /fincas
  Future<Finca> registrar({required String nombre, required String ubicacion, String descripcion = ''});

  /// POST /fincas/{id}/desactivar
  Future<CambioEstado> desactivar(String fincaId, {bool confirmar = true});
}

// features/lotes/domain/lotes_repository.dart
abstract interface class LotesRepository {
  /// GET /fincas/{fincaId}/lotes
  Future<List<Lote>> listarPorFinca(String fincaId);

  /// POST /fincas/{fincaId}/lotes
  Future<Lote> agregar(String fincaId, {required String nombre, required double area, String descripcion = ''});

  /// POST /lotes/{loteId}/eliminar
  Future<CambioEstado> eliminar(String loteId);
}

// features/muestras/domain/muestras_repository.dart
abstract interface class MuestrasRepository {
  /// loteId != null → GET /fincas/{fincaId}/lotes/{loteId}/muestras
  /// loteId == null → GET /fincas/{fincaId}/muestras
  Future<List<Muestra>> listar({required String fincaId, String? loteId});

  /// POST a las mismas rutas, body {latitud, longitud}
  Future<Muestra> tomar({
    required String fincaId,
    String? loteId,
    required double latitud,
    required double longitud,
  });
}

// features/diagnosticos/domain/diagnosticos_repository.dart
abstract interface class DiagnosticosRepository {
  /// POST YOLO /api/v1/diagnostico — multipart, campo 'archivo'
  Future<AnalisisYolo> analizarImagen(String rutaArchivo);

  /// POST /muestras/{muestraId}/diagnosticos/manual — body {imageURL}
  Future<SolicitudDiagnostico> solicitarManual(String muestraId, {required String imageUrl});

  /// POST /muestras/{muestraId}/diagnosticos/manual/resultado
  /// body {imageURL, tieneClorosis, confianza, procesadoAt(ISO8601 UTC)}
  Future<Diagnostico> guardarResultado(
    String muestraId, {
    required String imageUrl,
    required bool tieneClorosis,
    required double confianza,
    required DateTime procesadoAt,
  });

  /// POST /diagnosticos/{id}/aceptar
  Future<CambioEstado> aceptar(String diagnosticoId);

  /// POST /diagnosticos/{id}/rechazar — body {motivo}
  Future<CambioEstado> rechazar(String diagnosticoId, {String motivo = ''});
}

// features/reportes/domain/reportes_repository.dart  (fase 5, reservado)
abstract interface class ReportesRepository {
  /// GET /fincas/{fincaId}/lotes/{loteId}/reporte
  Future<ReporteLote> generar({required String fincaId, required String loteId});
}
```

Entidades compartidas: `CambioEstado {id, estado, motivo?, updatedAt}` (respuesta
`EstadoCambioResponse` de fincas). `AnalisisYolo {feedback, numDetections,
avgConfidence, detections[], imageBase64}` (respuesta YOLO). Estados como enums
Dart con valor `desconocido` de fallback (ver §5.7).

Casos de uso (uno por operación, wrapper fino sobre el repo): auth → `Login`,
`RestoreSession`, `Logout`, `SwitchTenant`, `CargarContextoUsuario`; fincas →
`ListarFincas`, `RegistrarFinca`, `DesactivarFinca`; lotes → `ListarLotes`,
`AgregarLote`, `EliminarLote`; muestras → `ListarMuestras`, `TomarMuestra`;
diagnósticos → `AnalizarImagen`, `VincularDiagnostico` (orquesta solicitarManual
+ guardarResultado, como hace la web), `AceptarDiagnostico`, `RechazarDiagnostico`.

**`CargarContextoUsuario`, corrección post-implementación:** `tenants` y
`permisos` degradan a vacío (`MisTenants(tenants: [], propioId: '')` /
`<Permiso>[]`) si su respectiva llamada lanza `AppException`, en vez de
abortar todo el bootstrap — igual que `fetchMisTenants`/`fetchMisPermisos` en
`AuthContext.jsx` del frontend web. Motivo confirmado contra el backend real:
`GET /mis-permisos` devuelve 401 para usuarios `sys_admin` (sin tenant, sin
claim `rol` en el JWT) con el MISMO token que `GET /mi-perfil` acepta
segundos antes — no es una sesión inválida, es una particularidad del
backend. Ver nota de deadlock en §4.

## 3. Manejo de errores unificado

### Modelo

```dart
sealed class AppException implements Exception {
  const AppException(this.message, {this.statusCode, this.cause});
  final String message;      // texto listo para UI, en español
  final int? statusCode;
  final Object? cause;       // DioException original, para logging
}

final class NetworkException extends AppException {}        // sin conexión, timeout, DNS, conn refused
final class SessionExpiredException extends AppException {} // 401 definitivo (refresh agotado)
final class ForbiddenException extends AppException {}      // 403 — NUNCA dispara logout
final class ValidationException extends AppException {      // 400 / 422
  final Map<String, String> fieldErrors;                    // campo → mensaje (de errors[] RFC 9457)
}
final class NotFoundException extends AppException {}       // 404
final class ConflictException extends AppException {}       // 409
final class RateLimitException extends AppException {}      // 429
final class ServerException extends AppException {}         // 5xx o body con formato inesperado
```

### Dónde se mapea

En un `ErrorInterceptor` por instancia Dio, construido con el `ErrorMapper` del
servicio. Así datasources y repositorios solo ven `AppException`; jamás un
`DioException` cruza a domain/presentation.

```dart
abstract interface class ErrorMapper {
  AppException map(DioException e);
}
```

Reglas comunes (ambos mappers, antes de mirar el body):
- `connectionTimeout | sendTimeout | receiveTimeout | connectionError | unknown sin response` → `NetworkException('No se pudo conectar con el servidor')`.
- `response == null` → `NetworkException`.
- Status → subtipo: 400/422→Validation, 403→Forbidden, 404→NotFound, 409→Conflict, 429→RateLimit, ≥500→Server.
- 401 solo llega aquí si el RefreshInterceptor ya se rindió (marca `extra['sessionExpired']`) → `SessionExpiredException`.

`IdentidadErrorMapper` (RFC 9457 `{title, status, detail, errors:[{field,message}]}`):
- `message = detail ?? title ?? 'Error inesperado'`.
- `errors[]` → `fieldErrors {field: message}` en `ValidationException` (422 sin
  `errors[]` = regla de negocio: cuenta bloqueada, correo no verificado, etc. —
  `fieldErrors` vacío, `detail` como mensaje).

`FincasErrorMapper` (`{error, detalle}`, a veces `gin.H{"error": ...}` o string plano):
- body Map → `message = detalle ?? error ?? 'Error {status}'`.
- body String no vacío → ese string como message.
- Cualquier otra cosa → `ServerException` genérica.

YOLO usa `FincasErrorMapper`-like mínimo (red/5xx; no tiene contrato de error estable).

### Presentación

Controllers usan `AsyncValue.guard`; la UI muestra `AppException.message`
directamente (ya viene en español). `ValidationException.fieldErrors` alimenta
errores inline de formularios. Política fija:
- `SessionExpiredException` → nunca se muestra como error de pantalla: el
  `authController` (que escucha `session_events`) ya transicionó a
  no-autenticado y GoRouter redirigió a `/login` con snackbar "Tu sesión expiró".
- `ForbiddenException` → banner/estado "No tienes permiso para esta acción",
  **sin** tocar la sesión.
- `RateLimitException` → snackbar "Demasiadas solicitudes, intenta en unos segundos".

## 4. Interceptores Dio y refresh con rotación

### Cadena por cliente (orden de registro — importa)

| Cliente | Interceptores |
|---|---|
| identidadDio | `AuthInterceptor` → `RefreshInterceptor` → `ErrorInterceptor(IdentidadErrorMapper)` (+ `LogInterceptor` solo debug) |
| fincasDio | `AuthInterceptor` → `RefreshInterceptor` → `ErrorInterceptor(FincasErrorMapper)` |
| yoloDio | `ErrorInterceptor` (sin auth: YOLO no recibe JWT; timeout receive 60 s por subida de fotos) |

Timeouts identidad/fincas: connect 10 s, receive 20 s.

`AuthInterceptor.onRequest`: lee `tokenStore.current` (memoria, síncrono) y
adjunta `Authorization: Bearer …`. Omite rutas públicas: `/auth/login`,
`/auth/register`, `/auth/refresh`, `/recuperacion/*`, `/verificacion/confirmar`,
`/health`.

### TokenRefreshCoordinator (núcleo anti-carrera, compartido entre AMBOS clientes)

Equivalente exacto del `isRefreshing` + `failedQueue` del frontend, elevado a
proceso completo porque aquí hay dos clientes que pueden recibir 401 a la vez:

```dart
class TokenRefreshCoordinator {
  Future<AuthTokens>? _inFlight;          // single-flight process-wide

  Future<AuthTokens> refresh({required String failedAccessToken}) {
    final current = _store.current;
    // (a) Otro request ya rotó los tokens mientras este fallaba → no refrescar de nuevo.
    if (current != null && current.accessToken != failedAccessToken) {
      return Future.value(current);
    }
    // (b) Refresh ya en curso → engancharse (la "cola" del frontend).
    if (_inFlight != null) return _inFlight!;
    // (c) Iniciar refresh único.
    _inFlight = _doRefresh().whenComplete(() => _inFlight = null);
    return _inFlight!;
  }

  Future<AuthTokens> _doRefresh() async {
    // Dio "pelado" SIN interceptores → imposible el bucle refresh→401→refresh.
    // POST /api/v1/identidad/auth/refresh {refresh_token}
    // Éxito: _store.save(nuevos)  ← guarda AMBOS tokens (rotación) memoria→disco, y retorna.
    // Fallo 4xx: _store.clear(); _sessionEvents.emitExpired(); throw SessionExpiredException.
    // Fallo de red: throw NetworkException (NO limpia tokens: reintentable).
  }
}
```

`RefreshInterceptor extends Interceptor` (**no** `QueuedInterceptor` — ver
nota de corrección abajo). El single-flight real entre clientes lo garantiza
por completo el coordinator, vía `_inFlight`:

```
onError(e):
  si status != 401                          → next(e)
  si e.requestOptions es ruta pública       → next(e)
  si extra['retried'] == true               → marcar sessionExpired, next(e)
  tokenUsado = header Authorization del request fallido
  try:
    nuevos = coordinator.refresh(failedAccessToken: tokenUsado)
    opts = e.requestOptions
      ..headers['Authorization'] = 'Bearer ${nuevos.accessToken}'
      ..extra['retried'] = true
    try:
      handler.resolve(await dio.fetch(opts))   // reintento único
    on DioException catch retryError:
      handler.reject(retryError)               // el reintento falló por su cuenta, no colgar la request
  catch SessionExpiredException:
    e.requestOptions.extra['sessionExpired'] = true
    next(e)                                  // ErrorInterceptor lo convierte
  catch NetworkException:
    next(e)                                  // el caller ve fallo de red, sesión intacta
```

> **Corrección post-implementación (fase 1):** el diseño original de este
> documento marcaba `RefreshInterceptor extends QueuedInterceptor` para que
> Dio serializara `onError` por instancia. Esto causaba un **deadlock real**:
> el reintento (`dio.fetch(opts)`) reutiliza la MISMA instancia de Dio, con
> el mismo `RefreshInterceptor` en su cadena; si ese reintento también recibe
> 401, necesita volver a pasar por `onError` de esa instancia — pero
> `QueuedInterceptor` encola esa invocación interna detrás de la externa, que
> a su vez está esperándola (`await`). Se reprodujo en la práctica con el
> usuario `sys_admin`: su JWT no lleva claim `rol` (no tiene tenant), y
> `GET /mis-permisos` trata `rol` vacío como "sin token" → 401 con el MISMO
> token que `GET /mi-perfil` acepta segundos antes — el refresh sale bien,
> pero el reintento vuelve a fallar sí o sí, disparando el deadlock. Como el
> single-flight de refresh ya lo garantiza `TokenRefreshCoordinator._inFlight`
> (no la serialización de Dio), `Interceptor` simple es correcto y suficiente.
> Ver también el fallback de `tenants`/`permisos` en
> `CargarContextoUsuario` (§2): un 401 aislado en esos dos endpoints, con un
> token que `mi-perfil` ya validó, no es una sesión inválida.

Garantías del diseño:
1. **Un solo refresh en vuelo en todo el proceso** (aunque identidad y fincas
   fallen 401 simultáneamente) — crítico porque el refresh **rota**: un segundo
   refresh con el token viejo mataría la sesión.
2. **Un solo reintento por request** (`extra['retried']`) — sin bucles.
3. La comparación `failedAccessToken != current` evita refrescar cuando otro
   request ya rotó (equivale a "resolver desde la cola" del frontend).
4. El request de refresh usa un Dio sin interceptores (equivale a la exclusión
   `url.includes('/auth/refresh')` del frontend).
5. Sesión server-side muerta (inactividad 30 min) ⇒ el refresh devuelve 401
   aunque el refresh token no haya expirado ⇒ ruta de expiración limpia.

Refresh reactivo únicamente (no proactivo) en fase 1; `restoreSession()` al
arrancar sí refresca si `exp` del access ya pasó (decodificado localmente, sin
verificar firma — igual que `parseJWT` del frontend). Nunca se hardcodea la
duración del token: se usa `expires_in`/claim `exp`.

## 5. Estrategia JSON (naming mixto de fincas)

1. **El wire format vive solo en DTOs** (`features/*/data/dtos/`). Entidades de
   domain: Dart camelCase puro, sin anotaciones JSON.
2. `build.yaml`: sin `field_rename` global (default `none`) y
   `explicit_to_json: true`. El naming se declara **siempre** por clase.
3. **DTOs de identidad** (snake_case consistente): anotación a nivel de factory
   Freezed —
   ```dart
   @freezed
   abstract class LoginResponseDto with _$LoginResponseDto {
     @JsonSerializable(fieldRename: FieldRename.snake)
     const factory LoginResponseDto({
       required String accessToken,    // ← access_token
       required String refreshToken,
       required int expiresIn,
       required String tokenType,
       required String usuarioId,      // ← usuario_id
       required String tenantId,
       required String rol,
     }) = _LoginResponseDto;
     factory LoginResponseDto.fromJson(Map<String, dynamic> json) => _$LoginResponseDtoFromJson(json);
   }
   ```
4. **DTOs de fincas, familia camelCase** (Finca/Lote/Muestra/Diagnóstico/Reporte):
   sin rename global; `@JsonKey(name:)` explícito solo donde el wire no coincide
   con lowerCamelCase Dart — en la práctica, los sufijos ID en mayúscula y URL:
   `@JsonKey(name: 'fincaID')`, `'loteID'`, `'muestraID'`, `'diagnosticoID'`,
   `'solicitudID'`, `'imageURL'`. (`createdAt`, `tieneClorosis`, `confianza`,
   `procesadoAt`, `areaTotal`, `radioMts`… coinciden solos.)
5. **DTOs de fincas, familia nodos** (snake_case: `finca_id`, `node_key`,
   `creado_en`, `tenant_id`, `lote_id`): `FieldRename.snake` a nivel factory,
   como identidad. (Módulo diferido, pero la regla queda fijada.)
6. **Regla dura para el implementador**: cada campo de cada DTO se coteja contra
   el struct Go correspondiente antes de generar — prohibido asumir por
   convención. Fuentes de verdad (solo lectura):
   `identidad/internal/presentation/dto/*.go`,
   `fincas/internal/presentation/dto/*.go`.
7. **Envelope**: `ApiEnvelope<T>` genérico (`genericArgumentFactories`) que lee
   solo la clave `data` (identidad usa `_links`, fincas `links`; se ignoran).
   Excepción documentada: `GET /nodos` viene doble-anidado (`data.data` +
   paginación) — se desenvuelve a mano en el datasource de nodos.
8. **Fechas**: ISO 8601 ↔ `DateTime` con el converter por defecto; al enviar
   (`procesadoAt`) siempre `DateTime.now().toUtc().toIso8601String()`.
9. **Estados** (`ACTIVO`, `PENDIENTE`, `ACEPTADO`, …): String en el DTO; el
   mapper `toDomain()` los convierte a enums Dart con fallback `desconocido`
   para no romper ante valores nuevos del backend.
10. **Errores RFC 9457 y `{error, detalle}`**: parseados a mano en los mappers de
   §3 (Map crudo), sin DTOs generados — son formatos pequeños y semiestables.

## 6. Dependencias (pubspec.yaml)

Los majors son vinculantes; el patch exacto se resuelve con `flutter pub add` al
implementar.

```yaml
environment:
  sdk: ^3.12.2

dependencies:
  flutter:
    sdk: flutter
  cupertino_icons: ^1.0.8
  # Estado
  flutter_riverpod: ^3.0.0
  riverpod_annotation: ^3.0.0
  # Navegación
  go_router: ^16.0.0
  # Red (FormData/multipart de YOLO incluido — sin dependencia extra)
  dio: ^5.9.0
  # Modelos inmutables + JSON
  freezed_annotation: ^3.0.0
  json_annotation: ^4.9.0
  # Persistencia de tokens (lo ÚNICO persistido)
  flutter_secure_storage: ^9.2.4
  # Formato de fechas/números en UI
  intl: ^0.20.0

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^6.0.0
  build_runner: ^2.4.15
  freezed: ^3.0.0
  json_serializable: ^6.9.0
  riverpod_generator: ^3.0.0
  custom_lint: ^0.7.0
  riverpod_lint: ^3.0.0
  mocktail: ^1.0.4
```

Deliberadamente fuera: paquete de JWT decode (util propio de ~15 líneas, igual
que `parseJWT` del frontend), fpdart/dartz (ver §0), cualquier base de datos
local (requisito: sin copia local de datos de negocio).

Previstas para fases posteriores (NO agregar aún): `geolocator` (fase 3 —
muestras GPS), `image_picker` (fase 4 — cámara), `flutter_map` + `latlong2`
(fase 5 — mapa de reporte), `cached_network_image` (imágenes de diagnósticos).

## 7. Navegación, estado de sesión y notas de implementación (fase 1)

**Rutas fase 1**: `/splash` (restore de sesión) · `/login` · shell con tabs:
`/fincas` (placeholder fase 1), `/diagnostico` (placeholder), `/perfil` (datos de
`mi-perfil` + selector de tenant si `misTenants.length > 1` + logout).
Fases 2–3 agregan `/fincas/:id` (lotes) y `/fincas/:id/lotes/:loteId` (muestras).

**AuthState** (`authController`, `AsyncNotifier` keepAlive):
`unknown` (splash) | `unauthenticated` | `authenticated {session, perfil, tenants, permisos}`.
GoRouter usa `refreshListenable` puenteado al provider; `redirect`: no autenticado
y ruta privada → `/login`; autenticado en `/login` o `/splash` → `/fincas`.
El controller escucha `session_events` para la expiración forzada.
`switchTenant` reejecuta `CargarContextoUsuario` completo (perfil/tenants/permisos).

**Permisos en UI**: helper `puede(codigo)` — `rol == 'sys_admin'` ⇒ true, si no
`permisos.contains(codigo)` — espejo de `usePermisos` web. Códigos relevantes de
fincas: `fincas:finca:crear`, `fincas:finca:desactivar`, `fincas:lote:crear`,
`fincas:lote:eliminar`, `fincas:muestra:crear`, `fincas:muestra:consultar`,
`fincas:diagnostico:solicitar|aceptar|rechazar`, `fincas:reporte:generar`.

**Caché en memoria**: providers de lectura `autoDispose` (default codegen);
mutaciones invalidan (`ref.invalidate(fincasProvider)`); pull-to-refresh =
`ref.refresh`. Nada de datos de negocio persiste al cerrar la app.

**Android (cambios dentro de ./mobile al implementar fase 1)**:
1. Permiso `INTERNET` en `android/app/src/main/AndroidManifest.xml` (el template
   solo lo trae en debug/profile).
2. Network security config permitiendo cleartext solo hacia `10.0.2.2` (Android
   9+ bloquea HTTP plano por defecto). No habilitar cleartext global.

**Tests mínimos exigidos por fase 1**: unit de ambos `ErrorMapper` (con bodies
reales de los backends), unit del `TokenRefreshCoordinator` (dos 401 concurrentes
⇒ un solo POST de refresh; token ya rotado ⇒ cero POSTs), unit de
`AuthRepositoryImpl` con mocktail, widget test del login (éxito y credenciales
inválidas).

**Verificación de fase 1 contra backend real** (docker-compose.dev.yml arriba):
login con usuario seed → tabs visibles → matar access token (esperar o borrar en
memoria) → cualquier request dispara refresh transparente → logout limpia secure
storage.
