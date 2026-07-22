import 'package:dio/dio.dart';

import '../../session/token_store.dart';

/// Adjunta `Authorization: Bearer <token>` a toda request excepto las rutas
/// públicas (login, register, refresh, recuperación, verificación, health).
final class AuthInterceptor extends Interceptor {
  AuthInterceptor(this._tokenStore);

  final TokenStore _tokenStore;

  /// Fragmentos de path que identifican rutas públicas — compartido con
  /// `RefreshInterceptor` para que ambos apliquen exactamente la misma regla.
  static const publicPathFragments = <String>[
    '/auth/login',
    '/auth/register',
    '/auth/refresh',
    '/recuperacion/',
    '/verificacion/confirmar',
    '/health',
  ];

  static bool isPublicPath(String path) =>
      publicPathFragments.any(path.contains);

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) {
    if (!isPublicPath(options.path)) {
      final token = _tokenStore.current?.accessToken;
      if (token != null) {
        options.headers['Authorization'] = 'Bearer $token';
      }
    }
    handler.next(options);
  }
}
