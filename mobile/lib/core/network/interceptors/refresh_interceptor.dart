import 'package:dio/dio.dart';

import '../../error/app_exception.dart';
import '../../session/token_refresh_coordinator.dart';
import 'auth_interceptor.dart';

/// Intercepta 401 y refresca el access token vía [TokenRefreshCoordinator]
/// antes de reintentar la request original UNA sola vez
/// (`extra['retried']`).
///
/// A propósito NO es `QueuedInterceptor`: si el reintento post-refresh
/// vuelve a fallar con 401 (p. ej. porque el 401 original no era por token
/// expirado sino por otra razón que persiste tras refrescar — caso real:
/// `/mis-permisos` devuelve 401 para usuarios `sys_admin` sin importar qué
/// tan fresco sea el token, porque su JWT no lleva claim `rol`), ese segundo
/// error necesita volver a pasar por `onError` de esta MISMA instancia
/// (`_dio.fetch(opts)` reutiliza `_dio`, con este interceptor en su cadena).
/// `QueuedInterceptor` serializa `onError` por instancia — la invocación
/// interna quedaría en cola detrás de la externa, que a su vez está
/// esperándola: deadlock. El single-flight real (evitar refrescos
/// concurrentes) ya lo garantiza `TokenRefreshCoordinator._inFlight`, no
/// esto — por eso `Interceptor` simple es suficiente y correcto.
final class RefreshInterceptor extends Interceptor {
  RefreshInterceptor({required TokenRefreshCoordinator coordinator, required Dio dio})
    : _coordinator = coordinator,
      _dio = dio;

  final TokenRefreshCoordinator _coordinator;
  final Dio _dio;

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    final requestOptions = err.requestOptions;

    final noEs401 = err.response?.statusCode != 401;
    final esRutaPublica = AuthInterceptor.isPublicPath(requestOptions.path);
    if (noEs401 || esRutaPublica) {
      handler.next(err);
      return;
    }

    if (requestOptions.extra['retried'] == true) {
      // Ya reintentamos una vez y volvió a fallar con 401: no hay más que
      // hacer, es una sesión definitivamente inválida.
      requestOptions.extra = {...requestOptions.extra, 'sessionExpired': true};
      handler.next(err);
      return;
    }

    final authHeader = requestOptions.headers['Authorization'] as String?;
    final failedAccessToken =
        authHeader?.replaceFirst('Bearer ', '') ?? '';

    try {
      final nuevosTokens = await _coordinator.refresh(
        failedAccessToken: failedAccessToken,
      );

      final opts = requestOptions
        ..headers['Authorization'] = 'Bearer ${nuevosTokens.accessToken}'
        ..extra = {...requestOptions.extra, 'retried': true};

      // Try/catch anidado a propósito: si el reintento vuelve a fallar, Dio
      // lanza un DioException (no una AppException — esa la lanza el
      // coordinator), y debe propagarse tal cual via `reject`, no quedar
      // sin capturar (lo que colgaría la request original para siempre).
      try {
        final reintento = await _dio.fetch(opts);
        handler.resolve(reintento);
      } on DioException catch (retryError) {
        handler.reject(retryError);
      }
    } on SessionExpiredException {
      requestOptions.extra = {...requestOptions.extra, 'sessionExpired': true};
      handler.next(err);
    } on NetworkException {
      // La sesión sigue intacta (no se tocaron los tokens); el caller
      // simplemente ve un fallo de red en su request original.
      handler.next(err);
    }
  }
}
