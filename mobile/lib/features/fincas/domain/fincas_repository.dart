import '../../../core/domain/cambio_estado.dart';
import 'entities/finca.dart';

/// Todos los métodos lanzan `AppException`. Ver ARQUITECTURA.md §2.
abstract interface class FincasRepository {
  /// GET /api/v1/fincas/fincas (el backend ya filtra por tenant del token)
  Future<List<Finca>> listar();

  /// POST /fincas
  Future<Finca> registrar({
    required String nombre,
    required String ubicacion,
    String descripcion = '',
  });

  /// POST /fincas/{id}/desactivar
  Future<CambioEstado> desactivar(String fincaId, {bool confirmar = true});
}
