import '../auth_repository.dart';
import '../entities/auth_session.dart';

final class SwitchTenant {
  const SwitchTenant(this._repo);

  final AuthRepository _repo;

  Future<AuthSession> call(String tenantId) => _repo.switchTenant(tenantId);
}
