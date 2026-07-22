import 'dart:async';

/// Eventos de sesión emitidos desde infraestructura de red (p. ej. el
/// `TokenRefreshCoordinator` cuando el refresh token ya no es válido) que
/// `AuthController` escucha para forzar el logout local y que GoRouter
/// redirija a `/login`.
enum SessionEventType { expired }

final class SessionEvents {
  final StreamController<SessionEventType> _controller =
      StreamController<SessionEventType>.broadcast();

  Stream<SessionEventType> get stream => _controller.stream;

  void emitExpired() => _controller.add(SessionEventType.expired);

  void dispose() => _controller.close();
}
