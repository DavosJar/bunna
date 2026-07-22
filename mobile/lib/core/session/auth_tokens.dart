import 'package:freezed_annotation/freezed_annotation.dart';

part 'auth_tokens.freezed.dart';
part 'auth_tokens.g.dart';

/// Par de tokens JWT persistido en Flutter Secure Storage. Es el ÚNICO dato
/// que la app guarda en disco (ver ARQUITECTURA.md, encabezado).
///
/// No es un DTO de wire: es el modelo interno de `core/session`, por eso usa
/// naming Dart plano (sin `FieldRename`) — el JSON que produce es el formato
/// de persistencia local, no un contrato con el backend.
@freezed
abstract class AuthTokens with _$AuthTokens {
  const factory AuthTokens({
    required String accessToken,
    required String refreshToken,
    required DateTime accessExpiresAt,
  }) = _AuthTokens;

  const AuthTokens._();

  factory AuthTokens.fromJson(Map<String, dynamic> json) =>
      _$AuthTokensFromJson(json);

  bool get isAccessExpired => DateTime.now().isAfter(accessExpiresAt);
}
