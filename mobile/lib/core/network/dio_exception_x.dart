import 'package:dio/dio.dart';

import '../error/app_exception.dart';

/// Puente entre Dio y el modelo de errores de la app: el `ErrorInterceptor`
/// deja la `AppException` ya calculada en `DioException.error`; esta
/// extensión es lo único que los datasources necesitan para extraerla sin
/// volver a tocar Dio.
extension DioExceptionToAppException on DioException {
  AppException toAppException() {
    final err = error;
    if (err is AppException) return err;
    return NetworkException('Error de red inesperado', this);
  }
}
