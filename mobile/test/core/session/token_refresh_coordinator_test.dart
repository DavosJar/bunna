import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/error/app_exception.dart';
import 'package:mobile/core/session/auth_tokens.dart';
import 'package:mobile/core/session/session_events.dart';
import 'package:mobile/core/session/token_refresh_coordinator.dart';
import 'package:mobile/core/session/token_store.dart';
import 'package:mocktail/mocktail.dart';

class _MockSecureStorage extends Mock implements FlutterSecureStorage {}

class _MockDio extends Mock implements Dio {}

void main() {
  late _MockSecureStorage secureStorage;
  late TokenStore tokenStore;
  late SessionEvents sessionEvents;
  late _MockDio bareDio;
  late TokenRefreshCoordinator coordinator;

  final tokensViejos = AuthTokens(
    accessToken: 'access-viejo',
    refreshToken: 'refresh-viejo',
    accessExpiresAt: DateTime.now().add(const Duration(hours: 1)),
  );

  setUp(() async {
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
    await tokenStore.save(tokensViejos);

    sessionEvents = SessionEvents();
    bareDio = _MockDio();

    coordinator = TokenRefreshCoordinator(
      tokenStore: tokenStore,
      sessionEvents: sessionEvents,
      identidadBaseUrl: 'http://10.0.2.2:8080',
      bareDioOverride: bareDio,
    );
  });

  tearDown(() {
    sessionEvents.dispose();
  });

  test(
    'dos refresh concurrentes con el mismo access token expirado disparan un solo POST',
    () async {
      var postCount = 0;
      when(
        () => bareDio.post<Map<String, dynamic>>(
          any(),
          data: any(named: 'data'),
        ),
      ).thenAnswer((_) async {
        postCount++;
        // Latencia simulada para forzar el solapamiento de las dos llamadas.
        await Future<void>.delayed(const Duration(milliseconds: 30));
        return Response(
          requestOptions: RequestOptions(
            path: '/api/v1/identidad/auth/refresh',
          ),
          statusCode: 200,
          data: {
            'data': {
              'access_token': 'access-nuevo',
              'refresh_token': 'refresh-nuevo',
              'expires_in': 3600,
            },
          },
        );
      });

      final resultados = await Future.wait([
        coordinator.refresh(failedAccessToken: 'access-viejo'),
        coordinator.refresh(failedAccessToken: 'access-viejo'),
      ]);

      expect(postCount, 1, reason: 'debe haber un único POST /auth/refresh');
      expect(resultados[0].accessToken, 'access-nuevo');
      expect(resultados[1].accessToken, 'access-nuevo');
      expect(tokenStore.current?.accessToken, 'access-nuevo');
      expect(tokenStore.current?.refreshToken, 'refresh-nuevo');
    },
  );

  test(
    'si otro request ya rotó los tokens mientras este esperaba, no dispara ningún POST',
    () async {
      final tokensYaRotados = AuthTokens(
        accessToken: 'access-ya-rotado',
        refreshToken: 'refresh-ya-rotado',
        accessExpiresAt: DateTime.now().add(const Duration(hours: 1)),
      );
      await tokenStore.save(tokensYaRotados);

      final resultado = await coordinator.refresh(
        failedAccessToken: 'access-viejo', // el token que causó el 401 original
      );

      expect(resultado.accessToken, 'access-ya-rotado');
      verifyNever(
        () => bareDio.post<Map<String, dynamic>>(
          any(),
          data: any(named: 'data'),
        ),
      );
    },
  );

  test(
    'refresh token inválido (400) limpia el TokenStore y emite sessionExpired',
    () async {
      when(
        () => bareDio.post<Map<String, dynamic>>(
          any(),
          data: any(named: 'data'),
        ),
      ).thenThrow(
        DioException(
          requestOptions: RequestOptions(
            path: '/api/v1/identidad/auth/refresh',
          ),
          response: Response(
            requestOptions: RequestOptions(
              path: '/api/v1/identidad/auth/refresh',
            ),
            statusCode: 400,
          ),
        ),
      );

      final eventosFuturos = sessionEvents.stream.toList();

      await expectLater(
        () => coordinator.refresh(failedAccessToken: 'access-viejo'),
        throwsA(isA<SessionExpiredException>()),
      );

      expect(tokenStore.current, isNull);
      sessionEvents.dispose();
      final eventos = await eventosFuturos;
      expect(eventos, [SessionEventType.expired]);
    },
  );

  test('fallo de red durante el refresh NO limpia el TokenStore', () async {
    when(
      () => bareDio.post<Map<String, dynamic>>(
        any(),
        data: any(named: 'data'),
      ),
    ).thenThrow(
      DioException(
        requestOptions: RequestOptions(path: '/api/v1/identidad/auth/refresh'),
        type: DioExceptionType.connectionError,
      ),
    );

    await expectLater(
      () => coordinator.refresh(failedAccessToken: 'access-viejo'),
      throwsA(isA<NetworkException>()),
    );

    expect(tokenStore.current?.accessToken, 'access-viejo');
  });
}
