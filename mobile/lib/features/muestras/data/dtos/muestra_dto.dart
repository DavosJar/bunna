import 'package:freezed_annotation/freezed_annotation.dart';

import '../../domain/entities/muestra.dart';

part 'muestra_dto.freezed.dart';
part 'muestra_dto.g.dart';

/// `MuestraResponse` / `MuestraItemResponse` de fincas (idénticos). Familia
/// camelCase; sufijos ID en mayúscula requieren `@JsonKey`. Cotejado contra
/// `fincas/internal/presentation/dto/muestra.go`.
@freezed
abstract class MuestraDto with _$MuestraDto {
  const factory MuestraDto({
    required String id,
    @JsonKey(name: 'fincaID') required String fincaId,
    @JsonKey(name: 'loteID') required String loteId,
    required double latitud,
    required double longitud,
    required DateTime createdAt,
  }) = _MuestraDto;

  const MuestraDto._();

  factory MuestraDto.fromJson(Map<String, dynamic> json) =>
      _$MuestraDtoFromJson(json);

  Muestra toDomain() => Muestra(
    id: id,
    fincaId: fincaId,
    loteId: loteId,
    latitud: latitud,
    longitud: longitud,
    createdAt: createdAt,
  );
}
