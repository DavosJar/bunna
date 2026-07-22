import 'package:freezed_annotation/freezed_annotation.dart';

part 'mis_tenants_dto.freezed.dart';
part 'mis_tenants_dto.g.dart';

/// GET /api/v1/identidad/tenants/mis-tenants — cotejado contra
/// `ListarMisTenantsResponse` / `TenantConRolDTO` en
/// `identidad/internal/presentation/dto/tenant_dto.go`.
@freezed
abstract class TenantConRolDto with _$TenantConRolDto {
  @JsonSerializable(fieldRename: FieldRename.snake)
  const factory TenantConRolDto({
    required String id,
    required String nombre,
    required String slug,
    required String rol,
    required bool esPropio,
  }) = _TenantConRolDto;

  factory TenantConRolDto.fromJson(Map<String, dynamic> json) =>
      _$TenantConRolDtoFromJson(json);
}

@freezed
abstract class MisTenantsDto with _$MisTenantsDto {
  @JsonSerializable(fieldRename: FieldRename.snake)
  const factory MisTenantsDto({
    required List<TenantConRolDto> tenants,
    required String propioId,
  }) = _MisTenantsDto;

  factory MisTenantsDto.fromJson(Map<String, dynamic> json) =>
      _$MisTenantsDtoFromJson(json);
}
