import 'package:freezed_annotation/freezed_annotation.dart';

part 'cambio_estado.freezed.dart';

/// Entidad compartida — respuesta de todas las mutaciones de estado de fincas
/// (`EstadoCambioResponse`): desactivar finca, eliminar lote, aceptar/rechazar
/// diagnóstico.
@freezed
abstract class CambioEstado with _$CambioEstado {
  const factory CambioEstado({
    required String id,
    required String estado,
    String? motivo,
    required DateTime updatedAt,
  }) = _CambioEstado;
}
