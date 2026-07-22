// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'diagnostico_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_DiagnosticoDto _$DiagnosticoDtoFromJson(Map<String, dynamic> json) =>
    _DiagnosticoDto(
      id: json['id'] as String,
      muestraId: json['muestraID'] as String,
      estado: json['estado'] as String,
      tieneClorosis: json['tieneClorosis'] as bool,
      confianza: (json['confianza'] as num).toDouble(),
      imageUrl: json['imageURL'] as String?,
      imageBase64: json['imageBase64'] as String?,
    );

Map<String, dynamic> _$DiagnosticoDtoToJson(_DiagnosticoDto instance) =>
    <String, dynamic>{
      'id': instance.id,
      'muestraID': instance.muestraId,
      'estado': instance.estado,
      'tieneClorosis': instance.tieneClorosis,
      'confianza': instance.confianza,
      'imageURL': instance.imageUrl,
      'imageBase64': instance.imageBase64,
    };
