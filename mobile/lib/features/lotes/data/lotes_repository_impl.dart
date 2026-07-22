import '../../../core/domain/cambio_estado.dart';
import '../domain/entities/lote.dart';
import '../domain/lotes_repository.dart';
import 'lotes_api.dart';

final class LotesRepositoryImpl implements LotesRepository {
  LotesRepositoryImpl(this._api);

  final LotesApi _api;

  @override
  Future<List<Lote>> listarPorFinca(String fincaId) async {
    final dtos = await _api.listarPorFinca(fincaId);
    return dtos.map((d) => d.toDomain()).toList(growable: false);
  }

  @override
  Future<Lote> agregar(
    String fincaId, {
    required String nombre,
    required double area,
    String descripcion = '',
  }) async {
    final dto = await _api.agregar(
      fincaId,
      nombre: nombre,
      area: area,
      descripcion: descripcion,
    );
    return dto.toDomain();
  }

  @override
  Future<CambioEstado> eliminar(String loteId) async {
    final dto = await _api.eliminar(loteId);
    return dto.toDomain();
  }
}
