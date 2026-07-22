import 'package:freezed_annotation/freezed_annotation.dart';

import '../../domain/entities/lote.dart';

part 'lote_dto.freezed.dart';
part 'lote_dto.g.dart';

/// `LoteResponse` de fincas (familia camelCase). El único campo que no
/// coincide con lowerCamel Dart es el sufijo ID en mayúscula (`fincaID`).
/// Cotejado contra `fincas/internal/presentation/dto/lote.go`.
@freezed
abstract class LoteDto with _$LoteDto {
  const factory LoteDto({
    required String id,
    @JsonKey(name: 'fincaID') required String fincaId,
    required String nombre,
    required double area,
    required String descripcion,
    required String estado,
    required DateTime createdAt,
  }) = _LoteDto;

  const LoteDto._();

  factory LoteDto.fromJson(Map<String, dynamic> json) =>
      _$LoteDtoFromJson(json);

  Lote toDomain() => Lote(
    id: id,
    fincaId: fincaId,
    nombre: nombre,
    area: area,
    descripcion: descripcion,
    estado: estado,
    createdAt: createdAt,
  );
}
