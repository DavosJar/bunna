// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'muestra_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_MuestraDto _$MuestraDtoFromJson(Map<String, dynamic> json) => _MuestraDto(
  id: json['id'] as String,
  fincaId: json['fincaID'] as String,
  loteId: json['loteID'] as String,
  latitud: (json['latitud'] as num).toDouble(),
  longitud: (json['longitud'] as num).toDouble(),
  createdAt: DateTime.parse(json['createdAt'] as String),
);

Map<String, dynamic> _$MuestraDtoToJson(_MuestraDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'fincaID': instance.fincaId,
      'loteID': instance.loteId,
      'latitud': instance.latitud,
      'longitud': instance.longitud,
      'createdAt': instance.createdAt.toIso8601String(),
    };
