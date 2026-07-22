import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'app/app.dart';
import 'core/network/dio_providers.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // El TokenStore debe tener los tokens cargados en memoria ANTES de que
  // AuthController.build() corra restoreSession() — por eso se resuelve acá,
  // fuera del árbol de widgets, con un ProviderContainer explícito.
  final container = ProviderContainer();
  await container.read(tokenStoreProvider).loadFromDisk();

  runApp(
    UncontrolledProviderScope(container: container, child: const BunnaApp()),
  );
}
