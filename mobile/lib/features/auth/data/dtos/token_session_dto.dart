import 'package:freezed_annotation/freezed_annotation.dart';

part 'token_session_dto.freezed.dart';
part 'token_session_dto.g.dart';

/// DTO compartido por los tres endpoints que devuelven un par de tokens con
/// la misma forma exacta: `POST /auth/login`, `POST /auth/refresh` y
/// `POST /auth/switch-tenant` (cotejado contra `login_dto.go` y
/// `sesion_dto.go` en identidad — los tres structs Go son idénticos campo a
/// campo). Wire: snake_case.
@freezed
abstract class TokenSessionDto with _$TokenSessionDto {
  @JsonSerializable(fieldRename: FieldRename.snake)
  const factory TokenSessionDto({
    required String accessToken,
    required String refreshToken,
    required int expiresIn,
    required String tokenType,
    required String usuarioId,
    required String tenantId,
    required String rol,
  }) = _TokenSessionDto;

  factory TokenSessionDto.fromJson(Map<String, dynamic> json) =>
      _$TokenSessionDtoFromJson(json);
}
