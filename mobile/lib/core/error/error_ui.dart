import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import 'app_exception.dart';

/// Capa de presentación de errores: traduce un [AppException] (cuyo `message`
/// puede venir crudo del backend, p. ej. `"permiso denegado"`) a texto claro
/// para el usuario, con título e icono acordes.
///
/// Regla clave: para un 403 NUNCA mostramos el detalle técnico del backend —
/// el usuario no tiene por qué entender "permiso denegado"; le decimos qué
/// pasó y qué hacer.
extension AppExceptionUi on AppException {
  bool get esPermiso => this is ForbiddenException;

  /// Título corto para el banner.
  String get tituloUsuario => switch (this) {
    ForbiddenException() => 'No tienes permiso',
    SessionExpiredException() => 'Sesión expirada',
    NetworkException() => 'Sin conexión',
    ServerException() => 'Error del servidor',
    NotFoundException() => 'No encontrado',
    _ => 'No se pudo completar',
  };

  /// Mensaje claro y accionable. Para 403 ignoramos el detalle crudo del
  /// backend; para validación/conflicto sí lo mostramos porque es informativo.
  String get mensajeUsuario => switch (this) {
    ForbiddenException() =>
      'Tu cuenta no tiene permiso para realizar esta acción. Pídele a un '
          'administrador de la finca que te asigne el permiso necesario.',
    SessionExpiredException() =>
      'Tu sesión expiró. Inicia sesión nuevamente para continuar.',
    NetworkException() =>
      'No se pudo conectar con el servidor. Revisa tu conexión e inténtalo de '
          'nuevo.',
    ServerException() =>
      'Ocurrió un problema en el servidor. Inténtalo de nuevo en unos '
          'momentos.',
    // Validación (400/422), conflicto (409), rate limit, not found: el mensaje
    // del backend es específico y útil, se muestra tal cual.
    _ => message,
  };

  IconData get iconoUsuario => switch (this) {
    ForbiddenException() => Icons.lock_outline,
    SessionExpiredException() => Icons.lock_clock_outlined,
    NetworkException() => Icons.wifi_off_rounded,
    NotFoundException() => Icons.search_off_rounded,
    ServerException() => Icons.cloud_off_rounded,
    _ => Icons.error_outline_rounded,
  };
}

/// SnackBar temático para errores puntuales de una acción (p. ej. "tomar
/// muestra"), donde no cabe un banner persistente.
SnackBar appErrorSnackBar(AppException e) {
  final permiso = e.esPermiso;
  return SnackBar(
    behavior: SnackBarBehavior.floating,
    backgroundColor: permiso ? AppColors.cafe : AppColors.error,
    content: Row(
      children: [
        Icon(e.iconoUsuario, color: AppColors.blanco, size: 20),
        const SizedBox(width: 10),
        Expanded(
          child: Text(
            e.mensajeUsuario,
            style: const TextStyle(color: AppColors.blanco),
          ),
        ),
      ],
    ),
  );
}
