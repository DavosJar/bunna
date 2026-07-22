// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'muestra.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$Muestra {

 String get id; String get fincaId; String get loteId; double get latitud; double get longitud; DateTime get createdAt;
/// Create a copy of Muestra
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MuestraCopyWith<Muestra> get copyWith => _$MuestraCopyWithImpl<Muestra>(this as Muestra, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Muestra&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.loteId, loteId) || other.loteId == loteId)&&(identical(other.latitud, latitud) || other.latitud == latitud)&&(identical(other.longitud, longitud) || other.longitud == longitud)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,fincaId,loteId,latitud,longitud,createdAt);

@override
String toString() {
  return 'Muestra(id: $id, fincaId: $fincaId, loteId: $loteId, latitud: $latitud, longitud: $longitud, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $MuestraCopyWith<$Res>  {
  factory $MuestraCopyWith(Muestra value, $Res Function(Muestra) _then) = _$MuestraCopyWithImpl;
@useResult
$Res call({
 String id, String fincaId, String loteId, double latitud, double longitud, DateTime createdAt
});




}
/// @nodoc
class _$MuestraCopyWithImpl<$Res>
    implements $MuestraCopyWith<$Res> {
  _$MuestraCopyWithImpl(this._self, this._then);

  final Muestra _self;
  final $Res Function(Muestra) _then;

/// Create a copy of Muestra
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? fincaId = null,Object? loteId = null,Object? latitud = null,Object? longitud = null,Object? createdAt = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,fincaId: null == fincaId ? _self.fincaId : fincaId // ignore: cast_nullable_to_non_nullable
as String,loteId: null == loteId ? _self.loteId : loteId // ignore: cast_nullable_to_non_nullable
as String,latitud: null == latitud ? _self.latitud : latitud // ignore: cast_nullable_to_non_nullable
as double,longitud: null == longitud ? _self.longitud : longitud // ignore: cast_nullable_to_non_nullable
as double,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [Muestra].
extension MuestraPatterns on Muestra {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Muestra value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Muestra() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Muestra value)  $default,){
final _that = this;
switch (_that) {
case _Muestra():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Muestra value)?  $default,){
final _that = this;
switch (_that) {
case _Muestra() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String fincaId,  String loteId,  double latitud,  double longitud,  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Muestra() when $default != null:
return $default(_that.id,_that.fincaId,_that.loteId,_that.latitud,_that.longitud,_that.createdAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String fincaId,  String loteId,  double latitud,  double longitud,  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _Muestra():
return $default(_that.id,_that.fincaId,_that.loteId,_that.latitud,_that.longitud,_that.createdAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String fincaId,  String loteId,  double latitud,  double longitud,  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _Muestra() when $default != null:
return $default(_that.id,_that.fincaId,_that.loteId,_that.latitud,_that.longitud,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc


class _Muestra implements Muestra {
  const _Muestra({required this.id, required this.fincaId, required this.loteId, required this.latitud, required this.longitud, required this.createdAt});
  

@override final  String id;
@override final  String fincaId;
@override final  String loteId;
@override final  double latitud;
@override final  double longitud;
@override final  DateTime createdAt;

/// Create a copy of Muestra
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MuestraCopyWith<_Muestra> get copyWith => __$MuestraCopyWithImpl<_Muestra>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Muestra&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.loteId, loteId) || other.loteId == loteId)&&(identical(other.latitud, latitud) || other.latitud == latitud)&&(identical(other.longitud, longitud) || other.longitud == longitud)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,fincaId,loteId,latitud,longitud,createdAt);

@override
String toString() {
  return 'Muestra(id: $id, fincaId: $fincaId, loteId: $loteId, latitud: $latitud, longitud: $longitud, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$MuestraCopyWith<$Res> implements $MuestraCopyWith<$Res> {
  factory _$MuestraCopyWith(_Muestra value, $Res Function(_Muestra) _then) = __$MuestraCopyWithImpl;
@override @useResult
$Res call({
 String id, String fincaId, String loteId, double latitud, double longitud, DateTime createdAt
});




}
/// @nodoc
class __$MuestraCopyWithImpl<$Res>
    implements _$MuestraCopyWith<$Res> {
  __$MuestraCopyWithImpl(this._self, this._then);

  final _Muestra _self;
  final $Res Function(_Muestra) _then;

/// Create a copy of Muestra
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? fincaId = null,Object? loteId = null,Object? latitud = null,Object? longitud = null,Object? createdAt = null,}) {
  return _then(_Muestra(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,fincaId: null == fincaId ? _self.fincaId : fincaId // ignore: cast_nullable_to_non_nullable
as String,loteId: null == loteId ? _self.loteId : loteId // ignore: cast_nullable_to_non_nullable
as String,latitud: null == latitud ? _self.latitud : latitud // ignore: cast_nullable_to_non_nullable
as double,longitud: null == longitud ? _self.longitud : longitud // ignore: cast_nullable_to_non_nullable
as double,createdAt: null == createdAt ? _self.createdAt : createdAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}

// dart format on
