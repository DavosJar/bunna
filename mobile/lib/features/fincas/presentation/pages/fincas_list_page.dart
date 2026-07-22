import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../../app/theme/app_colors.dart';
import '../../../../core/widgets/circle_icon.dart';
import '../../../../core/widgets/error_retry_view.dart';
import '../../../../core/widgets/loading_view.dart';
import '../../domain/entities/finca.dart';
import '../fincas_controller.dart';
import '../widgets/finca_form_sheet.dart';

class FincasListPage extends ConsumerWidget {
  const FincasListPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final fincasAsync = ref.watch(fincasControllerProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Mis fincas')),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _crearFinca(context, ref),
        icon: const Icon(Icons.add),
        label: const Text('Finca'),
      ),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(fincasControllerProvider),
        child: fincasAsync.when(
          loading: () => const LoadingView(),
          error: (e, _) => ErrorRetryView(
            error: e,
            onRetry: () => ref.invalidate(fincasControllerProvider),
          ),
          data: (fincas) => fincas.isEmpty
              ? _EmptyFincas(onCreate: () => _crearFinca(context, ref))
              : ListView.separated(
                  padding: const EdgeInsets.fromLTRB(16, 16, 16, 96),
                  itemCount: fincas.length,
                  separatorBuilder: (_, _) => const SizedBox(height: 12),
                  itemBuilder: (_, i) => _FincaCard(finca: fincas[i]),
                ),
        ),
      ),
    );
  }

  Future<void> _crearFinca(BuildContext context, WidgetRef ref) async {
    final creada = await FincaFormSheet.show(context);
    if (creada == true && context.mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Finca registrada')),
      );
    }
  }
}

class _FincaCard extends StatelessWidget {
  const _FincaCard({required this.finca});
  final Finca finca;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: InkWell(
        borderRadius: BorderRadius.circular(20),
        onTap: () => context.push('/fincas/${finca.id}', extra: finca),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              const CircleIcon(Icons.agriculture, size: 48),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      finca.nombre,
                      style: Theme.of(context).textTheme.titleMedium
                          ?.copyWith(fontWeight: FontWeight.w700),
                    ),
                    const SizedBox(height: 2),
                    Row(
                      children: [
                        const Icon(
                          Icons.location_on_outlined,
                          size: 14,
                          color: AppColors.textoTenue,
                        ),
                        const SizedBox(width: 3),
                        Expanded(
                          child: Text(
                            finca.ubicacion,
                            style: const TextStyle(color: AppColors.textoTenue),
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              _EstadoChip(finca.estado),
            ],
          ),
        ),
      ),
    );
  }
}

class _EstadoChip extends StatelessWidget {
  const _EstadoChip(this.estado);
  final String estado;

  @override
  Widget build(BuildContext context) {
    final activa = estado.toUpperCase() == 'ACTIVA';
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: activa ? AppColors.verdeSuave : AppColors.superficieTenue,
        borderRadius: BorderRadius.circular(999),
      ),
      child: Text(
        activa ? 'Activa' : estado,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: activa ? AppColors.verdeOscuro : AppColors.textoTenue,
        ),
      ),
    );
  }
}

class _EmptyFincas extends StatelessWidget {
  const _EmptyFincas({required this.onCreate});
  final VoidCallback onCreate;

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: [
        const SizedBox(height: 80),
        const Center(child: CircleIcon(Icons.agriculture, size: 72)),
        const SizedBox(height: 16),
        Text(
          'Aún no tienes fincas',
          textAlign: TextAlign.center,
          style: Theme.of(context).textTheme.titleMedium,
        ),
        const SizedBox(height: 6),
        const Padding(
          padding: EdgeInsets.symmetric(horizontal: 40),
          child: Text(
            'Registra tu primera finca para empezar a monitorear tus lotes.',
            textAlign: TextAlign.center,
            style: TextStyle(color: AppColors.textoTenue),
          ),
        ),
        const SizedBox(height: 24),
        Center(
          child: FilledButton.icon(
            onPressed: onCreate,
            icon: const Icon(Icons.add),
            label: const Text('Registrar finca'),
          ),
        ),
      ],
    );
  }
}
