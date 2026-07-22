import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../data/muestras_providers.dart';
import '../domain/entities/muestra.dart';

part 'muestras_controller.g.dart';

/// Muestras de un lote (family por fincaId + loteId).
@riverpod
class MuestrasController extends _$MuestrasController {
  @override
  Future<List<Muestra>> build(String fincaId, String loteId) {
    return ref
        .watch(muestrasRepositoryProvider)
        .listar(fincaId: fincaId, loteId: loteId);
  }

  Future<Muestra> tomar({
    required double latitud,
    required double longitud,
  }) async {
    final muestra = await ref.read(muestrasRepositoryProvider).tomar(
      fincaId: fincaId,
      loteId: loteId,
      latitud: latitud,
      longitud: longitud,
    );
    ref.invalidateSelf();
    await future;
    return muestra;
  }
}
