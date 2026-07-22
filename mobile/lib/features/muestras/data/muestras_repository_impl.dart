import '../domain/entities/muestra.dart';
import '../domain/muestras_repository.dart';
import 'muestras_api.dart';

final class MuestrasRepositoryImpl implements MuestrasRepository {
  MuestrasRepositoryImpl(this._api);

  final MuestrasApi _api;

  @override
  Future<List<Muestra>> listar({
    required String fincaId,
    required String loteId,
  }) async {
    final dtos = await _api.listar(fincaId: fincaId, loteId: loteId);
    return dtos.map((d) => d.toDomain()).toList(growable: false);
  }

  @override
  Future<Muestra> tomar({
    required String fincaId,
    required String loteId,
    required double latitud,
    required double longitud,
  }) async {
    final dto = await _api.tomar(
      fincaId: fincaId,
      loteId: loteId,
      latitud: latitud,
      longitud: longitud,
    );
    return dto.toDomain();
  }
}
