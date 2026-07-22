import 'dart:typed_data';

import '../../../core/domain/cambio_estado.dart';
import '../domain/diagnosticos_repository.dart';
import '../domain/entities/analisis_yolo.dart';
import '../domain/entities/diagnostico.dart';
import 'diagnosticos_api.dart';
import 'yolo_api.dart';

final class DiagnosticosRepositoryImpl implements DiagnosticosRepository {
  DiagnosticosRepositoryImpl({
    required YoloApi yoloApi,
    required DiagnosticosApi diagnosticosApi,
  }) : _yolo = yoloApi,
       _fincas = diagnosticosApi;

  final YoloApi _yolo;
  final DiagnosticosApi _fincas;

  @override
  Future<AnalisisYolo> analizarImagen(String rutaArchivo) =>
      _yolo.analizar(rutaArchivo);

  /// Ruta alternativa para bytes en memoria (web / XFile sin path).
  Future<AnalisisYolo> analizarBytes(Uint8List bytes, {String? filename}) =>
      _yolo.analizarBytes(bytes, filename: filename ?? 'muestra.jpg');

  @override
  Future<Diagnostico> guardarResultado(
    String muestraId, {
    required String imageUrl,
    required bool tieneClorosis,
    required double confianza,
    required DateTime procesadoAt,
  }) async {
    final dto = await _fincas.guardarResultado(
      muestraId,
      imageUrl: imageUrl,
      tieneClorosis: tieneClorosis,
      confianza: confianza,
      procesadoAt: procesadoAt,
    );
    return dto.toDomain();
  }

  @override
  Future<CambioEstado> aceptar(String diagnosticoId) async =>
      (await _fincas.aceptar(diagnosticoId)).toDomain();

  @override
  Future<CambioEstado> rechazar(String diagnosticoId, {String motivo = ''}) async =>
      (await _fincas.rechazar(diagnosticoId, motivo: motivo)).toDomain();
}
