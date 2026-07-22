import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';

import '../../../../app/theme/app_colors.dart';
import '../../../../core/widgets/circle_icon.dart';
import '../../../../core/widgets/donut_metric.dart';
import '../../domain/entities/analisis_yolo.dart';
import '../quick_scan_controller.dart';

/// Tab central "Escanear": análisis YOLO rápido de una foto de hoja, sin
/// vincular a una muestra.
class QuickScanPage extends ConsumerWidget {
  const QuickScanPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(quickScanControllerProvider);
    final controller = ref.read(quickScanControllerProvider.notifier);

    Future<void> pick(ImageSource source) async {
      final file = await ImagePicker().pickImage(
        source: source,
        maxWidth: 1600,
        imageQuality: 85,
      );
      if (file == null) return;
      final bytes = await file.readAsBytes();
      final path = file.path.startsWith('blob:') || file.path.isEmpty
          ? null
          : file.path;
      await controller.analizar(bytes: bytes, path: path, filename: file.name);
    }

    return Scaffold(
      appBar: AppBar(title: const Text('Escanear hoja')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          AspectRatio(
            aspectRatio: 4 / 3,
            child: Container(
              decoration: BoxDecoration(
                color: AppColors.superficieTenue,
                borderRadius: BorderRadius.circular(20),
              ),
              clipBehavior: Clip.antiAlias,
              child: state.previewBytes != null
                  ? Image.memory(state.previewBytes!, fit: BoxFit.cover)
                  : const Center(child: CircleIcon(Icons.eco, size: 80)),
            ),
          ),
          const SizedBox(height: 16),
          if (state.error != null) ...[
            Container(
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
                      state.error!,
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
            ),
            const SizedBox(height: 12),
          ],
          if (state.analizando)
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 24),
              child: Center(child: CircularProgressIndicator()),
            )
          else if (state.analisis != null) ...[
            _QuickResultado(analisis: state.analisis!),
            const SizedBox(height: 16),
          ],
          FilledButton.icon(
            onPressed: state.analizando ? null : () => pick(ImageSource.camera),
            icon: const Icon(Icons.camera_alt),
            label: const Text('Tomar foto'),
          ),
          const SizedBox(height: 10),
          OutlinedButton.icon(
            onPressed: state.analizando ? null : () => pick(ImageSource.gallery),
            icon: const Icon(Icons.photo_library_outlined),
            label: const Text('Elegir de la galería'),
          ),
          if (state.analisis != null) ...[
            const SizedBox(height: 10),
            TextButton(
              onPressed: controller.reiniciar,
              child: const Text('Limpiar'),
            ),
          ],
        ],
      ),
    );
  }
}

class _QuickResultado extends StatelessWidget {
  const _QuickResultado({required this.analisis});
  final AnalisisYolo analisis;

  @override
  Widget build(BuildContext context) {
    final clorosis = analisis.tieneClorosis;
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
            DonutMetric(
              value: analisis.avgConfidence,
              label: 'Confianza del modelo',
            ),
            const SizedBox(height: 12),
            Text(
              analisis.feedback?.label ?? 'Análisis completado',
              style: Theme.of(context).textTheme.titleMedium
                  ?.copyWith(fontWeight: FontWeight.w700),
            ),
            const SizedBox(height: 4),
            Text(
              textoClorosis,
              style: const TextStyle(color: AppColors.textoTenue),
            ),
            if (analisis.feedback?.recommendation != null) ...[
              const Divider(height: 24),
              Text(analisis.feedback!.recommendation!),
            ],
          ],
        ),
      ),
    );
  }
}
