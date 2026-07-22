import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/error/app_exception.dart';
import 'package:mobile/features/auth/domain/auth_repository.dart';
import 'package:mobile/features/auth/domain/entities/mis_tenants.dart';
import 'package:mobile/features/auth/domain/entities/perfil.dart';
import 'package:mobile/features/auth/domain/entities/permiso.dart';
import 'package:mobile/features/auth/domain/usecases/cargar_contexto_usuario.dart';
import 'package:mocktail/mocktail.dart';

class _MockAuthRepository extends Mock implements AuthRepository {}

Perfil _perfil() => Perfil(
  id: 'u1',
  correo: 'admin@bunna.com',
  nombre: 'Admin',
  apellido: 'Bunna',
  telefono: '',
  estado: 'ACTIVO',
  creadoEn: DateTime.utc(2026),
);

void main() {
  late _MockAuthRepository repo;
  late CargarContextoUsuario cargarContexto;

  setUp(() {
    repo = _MockAuthRepository();
    cargarContexto = CargarContextoUsuario(repo);
    when(() => repo.getMiPerfil()).thenAnswer((_) async => _perfil());
  });

  test(
    'caso sys_admin real: mis-permisos 401 con el mismo token que mi-perfil '
    'acepta ⇒ degrada a lista vacía, NO aborta el bootstrap',
    () async {
      when(() => repo.getMisTenants()).thenAnswer(
        (_) async => const MisTenants(tenants: [], propioId: ''),
      );
      when(
        () => repo.getMisPermisos(),
      ).thenThrow(const SessionExpiredException('token requerido'));

      final contexto = await cargarContexto();

      expect(contexto.perfil.correo, 'admin@bunna.com');
      expect(contexto.permisos, isEmpty);
      expect(contexto.tenants.tenants, isEmpty);
    },
  );

  test('mis-tenants también degrada a vacío si falla', () async {
    when(
      () => repo.getMisTenants(),
    ).thenThrow(const NetworkException('sin conexión'));
    when(() => repo.getMisPermisos()).thenAnswer((_) async => const []);

    final contexto = await cargarContexto();

    expect(contexto.tenants, const MisTenants(tenants: [], propioId: ''));
  });

  test('camino feliz: perfil, tenants y permisos se propagan tal cual', () async {
    when(() => repo.getMisTenants()).thenAnswer(
      (_) async => const MisTenants(
        tenants: [],
        propioId: 't1',
      ),
    );
    when(() => repo.getMisPermisos()).thenAnswer(
      (_) async => const [
        Permiso(
          codigo: 'fincas:finca:crear',
          nombre: 'Crear Finca',
          descripcion: '',
          modulo: 'fincas',
        ),
      ],
    );

    final contexto = await cargarContexto();

    expect(contexto.tenants.propioId, 't1');
    expect(contexto.permisos.single.codigo, 'fincas:finca:crear');
  });

  test(
    'si getMiPerfil() falla, el error se propaga sin degradar (sesión sí inválida)',
    () async {
      when(() => repo.getMiPerfil()).thenThrow(
        const SessionExpiredException('token inválido'),
      );

      expect(
        () => cargarContexto(),
        throwsA(isA<SessionExpiredException>()),
      );
    },
  );
}
