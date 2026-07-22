import 'package:freezed_annotation/freezed_annotation.dart';

import '../../domain/entities/diagnostico.dart';

part 'diagnostico_dto.freezed.dart';
part 'diagnostico_dto.g.dart';

/// `DiagnosticoResponse` de fincas (familia camelCase; sufijos ID/URL en
/// mayúscula vía `@JsonKey`). Cotejado contra
/// `fincas/internal/presentation/dto/diagnostico.go`.
@freezed
abstract class DiagnosticoDto with _$DiagnosticoDto {
  const factory DiagnosticoDto({
    required String id,
    @JsonKey(name: 'muestraID') required String muestraId,
    required String estado,
    required bool tieneClorosis,
    required double confianza,
    @JsonKey(name: 'imageURL') String? imageUrl,
    String? imageBase64,
  }) = _DiagnosticoDto;

  const DiagnosticoDto._();

  factory DiagnosticoDto.fromJson(Map<String, dynamic> json) =>
      _$DiagnosticoDtoFromJson(json);

  Diagnostico toDomain() => Diagnostico(
    id: id,
    muestraId: muestraId,
    estado: estado,
    tieneClorosis: tieneClorosis,
    confianza: confianza,
    imageUrl: imageUrl,
    imageBase64: imageBase64,
  );
}
