import 'package:freezed_annotation/freezed_annotation.dart';

part 'analisis_yolo.freezed.dart';

@freezed
abstract class YoloFeedback with _$YoloFeedback {
  const factory YoloFeedback({
    String? label,
    String? level, // low | medium | high
    String? recommendation,
  }) = _YoloFeedback;
}

/// Resultado del análisis YOLO de una foto de hoja. `imageBase64` es la
/// imagen anotada devuelta por la API (data URL o base64).
@freezed
abstract class AnalisisYolo with _$AnalisisYolo {
  const factory AnalisisYolo({
    YoloFeedback? feedback,
    required int numDetections,
    required double avgConfidence,
    required List<String> clasesDetectadas,
    String? imageBase64,
  }) = _AnalisisYolo;

  const AnalisisYolo._();

  /// Deriva si hay clorosis (deficiencia de nitrógeno), replicando
  /// `tieneClorosisFromYolo` del frontend web. `null` = indeterminado.
  bool? get tieneClorosis {
    if (clasesDetectadas.contains('deficiencia_nitrogeno')) return true;
    if (numDetections == 0) return null;
    switch (feedback?.level) {
      case 'low':
      case 'medium':
        return true;
      case 'high':
        return false;
      default:
        return null;
    }
  }
}
