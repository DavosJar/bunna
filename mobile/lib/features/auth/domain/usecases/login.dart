import '../auth_repository.dart';
import '../entities/auth_session.dart';

/// Domain puro: sin Riverpod ni ningún import de framework — la regla de
/// capas (`domain` no depende de `data` ni de `presentation`) se cumple
/// literalmente aquí. El wiring con Riverpod vive en
/// `presentation/auth_usecases_providers.dart`.
final class Login {
  const Login(this._repo);

  final AuthRepository _repo;

  Future<AuthSession> call({
    required String correo,
    required String password,
  }) => _repo.login(correo: correo, password: password);
}
