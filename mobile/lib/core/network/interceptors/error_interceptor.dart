import 'package:dio/dio.dart';

import '../../error/error_mapper.dart';

/// Último interceptor de la cadena: convierte el `DioException` crudo en una
/// `AppException` usando el `ErrorMapper` del servicio, y la deja en
/// `DioException.error` para que `DioExceptionToAppException` la extraiga.
/// A partir de aquí ningún código de dominio vuelve a ver un `DioException`
/// sin procesar.
final class ErrorInterceptor extends Interceptor {
  const ErrorInterceptor(this._mapper);

  final ErrorMapper _mapper;

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final appException = _mapper.map(err);
    handler.reject(
      DioException(
        requestOptions: err.requestOptions,
        response: err.response,
        type: err.type,
        error: appException,
        message: appException.message,
      ),
    );
  }
}
