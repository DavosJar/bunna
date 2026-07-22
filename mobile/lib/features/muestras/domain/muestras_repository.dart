import 'entities/muestra.dart';

abstract interface class MuestrasRepository {
  /// GET /fincas/{fincaId}/lotes/{loteId}/muestras
  Future<List<Muestra>> listar({
    required String fincaId,
    required String loteId,
  });

  /// POST /fincas/{fincaId}/lotes/{loteId}/muestras — body {latitud, longitud}
  Future<Muestra> tomar({
    required String fincaId,
    required String loteId,
    required double latitud,
    required double longitud,
  });
}
