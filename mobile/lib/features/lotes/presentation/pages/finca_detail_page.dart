import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../app/theme/app_colors.dart';
import '../../../../core/widgets/circle_icon.dart';
import '../../../../core/widgets/error_retry_view.dart';
import '../../../../core/widgets/loading_view.dart';
import '../../../fincas/domain/entities/finca.dart';
import '../../domain/entities/lote.dart';
import '../lotes_controller.dart';
import '../widgets/lote_form_sheet.dart';

/// Detalle de finca: cabecera + lista de lotes. `finca` llega por `extra` de
/// GoRouter (evita un refetch); si no viene (deep link directo), se muestra
/// solo el id.
class FincaDetailPage extends ConsumerWidget {
  const FincaDetailPage({super.key, required this.fincaId, this.finca});

  final String fincaId;
  final Finca? finca;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final lotesAsync = ref.watch(lotesControllerProvider(fincaId));

    return Scaffold(
      appBar: AppBar(title: Text(finca?.nombre ?? 'Finca')),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _crearLote(context, ref),
        icon: const Icon(Icons.add),
        label: const Text('Lote'),
      ),
      body: RefreshIndicator(
        onRefresh: () async =>
            ref.invalidate(lotesControllerProvider(fincaId)),
        child: lotesAsync.when(
          loading: () => const LoadingView(),
          error: (e, _) => ErrorRetryView(
            error: e,
            onRetry: () => ref.invalidate(lotesControllerProvider(fincaId)),
          ),
          data: (lotes) => ListView(
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 96),
            children: [
              if (finca != null) _FincaHeader(finca: finca!),
              const SizedBox(height: 16),
              Text(
                'Lotes (${lotes.length})',
                style: Theme.of(context).textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.w700),
              ),
              const SizedBox(height: 8),
              if (lotes.isEmpty)
                _EmptyLotes(onCreate: () => _crearLote(context, ref))
              else
                ...lotes.map(
                  (l) => Padding(
                    padding: const EdgeInsets.only(bottom: 12),
                    child: _LoteCard(fincaId: fincaId, lote: l),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _crearLote(BuildContext context, WidgetRef ref) async {
    final creado = await LoteFormSheet.show(context, fincaId);
    if (creado == true && context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Lote agregado')),
      );
    }
  }
}

class _FincaHeader extends StatelessWidget {
  const _FincaHeader({required this.finca});
  final Finca finca;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            const CircleIcon(Icons.agriculture, size: 52),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    finca.nombre,
                    style: Theme.of(context).textTheme.titleLarge
                        ?.copyWith(fontWeight: FontWeight.w700),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    finca.ubicacion,
                    style: const TextStyle(color: AppColors.textoTenue),
                  ),
                  if (finca.descripcion.isNotEmpty) ...[
                    const SizedBox(height: 6),
                    Text(finca.descripcion),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _LoteCard extends StatelessWidget {
  const _LoteCard({required this.fincaId, required this.lote});
  final String fincaId;
  final Lote lote;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: InkWell(
        borderRadius: BorderRadius.circular(20),
        onTap: () => context.push(
          '/fincas/$fincaId/lotes/${lote.id}',
          extra: lote,
        ),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              const CircleIcon(
                Icons.grid_view_rounded,
                background: AppColors.superficieTenue,
                foreground: AppColors.cafe,
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      lote.nombre,
                      style: const TextStyle(fontWeight: FontWeight.w700),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${lote.area.toStringAsFixed(2)} ha',
                      style: const TextStyle(color: AppColors.textoTenue),
                    ),
                  ],
                ),
              ),
              const Icon(Icons.chevron_right, color: AppColors.textoTenue),
            ],
          ),
        ),
      ),
    );
  }
}

class _EmptyLotes extends StatelessWidget {
  const _EmptyLotes({required this.onCreate});
  final VoidCallback onCreate;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            const CircleIcon(Icons.grid_view_rounded, size: 56),
            const SizedBox(height: 12),
            const Text(
              'Esta finca no tiene lotes todavía',
              textAlign: TextAlign.center,
              style: TextStyle(color: AppColors.textoTenue),
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: onCreate,
              icon: const Icon(Icons.add),
              label: const Text('Agregar lote'),
            ),
          ],
        ),
      ),
    );
  }
}
