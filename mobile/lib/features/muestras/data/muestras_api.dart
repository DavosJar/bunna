import 'package:dio/dio.dart';

import '../../../core/network/dio_exception_x.dart';
import 'dtos/muestra_dto.dart';

class MuestrasApi {
  MuestrasApi(this._dio);

  final Dio _dio;
  static const _base = '/api/v1/fincas';

  String _path(String fincaId, String loteId) =>
      '$_base/fincas/$fincaId/lotes/$loteId/muestras';

  Future<List<MuestraDto>> listar({
    required String fincaId,
    required String loteId,
  }) async {
    try {
      final res = await _dio.get<Map<String, dynamic>>(_path(fincaId, loteId));
      final data = res.data?['data'] as List<dynamic>? ?? const [];
      return data
          .map((e) => MuestraDto.fromJson(e as Map<String, dynamic>))
          .toList(growable: false);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<MuestraDto> tomar({
    required String fincaId,
    required String loteId,
    required double latitud,
    required double longitud,
  }) async {
    try {
      final res = await _dio.post<Map<String, dynamic>>(
        _path(fincaId, loteId),
        data: {'latitud': latitud, 'longitud': longitud},
      );
      return MuestraDto.fromJson(res.data!['data'] as Map<String, dynamic>);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }
}
