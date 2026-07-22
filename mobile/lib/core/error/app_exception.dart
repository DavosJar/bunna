/// Modelo unificado de errores de la app. Ningún `DioException` ni error de
/// backend cruza más allá de `core/network` — todo llega a domain/presentation
/// ya convertido a uno de estos subtipos (ver ARQUITECTURA.md §3).
sealed class AppException implements Exception {
  const AppException(this.message, {this.statusCode, this.cause});

  /// Mensaje listo para mostrar en UI, en español.
  final String message;

  final int? statusCode;

  /// Excepción/objeto original, solo para logging — nunca se muestra en UI.
  final Object? cause;

  @override
  String toString() => '$runtimeType(status: $statusCode, message: $message)';
}

/// Sin conexión, timeout, DNS, connection refused, o sin `response` en absoluto.
final class NetworkException extends AppException {
  const NetworkException([
    String message = 'No se pudo conectar con el servidor',
    Object? cause,
  ]) : super(message, cause: cause);
}

/// 401 definitivo: el refresh ya se intentó y falló (token/sesión inválidos).
final class SessionExpiredException extends AppException {
  const SessionExpiredException([
    String message = 'Tu sesión expiró, inicia sesión nuevamente',
    Object? cause,
  ]) : super(message, statusCode: 401, cause: cause);
}

/// 403 — el usuario está autenticado pero no tiene el permiso requerido.
/// NUNCA debe disparar logout (ver riesgo #6 del análisis: fincas puede tener
/// la cache de permisos vacía por falta de sync de Kafka aunque el JWT sea
/// válido).
final class ForbiddenException extends AppException {
  const ForbiddenException([
    String message = 'No tienes permiso para realizar esta acción',
    Object? cause,
  ]) : super(message, statusCode: 403, cause: cause);
}

/// 400 / 422 con errores de campo (validación de formulario o regla de
/// negocio). `fieldErrors` puede venir vacío cuando el 422 es una regla de
/// negocio sin campo asociado (p. ej. "cuenta bloqueada").
final class ValidationException extends AppException {
  const ValidationException(
    super.message, {
    this.fieldErrors = const {},
    super.statusCode,
    super.cause,
  });

  final Map<String, String> fieldErrors;
}

final class NotFoundException extends AppException {
  const NotFoundException([
    String message = 'Recurso no encontrado',
    Object? cause,
  ]) : super(message, statusCode: 404, cause: cause);
}

final class ConflictException extends AppException {
  const ConflictException([
    String message = 'Conflicto con el estado actual del recurso',
    Object? cause,
  ]) : super(message, statusCode: 409, cause: cause);
}

final class RateLimitException extends AppException {
  const RateLimitException([
    String message = 'Demasiadas solicitudes, intenta en unos segundos',
    Object? cause,
  ]) : super(message, statusCode: 429, cause: cause);
}

/// 5xx, o body con formato inesperado que no se pudo interpretar.
final class ServerException extends AppException {
  const ServerException([
    String message = 'Error del servidor, intenta más tarde',
    int? statusCode,
    Object? cause,
  ]) : super(message, statusCode: statusCode, cause: cause);
}
