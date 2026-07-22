import 'dart:typed_data';

import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/error/app_exception.dart';
import '../../../core/error/error_ui.dart';
import '../data/diagnosticos_providers.dart';
import '../domain/entities/analisis_yolo.dart';

part 'quick_scan_controller.freezed.dart';
part 'quick_scan_controller.g.dart';

@freezed
abstract class QuickScanState with _$QuickScanState {
  const factory QuickScanState({
    @Default(false) bool analizando,
    Uint8List? previewBytes,
    AnalisisYolo? analisis,
    String? error,
  }) = _QuickScanState;
}

/// Escaneo YOLO independiente (tab central "Escanear"): analiza una foto sin
/// vincularla a ninguna muestra. Para vincular a una muestra concreta se usa
/// el flujo dentro de un lote.
@riverpod
class QuickScanController extends _$QuickScanController {
  @override
  QuickScanState build() => const QuickScanState();

  Future<void> analizar({
    required Uint8List bytes,
    String? path,
    String? filename,
  }) async {
    state = QuickScanState(analizando: true, previewBytes: bytes);
    try {
      final repo = ref.read(diagnosticosRepositoryProvider);
      final analisis = path != null
          ? await repo.analizarImagen(path)
          : await repo.analizarBytes(bytes, filename: filename);
      state = state.copyWith(analizando: false, analisis: analisis);
    } on AppException catch (e) {
      state = state.copyWith(analizando: false, error: e.mensajeUsuario);
    }
  }

  void reiniciar() => state = const QuickScanState();
}
