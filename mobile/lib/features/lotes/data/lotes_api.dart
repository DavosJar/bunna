import 'package:dio/dio.dart';

import '../../../core/data/cambio_estado_dto.dart';
import '../../../core/network/dio_exception_x.dart';
import 'dtos/lote_dto.dart';

class LotesApi {
  LotesApi(this._dio);

  final Dio _dio;
  static const _base = '/api/v1/fincas';

  Future<List<LoteDto>> listarPorFinca(String fincaId) async {
    try {
      final res = await _dio.get<Map<String, dynamic>>(
        '$_base/fincas/$fincaId/lotes',
      );
      final data = res.data?['data'] as List<dynamic>? ?? const [];
      return data
          .map((e) => LoteDto.fromJson(e as Map<String, dynamic>))
          .toList(growable: false);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<LoteDto> agregar(
    String fincaId, {
    required String nombre,
    required double area,
    required String descripcion,
  }) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        '$_base/fincas/$fincaId/lotes',
        data: {'nombre': nombre, 'area': area, 'descripcion': descripcion},
      );
      return LoteDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<CambioEstadoDto> eliminar(String loteId) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        '$_base/lotes/$loteId/eliminar',
      );
      return CambioEstadoDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }
}
