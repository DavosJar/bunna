// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'cambio_estado.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$CambioEstado {

 String get id; String get estado; String? get motivo; DateTime get updatedAt;
/// Create a copy of CambioEstado
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$CambioEstadoCopyWith<CambioEstado> get copyWith => _$CambioEstadoCopyWithImpl<CambioEstado>(this as CambioEstado, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is CambioEstado&&(identical(other.id, id) || other.id == id)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.motivo, motivo) || other.motivo == motivo)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,estado,motivo,updatedAt);

@override
String toString() {
  return 'CambioEstado(id: $id, estado: $estado, motivo: $motivo, updatedAt: $updatedAt)';
}


}

/// @nodoc
abstract mixin class $CambioEstadoCopyWith<$Res>  {
  factory $CambioEstadoCopyWith(CambioEstado value, $Res Function(CambioEstado) _then) = _$CambioEstadoCopyWithImpl;
@useResult
$Res call({
 String id, String estado, String? motivo, DateTime updatedAt
});




}
/// @nodoc
class _$CambioEstadoCopyWithImpl<$Res>
    implements $CambioEstadoCopyWith<$Res> {
  _$CambioEstadoCopyWithImpl(this._self, this._then);

  final CambioEstado _self;
  final $Res Function(CambioEstado) _then;

/// Create a copy of CambioEstado
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? estado = null,Object? motivo = freezed,Object? updatedAt = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,motivo: freezed == motivo ? _self.motivo : motivo // ignore: cast_nullable_to_non_nullable
as String?,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [CambioEstado].
extension CambioEstadoPatterns on CambioEstado {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _CambioEstado value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _CambioEstado() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _CambioEstado value)  $default,){
final _that = this;
switch (_that) {
case _CambioEstado():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _CambioEstado value)?  $default,){
final _that = this;
switch (_that) {
case _CambioEstado() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String estado,  String? motivo,  DateTime updatedAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _CambioEstado() when $default != null:
return $default(_that.id,_that.estado,_that.motivo,_that.updatedAt);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String estado,  String? motivo,  DateTime updatedAt)  $default,) {final _that = this;
switch (_that) {
case _CambioEstado():
return $default(_that.id,_that.estado,_that.motivo,_that.updatedAt);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String estado,  String? motivo,  DateTime updatedAt)?  $default,) {final _that = this;
switch (_that) {
case _CambioEstado() when $default != null:
return $default(_that.id,_that.estado,_that.motivo,_that.updatedAt);case _:
  return null;

}
}

}

/// @nodoc


class _CambioEstado implements CambioEstado {
  const _CambioEstado({required this.id, required this.estado, this.motivo, required this.updatedAt});
  

@override final  String id;
@override final  String estado;
@override final  String? motivo;
@override final  DateTime updatedAt;

/// Create a copy of CambioEstado
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$CambioEstadoCopyWith<_CambioEstado> get copyWith => __$CambioEstadoCopyWithImpl<_CambioEstado>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _CambioEstado&&(identical(other.id, id) || other.id == id)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.motivo, motivo) || other.motivo == motivo)&&(identical(other.updatedAt, updatedAt) || other.updatedAt == updatedAt));
}


@override
int get hashCode => Object.hash(runtimeType,id,estado,motivo,updatedAt);

@override
String toString() {
  return 'CambioEstado(id: $id, estado: $estado, motivo: $motivo, updatedAt: $updatedAt)';
}


}

/// @nodoc
abstract mixin class _$CambioEstadoCopyWith<$Res> implements $CambioEstadoCopyWith<$Res> {
  factory _$CambioEstadoCopyWith(_CambioEstado value, $Res Function(_CambioEstado) _then) = __$CambioEstadoCopyWithImpl;
@override @useResult
$Res call({
 String id, String estado, String? motivo, DateTime updatedAt
});




}
/// @nodoc
class __$CambioEstadoCopyWithImpl<$Res>
    implements _$CambioEstadoCopyWith<$Res> {
  __$CambioEstadoCopyWithImpl(this._self, this._then);

  final _CambioEstado _self;
  final $Res Function(_CambioEstado) _then;

/// Create a copy of CambioEstado
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? estado = null,Object? motivo = freezed,Object? updatedAt = null,}) {
  return _then(_CambioEstado(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,motivo: freezed == motivo ? _self.motivo : motivo // ignore: cast_nullable_to_non_nullable
as String?,updatedAt: null == updatedAt ? _self.updatedAt : updatedAt // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}

// dart format on
