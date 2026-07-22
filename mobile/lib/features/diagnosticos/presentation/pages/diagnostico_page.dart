import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';

import '../../../../app/theme/app_colors.dart';
import '../../../../core/widgets/circle_icon.dart';
import '../../../../core/widgets/donut_metric.dart';
import '../../domain/entities/analisis_yolo.dart';
import '../diagnostico_flow_controller.dart';
import '../diagnostico_flow_state.dart';

/// Flujo de diagnóstico de una muestra: tomar/elegir foto → análisis YOLO →
/// vincular como diagnóstico → aceptar/rechazar.
class DiagnosticoPage extends ConsumerWidget {
  const DiagnosticoPage({super.key, required this.muestraId});

  final String muestraId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(diagnosticoFlowControllerProvider(muestraId));
    final controller =
        ref.read(diagnosticoFlowControllerProvider(muestraId).notifier);

    return Scaffold(
      appBar: AppBar(title: const Text('Diagnóstico de muestra')),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            _ImagenPreview(state: state),
            const SizedBox(height: 16),
            if (state.error != null) ...[
              _ErrorBanner(state.error!),
              const SizedBox(height: 12),
            ],
            ..._buildFase(context, ref, state, controller),
          ],
        ),
      ),
    );
  }

  List<Widget> _buildFase(
    BuildContext context,
    WidgetRef ref,
    DiagnosticoFlowState state,
    DiagnosticoFlowController controller,
  ) {
    switch (state.fase) {
      case FlowFase.inicial:
        return [_PickerButtons(muestraId: muestraId)];
      case FlowFase.analizando:
        return const [_Cargando('Analizando la hoja con IA…')];
      case FlowFase.analizado:
        return [
          _ResultadoYolo(analisis: state.analisis!),
          const SizedBox(height: 16),
          FilledButton.icon(
            onPressed: controller.vincular,
            icon: const Icon(Icons.link),
            label: const Text('Vincular diagnóstico a la muestra'),
          ),
          const SizedBox(height: 8),
          _PickerButtons(muestraId: muestraId, label: 'Volver a tomar foto'),
        ];
      case FlowFase.vinculando:
        return const [_Cargando('Guardando el diagnóstico…')];
      case FlowFase.vinculado:
        return [
          _ResultadoYolo(analisis: state.analisis!),
          const SizedBox(height: 16),
          const _InfoBanner(
            'Diagnóstico vinculado. Revisa el resultado y acéptalo o recházalo.',
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: OutlinedButton.icon(
                  onPressed: () => _rechazar(context, controller),
                  icon: const Icon(Icons.close),
                  label: const Text('Rechazar'),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: FilledButton.icon(
                  onPressed: controller.aceptar,
                  icon: const Icon(Icons.check),
                  label: const Text('Aceptar'),
                ),
              ),
            ],
          ),
        ];
      case FlowFase.procesando:
        return const [_Cargando('Aplicando decisión…')];
      case FlowFase.resuelto:
        return [
          _Resuelto(estado: state.resueltoEstado ?? ''),
          const SizedBox(height: 16),
          OutlinedButton(
            onPressed: controller.reiniciar,
            child: const Text('Nuevo diagnóstico'),
          ),
        ];
    }
  }

  Future<void> _rechazar(
    BuildContext context,
    DiagnosticoFlowController controller,
  ) async {
    final motivo = await showDialog<String>(
      context: context,
      builder: (_) => const _MotivoRechazoDialog(),
    );
    if (motivo != null) controller.rechazar(motivo: motivo);
  }
}

class _ImagenPreview extends StatelessWidget {
  const _ImagenPreview({required this.state});
  final DiagnosticoFlowState state;

  @override
  Widget build(BuildContext context) {
    return AspectRatio(
      aspectRatio: 4 / 3,
      child: Container(
        decoration: BoxDecoration(
          color: AppColors.superficieTenue,
          borderRadius: BorderRadius.circular(20),
        ),
        clipBehavior: Clip.antiAlias,
        child: state.previewBytes != null
            ? Image.memory(state.previewBytes!, fit: BoxFit.cover)
            : const Center(
                child: CircleIcon(Icons.camera_alt, size: 72),
              ),
      ),
    );
  }
}

class _PickerButtons extends ConsumerWidget {
  const _PickerButtons({required this.muestraId, this.label});
  final String muestraId;
  final String? label;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    Future<void> pick(ImageSource source) async {
      final picker = ImagePicker();
      final XFile? file = await picker.pickImage(
        source: source,
        maxWidth: 1600,
        imageQuality: 85,
      );
      if (file == null) return;
      final bytes = await file.readAsBytes();
      // En web XFile.path es un blob URL inservible para File(); enviamos
      // bytes. En móvil el path funciona pero bytes también, así que usamos
      // bytes siempre por uniformidad, salvo que exista un path real de disco.
      final path = file.path.startsWith('blob:') || file.path.isEmpty
          ? null
          : file.path;
      await ref
          .read(diagnosticoFlowControllerProvider(muestraId).notifier)
          .analizar(bytes: bytes, path: path, filename: file.name);
    }

    return Column(
      children: [
        FilledButton.icon(
          onPressed: () => pick(ImageSource.camera),
          icon: const Icon(Icons.camera_alt),
          label: Text(label ?? 'Tomar foto de la hoja'),
        ),
        const SizedBox(height: 10),
        OutlinedButton.icon(
          onPressed: () => pick(ImageSource.gallery),
          icon: const Icon(Icons.photo_library_outlined),
          label: const Text('Elegir de la galería'),
        ),
      ],
    );
  }
}

class _ResultadoYolo extends StatelessWidget {
  const _ResultadoYolo({required this.analisis});
  final AnalisisYolo analisis;

  @override
  Widget build(BuildContext context) {
    final clorosis = analisis.tieneClorosis;
    final colorClorosis = clorosis == true
        ? AppColors.advertencia
        : clorosis == false
        ? AppColors.verde
        : AppColors.textoTenue;
    final textoClorosis = clorosis == true
        ? 'Deficiencia de nitrógeno detectada'
        : clorosis == false
        ? 'Sin deficiencia aparente'
        : 'Resultado indeterminado';

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              children: [
                DonutMetric(
                  value: analisis.avgConfidence,
                  label: 'Confianza',
                  color: AppColors.verde,
                  size: 110,
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        analisis.feedback?.label ?? 'Análisis YOLO',
                        style: Theme.of(context).textTheme.titleMedium
                            ?.copyWith(fontWeight: FontWeight.w700),
                      ),
                      const SizedBox(height: 6),
                      Row(
                        children: [
                          Icon(Icons.eco, size: 16, color: colorClorosis),
                          const SizedBox(width: 4),
                          Expanded(
                            child: Text(
                              textoClorosis,
                              style: TextStyle(
                                color: colorClorosis,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '${analisis.numDetections} detecciones',
                        style: const TextStyle(color: AppColors.textoTenue),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            if (analisis.feedback?.recommendation != null) ...[
              const Divider(height: 24),
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Icon(
                    Icons.tips_and_updates_outlined,
                    size: 18,
                    color: AppColors.cafe,
                  ),
                  const SizedBox(width: 8),
                  Expanded(child: Text(analisis.feedback!.recommendation!)),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _Resuelto extends StatelessWidget {
  const _Resuelto({required this.estado});
  final String estado;

  @override
  Widget build(BuildContext context) {
    final aceptado = estado.toUpperCase() == 'ACEPTADO';
    final color = aceptado ? AppColors.verde : AppColors.error;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            CircleIcon(
              aceptado ? Icons.check_circle : Icons.cancel,
              size: 64,
              background: color.withValues(alpha: 0.12),
              foreground: color,
            ),
            const SizedBox(height: 12),
            Text(
              aceptado ? 'Diagnóstico aceptado' : 'Diagnóstico rechazado',
              style: Theme.of(context).textTheme.titleMedium
                  ?.copyWith(fontWeight: FontWeight.w700),
            ),
          ],
        ),
      ),
    );
  }
}

class _Cargando extends StatelessWidget {
  const _Cargando(this.mensaje);
  final String mensaje;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 32),
      child: Column(
        children: [
          const CircularProgressIndicator(),
          const SizedBox(height: 16),
          Text(mensaje, style: const TextStyle(color: AppColors.textoTenue)),
        ],
      ),
    );
  }
}

class _ErrorBanner extends StatelessWidget {
  const _ErrorBanner(this.message);
  final String message;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(13),
      decoration: BoxDecoration(
        color: AppColors.errorBg,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.errorBorder),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(
            Icons.error_outline_rounded,
            color: AppColors.errorIcon,
            size: 20,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Text(
              message,
              style: const TextStyle(
                color: AppColors.errorText,
                fontSize: 13.5,
                fontWeight: FontWeight.w500,
                height: 1.35,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _InfoBanner extends StatelessWidget {
  const _InfoBanner(this.message);
  final String message;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(13),
      decoration: BoxDecoration(
        color: AppColors.successBg,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: AppColors.successBorder),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(
            Icons.check_circle_outline_rounded,
            color: AppColors.successText,
            size: 20,
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Text(
              message,
              style: const TextStyle(
                color: AppColors.successText,
                fontSize: 13.5,
                height: 1.35,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _MotivoRechazoDialog extends StatefulWidget {
  const _MotivoRechazoDialog();

  @override
  State<_MotivoRechazoDialog> createState() => _MotivoRechazoDialogState();
}

class _MotivoRechazoDialogState extends State<_MotivoRechazoDialog> {
  final _controller = TextEditingController();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Rechazar diagnóstico'),
      content: TextField(
        controller: _controller,
        decoration: const InputDecoration(labelText: 'Motivo (opcional)'),
        maxLines: 2,
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancelar'),
        ),
        FilledButton(
          onPressed: () => Navigator.pop(context, _controller.text.trim()),
          child: const Text('Rechazar'),
        ),
      ],
    );
  }
}
