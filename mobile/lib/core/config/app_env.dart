/// Configuración de entorno: base URLs de los tres backends.
///
/// Los defaults apuntan al emulador Android (`10.0.2.2` = loopback del host
/// visto desde el emulador). Para dispositivo físico o iOS simulator,
/// sobrescribir con `--dart-define`:
///   flutter run \
///     --dart-define=IDENTIDAD_BASE_URL=http://192.168.1.50:8080 \
///     --dart-define=FINCAS_BASE_URL=http://192.168.1.50:8082
abstract final class AppEnv {
  static const String identidadBaseUrl = String.fromEnvironment(
    'IDENTIDAD_BASE_URL',
    defaultValue: 'http://10.0.2.2:8080',
  );

  static const String fincasBaseUrl = String.fromEnvironment(
    'FINCAS_BASE_URL',
    defaultValue: 'http://10.0.2.2:8082',
  );

  static const String yoloBaseUrl = String.fromEnvironment(
    'YOLO_BASE_URL',
    defaultValue: 'https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com',
  );
}
