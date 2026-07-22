import 'package:freezed_annotation/freezed_annotation.dart';

part 'auth_session.freezed.dart';

/// Sesión activa: lo que se obtiene de login/refresh/switch-tenant, más el
/// `sesionId` decodificado del claim `sid` del access token.
@freezed
abstract class AuthSession with _$AuthSession {
  const factory AuthSession({
    required String usuarioId,
    required String tenantId,
    required String rol,
    String? sesionId,
  }) = _AuthSession;
}
