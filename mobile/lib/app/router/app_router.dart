import 'package:flutter/foundation.dart';
import 'package:go_router/go_router.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../features/auth/presentation/auth_controller.dart';
import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/auth/presentation/pages/perfil_page.dart';
import '../../features/auth/presentation/pages/splash_page.dart';
import '../../features/diagnosticos/presentation/pages/diagnostico_page.dart';
import '../../features/diagnosticos/presentation/pages/quick_scan_page.dart';
import '../../features/fincas/domain/entities/finca.dart';
import '../../features/fincas/presentation/pages/fincas_list_page.dart';
import '../../features/lotes/domain/entities/lote.dart';
import '../../features/lotes/presentation/pages/finca_detail_page.dart';
import '../../features/muestras/presentation/pages/lote_detail_page.dart';
import 'routes.dart';
import 'shell_scaffold.dart';

part 'app_router.g.dart';

/// Puente entre `AuthController` (Riverpod) y `GoRouter.refreshListenable`.
final class _AuthRefreshNotifier extends ChangeNotifier {
  _AuthRefreshNotifier(Ref ref) {
    ref.listen(authControllerProvider, (_, _) => notifyListeners());
  }
}

@Riverpod(keepAlive: true)
GoRouter appRouter(Ref ref) {
  final refreshNotifier = _AuthRefreshNotifier(ref);
  ref.onDispose(refreshNotifier.dispose);

  return GoRouter(
    initialLocation: AppRoutes.splash,
    refreshListenable: refreshNotifier,
    redirect: (context, state) => _redirect(ref, state),
    routes: [
      GoRoute(
        path: AppRoutes.splash,
        builder: (context, state) => const SplashPage(),
      ),
      GoRoute(
        path: AppRoutes.login,
        builder: (context, state) => const LoginPage(),
      ),
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) =>
            ShellScaffold(navigationShell: navigationShell),
        branches: [
          // Rama 0 — Fincas (con sub-rutas anidadas: finca → lote → diagnóstico)
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.fincas,
                builder: (context, state) => const FincasListPage(),
                routes: [
                  GoRoute(
                    path: ':fincaId',
                    builder: (context, state) => FincaDetailPage(
                      fincaId: state.pathParameters['fincaId']!,
                      finca: state.extra as Finca?,
                    ),
                    routes: [
                      GoRoute(
                        path: 'lotes/:loteId',
                        builder: (context, state) => LoteDetailPage(
                          fincaId: state.pathParameters['fincaId']!,
                          loteId: state.pathParameters['loteId']!,
                          lote: state.extra as Lote?,
                        ),
                        routes: [
                          GoRoute(
                            path: 'muestras/:muestraId/diagnostico',
                            builder: (context, state) => DiagnosticoPage(
                              muestraId: state.pathParameters['muestraId']!,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ),
          // Rama 1 — Escanear (YOLO rápido)
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.escanear,
                builder: (context, state) => const QuickScanPage(),
              ),
            ],
          ),
          // Rama 2 — Perfil
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: AppRoutes.perfil,
                builder: (context, state) => const PerfilPage(),
              ),
            ],
          ),
        ],
      ),
    ],
  );
}

String? _redirect(Ref ref, GoRouterState routerState) {
  final authAsync = ref.read(authControllerProvider);
  final location = routerState.matchedLocation;
  final onSplash = location == AppRoutes.splash;
  final onLogin = location == AppRoutes.login;

  if (authAsync.isLoading) {
    return onSplash ? null : AppRoutes.splash;
  }

  final authState = authAsync.value ?? const AuthState.unauthenticated();

  return switch (authState) {
    AuthUnknown() => onSplash ? null : AppRoutes.splash,
    AuthUnauthenticated() => onLogin ? null : AppRoutes.login,
    AuthAuthenticated() => (onLogin || onSplash) ? AppRoutes.fincas : null,
  };
}
