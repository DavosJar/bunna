// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'finca_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_FincaDto _$FincaDtoFromJson(Map<String, dynamic> json) => _FincaDto(
  id: json['id'] as String,
  nombre: json['nombre'] as String,
  ubicacion: json['ubicacion'] as String,
  descripcion: json['descripcion'] as String,
  estado: json['estado'] as String,
  createdAt: DateTime.parse(json['createdAt'] as String),
);

Map<String, dynamic> _$FincaDtoToJson(_FincaDto instance) => <String, dynamic>{
  'id': instance.id,
  'nombre': instance.nombre,
  'ubicacion': instance.ubicacion,
  'descripcion': instance.descripcion,
  'estado': instance.estado,
  'createdAt': instance.createdAt.toIso8601String(),
};
