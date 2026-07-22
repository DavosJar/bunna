import '../../../core/domain/cambio_estado.dart';
import '../domain/entities/finca.dart';
import '../domain/fincas_repository.dart';
import 'fincas_api.dart';

final class FincasRepositoryImpl implements FincasRepository {
  FincasRepositoryImpl(this._api);

  final FincasApi _api;

  @override
  Future<List<Finca>> listar() async {
    final dtos = await _api.listar();
    return dtos.map((d) => d.toDomain()).toList(growable: false);
  }

  @override
  Future<Finca> registrar({
    required String nombre,
    required String ubicacion,
    String descripcion = '',
  }) async {
    final dto = await _api.registrar(
      nombre: nombre,
      ubicacion: ubicacion,
      descripcion: descripcion,
    );
    return dto.toDomain();
  }

  @override
  Future<CambioEstado> desactivar(String fincaId, {bool confirmar = true}) async {
    final dto = await _api.desactivar(fincaId, confirmar: confirmar);
    return dto.toDomain();
  }
}
