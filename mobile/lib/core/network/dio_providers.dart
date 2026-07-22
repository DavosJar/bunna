import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../config/app_env.dart';
import '../error/fincas_error_mapper.dart';
import '../error/identidad_error_mapper.dart';
import '../session/session_events.dart';
import '../session/token_refresh_coordinator.dart';
import '../session/token_store.dart';
import 'interceptors/auth_interceptor.dart';
import 'interceptors/error_interceptor.dart';
import 'interceptors/refresh_interceptor.dart';

part 'dio_providers.g.dart';

const _connectTimeout = Duration(seconds: 10);
const _receiveTimeout = Duration(seconds: 20);
const _yoloReceiveTimeout = Duration(seconds: 60);

/// Singleton de proceso: memoria + secure storage de los tokens JWT.
/// `main()` llama a `loadFromDisk()` antes de `runApp` (ver ARQUITECTURA.md
/// §1, `main.dart`).
@Riverpod(keepAlive: true)
TokenStore tokenStore(Ref ref) => TokenStore();

@Riverpod(keepAlive: true)
SessionEvents sessionEvents(Ref ref) {
  final events = SessionEvents();
  ref.onDispose(events.dispose);
  return events;
}

/// Single-flight de refresh compartido por AMBOS clientes Dio — es lo que
/// evita el doble refresh cuando identidad y fincas dan 401 a la vez.
@Riverpod(keepAlive: true)
TokenRefreshCoordinator tokenRefreshCoordinator(Ref ref) {
  return TokenRefreshCoordinator(
    tokenStore: ref.watch(tokenStoreProvider),
    sessionEvents: ref.watch(sessionEventsProvider),
    identidadBaseUrl: AppEnv.identidadBaseUrl,
  );
}

@Riverpod(keepAlive: true)
Dio identidadDio(Ref ref) {
  final dio = Dio(
    BaseOptions(
      baseUrl: AppEnv.identidadBaseUrl,
      connectTimeout: _connectTimeout,
      receiveTimeout: _receiveTimeout,
      headers: const {'Content-Type': 'application/json'},
    ),
  );

  dio.interceptors.addAll([
    AuthInterceptor(ref.watch(tokenStoreProvider)),
    RefreshInterceptor(
      coordinator: ref.watch(tokenRefreshCoordinatorProvider),
      dio: dio,
    ),
    const ErrorInterceptor(IdentidadErrorMapper()),
    if (kDebugMode)
      LogInterceptor(requestBody: false, responseBody: false),
  ]);

  return dio;
}

@Riverpod(keepAlive: true)
Dio fincasDio(Ref ref) {
  final dio = Dio(
    BaseOptions(
      baseUrl: AppEnv.fincasBaseUrl,
      connectTimeout: _connectTimeout,
      receiveTimeout: _receiveTimeout,
      headers: const {'Content-Type': 'application/json'},
    ),
  );

  dio.interceptors.addAll([
    AuthInterceptor(ref.watch(tokenStoreProvider)),
    RefreshInterceptor(
      coordinator: ref.watch(tokenRefreshCoordinatorProvider),
      dio: dio,
    ),
    const ErrorInterceptor(FincasErrorMapper()),
    if (kDebugMode)
      LogInterceptor(requestBody: false, responseBody: false),
  ]);

  return dio;
}

/// YOLO no recibe JWT (ver ARQUITECTURA.md §4) y no tiene contrato de error
/// estable, por eso reutiliza el mapper tolerante de fincas en vez de tener
/// uno propio.
@Riverpod(keepAlive: true)
Dio yoloDio(Ref ref) {
  final dio = Dio(
    BaseOptions(
      baseUrl: AppEnv.yoloBaseUrl,
      connectTimeout: _connectTimeout,
      receiveTimeout: _yoloReceiveTimeout,
    ),
  );

  dio.interceptors.addAll([
    const ErrorInterceptor(FincasErrorMapper()),
    if (kDebugMode)
      LogInterceptor(requestBody: false, responseBody: false),
  ]);

  return dio;
}
