import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/error/app_exception.dart';
import 'package:mobile/core/error/fincas_error_mapper.dart';

DioException _errorConBody(int status, Object? body) {
  final requestOptions = RequestOptions(path: '/api/v1/fincas/fincas');
  return DioException(
    requestOptions: requestOptions,
    response: Response(
      requestOptions: requestOptions,
      statusCode: status,
      data: body,
    ),
  );
}

void main() {
  const mapper = FincasErrorMapper();

  test('sin response (connection refused) ⇒ NetworkException', () {
    final err = DioException(
      requestOptions: RequestOptions(path: '/api/v1/fincas/fincas'),
      type: DioExceptionType.connectionError,
    );

    expect(mapper.map(err), isA<NetworkException>());
  });

  test('{error, detalle} ⇒ usa detalle como mensaje', () {
    final err = _errorConBody(422, {
      'error': 'Validación fallida',
      'detalle': 'el nombre de la finca es requerido',
    });

    final result = mapper.map(err);

    expect(result, isA<ValidationException>());
    expect(result.message, 'el nombre de la finca es requerido');
  });

  test('solo {error} (gin.H) ⇒ usa error como mensaje', () {
    final err = _errorConBody(401, {'error': 'no autorizado'});

    // 401 sin marcar sessionExpired: el mapper lo trata como sesión
    // expirada de todas formas (red de seguridad defensiva).
    final result = mapper.map(err);
    expect(result, isA<SessionExpiredException>());
  });

  test('403 con {error, detalle} ⇒ ForbiddenException (nunca desloguea)', () {
    final err = _errorConBody(403, {
      'error': 'Forbidden',
      'detalle': 'permiso denegado',
    });

    final result = mapper.map(err);
    expect(result, isA<ForbiddenException>());
    expect(result.message, 'permiso denegado');
  });

  test('body como string plano ⇒ se usa tal cual como mensaje', () {
    final err = _errorConBody(500, 'el servicio de fincas no está disponible');

    final result = mapper.map(err);
    expect(result, isA<ServerException>());
    expect(result.message, 'el servicio de fincas no está disponible');
  });

  test('body vacío/nulo ⇒ mensaje genérico con status', () {
    final err = _errorConBody(500, null);

    final result = mapper.map(err);
    expect(result, isA<ServerException>());
    expect(result.message, 'Error 500 del servidor');
  });

  test('409 ⇒ ConflictException', () {
    final err = _errorConBody(409, {
      'error': 'conflicto',
      'detalle': 'el lote ya fue eliminado',
    });

    expect(mapper.map(err), isA<ConflictException>());
  });
}
