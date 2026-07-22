import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/error/app_exception.dart';
import 'package:mobile/features/auth/data/auth_providers.dart';
import 'package:mobile/features/auth/domain/auth_repository.dart';
import 'package:mobile/features/auth/domain/entities/auth_session.dart';
import 'package:mobile/features/auth/domain/entities/mis_tenants.dart';
import 'package:mobile/features/auth/domain/entities/perfil.dart';
import 'package:mobile/features/auth/presentation/pages/login_page.dart';
import 'package:mocktail/mocktail.dart';

class _MockAuthRepository extends Mock implements AuthRepository {}

void main() {
  late _MockAuthRepository repo;

  setUp(() {
    repo = _MockAuthRepository();
    // AuthController.build() llama restoreSession() al montar: sin sesión
    // previa para que el widget arranque en estado no-autenticado.
    when(() => repo.restoreSession()).thenAnswer((_) async => null);
  });

  Future<void> pumpLoginPage(WidgetTester tester) async {
    await tester.pumpWidget(
      ProviderScope(
        overrides: [authRepositoryProvider.overrideWithValue(repo)],
        child: const MaterialApp(home: LoginPage()),
      ),
    );
    await tester.pumpAndSettle();
  }

  testWidgets('login exitoso: llama al repositorio con las credenciales', (
    tester,
  ) async {
    when(
      () => repo.login(correo: 'juan@correo.com', password: 'secreto123'),
    ).thenAnswer(
      (_) async => const AuthSession(
        usuarioId: 'u1',
        tenantId: 't1',
        rol: 'caficultor',
        sesionId: 's1',
      ),
    );
    when(() => repo.getMiPerfil()).thenAnswer(
      (_) async => Perfil(
        id: 'u1',
        correo: 'juan@correo.com',
        nombre: 'Juan',
        apellido: 'Pérez',
        telefono: '',
        estado: 'ACTIVO',
        creadoEn: DateTime.utc(2026),
      ),
    );
    when(
      () => repo.getMisTenants(),
    ).thenAnswer((_) async => const MisTenants(tenants: [], propioId: 't1'));
    when(() => repo.getMisPermisos()).thenAnswer((_) async => []);

    await pumpLoginPage(tester);

    await tester.enterText(
      find.widgetWithText(TextFormField, 'Correo electrónico'),
      'juan@correo.com',
    );
    await tester.enterText(
      find.widgetWithText(TextFormField, 'Contraseña'),
      'secreto123',
    );
    await tester.tap(find.widgetWithText(FilledButton, 'Ingresar'));
    await tester.pumpAndSettle();

    verify(
      () => repo.login(correo: 'juan@correo.com', password: 'secreto123'),
    ).called(1);
    // Sin banner de error tras un login exitoso.
    expect(find.byIcon(Icons.error_outline), findsNothing);
  });

  testWidgets('credenciales inválidas: muestra el mensaje de error inline', (
    tester,
  ) async {
    when(
      () => repo.login(
        correo: any(named: 'correo'),
        password: any(named: 'password'),
      ),
    ).thenThrow(const ValidationException('credenciales inválidas'));

    await pumpLoginPage(tester);

    await tester.enterText(
      find.widgetWithText(TextFormField, 'Correo electrónico'),
      'juan@correo.com',
    );
    await tester.enterText(
      find.widgetWithText(TextFormField, 'Contraseña'),
      'incorrecta',
    );
    await tester.tap(find.widgetWithText(FilledButton, 'Ingresar'));
    await tester.pumpAndSettle();

    expect(find.text('credenciales inválidas'), findsOneWidget);
  });

  testWidgets('no envía el formulario si los campos están vacíos', (
    tester,
  ) async {
    await pumpLoginPage(tester);

    await tester.tap(find.widgetWithText(FilledButton, 'Ingresar'));
    await tester.pumpAndSettle();

    verifyNever(
      () => repo.login(
        correo: any(named: 'correo'),
        password: any(named: 'password'),
      ),
    );
    expect(find.text('Ingresa tu correo'), findsOneWidget);
    expect(find.text('Ingresa tu contraseña'), findsOneWidget);
  });
}
