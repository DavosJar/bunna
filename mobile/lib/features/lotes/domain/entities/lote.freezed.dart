// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'lote.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$Lote {

 String get id; String get fincaId; String get nombre; double get area; String get descripcion; String get estado; DateTime get createdAt;
/// Create a copy of Lote
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$LoteCopyWith<Lote> get copyWith => _$LoteCopyWithImpl<Lote>(this as Lote, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Lote&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.area, area) || other.area == area)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,fincaId,nombre,area,descripcion,estado,createdAt);

@override
String toString() {
  return 'Lote(id: $id, fincaId: $fincaId, nombre: $nombre, area: $area, descripcion: $descripcion, estado: $estado, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $LoteCopyWith<$Res>  {
  factory $LoteCopyWith(Lote value, $Res Function(Lote) _then) = _$LoteCopyWithImpl;
@useResult
$Res call({
 String id, String fincaId, String nombre, double area, String descripcion, String estado, DateTime createdAt
});




}
/// @nodoc
class _$LoteCopyWithImpl<$Res>
    implements $LoteCopyWith<$Res> {
  _$LoteCopyWithImpl(this._self, this._then);

  final Lote _self;
  final $Res Function(Lote) _then;

/// Create a copy of Lote
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? fincaId = null,Object? nombre = null,Object? area = null,Object? descripcion = null,Object? estado = null,Object? createdAt = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,fincaId: null == fincaId ? _self.fincaId : fincaId // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,area: null == area ? _self.area : area // ignore: cast_nullable_to_non_nullable
as double,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [Lote].
extension LotePatterns on Lote {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Lote value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Lote() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Lote value)  $default,){
final _that = this;
switch (_that) {
case _Lote():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Lote value)?  $default,){
final _that = this;
switch (_that) {
case _Lote() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String fincaId,  String nombre,  double area,  String descripcion,  String estado,  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Lote() when $default != null:
return $default(_that.id,_that.fincaId,_that.nombre,_that.area,_that.descripcion,_that.estado,_that.createdAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String fincaId,  String nombre,  double area,  String descripcion,  String estado,  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _Lote():
return $default(_that.id,_that.fincaId,_that.nombre,_that.area,_that.descripcion,_that.estado,_that.createdAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String fincaId,  String nombre,  double area,  String descripcion,  String estado,  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _Lote() when $default != null:
return $default(_that.id,_that.fincaId,_that.nombre,_that.area,_that.descripcion,_that.estado,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc


class _Lote extends Lote {
  const _Lote({required this.id, required this.fincaId, required this.nombre, required this.area, required this.descripcion, required this.estado, required this.createdAt}): super._();
  

@override final  String id;
@override final  String fincaId;
@override final  String nombre;
@override final  double area;
@override final  String descripcion;
@override final  String estado;
@override final  DateTime createdAt;

/// Create a copy of Lote
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$LoteCopyWith<_Lote> get copyWith => __$LoteCopyWithImpl<_Lote>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Lote&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.area, area) || other.area == area)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,fincaId,nombre,area,descripcion,estado,createdAt);

@override
String toString() {
  return 'Lote(id: $id, fincaId: $fincaId, nombre: $nombre, area: $area, descripcion: $descripcion, estado: $estado, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$LoteCopyWith<$Res> implements $LoteCopyWith<$Res> {
  factory _$LoteCopyWith(_Lote value, $Res Function(_Lote) _then) = __$LoteCopyWithImpl;
@override @useResult
$Res call({
 String id, String fincaId, String nombre, double area, String descripcion, String estado, DateTime createdAt
});




}
/// @nodoc
class __$LoteCopyWithImpl<$Res>
    implements _$LoteCopyWith<$Res> {
  __$LoteCopyWithImpl(this._self, this._then);

  final _Lote _self;
  final $Res Function(_Lote) _then;

/// Create a copy of Lote
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? fincaId = null,Object? nombre = null,Object? area = null,Object? descripcion = null,Object? estado = null,Object? createdAt = null,}) {
  return _then(_Lote(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,fincaId: null == fincaId ? _self.fincaId : fincaId // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,area: null == area ? _self.area : area // ignore: cast_nullable_to_non_nullable
as double,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}

// dart format on
