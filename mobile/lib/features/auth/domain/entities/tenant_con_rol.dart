import 'package:freezed_annotation/freezed_annotation.dart';

part 'tenant_con_rol.freezed.dart';

@freezed
abstract class TenantConRol with _$TenantConRol {
  const factory TenantConRol({
    required String id,
    required String nombre,
    required String slug,
    required String rol,
    required bool esPropio,
  }) = _TenantConRol;
}
