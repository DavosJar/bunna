import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../data/fincas_providers.dart';
import '../domain/entities/finca.dart';

part 'fincas_controller.g.dart';

/// Lista de fincas del tenant activo. `autoDispose` (default de `@riverpod`):
/// la caché vive solo mientras la pantalla está montada — sin persistencia
/// de datos de negocio, se relee del backend en cada entrada. Las mutaciones
/// invalidan el provider para forzar un refetch.
@riverpod
class FincasController extends _$FincasController {
  @override
  Future<List<Finca>> build() {
    return ref.watch(fincasRepositoryProvider).listar();
  }

  /// Registra una finca y refresca la lista. Lanza `AppException` en error
  /// (la UI lo captura para mostrarlo).
  Future<void> registrar({
    required String nombre,
    required String ubicacion,
    String descripcion = '',
  }) async {
    await ref.read(fincasRepositoryProvider).registrar(
      nombre: nombre,
      ubicacion: ubicacion,
      descripcion: descripcion,
    );
    ref.invalidateSelf();
    await future;
  }

  Future<void> desactivar(String fincaId) async {
    await ref.read(fincasRepositoryProvider).desactivar(fincaId);
    ref.invalidateSelf();
    await future;
  }
}
