import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/network/dio_providers.dart';
import '../domain/diagnosticos_repository.dart';
import 'diagnosticos_api.dart';
import 'diagnosticos_repository_impl.dart';
import 'yolo_api.dart';

part 'diagnosticos_providers.g.dart';

@Riverpod(keepAlive: true)
YoloApi yoloApi(Ref ref) => YoloApi(ref.watch(yoloDioProvider));

@Riverpod(keepAlive: true)
DiagnosticosApi diagnosticosApi(Ref ref) =>
    DiagnosticosApi(ref.watch(fincasDioProvider));

/// Se expone como `DiagnosticosRepositoryImpl` (no la interfaz) porque la UI
/// necesita el método extra `analizarBytes` para el caso web/bytes en
/// memoria, que no está en el contrato `DiagnosticosRepository`.
@Riverpod(keepAlive: true)
DiagnosticosRepositoryImpl diagnosticosRepository(Ref ref) =>
    DiagnosticosRepositoryImpl(
      yoloApi: ref.watch(yoloApiProvider),
      diagnosticosApi: ref.watch(diagnosticosApiProvider),
    );

/// Alias al contrato para código que solo necesita la interfaz.
@Riverpod(keepAlive: true)
DiagnosticosRepository diagnosticosRepositoryContract(Ref ref) =>
    ref.watch(diagnosticosRepositoryProvider);
