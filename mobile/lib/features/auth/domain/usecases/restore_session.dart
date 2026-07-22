import '../auth_repository.dart';
import '../entities/auth_session.dart';

final class RestoreSession {
  const RestoreSession(this._repo);

  final AuthRepository _repo;

  Future<AuthSession?> call() => _repo.restoreSession();
}
