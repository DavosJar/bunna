/// Desenvuelve el sobre `{"data": ..., "_links"/"links": ...}` que ambos
/// backends usan en sus respuestas exitosas (identidad: `_links`; fincas:
/// `links`). Los links HATEOAS se ignoran a propósito — la app navega con
/// GoRouter, no los necesita.
///
/// Excepción documentada: `GET /nodos` de fincas viene doble-anidado
/// (`data.data`) — ese caso se desenvuelve a mano en su propio datasource,
/// no aquí.
final class ApiEnvelope {
  const ApiEnvelope._();

  static T unwrap<T>(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) {
    return fromJsonT(json['data']);
  }
}
