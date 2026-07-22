// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'cambio_estado_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_CambioEstadoDto _$CambioEstadoDtoFromJson(Map<String, dynamic> json) =>
    _CambioEstadoDto(
      id: json['id'] as String,
      estado: json['estado'] as String,
      motivo: json['motivo'] as String?,
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );

Map<String, dynamic> _$CambioEstadoDtoToJson(_CambioEstadoDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'estado': instance.estado,
      'motivo': instance.motivo,
      'updatedAt': instance.updatedAt.toIso8601String(),
    };
