// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'permiso.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$Permiso {

 String get codigo; String get nombre; String get descripcion; String get modulo;
/// Create a copy of Permiso
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PermisoCopyWith<Permiso> get copyWith => _$PermisoCopyWithImpl<Permiso>(this as Permiso, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Permiso&&(identical(other.codigo, codigo) || other.codigo == codigo)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.modulo, modulo) || other.modulo == modulo));
}


@override
int get hashCode => Object.hash(runtimeType,codigo,nombre,descripcion,modulo);

@override
String toString() {
  return 'Permiso(codigo: $codigo, nombre: $nombre, descripcion: $descripcion, modulo: $modulo)';
}


}

/// @nodoc
abstract mixin class $PermisoCopyWith<$Res>  {
  factory $PermisoCopyWith(Permiso value, $Res Function(Permiso) _then) = _$PermisoCopyWithImpl;
@useResult
$Res call({
 String codigo, String nombre, String descripcion, String modulo
});




}
/// @nodoc
class _$PermisoCopyWithImpl<$Res>
    implements $PermisoCopyWith<$Res> {
  _$PermisoCopyWithImpl(this._self, this._then);

  final Permiso _self;
  final $Res Function(Permiso) _then;

/// Create a copy of Permiso
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? codigo = null,Object? nombre = null,Object? descripcion = null,Object? modulo = null,}) {
  return _then(_self.copyWith(
codigo: null == codigo ? _self.codigo : codigo // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,modulo: null == modulo ? _self.modulo : modulo // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [Permiso].
extension PermisoPatterns on Permiso {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Permiso value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Permiso() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Permiso value)  $default,){
final _that = this;
switch (_that) {
case _Permiso():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Permiso value)?  $default,){
final _that = this;
switch (_that) {
case _Permiso() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String codigo,  String nombre,  String descripcion,  String modulo)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Permiso() when $default != null:
return $default(_that.codigo,_that.nombre,_that.descripcion,_that.modulo);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String codigo,  String nombre,  String descripcion,  String modulo)  $default,) {final _that = this;
switch (_that) {
case _Permiso():
return $default(_that.codigo,_that.nombre,_that.descripcion,_that.modulo);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String codigo,  String nombre,  String descripcion,  String modulo)?  $default,) {final _that = this;
switch (_that) {
case _Permiso() when $default != null:
return $default(_that.codigo,_that.nombre,_that.descripcion,_that.modulo);case _:
  return null;

}
}

}

/// @nodoc


class _Permiso implements Permiso {
  const _Permiso({required this.codigo, required this.nombre, required this.descripcion, required this.modulo});
  

@override final  String codigo;
@override final  String nombre;
@override final  String descripcion;
@override final  String modulo;

/// Create a copy of Permiso
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PermisoCopyWith<_Permiso> get copyWith => __$PermisoCopyWithImpl<_Permiso>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Permiso&&(identical(other.codigo, codigo) || other.codigo == codigo)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.modulo, modulo) || other.modulo == modulo));
}


@override
int get hashCode => Object.hash(runtimeType,codigo,nombre,descripcion,modulo);

@override
String toString() {
  return 'Permiso(codigo: $codigo, nombre: $nombre, descripcion: $descripcion, modulo: $modulo)';
}


}

/// @nodoc
abstract mixin class _$PermisoCopyWith<$Res> implements $PermisoCopyWith<$Res> {
  factory _$PermisoCopyWith(_Permiso value, $Res Function(_Permiso) _then) = __$PermisoCopyWithImpl;
@override @useResult
$Res call({
 String codigo, String nombre, String descripcion, String modulo
});




}
/// @nodoc
class __$PermisoCopyWithImpl<$Res>
    implements _$PermisoCopyWith<$Res> {
  __$PermisoCopyWithImpl(this._self, this._then);

  final _Permiso _self;
  final $Res Function(_Permiso) _then;

/// Create a copy of Permiso
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? codigo = null,Object? nombre = null,Object? descripcion = null,Object? modulo = null,}) {
  return _then(_Permiso(
codigo: null == codigo ? _self.codigo : codigo // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,modulo: null == modulo ? _self.modulo : modulo // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}

// dart format on
