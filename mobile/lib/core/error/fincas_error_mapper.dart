import 'package:dio/dio.dart';

import 'app_exception.dart';
import 'error_mapper.dart';

/// Interpreta errores del servicio `fincas`. Formato dominante
/// `{error, detalle}` (a veces solo `{error}` vía `gin.H`), pero algunos
/// handlers devuelven el body como string plano — se toleran ambos.
final class FincasErrorMapper extends BaseErrorMapper {
  const FincasErrorMapper();

  @override
  AppException mapStatus(int status, dynamic data, DioException cause) {
    if (data is Map) {
      final body = data.cast<String, dynamic>();
      final detalle = body['detalle'] as String?;
      final error = body['error'] as String?;
      final message = detalle ?? error ?? 'Error $status';
      return porStatus(status, message, cause: cause);
    }

    if (data is String && data.trim().isNotEmpty) {
      return porStatus(status, data.trim(), cause: cause);
    }

    return ServerException('Error $status del servidor', status, cause);
  }
}
