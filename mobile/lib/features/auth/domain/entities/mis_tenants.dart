import 'package:freezed_annotation/freezed_annotation.dart';

import 'tenant_con_rol.dart';

part 'mis_tenants.freezed.dart';

@freezed
abstract class MisTenants with _$MisTenants {
  const factory MisTenants({
    required List<TenantConRol> tenants,
    required String propioId,
  }) = _MisTenants;
}
