// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'perfil_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_PerfilDto _$PerfilDtoFromJson(Map<String, dynamic> json) => _PerfilDto(
  id: json['id'] as String,
  correo: json['correo'] as String,
  nombre: json['nombre'] as String,
  apellido: json['apellido'] as String,
  telefono: json['telefono'] as String,
  estado: json['estado'] as String,
  creadoEn: DateTime.parse(json['creado_en'] as String),
);

Map<String, dynamic> _$PerfilDtoToJson(_PerfilDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'correo': instance.correo,
      'nombre': instance.nombre,
      'apellido': instance.apellido,
      'telefono': instance.telefono,
      'estado': instance.estado,
      'creado_en': instance.creadoEn.toIso8601String(),
    };
