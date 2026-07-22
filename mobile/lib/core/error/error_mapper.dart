import 'package:dio/dio.dart';

import 'app_exception.dart';

/// Traduce un `DioException` de un servicio concreto a `AppException`.
abstract interface class ErrorMapper {
  AppException map(DioException e);
}

/// Reglas comunes a ambos backends (ARQUITECTURA.md §3): conectividad,
/// ausencia de `response`, y expiración de sesión ya resuelta por el
/// `RefreshInterceptor`. Las subclases solo interpretan el body específico
/// de su servicio vía [mapStatus].
abstract class BaseErrorMapper implements ErrorMapper {
  const BaseErrorMapper();

  @override
  AppException map(DioException e) {
    if (e.requestOptions.extra['sessionExpired'] == true) {
      return SessionExpiredException(
        'Tu sesión expiró, inicia sesión nuevamente',
        e,
      );
    }

    if (_esErrorDeConectividad(e) || e.response == null) {
      return NetworkException('No se pudo conectar con el servidor', e);
    }

    final status = e.response!.statusCode ?? 0;

    // El RefreshInterceptor debe interceptar todo 401 antes de que llegue
    // aquí; esto es una red de seguridad defensiva, no el camino esperado.
    if (status == 401) {
      return SessionExpiredException(
        'Tu sesión expiró, inicia sesión nuevamente',
        e,
      );
    }

    return mapStatus(status, e.response!.data, e);
  }

  /// Interpreta el body específico del servicio para `status` dado.
  AppException mapStatus(int status, dynamic data, DioException cause);

  bool _esErrorDeConectividad(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.connectionError:
        return true;
      case DioExceptionType.unknown:
        return e.response == null;
      default:
        return false;
    }
  }

  /// Construye el subtipo de `AppException` correcto según el código de
  /// estado HTTP, con el `message` ya extraído del body por la subclase.
  AppException porStatus(
    int status,
    String message, {
    Map<String, String> fieldErrors = const {},
    Object? cause,
  }) {
    switch (status) {
      case 400:
      case 422:
        return ValidationException(
          message,
          fieldErrors: fieldErrors,
          statusCode: status,
          cause: cause,
        );
      case 403:
        return ForbiddenException(message, cause);
      case 404:
        return NotFoundException(message, cause);
      case 409:
        return ConflictException(message, cause);
      case 429:
        return RateLimitException(message, cause);
      default:
        return ServerException(message, status, cause);
    }
  }
}
