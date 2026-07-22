import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/error/app_exception.dart';
import 'package:mobile/core/session/auth_tokens.dart';
import 'package:mobile/core/session/session_events.dart';
import 'package:mobile/core/session/token_refresh_coordinator.dart';
import 'package:mobile/core/session/token_store.dart';
import 'package:mobile/features/auth/data/auth_api.dart';
import 'package:mobile/features/auth/data/auth_repository_impl.dart';
import 'package:mobile/features/auth/data/dtos/mis_tenants_dto.dart';
import 'package:mobile/features/auth/data/dtos/perfil_dto.dart';
import 'package:mobile/features/auth/data/dtos/permiso_dto.dart';
import 'package:mobile/features/auth/data/dtos/token_session_dto.dart';
import 'package:mocktail/mocktail.dart';

class _MockAuthApi extends Mock implements AuthApi {}

class _MockSecureStorage extends Mock implements FlutterSecureStorage {}

/// JWT sin firma real (no se verifica en el cliente, solo se decodifica el
/// payload) — igual que `parseJWT` del frontend web.
String _fakeJwt(Map<String, dynamic> payload) {
  String b64(Map<String, dynamic> json) =>
      base64Url.encode(utf8.encode(jsonEncode(json))).replaceAll('=', '');
  return '${b64({'alg': 'HS256'})}.${b64(payload)}.firma-no-verificada';
}

AuthTokens _tokensVigentes(String accessToken, {String refreshToken = 'r'}) =>
    AuthTokens(
      accessToken: accessToken,
      refreshToken: refreshToken,
      accessExpiresAt: DateTime.now().add(const Duration(hours: 1)),
    );

void main() {
  late _MockAuthApi api;
  late _MockSecureStorage secureStorage;
  late TokenStore tokenStore;
  late SessionEvents sessionEvents;
  late TokenRefreshCoordinator refreshCoordinator;
  late AuthRepositoryImpl repo;

  setUp(() {
    api = _MockAuthApi();
    secureStorage = _MockSecureStorage();
    when(
      () => secureStorage.write(
        key: any(named: 'key'),
        value: any(named: 'value'),
      ),
    ).thenAnswer((_) async {});
    when(
      () => secureStorage.delete(key: any(named: 'key')),
    ).thenAnswer((_) async {});

    tokenStore = TokenStore(storage: secureStorage);
    sessionEvents = SessionEvents();
    refreshCoordinator = TokenRefreshCoordinator(
      tokenStore: tokenStore,
      sessionEvents: sessionEvents,
      identidadBaseUrl: 'http://10.0.2.2:8080',
    );

    repo = AuthRepositoryImpl(
      api: api,
      tokenStore: tokenStore,
      refreshCoordinator: refreshCoordinator,
    );
  });

  tearDown(() => sessionEvents.dispose());

  group('login', () {
    test(
      'persiste los tokens y construye la sesión desde la respuesta',
      () async {
        final jwt = _fakeJwt({
          'sub': 'usuario-1',
          'sid': 'sesion-1',
          'tenant_id': 'tenant-1',
          'rol': 'administrador',
          'exp': 9999999999,
        });

        when(
          () => api.login(correo: 'juan@correo.com', password: 'secreto123'),
        ).thenAnswer(
          (_) async => TokenSessionDto(
            accessToken: jwt,
            refreshToken: 'refresh-token-1',
            expiresIn: 3600,
            tokenType: 'Bearer',
            usuarioId: 'usuario-1',
            tenantId: 'tenant-1',
            rol: 'administrador',
          ),
        );

        final session = await repo.login(
          correo: 'juan@correo.com',
          password: 'secreto123',
        );

        expect(session.usuarioId, 'usuario-1');
        expect(session.tenantId, 'tenant-1');
        expect(session.rol, 'administrador');
        expect(session.sesionId, 'sesion-1');
        expect(tokenStore.current?.accessToken, jwt);
        expect(tokenStore.current?.refreshToken, 'refresh-token-1');
      },
    );

    test('propaga la AppException del datasource sin envolverla', () async {
      when(
        () => api.login(
          correo: any(named: 'correo'),
          password: any(named: 'password'),
        ),
      ).thenThrow(const ValidationException('credenciales inválidas'));

      expect(
        () => repo.login(correo: 'x@x.com', password: 'malo'),
        throwsA(isA<ValidationException>()),
      );
    });
  });

  group('restoreSession', () {
    test('sin tokens guardados ⇒ null', () async {
      final result = await repo.restoreSession();
      expect(result, isNull);
    });

    test('access token vigente ⇒ reconstruye sesión sin refrescar', () async {
      final jwt = _fakeJwt({
        'sub': 'usuario-2',
        'sid': 'sesion-2',
        'tenant_id': 'tenant-2',
        'rol': 'caficultor',
        'exp': 9999999999,
      });
      await tokenStore.save(_tokensVigentes(jwt));

      final result = await repo.restoreSession();

      expect(result?.usuarioId, 'usuario-2');
      expect(result?.tenantId, 'tenant-2');
      expect(result?.rol, 'caficultor');

      // No debió intentar refrescar: el access token seguía vigente.
      verifyNever(
        () => api.login(
          correo: any(named: 'correo'),
          password: any(named: 'password'),
        ),
      );
    });
  });

  group('logout', () {
    test('limpia el TokenStore aunque el servidor falle', () async {
      final jwt = _fakeJwt({'sub': 'u', 'exp': 9999999999});
      await tokenStore.save(_tokensVigentes(jwt));

      when(() => api.logout()).thenThrow(const NetworkException());

      await repo.logout();

      expect(tokenStore.current, isNull);
    });
  });

  group('getMiPerfil / getMisTenants / getMisPermisos', () {
    test('mapean el DTO a la entidad de dominio', () async {
      when(() => api.getMiPerfil()).thenAnswer(
        (_) async => PerfilDto(
          id: 'u1',
          correo: 'juan@correo.com',
          nombre: 'Juan',
          apellido: 'Pérez',
          telefono: '0999999999',
          estado: 'ACTIVO',
          creadoEn: DateTime.utc(2026, 5, 23, 12),
        ),
      );
      when(() => api.getMisTenants()).thenAnswer(
        (_) async => const MisTenantsDto(
          tenants: [
            TenantConRolDto(
              id: 't1',
              nombre: 'Mi Finca',
              slug: 'mi-finca',
              rol: 'administrador',
              esPropio: true,
            ),
          ],
          propioId: 't1',
        ),
      );
      when(() => api.getMisPermisos()).thenAnswer(
        (_) async => const MisPermisosDto(
          permisos: [
            PermisoDto(
              codigo: 'fincas:finca:crear',
              nombre: 'Crear Finca',
              descripcion: 'Registrar una nueva finca',
              modulo: 'fincas',
            ),
          ],
        ),
      );

      final perfil = await repo.getMiPerfil();
      final tenants = await repo.getMisTenants();
      final permisos = await repo.getMisPermisos();

      expect(perfil.nombre, 'Juan');
      expect(tenants.propioId, 't1');
      expect(tenants.tenants.single.esPropio, isTrue);
      expect(permisos.single.codigo, 'fincas:finca:crear');
    });
  });
}
