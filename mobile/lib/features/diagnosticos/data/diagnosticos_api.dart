import 'package:dio/dio.dart';

import '../../../core/data/cambio_estado_dto.dart';
import '../../../core/network/dio_exception_x.dart';
import 'dtos/diagnostico_dto.dart';

/// Datasource de diagnósticos contra el servicio fincas (fincasDio).
class DiagnosticosApi {
  DiagnosticosApi(this._dio);

  final Dio _dio;
  static const _base = '/api/v1/fincas';

  Future<DiagnosticoDto> guardarResultado(
    String muestraId, {
    required String imageUrl,
    required bool tieneClorosis,
    required double confianza,
    required DateTime procesadoAt,
  }) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        '$_base/muestras/$muestraId/diagnosticos/manual/resultado',
        data: {
          'imageURL': imageUrl,
          'tieneClorosis': tieneClorosis,
          'confianza': confianza,
          'procesadoAt': procesadoAt.toUtc().toIso8601String(),
        },
      );
      return DiagnosticoDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<CambioEstadoDto> aceptar(String diagnosticoId) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        '$_base/diagnosticos/$diagnosticoId/aceptar',
      );
      return CambioEstadoDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<CambioEstadoDto> rechazar(
    String diagnosticoId, {
    required String motivo,
  }) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        '$_base/diagnosticos/$diagnosticoId/rechazar',
        data: {'motivo': motivo},
      );
      return CambioEstadoDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }
}
