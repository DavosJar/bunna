import 'package:dio/dio.dart';

import '../error/app_exception.dart';
import 'auth_tokens.dart';
import 'session_events.dart';
import 'token_store.dart';

/// Serializa el refresh de tokens en todo el proceso, entre los clientes de
/// identidad y fincas a la vez (ARQUITECTURA.md §4). El refresh **rota**
/// ambos tokens en el backend: un segundo refresh concurrente con el
/// `refresh_token` ya invalidado mataría la sesión, así que aquí solo puede
/// haber un POST /auth/refresh en vuelo por proceso.
///
/// Usa un Dio "pelado", sin interceptores, para que la propia llamada de
/// refresh no pueda disparar el `RefreshInterceptor` (equivalente a la
/// exclusión `url.includes('/auth/refresh')` del interceptor axios web).
///
/// Nota de diseño: parsea el JSON de la respuesta a mano en vez de
/// reutilizar los DTOs de `features/auth/data` — `core/` no puede depender
/// de `features/` (regla de capas del documento de arquitectura), y el
/// refresh es infraestructura transversal, no una operación de negocio.
final class TokenRefreshCoordinator {
  TokenRefreshCoordinator({
    required TokenStore tokenStore,
    required SessionEvents sessionEvents,
    required String identidadBaseUrl,
    Dio? bareDioOverride,
  }) : _tokenStore = tokenStore,
       _sessionEvents = sessionEvents,
       _bareDio =
           bareDioOverride ??
           Dio(
             BaseOptions(
               baseUrl: identidadBaseUrl,
               connectTimeout: const Duration(seconds: 10),
               receiveTimeout: const Duration(seconds: 20),
             ),
           );

  static const _refreshPath = '/api/v1/identidad/auth/refresh';

  final TokenStore _tokenStore;
  final SessionEvents _sessionEvents;
  final Dio _bareDio;

  Future<AuthTokens>? _inFlight;

  /// [failedAccessToken] es el access token que causó el 401 que disparó
  /// este refresh. Se usa para distinguir "necesito refrescar" de "otro
  /// request ya refrescó mientras yo esperaba".
  Future<AuthTokens> refresh({required String failedAccessToken}) {
    final current = _tokenStore.current;

    // (a) Otro request ya rotó los tokens mientras este fallaba: no hace
    // falta pegarle a la red de nuevo, ya hay tokens nuevos disponibles.
    if (current != null && current.accessToken != failedAccessToken) {
      return Future.value(current);
    }

    // (b) Refresh ya en curso → engancharse (la "cola" del interceptor web).
    final inFlight = _inFlight;
    if (inFlight != null) return inFlight;

    // (c) Único refresh en vuelo.
    final future = _doRefresh();
    _inFlight = future;
    // `.ignore()`: el future derivado de `whenComplete` también carga el
    // error si `future` rechaza, y nadie más lo escucha — sin esto, Dart
    // lo reporta como un error async sin manejar (side-effect-only future).
    future.whenComplete(() => _inFlight = null).ignore();
    return future;
  }

  Future<AuthTokens> _doRefresh() async {
    final current = _tokenStore.current;
    if (current == null) {
      const ex = SessionExpiredException('No hay sesión activa');
      throw ex;
    }

    try {
      final response = await _bareDio.post<Map<String, dynamic>>(
        _refreshPath,
        data: {'refresh_token': current.refreshToken},
      );

      final data = response.data?['data'] as Map<String, dynamic>?;
      final accessToken = data?['access_token'] as String?;
      final refreshToken = data?['refresh_token'] as String?;
      final expiresIn = data?['expires_in'] as int?;
      if (data == null ||
          accessToken == null ||
          refreshToken == null ||
          expiresIn == null) {
        throw const ServerException('Respuesta de refresh con formato inesperado');
      }

      final nuevosTokens = AuthTokens(
        accessToken: accessToken,
        refreshToken: refreshToken,
        accessExpiresAt: DateTime.now().add(Duration(seconds: expiresIn)),
      );
      await _tokenStore.save(nuevosTokens);
      return nuevosTokens;
    } on DioException catch (e) {
      final status = e.response?.statusCode;
      if (status != null && status >= 400 && status < 500) {
        // El refresh token ya no sirve (expirado, revocado, rotado por otra
        // sesión, o la sesión server-side murió por inactividad): no hay
        // forma de recuperar la sesión, hay que desloguear.
        await _tokenStore.clear();
        _sessionEvents.emitExpired();
        throw SessionExpiredException(
          'Tu sesión expiró, inicia sesión nuevamente',
          e,
        );
      }
      // Fallo de red: NO se limpia la sesión, es reintentable.
      throw NetworkException('No se pudo conectar con el servidor', e);
    }
  }
}
