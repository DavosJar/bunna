import 'package:dio/dio.dart';

import 'app_exception.dart';
import 'error_mapper.dart';

/// Interpreta errores del servicio `identidad`, formato RFC 9457 Problem
/// Details: `{title, status, detail, errors: [{field, message}]}`.
final class IdentidadErrorMapper extends BaseErrorMapper {
  const IdentidadErrorMapper();

  @override
  AppException mapStatus(int status, dynamic data, DioException cause) {
    if (data is! Map) {
      return porStatus(status, 'Error inesperado', cause: cause);
    }

    final body = data.cast<String, dynamic>();
    final detail = body['detail'] as String?;
    final title = body['title'] as String?;
    final message = detail ?? title ?? 'Error inesperado';

    final fieldErrors = <String, String>{};
    final errors = body['errors'];
    if (errors is List) {
      for (final item in errors) {
        if (item is Map) {
          final field = item['field'] as String?;
          final msg = item['message'] as String?;
          if (field != null && msg != null) {
            fieldErrors[field] = msg;
          }
        }
      }
    }

    return porStatus(
      status,
      message,
      fieldErrors: fieldErrors,
      cause: cause,
    );
  }
}
