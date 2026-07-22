// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'finca.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$Finca {

 String get id; String get nombre; String get ubicacion; String get descripcion; String get estado; DateTime get createdAt;
/// Create a copy of Finca
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$FincaCopyWith<Finca> get copyWith => _$FincaCopyWithImpl<Finca>(this as Finca, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Finca&&(identical(other.id, id) || other.id == id)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.ubicacion, ubicacion) || other.ubicacion == ubicacion)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,nombre,ubicacion,descripcion,estado,createdAt);

@override
String toString() {
  return 'Finca(id: $id, nombre: $nombre, ubicacion: $ubicacion, descripcion: $descripcion, estado: $estado, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $FincaCopyWith<$Res>  {
  factory $FincaCopyWith(Finca value, $Res Function(Finca) _then) = _$FincaCopyWithImpl;
@useResult
$Res call({
 String id, String nombre, String ubicacion, String descripcion, String estado, DateTime createdAt
});




}
/// @nodoc
class _$FincaCopyWithImpl<$Res>
    implements $FincaCopyWith<$Res> {
  _$FincaCopyWithImpl(this._self, this._then);

  final Finca _self;
  final $Res Function(Finca) _then;

/// Create a copy of Finca
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? nombre = null,Object? ubicacion = null,Object? descripcion = null,Object? estado = null,Object? createdAt = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,ubicacion: null == ubicacion ? _self.ubicacion : ubicacion // ignore: cast_nullable_to_non_nullable
as String,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [Finca].
extension FincaPatterns on Finca {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Finca value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Finca() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Finca value)  $default,){
final _that = this;
switch (_that) {
case _Finca():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Finca value)?  $default,){
final _that = this;
switch (_that) {
case _Finca() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String nombre,  String ubicacion,  String descripcion,  String estado,  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Finca() when $default != null:
return $default(_that.id,_that.nombre,_that.ubicacion,_that.descripcion,_that.estado,_that.createdAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String nombre,  String ubicacion,  String descripcion,  String estado,  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _Finca():
return $default(_that.id,_that.nombre,_that.ubicacion,_that.descripcion,_that.estado,_that.createdAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String nombre,  String ubicacion,  String descripcion,  String estado,  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _Finca() when $default != null:
return $default(_that.id,_that.nombre,_that.ubicacion,_that.descripcion,_that.estado,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc


class _Finca extends Finca {
  const _Finca({required this.id, required this.nombre, required this.ubicacion, required this.descripcion, required this.estado, required this.createdAt}): super._();
  

@override final  String id;
@override final  String nombre;
@override final  String ubicacion;
@override final  String descripcion;
@override final  String estado;
@override final  DateTime createdAt;

/// Create a copy of Finca
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$FincaCopyWith<_Finca> get copyWith => __$FincaCopyWithImpl<_Finca>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Finca&&(identical(other.id, id) || other.id == id)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.ubicacion, ubicacion) || other.ubicacion == ubicacion)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,nombre,ubicacion,descripcion,estado,createdAt);

@override
String toString() {
  return 'Finca(id: $id, nombre: $nombre, ubicacion: $ubicacion, descripcion: $descripcion, estado: $estado, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$FincaCopyWith<$Res> implements $FincaCopyWith<$Res> {
  factory _$FincaCopyWith(_Finca value, $Res Function(_Finca) _then) = __$FincaCopyWithImpl;
@override @useResult
$Res call({
 String id, String nombre, String ubicacion, String descripcion, String estado, DateTime createdAt
});




}
/// @nodoc
class __$FincaCopyWithImpl<$Res>
    implements _$FincaCopyWith<$Res> {
  __$FincaCopyWithImpl(this._self, this._then);

  final _Finca _self;
  final $Res Function(_Finca) _then;

/// Create a copy of Finca
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? nombre = null,Object? ubicacion = null,Object? descripcion = null,Object? estado = null,Object? createdAt = null,}) {
  return _then(_Finca(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,ubicacion: null == ubicacion ? _self.ubicacion : ubicacion // ignore: cast_nullable_to_non_nullable
as String,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}

// dart format on
