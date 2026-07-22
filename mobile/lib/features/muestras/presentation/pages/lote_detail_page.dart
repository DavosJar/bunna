import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../../../app/theme/app_colors.dart';
import '../../../../core/error/app_exception.dart';
import '../../../../core/error/error_ui.dart';
import '../../../../core/location/location_service.dart';
import '../../../../core/widgets/circle_icon.dart';
import '../../../../core/widgets/error_retry_view.dart';
import '../../../../core/widgets/loading_view.dart';
import '../../../lotes/domain/entities/lote.dart';
import '../../domain/entities/muestra.dart';
import '../muestras_controller.dart';

/// Detalle de lote: lista de muestras + acción "tomar muestra" que captura la
/// ubicación GPS del dispositivo. Cada muestra abre el flujo de diagnóstico.
class LoteDetailPage extends ConsumerStatefulWidget {
  const LoteDetailPage({
    super.key,
    required this.fincaId,
    required this.loteId,
    this.lote,
  });

  final String fincaId;
  final String loteId;
  final Lote? lote;

  @override
  ConsumerState<LoteDetailPage> createState() => _LoteDetailPageState();
}

class _LoteDetailPageState extends ConsumerState<LoteDetailPage> {
  bool _capturando = false;

  Future<void> _tomarMuestra() async {
    setState(() => _capturando = true);
    final messenger = ScaffoldMessenger.of(context);
    try {
      final coord = await ref
          .read(locationServiceProvider)
          .obtenerUbicacion();
      await ref
          .read(
            muestrasControllerProvider(widget.fincaId, widget.loteId).notifier,
          )
          .tomar(latitud: coord.latitud, longitud: coord.longitud);
      messenger.showSnackBar(
        SnackBar(
          content: Text(
            'Muestra registrada en '
            '${coord.latitud.toStringAsFixed(5)}, '
            '${coord.longitud.toStringAsFixed(5)}',
          ),
        ),
      );
    } on AppException catch (e) {
      messenger.showSnackBar(appErrorSnackBar(e));
    } finally {
      if (mounted) setState(() => _capturando = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final muestrasAsync = ref.watch(
      muestrasControllerProvider(widget.fincaId, widget.loteId),
    );

    return Scaffold(
      appBar: AppBar(title: Text(widget.lote?.nombre ?? 'Lote')),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _capturando ? null : _tomarMuestra,
        icon: _capturando
            ? const SizedBox(
                height: 18,
                width: 18,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  color: Colors.white,
                ),
              )
            : const Icon(Icons.my_location),
        label: Text(_capturando ? 'Ubicando…' : 'Tomar muestra'),
      ),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(
          muestrasControllerProvider(widget.fincaId, widget.loteId),
        ),
        child: muestrasAsync.when(
          loading: () => const LoadingView(),
          error: (e, _) => ErrorRetryView(
            error: e,
            onRetry: () => ref.invalidate(
              muestrasControllerProvider(widget.fincaId, widget.loteId),
            ),
          ),
          data: (muestras) => ListView(
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 96),
            children: [
              if (widget.lote != null) _LoteHeader(lote: widget.lote!),
              const SizedBox(height: 16),
              Text(
                'Muestras (${muestras.length})',
                style: Theme.of(context).textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.w700),
              ),
              const SizedBox(height: 8),
              if (muestras.isEmpty)
                const _EmptyMuestras()
              else
                ...muestras.map(
                  (m) => Padding(
                    padding: const EdgeInsets.only(bottom: 12),
                    child: _MuestraCard(
                      fincaId: widget.fincaId,
                      loteId: widget.loteId,
                      muestra: m,
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

class _LoteHeader extends StatelessWidget {
  const _LoteHeader({required this.lote});
  final Lote lote;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            const CircleIcon(
              Icons.grid_view_rounded,
              size: 52,
              foreground: AppColors.cafe,
            ),
            const SizedBox(width: 14),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    lote.nombre,
                    style: Theme.of(context).textTheme.titleLarge
                        ?.copyWith(fontWeight: FontWeight.w700),
                  ),
                  Text(
                    '${lote.area.toStringAsFixed(2)} ha',
                    style: const TextStyle(color: AppColors.textoTenue),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _MuestraCard extends StatelessWidget {
  const _MuestraCard({
    required this.fincaId,
    required this.loteId,
    required this.muestra,
  });

  final String fincaId;
  final String loteId;
  final Muestra muestra;

  @override
  Widget build(BuildContext context) {
    final fecha = DateFormat('d MMM, HH:mm', 'es').format(muestra.createdAt);
    return Card(
      child: InkWell(
        borderRadius: BorderRadius.circular(20),
        onTap: () => context.push(
          '/fincas/$fincaId/lotes/$loteId/muestras/${muestra.id}/diagnostico',
        ),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              const CircleIcon(Icons.location_on),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${muestra.latitud.toStringAsFixed(5)}, '
                      '${muestra.longitud.toStringAsFixed(5)}',
                      style: const TextStyle(fontWeight: FontWeight.w700),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      fecha,
                      style: const TextStyle(color: AppColors.textoTenue),
                    ),
                  ],
                ),
              ),
              Row(
                children: const [
                  Icon(Icons.camera_alt_outlined, color: AppColors.verde),
                  Icon(Icons.chevron_right, color: AppColors.textoTenue),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _EmptyMuestras extends StatelessWidget {
  const _EmptyMuestras();

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: const [
            CircleIcon(Icons.my_location, size: 56),
            SizedBox(height: 12),
            Text(
              'Aún no hay muestras en este lote.\nToma una muestra con el botón de abajo — se guardará con tu ubicación GPS.',
              textAlign: TextAlign.center,
              style: TextStyle(color: AppColors.textoTenue),
            ),
          ],
        ),
      ),
    );
  }
}
