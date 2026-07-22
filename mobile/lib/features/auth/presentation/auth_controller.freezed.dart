// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'auth_controller.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$AuthState {





@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is AuthState);
}


@override
int get hashCode => runtimeType.hashCode;

@override
String toString() {
  return 'AuthState()';
}


}

/// @nodoc
class $AuthStateCopyWith<$Res>  {
$AuthStateCopyWith(AuthState _, $Res Function(AuthState) __);
}


/// Adds pattern-matching-related methods to [AuthState].
extension AuthStatePatterns on AuthState {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>({TResult Function( AuthUnknown value)?  unknown,TResult Function( AuthUnauthenticated value)?  unauthenticated,TResult Function( AuthAuthenticated value)?  authenticated,required TResult orElse(),}){
final _that = this;
switch (_that) {
case AuthUnknown() when unknown != null:
return unknown(_that);case AuthUnauthenticated() when unauthenticated != null:
return unauthenticated(_that);case AuthAuthenticated() when authenticated != null:
return authenticated(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>({required TResult Function( AuthUnknown value)  unknown,required TResult Function( AuthUnauthenticated value)  unauthenticated,required TResult Function( AuthAuthenticated value)  authenticated,}){
final _that = this;
switch (_that) {
case AuthUnknown():
return unknown(_that);case AuthUnauthenticated():
return unauthenticated(_that);case AuthAuthenticated():
return authenticated(_that);}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>({TResult? Function( AuthUnknown value)?  unknown,TResult? Function( AuthUnauthenticated value)?  unauthenticated,TResult? Function( AuthAuthenticated value)?  authenticated,}){
final _that = this;
switch (_that) {
case AuthUnknown() when unknown != null:
return unknown(_that);case AuthUnauthenticated() when unauthenticated != null:
return unauthenticated(_that);case AuthAuthenticated() when authenticated != null:
return authenticated(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>({TResult Function()?  unknown,TResult Function()?  unauthenticated,TResult Function( AuthSession session,  Perfil perfil,  MisTenants tenants,  List<Permiso> permisos)?  authenticated,required TResult orElse(),}) {final _that = this;
switch (_that) {
case AuthUnknown() when unknown != null:
return unknown();case AuthUnauthenticated() when unauthenticated != null:
return unauthenticated();case AuthAuthenticated() when authenticated != null:
return authenticated(_that.session,_that.perfil,_that.tenants,_that.permisos);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>({required TResult Function()  unknown,required TResult Function()  unauthenticated,required TResult Function( AuthSession session,  Perfil perfil,  MisTenants tenants,  List<Permiso> permisos)  authenticated,}) {final _that = this;
switch (_that) {
case AuthUnknown():
return unknown();case AuthUnauthenticated():
return unauthenticated();case AuthAuthenticated():
return authenticated(_that.session,_that.perfil,_that.tenants,_that.permisos);}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>({TResult? Function()?  unknown,TResult? Function()?  unauthenticated,TResult? Function( AuthSession session,  Perfil perfil,  MisTenants tenants,  List<Permiso> permisos)?  authenticated,}) {final _that = this;
switch (_that) {
case AuthUnknown() when unknown != null:
return unknown();case AuthUnauthenticated() when unauthenticated != null:
return unauthenticated();case AuthAuthenticated() when authenticated != null:
return authenticated(_that.session,_that.perfil,_that.tenants,_that.permisos);case _:
  return null;

}
}

}

/// @nodoc


class AuthUnknown implements AuthState {
  const AuthUnknown();
  






@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is AuthUnknown);
}


@override
int get hashCode => runtimeType.hashCode;

@override
String toString() {
  return 'AuthState.unknown()';
}


}




/// @nodoc


class AuthUnauthenticated implements AuthState {
  const AuthUnauthenticated();
  






@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is AuthUnauthenticated);
}


@override
int get hashCode => runtimeType.hashCode;

@override
String toString() {
  return 'AuthState.unauthenticated()';
}


}




/// @nodoc


class AuthAuthenticated implements AuthState {
  const AuthAuthenticated({required this.session, required this.perfil, required this.tenants, required final  List<Permiso> permisos}): _permisos = permisos;
  

 final  AuthSession session;
 final  Perfil perfil;
 final  MisTenants tenants;
 final  List<Permiso> _permisos;
 List<Permiso> get permisos {
  if (_permisos is EqualUnmodifiableListView) return _permisos;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_permisos);
}


/// Create a copy of AuthState
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$AuthAuthenticatedCopyWith<AuthAuthenticated> get copyWith => _$AuthAuthenticatedCopyWithImpl<AuthAuthenticated>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is AuthAuthenticated&&(identical(other.session, session) || other.session == session)&&(identical(other.perfil, perfil) || other.perfil == perfil)&&(identical(other.tenants, tenants) || other.tenants == tenants)&&const DeepCollectionEquality().equals(other._permisos, _permisos));
}


@override
int get hashCode => Object.hash(runtimeType,session,perfil,tenants,const DeepCollectionEquality().hash(_permisos));

@override
String toString() {
  return 'AuthState.authenticated(session: $session, perfil: $perfil, tenants: $tenants, permisos: $permisos)';
}


}

/// @nodoc
abstract mixin class $AuthAuthenticatedCopyWith<$Res> implements $AuthStateCopyWith<$Res> {
  factory $AuthAuthenticatedCopyWith(AuthAuthenticated value, $Res Function(AuthAuthenticated) _then) = _$AuthAuthenticatedCopyWithImpl;
@useResult
$Res call({
 AuthSession session, Perfil perfil, MisTenants tenants, List<Permiso> permisos
});


$AuthSessionCopyWith<$Res> get session;$PerfilCopyWith<$Res> get perfil;$MisTenantsCopyWith<$Res> get tenants;

}
/// @nodoc
class _$AuthAuthenticatedCopyWithImpl<$Res>
    implements $AuthAuthenticatedCopyWith<$Res> {
  _$AuthAuthenticatedCopyWithImpl(this._self, this._then);

  final AuthAuthenticated _self;
  final $Res Function(AuthAuthenticated) _then;

/// Create a copy of AuthState
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') $Res call({Object? session = null,Object? perfil = null,Object? tenants = null,Object? permisos = null,}) {
  return _then(AuthAuthenticated(
session: null == session ? _self.session : session // ignore: cast_nullable_to_non_nullable
as AuthSession,perfil: null == perfil ? _self.perfil : perfil // ignore: cast_nullable_to_non_nullable
as Perfil,tenants: null == tenants ? _self.tenants : tenants // ignore: cast_nullable_to_non_nullable
as MisTenants,permisos: null == permisos ? _self._permisos : permisos // ignore: cast_nullable_to_non_nullable
as List<Permiso>,
  ));
}

/// Create a copy of AuthState
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$AuthSessionCopyWith<$Res> get session {
  
  return $AuthSessionCopyWith<$Res>(_self.session, (value) {
    return _then(_self.copyWith(session: value));
  });
}/// Create a copy of AuthState
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$PerfilCopyWith<$Res> get perfil {
  
  return $PerfilCopyWith<$Res>(_self.perfil, (value) {
    return _then(_self.copyWith(perfil: value));
  });
}/// Create a copy of AuthState
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$MisTenantsCopyWith<$Res> get tenants {
  
  return $MisTenantsCopyWith<$Res>(_self.tenants, (value) {
    return _then(_self.copyWith(tenants: value));
  });
}
}

// dart format on
