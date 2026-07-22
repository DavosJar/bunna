// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'permiso_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_PermisoDto _$PermisoDtoFromJson(Map<String, dynamic> json) => _PermisoDto(
  codigo: json['codigo'] as String,
  nombre: json['nombre'] as String,
  descripcion: json['descripcion'] as String,
  modulo: json['modulo'] as String,
);

Map<String, dynamic> _$PermisoDtoToJson(_PermisoDto instance) =>
    <String, dynamic>{
      'codigo': instance.codigo,
      'nombre': instance.nombre,
      'descripcion': instance.descripcion,
      'modulo': instance.modulo,
    };

_MisPermisosDto _$MisPermisosDtoFromJson(Map<String, dynamic> json) =>
    _MisPermisosDto(
      permisos: (json['permisos'] as List<dynamic>)
          .map((e) => PermisoDto.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$MisPermisosDtoToJson(_MisPermisosDto instance) =>
    <String, dynamic>{
      'permisos': instance.permisos.map((e) => e.toJson()).toList(),
    };
