import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/error/app_exception.dart';
import '../../../../core/widgets/empty_view.dart';
import '../../../../core/widgets/error_retry_view.dart';
import '../../../../core/widgets/loading_view.dart';
import '../auth_controller.dart';

/// A diferencia de `fincas`/`diagnostico`, este tab es funcional desde la
/// Fase 1: consume el `AuthController` (que ya trae perfil/tenants/permisos
/// del bootstrap post-login) sin necesitar su propio repositorio. La edición
/// de perfil (PUT /mi-perfil, /mi-password) queda para `features/perfil` en
/// una fase posterior — ver ARQUITECTURA.md §1.
class PerfilPage extends ConsumerWidget {
  const PerfilPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authControllerProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Mi perfil')),
      body: authState.when(
        loading: () => const LoadingView(),
        error: (error, _) => ErrorRetryView(
          error: error,
          onRetry: () => ref.invalidate(authControllerProvider),
        ),
        data: (state) => switch (state) {
          AuthAuthenticated() => _PerfilContent(state: state),
          _ => const EmptyView(message: 'No hay sesión activa'),
        },
      ),
    );
  }
}

class _PerfilContent extends ConsumerStatefulWidget {
  const _PerfilContent({required this.state});

  final AuthAuthenticated state;

  @override
  ConsumerState<_PerfilContent> createState() => _PerfilContentState();
}

class _PerfilContentState extends ConsumerState<_PerfilContent> {
  bool _cambiandoTenant = false;

  Future<void> _onTenantSeleccionado(String tenantId) async {
    if (tenantId == widget.state.session.tenantId) return;
    setState(() => _cambiandoTenant = true);
    try {
      await ref.read(authControllerProvider.notifier).switchTenant(tenantId);
    } on AppException catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(e.message)));
      }
    } finally {
      if (mounted) setState(() => _cambiandoTenant = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final perfil = widget.state.perfil;
    final session = widget.state.session;
    final tenants = widget.state.tenants.tenants;
    final iniciales = perfil.nombre.isNotEmpty
        ? perfil.nombre[0].toUpperCase()
        : '?';

    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        CircleAvatar(
          radius: 32,
          child: Text(iniciales, style: const TextStyle(fontSize: 24)),
        ),
        const SizedBox(height: 12),
        Center(
          child: Text(
            '${perfil.nombre} ${perfil.apellido}',
            style: theme.textTheme.titleLarge,
            textAlign: TextAlign.center,
          ),
        ),
        Center(
          child: Text(perfil.correo, style: theme.textTheme.bodyMedium),
        ),
        const SizedBox(height: 8),
        Center(child: Chip(label: Text(session.rol))),
        const SizedBox(height: 24),
        if (tenants.length > 1) ...[
          Text('Tenant activo', style: theme.textTheme.titleMedium),
          const SizedBox(height: 4),
          if (_cambiandoTenant) const LinearProgressIndicator(),
          RadioGroup<String>(
            groupValue: session.tenantId,
            onChanged: _cambiandoTenant
                ? (_) {}
                : (value) {
                    if (value != null) _onTenantSeleccionado(value);
                  },
            child: Column(
              children: tenants
                  .map(
                    (tenant) => RadioListTile<String>(
                      value: tenant.id,
                      title: Text(tenant.nombre),
                      subtitle: Text(tenant.rol),
                    ),
                  )
                  .toList(growable: false),
            ),
          ),
          const SizedBox(height: 24),
        ],
        OutlinedButton.icon(
          onPressed: () =>
              ref.read(authControllerProvider.notifier).logout(),
          icon: const Icon(Icons.logout),
          label: const Text('Cerrar sesión'),
        ),
      ],
    );
  }
}
