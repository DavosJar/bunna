import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/features/diagnosticos/domain/entities/analisis_yolo.dart';

/// La derivación de clorosis espeja `tieneClorosisFromYolo` del frontend web
/// (frontend/src/utils/yoloDiagnostico.js) — es lo que decide el resultado
/// que se guarda en el backend.
void main() {
  AnalisisYolo build({
    List<String> clases = const [],
    int numDet = 0,
    String? level,
  }) => AnalisisYolo(
    feedback: level == null ? null : YoloFeedback(level: level),
    numDetections: numDet,
    avgConfidence: 0.5,
    clasesDetectadas: clases,
  );

  test('detección explícita de deficiencia_nitrogeno ⇒ true', () {
    expect(
      build(clases: ['deficiencia_nitrogeno'], numDet: 1).tieneClorosis,
      isTrue,
    );
  });

  test('sin detecciones ⇒ indeterminado (null)', () {
    expect(build(numDet: 0).tieneClorosis, isNull);
  });

  test('feedback level low/medium ⇒ true', () {
    expect(build(numDet: 2, level: 'low').tieneClorosis, isTrue);
    expect(build(numDet: 2, level: 'medium').tieneClorosis, isTrue);
  });

  test('feedback level high ⇒ false (sano)', () {
    expect(build(numDet: 2, level: 'high').tieneClorosis, isFalse);
  });

  test('detecciones sin clase ni level conocido ⇒ null', () {
    expect(build(numDet: 3, level: 'otro').tieneClorosis, isNull);
  });
}
