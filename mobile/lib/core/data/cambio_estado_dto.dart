import 'package:freezed_annotation/freezed_annotation.dart';

import '../domain/cambio_estado.dart';

part 'cambio_estado_dto.freezed.dart';
part 'cambio_estado_dto.g.dart';

/// `EstadoCambioResponse` de fincas (familia camelCase). Cotejado contra
/// `fincas/internal/presentation/dto/diagnostico.go`.
@freezed
abstract class CambioEstadoDto with _$CambioEstadoDto {
  const factory CambioEstadoDto({
    required String id,
    required String estado,
    String? motivo,
    required DateTime updatedAt,
  }) = _CambioEstadoDto;

  const CambioEstadoDto._();

  factory CambioEstadoDto.fromJson(Map<String, dynamic> json) =>
      _$CambioEstadoDtoFromJson(json);

  CambioEstado toDomain() =>
      CambioEstado(id: id, estado: estado, motivo: motivo, updatedAt: updatedAt);
}
