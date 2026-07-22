import 'package:dio/dio.dart';

import '../../../core/data/cambio_estado_dto.dart';
import '../../../core/network/dio_exception_x.dart';
import 'dtos/finca_dto.dart';

/// Datasource remoto de fincas (fincasDio). Devuelve DTOs, lanza
/// `AppException`. Las respuestas de fincas envuelven en `{data, links}`.
class FincasApi {
  FincasApi(this._dio);

  final Dio _dio;

  static const _base = '/api/v1/fincas';

  Future<List<FincaDto>> listar() async {
    try {
      final res = await _dio.get<Map<String, dynamic>>('$_base/fincas');
      final data = res.data?['data'] as List<dynamic>? ?? const [];
      return data
          .map((e) => FincaDto.fromJson(e as Map<String, dynamic>))
          .toList(growable: false);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<FincaDto> registrar({
    required String nombre,
    required String ubicacion,
    required String descripcion,
  }) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        '$_base/fincas',
        data: {
          'nombre': nombre,
          'ubicacion': ubicacion,
          'descripcion': descripcion,
        },
      );
      return FincaDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<CambioEstadoDto> desactivar(
    String fincaId, {
    required bool confirmar,
  }) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        '$_base/fincas/$fincaId/desactivar',
        data: {'confirmar': confirmar},
      );
      return CambioEstadoDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }
}
