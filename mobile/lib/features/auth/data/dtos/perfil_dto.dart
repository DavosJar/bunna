import 'package:freezed_annotation/freezed_annotation.dart';

part 'perfil_dto.freezed.dart';
part 'perfil_dto.g.dart';

/// GET /api/v1/identidad/mi-perfil — cotejado contra
/// `VerMiPerfilResponse` en `identidad/internal/presentation/dto/usuario_dto.go`.
@freezed
abstract class PerfilDto with _$PerfilDto {
  @JsonSerializable(fieldRename: FieldRename.snake)
  const factory PerfilDto({
    required String id,
    required String correo,
    required String nombre,
    required String apellido,
    required String telefono,
    required String estado,
    required DateTime creadoEn,
  }) = _PerfilDto;

  factory PerfilDto.fromJson(Map<String, dynamic> json) =>
      _$PerfilDtoFromJson(json);
}
