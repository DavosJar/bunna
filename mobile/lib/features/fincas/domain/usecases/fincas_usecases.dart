import '../../../../core/domain/cambio_estado.dart';
import '../entities/finca.dart';
import '../fincas_repository.dart';

/// Casos de uso de fincas — wrappers finos sobre el repositorio (uno por
/// operación, ver ARQUITECTURA.md §2). Agrupados en un archivo por brevedad;
/// cada clase sigue siendo una unidad independiente.
final class ListarFincas {
  const ListarFincas(this._repo);
  final FincasRepository _repo;
  Future<List<Finca>> call() => _repo.listar();
}

final class RegistrarFinca {
  const RegistrarFinca(this._repo);
  final FincasRepository _repo;
  Future<Finca> call({
    required String nombre,
    required String ubicacion,
    String descripcion = '',
  }) => _repo.registrar(
    nombre: nombre,
    ubicacion: ubicacion,
    descripcion: descripcion,
  );
}

final class DesactivarFinca {
  const DesactivarFinca(this._repo);
  final FincasRepository _repo;
  Future<CambioEstado> call(String fincaId) => _repo.desactivar(fincaId);
}
