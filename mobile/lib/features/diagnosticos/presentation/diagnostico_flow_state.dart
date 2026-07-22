import 'dart:typed_data';

import 'package:freezed_annotation/freezed_annotation.dart';

import '../domain/entities/analisis_yolo.dart';
import '../domain/entities/diagnostico.dart';

part 'diagnostico_flow_state.freezed.dart';

enum FlowFase { inicial, analizando, analizado, vinculando, vinculado, procesando, resuelto }

@freezed
abstract class DiagnosticoFlowState with _$DiagnosticoFlowState {
  const factory DiagnosticoFlowState({
    @Default(FlowFase.inicial) FlowFase fase,
    Uint8List? previewBytes,
    AnalisisYolo? analisis,
    Diagnostico? diagnostico,
    String? resueltoEstado, // ACEPTADO | RECHAZADO
    String? error,
  }) = _DiagnosticoFlowState;

  const DiagnosticoFlowState._();

  bool get ocupado =>
      fase == FlowFase.analizando ||
      fase == FlowFase.vinculando ||
      fase == FlowFase.procesando;
}
