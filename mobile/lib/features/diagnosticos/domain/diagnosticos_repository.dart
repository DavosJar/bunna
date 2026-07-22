import '../../../core/domain/cambio_estado.dart';
import 'entities/analisis_yolo.dart';
import 'entities/diagnostico.dart';

abstract interface class DiagnosticosRepository {
  /// POST YOLO /api/v1/diagnostico — multipart, campo 'archivo'
  Future<AnalisisYolo> analizarImagen(String rutaArchivo);

  /// POST /muestras/{muestraId}/diagnosticos/manual/resultado
  /// Registra el resultado del análisis como diagnóstico de la muestra.
  Future<Diagnostico> guardarResultado(
    String muestraId, {
    required String imageUrl,
    required bool tieneClorosis,
    required double confianza,
    required DateTime procesadoAt,
  });

  /// POST /diagnosticos/{id}/aceptar
  Future<CambioEstado> aceptar(String diagnosticoId);

  /// POST /diagnosticos/{id}/rechazar — body {motivo}
  Future<CambioEstado> rechazar(String diagnosticoId, {String motivo = ''});
}
