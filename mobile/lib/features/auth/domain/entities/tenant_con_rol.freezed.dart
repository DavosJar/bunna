// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'tenant_con_rol.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$TenantConRol {

 String get id; String get nombre; String get slug; String get rol; bool get esPropio;
/// Create a copy of TenantConRol
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TenantConRolCopyWith<TenantConRol> get copyWith => _$TenantConRolCopyWithImpl<TenantConRol>(this as TenantConRol, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is TenantConRol&&(identical(other.id, id) || other.id == id)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.slug, slug) || other.slug == slug)&&(identical(other.rol, rol) || other.rol == rol)&&(identical(other.esPropio, esPropio) || other.esPropio == esPropio));
}


@override
int get hashCode => Object.hash(runtimeType,id,nombre,slug,rol,esPropio);

@override
String toString() {
  return 'TenantConRol(id: $id, nombre: $nombre, slug: $slug, rol: $rol, esPropio: $esPropio)';
}


}

/// @nodoc
abstract mixin class $TenantConRolCopyWith<$Res>  {
  factory $TenantConRolCopyWith(TenantConRol value, $Res Function(TenantConRol) _then) = _$TenantConRolCopyWithImpl;
@useResult
$Res call({
 String id, String nombre, String slug, String rol, bool esPropio
});




}
/// @nodoc
class _$TenantConRolCopyWithImpl<$Res>
    implements $TenantConRolCopyWith<$Res> {
  _$TenantConRolCopyWithImpl(this._self, this._then);

  final TenantConRol _self;
  final $Res Function(TenantConRol) _then;

/// Create a copy of TenantConRol
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? nombre = null,Object? slug = null,Object? rol = null,Object? esPropio = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,slug: null == slug ? _self.slug : slug // ignore: cast_nullable_to_non_nullable
as String,rol: null == rol ? _self.rol : rol // ignore: cast_nullable_to_non_nullable
as String,esPropio: null == esPropio ? _self.esPropio : esPropio // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [TenantConRol].
extension TenantConRolPatterns on TenantConRol {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _TenantConRol value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _TenantConRol() when $default != null:
return $default(_that);case _:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _TenantConRol value)  $default,){
final _that = this;
switch (_that) {
case _TenantConRol():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _TenantConRol value)?  $default,){
final _that = this;
switch (_that) {
case _TenantConRol() when $default != null:
return $default(_that);case _:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String nombre,  String slug,  String rol,  bool esPropio)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _TenantConRol() when $default != null:
return $default(_that.id,_that.nombre,_that.slug,_that.rol,_that.esPropio);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String nombre,  String slug,  String rol,  bool esPropio)  $default,) {final _that = this;
switch (_that) {
case _TenantConRol():
return $default(_that.id,_that.nombre,_that.slug,_that.rol,_that.esPropio);case _:
  throw StateError('Unexpected subclass');

}
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String nombre,  String slug,  String rol,  bool esPropio)?  $default,) {final _that = this;
switch (_that) {
case _TenantConRol() when $default != null:
return $default(_that.id,_that.nombre,_that.slug,_that.rol,_that.esPropio);case _:
  return null;

}
}

}

/// @nodoc


class _TenantConRol implements TenantConRol {
  const _TenantConRol({required this.id, required this.nombre, required this.slug, required this.rol, required this.esPropio});
  

@override final  String id;
@override final  String nombre;
@override final  String slug;
@override final  String rol;
@override final  bool esPropio;

/// Create a copy of TenantConRol
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TenantConRolCopyWith<_TenantConRol> get copyWith => __$TenantConRolCopyWithImpl<_TenantConRol>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _TenantConRol&&(identical(other.id, id) || other.id == id)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.slug, slug) || other.slug == slug)&&(identical(other.rol, rol) || other.rol == rol)&&(identical(other.esPropio, esPropio) || other.esPropio == esPropio));
}


@override
int get hashCode => Object.hash(runtimeType,id,nombre,slug,rol,esPropio);

@override
String toString() {
  return 'TenantConRol(id: $id, nombre: $nombre, slug: $slug, rol: $rol, esPropio: $esPropio)';
}


}

/// @nodoc
abstract mixin class _$TenantConRolCopyWith<$Res> implements $TenantConRolCopyWith<$Res> {
  factory _$TenantConRolCopyWith(_TenantConRol value, $Res Function(_TenantConRol) _then) = __$TenantConRolCopyWithImpl;
@override @useResult
$Res call({
 String id, String nombre, String slug, String rol, bool esPropio
});




}
/// @nodoc
class __$TenantConRolCopyWithImpl<$Res>
    implements _$TenantConRolCopyWith<$Res> {
  __$TenantConRolCopyWithImpl(this._self, this._then);

  final _TenantConRol _self;
  final $Res Function(_TenantConRol) _then;

/// Create a copy of TenantConRol
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? nombre = null,Object? slug = null,Object? rol = null,Object? esPropio = null,}) {
  return _then(_TenantConRol(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,slug: null == slug ? _self.slug : slug // ignore: cast_nullable_to_non_nullable
as String,rol: null == rol ? _self.rol : rol // ignore: cast_nullable_to_non_nullable
as String,esPropio: null == esPropio ? _self.esPropio : esPropio // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}

// dart format on
