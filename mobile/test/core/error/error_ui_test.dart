import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/error/app_exception.dart';
import 'package:mobile/core/error/error_ui.dart';

/// Fija la regla de la mejora de UX: el 403 NUNCA muestra el detalle crudo del
/// backend ("permiso denegado"); se traduce a un mensaje claro y accionable.
/// En cambio, validación/conflicto sí conservan el mensaje del backend porque
/// es informativo.
void main() {
  test('403 (permiso): oculta el detalle crudo y da mensaje accionable', () {
    // Así llega desde FincasErrorMapper: message = detalle del backend.
    const e = ForbiddenException('permiso denegado');

    expect(e.esPermiso, isTrue);
    expect(e.tituloUsuario, 'No tienes permiso');
    expect(e.mensajeUsuario, isNot(contains('permiso denegado')));
    expect(e.mensajeUsuario.toLowerCase(), contains('permiso'));
    expect(e.mensajeUsuario.toLowerCase(), contains('administrador'));
  });

  test('validación (400/422): conserva el mensaje del backend', () {
    const e = ValidationException('El nombre ya existe', statusCode: 422);
    expect(e.esPermiso, isFalse);
    expect(e.mensajeUsuario, 'El nombre ya existe');
  });

  test('conflicto (409): conserva el mensaje del backend', () {
    const e = ConflictException('La finca ya está inactiva');
    expect(e.mensajeUsuario, 'La finca ya está inactiva');
  });

  test('red y servidor: mensajes genéricos claros, no detalle técnico', () {
    expect(
      const NetworkException('SocketException: failed').mensajeUsuario,
      isNot(contains('SocketException')),
    );
    expect(
      const ServerException('Error 503 del servidor').mensajeUsuario.toLowerCase(),
      contains('servidor'),
    );
  });
}
