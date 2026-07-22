// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'token_session_dto.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_TokenSessionDto _$TokenSessionDtoFromJson(Map<String, dynamic> json) =>
    _TokenSessionDto(
      accessToken: json['access_token'] as String,
      refreshToken: json['refresh_token'] as String,
      expiresIn: (json['expires_in'] as num).toInt(),
      tokenType: json['token_type'] as String,
      usuarioId: json['usuario_id'] as String,
      tenantId: json['tenant_id'] as String,
      rol: json['rol'] as String,
    );

Map<String, dynamic> _$TokenSessionDtoToJson(_TokenSessionDto instance) =>
    <String, dynamic>{
      'access_token': instance.accessToken,
      'refresh_token': instance.refreshToken,
      'expires_in': instance.expiresIn,
      'token_type': instance.tokenType,
      'usuario_id': instance.usuarioId,
      'tenant_id': instance.tenantId,
      'rol': instance.rol,
    };
