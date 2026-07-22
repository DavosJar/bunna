import 'package:geolocator/geolocator.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../error/app_exception.dart';

part 'location_service.g.dart';

/// Coordenada capturada del dispositivo.
typedef Coordenada = ({double latitud, double longitud});

/// Envuelve geolocator y traduce sus fallos a `AppException` para que la UI
/// los muestre igual que cualquier otro error.
class LocationService {
  const LocationService();

  Future<Coordenada> obtenerUbicacion() async {
    final habilitado = await Geolocator.isLocationServiceEnabled();
    if (!habilitado) {
      throw const ValidationException(
        'El GPS está desactivado. Actívalo para tomar la muestra.',
      );
    }

    var permiso = await Geolocator.checkPermission();
    if (permiso == LocationPermission.denied) {
      permiso = await Geolocator.requestPermission();
    }
    if (permiso == LocationPermission.denied) {
      throw const ValidationException(
        'Permiso de ubicación denegado. Concédelo para continuar.',
      );
    }
    if (permiso == LocationPermission.deniedForever) {
      throw const ValidationException(
        'El permiso de ubicación está bloqueado. Habilítalo en ajustes.',
      );
    }

    try {
      final pos = await Geolocator.getCurrentPosition(
        locationSettings: const LocationSettings(
          accuracy: LocationAccuracy.high,
        ),
      );
      return (latitud: pos.latitude, longitud: pos.longitude);
    } catch (e) {
      throw NetworkException('No se pudo obtener la ubicación GPS', e);
    }
  }
}

@Riverpod(keepAlive: true)
LocationService locationService(Ref ref) => const LocationService();
