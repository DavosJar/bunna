import 'package:flutter/material.dart';

import '../../../../core/widgets/loading_view.dart';

/// Se muestra mientras `AuthController.build()` resuelve `restoreSession()`.
/// GoRouter redirige fuera de aquí en cuanto el estado deja de ser
/// `AsyncLoading`/`unknown` (ver `app/router/app_router.dart`).
class SplashPage extends StatelessWidget {
  const SplashPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(body: LoadingView(message: 'Cargando sesión…'));
  }
}
