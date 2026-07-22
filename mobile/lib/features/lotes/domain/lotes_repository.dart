import '../../../core/domain/cambio_estado.dart';
import 'entities/lote.dart';

abstract interface class LotesRepository {
  /// GET /fincas/{fincaId}/lotes
  Future<List<Lote>> listarPorFinca(String fincaId);

  /// POST /fincas/{fincaId}/lotes
  Future<Lote> agregar(
    String fincaId, {
    required String nombre,
    required double area,
    String descripcion = '',
  });

  /// POST /lotes/{loteId}/eliminar
  Future<CambioEstado> eliminar(String loteId);
}
