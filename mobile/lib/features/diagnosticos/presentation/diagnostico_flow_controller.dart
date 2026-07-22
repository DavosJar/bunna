import 'dart:typed_data';

import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/error/app_exception.dart';
import '../../../core/error/error_ui.dart';
import '../data/diagnosticos_providers.dart';
import 'diagnostico_flow_state.dart';

part 'diagnostico_flow_controller.g.dart';

/// Orquesta el flujo de diagnóstico de UNA muestra: foto → análisis YOLO →
/// vincular como diagnóstico en el backend → aceptar/rechazar. Family por
/// `muestraId`. Imperativo (Notifier), no AsyncNotifier: cada paso lo dispara
/// el usuario y el estado avanza por fases.
@riverpod
class DiagnosticoFlowController extends _$DiagnosticoFlowController {
  @override
  DiagnosticoFlowState build(String muestraId) =>
      const DiagnosticoFlowState();

  /// Analiza una imagen ya tomada. Pasar `bytes` siempre (para preview y para
  /// el envío multipart, que funciona igual en web y móvil). `path` es
  /// opcional: si viene, se usa para el envío desde disco (móvil).
  Future<void> analizar({
    required Uint8List bytes,
    String? path,
    String? filename,
  }) async {
    state = state.copyWith(
      fase: FlowFase.analizando,
      previewBytes: bytes,
      error: null,
      analisis: null,
      diagnostico: null,
      resueltoEstado: null,
    );
    try {
      final repo = ref.read(diagnosticosRepositoryProvider);
      final analisis = path != null
          ? await repo.analizarImagen(path)
          : await repo.analizarBytes(bytes, filename: filename);
      state = state.copyWith(fase: FlowFase.analizado, analisis: analisis);
    } on AppException catch (e) {
      state = state.copyWith(fase: FlowFase.inicial, error: e.mensajeUsuario);
    }
  }

  /// Registra el resultado del análisis como diagnóstico de la muestra
  /// (POST .../diagnosticos/manual/resultado). Espeja el flujo del web:
  /// `imageURL` = imagen anotada de YOLO; clorosis derivada; confianza = avg.
  Future<void> vincular() async {
    final analisis = state.analisis;
    if (analisis == null) return;

    state = state.copyWith(fase: FlowFase.vinculando, error: null);
    try {
      final diagnostico = await ref
          .read(diagnosticosRepositoryProvider)
          .guardarResultado(
            muestraId,
            imageUrl: analisis.imageBase64 ?? '',
            tieneClorosis: analisis.tieneClorosis ?? false,
            confianza: analisis.avgConfidence,
            procesadoAt: DateTime.now(),
          );
      state = state.copyWith(
        fase: FlowFase.vinculado,
        diagnostico: diagnostico,
      );
    } on AppException catch (e) {
      state = state.copyWith(fase: FlowFase.analizado, error: e.mensajeUsuario);
    }
  }

  Future<void> aceptar() => _resolver(aceptar: true);

  Future<void> rechazar({String motivo = ''}) =>
      _resolver(aceptar: false, motivo: motivo);

  Future<void> _resolver({required bool aceptar, String motivo = ''}) async {
    final diagnostico = state.diagnostico;
    if (diagnostico == null) return;

    state = state.copyWith(fase: FlowFase.procesando, error: null);
    try {
      final repo = ref.read(diagnosticosRepositoryProvider);
      final cambio = aceptar
          ? await repo.aceptar(diagnostico.id)
          : await repo.rechazar(diagnostico.id, motivo: motivo);
      state = state.copyWith(
        fase: FlowFase.resuelto,
        resueltoEstado: cambio.estado,
      );
    } on AppException catch (e) {
      state = state.copyWith(fase: FlowFase.vinculado, error: e.mensajeUsuario);
    }
  }

  void reiniciar() => state = const DiagnosticoFlowState();
}
