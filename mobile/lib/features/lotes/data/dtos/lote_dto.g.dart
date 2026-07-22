// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'lote_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_LoteDto _$LoteDtoFromJson(Map<String, dynamic> json) => _LoteDto(
  id: json['id'] as String,
  fincaId: json['fincaID'] as String,
  nombre: json['nombre'] as String,
  area: (json['area'] as num).toDouble(),
  descripcion: json['descripcion'] as String,
  estado: json['estado'] as String,
  createdAt: DateTime.parse(json['createdAt'] as String),
);

Map<String, dynamic> _$LoteDtoToJson(_LoteDto instance) => <String, dynamic>{
  'id': instance.id,
  'fincaID': instance.fincaId,
  'nombre': instance.nombre,
  'area': instance.area,
  'descripcion': instance.descripcion,
  'estado': instance.estado,
  'createdAt': instance.createdAt.toIso8601String(),
};
