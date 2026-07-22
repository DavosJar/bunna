import 'package:dio/dio.dart';

import '../../../core/network/dio_exception_x.dart';
import '../domain/entities/analisis_yolo.dart';

/// Cliente de la API YOLO (yoloDio, sin JWT). Envía la foto como multipart
/// campo `archivo` y parsea el resultado a mano — la respuesta YOLO no tiene
/// un contrato tan estable como para justificar un DTO generado.
class YoloApi {
  YoloApi(this._dio);

  final Dio _dio;

  Future<AnalisisYolo> analizar(String rutaArchivo) async {
    try {
      final formData = FormData.fromMap({
        'archivo': await MultipartFile.fromFile(rutaArchivo),
      });
      final res = await _dio.post<Map<String, dynamic>>(
        '/api/v1/diagnostico',
        data: formData,
      );
      return _parse(res.data ?? const {});
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  /// Variante para web / bytes en memoria (image_picker en Chrome no da path).
  Future<AnalisisYolo> analizarBytes(
    List<int> bytes, {
    String filename = 'muestra.jpg',
  }) async {
    try {
      final formData = FormData.fromMap({
        'archivo': MultipartFile.fromBytes(bytes, filename: filename),
      });
      final res = await _dio.post<Map<String, dynamic>>(
        '/api/v1/diagnostico',
        data: formData,
      );
      return _parse(res.data ?? const {});
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  AnalisisYolo _parse(Map<String, dynamic> json) {
    final feedbackJson = json['feedback'];
    YoloFeedback? feedback;
    if (feedbackJson is Map) {
      feedback = YoloFeedback(
        label: feedbackJson['label'] as String?,
        level: feedbackJson['level'] as String?,
        recommendation: feedbackJson['recommendation'] as String?,
      );
    }

    final detecciones = <String>[];
    final rawDet = json['detections'];
    if (rawDet is List) {
      for (final d in rawDet) {
        if (d is Map && d['class_name'] is String) {
          detecciones.add(d['class_name'] as String);
        }
      }
    }

    return AnalisisYolo(
      feedback: feedback,
      numDetections: (json['num_detections'] as num?)?.toInt() ?? 0,
      avgConfidence: (json['avg_confidence'] as num?)?.toDouble() ?? 0,
      clasesDetectadas: detecciones,
      imageBase64: json['image'] as String?,
    );
  }
}
