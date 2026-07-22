import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../data/lotes_providers.dart';
import '../domain/entities/lote.dart';

part 'lotes_controller.g.dart';

/// Lotes de una finca (family por `fincaId`). `autoDispose`: caché por
/// pantalla, sin persistencia.
@riverpod
class LotesController extends _$LotesController {
  @override
  Future<List<Lote>> build(String fincaId) {
    return ref.watch(lotesRepositoryProvider).listarPorFinca(fincaId);
  }

  Future<void> agregar({
    required String nombre,
    required double area,
    String descripcion = '',
  }) async {
    await ref.read(lotesRepositoryProvider).agregar(
      fincaId,
      nombre: nombre,
      area: area,
      descripcion: descripcion,
    );
    ref.invalidateSelf();
    await future;
  }

  Future<void> eliminar(String loteId) async {
    await ref.read(lotesRepositoryProvider).eliminar(loteId);
    ref.invalidateSelf();
    await future;
  }
}
