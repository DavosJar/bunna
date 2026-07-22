import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/error/app_exception.dart';
import 'package:mobile/core/error/identidad_error_mapper.dart';

DioException _errorConBody(
  int status,
  Object? body, {
  DioExceptionType type = DioExceptionType.badResponse,
}) {
  final requestOptions = RequestOptions(path: '/api/v1/identidad/usuarios');
  return DioException(
    requestOptions: requestOptions,
    type: type,
    response: Response(
      requestOptions: requestOptions,
      statusCode: status,
      data: body,
    ),
  );
}

void main() {
  const mapper = IdentidadErrorMapper();

  test('401 con extra sessionExpired ⇒ SessionExpiredException', () {
    final requestOptions = RequestOptions(path: '/api/v1/identidad/usuarios')
      ..extra['sessionExpired'] = true;
    final err = DioException(
      requestOptions: requestOptions,
      response: Response(requestOptions: requestOptions, statusCode: 401),
    );

    expect(mapper.map(err), isA<SessionExpiredException>());
  });

  test('sin response (timeout) ⇒ NetworkException', () {
    final err = DioException(
      requestOptions: RequestOptions(path: '/api/v1/identidad/auth/login'),
      type: DioExceptionType.connectionTimeout,
    );

    expect(mapper.map(err), isA<NetworkException>());
  });

  test('400 RFC 9457 con errors[] ⇒ ValidationException con fieldErrors', () {
    final err = _errorConBody(400, {
      'title': 'Bad Request',
      'status': 400,
      'detail': 'Datos inválidos',
      'errors': [
        {'field': 'correo', 'message': 'El correo no tiene un formato válido'},
      ],
    });

    final result = mapper.map(err);

    expect(result, isA<ValidationException>());
    final validation = result as ValidationException;
    expect(validation.message, 'Datos inválidos');
    expect(validation.statusCode, 400);
    expect(
      validation.fieldErrors['correo'],
      'El correo no tiene un formato válido',
    );
  });

  test(
    '422 sin errors[] (regla de negocio) ⇒ ValidationException con fieldErrors vacío',
    () {
      final err = _errorConBody(422, {
        'title': 'Unprocessable Entity',
        'status': 422,
        'detail': 'debes verificar tu correo electrónico antes de iniciar sesión',
      });

      final result = mapper.map(err) as ValidationException;

      expect(
        result.message,
        'debes verificar tu correo electrónico antes de iniciar sesión',
      );
      expect(result.fieldErrors, isEmpty);
    },
  );

  test('403 ⇒ ForbiddenException', () {
    final err = _errorConBody(403, {
      'title': 'Forbidden',
      'status': 403,
      'detail': 'no tienes permiso',
    });

    expect(mapper.map(err), isA<ForbiddenException>());
  });

  test('404 ⇒ NotFoundException', () {
    final err = _errorConBody(404, {
      'title': 'Not Found',
      'status': 404,
      'detail': 'usuario no encontrado',
    });

    expect(mapper.map(err), isA<NotFoundException>());
  });

  test('409 ⇒ ConflictException', () {
    final err = _errorConBody(409, {
      'title': 'Conflict',
      'status': 409,
      'detail': 'rol ya asignado',
    });

    expect(mapper.map(err), isA<ConflictException>());
  });

  test('429 ⇒ RateLimitException', () {
    final err = _errorConBody(429, {
      'title': 'Too Many Requests',
      'status': 429,
      'detail': 'demasiadas solicitudes',
    });

    expect(mapper.map(err), isA<RateLimitException>());
  });

  test('500 ⇒ ServerException', () {
    final err = _errorConBody(500, {
      'title': 'Internal Server Error',
      'status': 500,
      'detail': 'contacte al equipo de plataforma',
    });

    expect(mapper.map(err), isA<ServerException>());
  });

  test('body con formato inesperado (no Map) ⇒ mensaje genérico', () {
    final err = _errorConBody(500, '<html>502 Bad Gateway</html>');

    final result = mapper.map(err);

    expect(result, isA<ServerException>());
    expect(result.message, 'Error inesperado');
  });
}
