// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'mis_tenants_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_TenantConRolDto _$TenantConRolDtoFromJson(Map<String, dynamic> json) =>
    _TenantConRolDto(
      id: json['id'] as String,
      nombre: json['nombre'] as String,
      slug: json['slug'] as String,
      rol: json['rol'] as String,
      esPropio: json['es_propio'] as bool,
    );

Map<String, dynamic> _$TenantConRolDtoToJson(_TenantConRolDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'nombre': instance.nombre,
      'slug': instance.slug,
      'rol': instance.rol,
      'es_propio': instance.esPropio,
    };

_MisTenantsDto _$MisTenantsDtoFromJson(Map<String, dynamic> json) =>
    _MisTenantsDto(
      tenants: (json['tenants'] as List<dynamic>)
          .map((e) => TenantConRolDto.fromJson(e as Map<String, dynamic>))
          .toList(),
      propioId: json['propio_id'] as String,
    );

Map<String, dynamic> _$MisTenantsDtoToJson(_MisTenantsDto instance) =>
    <String, dynamic>{
      'tenants': instance.tenants.map((e) => e.toJson()).toList(),
      'propio_id': instance.propioId,
    };
