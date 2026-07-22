import '../auth_repository.dart';

final class Logout {
  const Logout(this._repo);

  final AuthRepository _repo;

  Future<void> call({bool todasLasSesiones = false}) =>
      _repo.logout(todasLasSesiones: todasLasSesiones);
}
