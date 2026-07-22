import 'package:freezed_annotation/freezed_annotation.dart';

part 'permiso_dto.freezed.dart';
part 'permiso_dto.g.dart';

/// GET /api/v1/identidad/mis-permisos — cotejado contra `misPermisoItem` /
/// `listarMisPermisosData` en
/// `identidad/internal/presentation/handlers/mis_permisos_handler.go`.
@freezed
abstract class PermisoDto with _$PermisoDto {
  @JsonSerializable(fieldRename: FieldRename.snake)
  const factory PermisoDto({
    required String codigo,
    required String nombre,
    required String descripcion,
    required String modulo,
  }) = _PermisoDto;

  factory PermisoDto.fromJson(Map<String, dynamic> json) =>
      _$PermisoDtoFromJson(json);
}

@freezed
abstract class MisPermisosDto with _$MisPermisosDto {
  @JsonSerializable(fieldRename: FieldRename.snake)
  const factory MisPermisosDto({required List<PermisoDto> permisos}) =
      _MisPermisosDto;

  factory MisPermisosDto.fromJson(Map<String, dynamic> json) =>
      _$MisPermisosDtoFromJson(json);
}
