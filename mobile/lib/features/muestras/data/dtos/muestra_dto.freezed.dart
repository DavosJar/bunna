// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'muestra_dto.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$MuestraDto {

 String get id;@JsonKey(name: 'fincaID') String get fincaId;@JsonKey(name: 'loteID') String get loteId; double get latitud; double get longitud; DateTime get createdAt;
/// Create a copy of MuestraDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MuestraDtoCopyWith<MuestraDto> get copyWith => _$MuestraDtoCopyWithImpl<MuestraDto>(this as MuestraDto, _$identity);

  /// Serializes this MuestraDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MuestraDto&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.loteId, loteId) || other.loteId == loteId)&&(identical(other.latitud, latitud) || other.latitud == latitud)&&(identical(other.longitud, longitud) || other.longitud == longitud)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,fincaId,loteId,latitud,longitud,createdAt);

@override
String toString() {
  return 'MuestraDto(id: $id, fincaId: $fincaId, loteId: $loteId, latitud: $latitud, longitud: $longitud, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $MuestraDtoCopyWith<$Res>  {
  factory $MuestraDtoCopyWith(MuestraDto value, $Res Function(MuestraDto) _then) = _$MuestraDtoCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'fincaID') String fincaId,@JsonKey(name: 'loteID') String loteId, double latitud, double longitud, DateTime createdAt
});




}
/// @nodoc
class _$MuestraDtoCopyWithImpl<$Res>
    implements $MuestraDtoCopyWith<$Res> {
  _$MuestraDtoCopyWithImpl(this._self, this._then);

  final MuestraDto _self;
  final $Res Function(MuestraDto) _then;

/// Create a copy of MuestraDto
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


/// Adds pattern-matching-related methods to [MuestraDto].
extension MuestraDtoPatterns on MuestraDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MuestraDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MuestraDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MuestraDto value)  $default,){
final _that = this;
switch (_that) {
case _MuestraDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MuestraDto value)?  $default,){
final _that = this;
switch (_that) {
case _MuestraDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'fincaID')  String fincaId, @JsonKey(name: 'loteID')  String loteId,  double latitud,  double longitud,  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MuestraDto() when $default != null:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'fincaID')  String fincaId, @JsonKey(name: 'loteID')  String loteId,  double latitud,  double longitud,  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _MuestraDto():
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'fincaID')  String fincaId, @JsonKey(name: 'loteID')  String loteId,  double latitud,  double longitud,  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _MuestraDto() when $default != null:
return $default(_that.id,_that.fincaId,_that.loteId,_that.latitud,_that.longitud,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _MuestraDto extends MuestraDto {
  const _MuestraDto({required this.id, @JsonKey(name: 'fincaID') required this.fincaId, @JsonKey(name: 'loteID') required this.loteId, required this.latitud, required this.longitud, required this.createdAt}): super._();
  factory _MuestraDto.fromJson(Map<String, dynamic> json) => _$MuestraDtoFromJson(json);

@override final  String id;
@override@JsonKey(name: 'fincaID') final  String fincaId;
@override@JsonKey(name: 'loteID') final  String loteId;
@override final  double latitud;
@override final  double longitud;
@override final  DateTime createdAt;

/// Create a copy of MuestraDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MuestraDtoCopyWith<_MuestraDto> get copyWith => __$MuestraDtoCopyWithImpl<_MuestraDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MuestraDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _MuestraDto&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.loteId, loteId) || other.loteId == loteId)&&(identical(other.latitud, latitud) || other.latitud == latitud)&&(identical(other.longitud, longitud) || other.longitud == longitud)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,fincaId,loteId,latitud,longitud,createdAt);

@override
String toString() {
  return 'MuestraDto(id: $id, fincaId: $fincaId, loteId: $loteId, latitud: $latitud, longitud: $longitud, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$MuestraDtoCopyWith<$Res> implements $MuestraDtoCopyWith<$Res> {
  factory _$MuestraDtoCopyWith(_MuestraDto value, $Res Function(_MuestraDto) _then) = __$MuestraDtoCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'fincaID') String fincaId,@JsonKey(name: 'loteID') String loteId, double latitud, double longitud, DateTime createdAt
});




}
/// @nodoc
class __$MuestraDtoCopyWithImpl<$Res>
    implements _$MuestraDtoCopyWith<$Res> {
  __$MuestraDtoCopyWithImpl(this._self, this._then);

  final _MuestraDto _self;
  final $Res Function(_MuestraDto) _then;

/// Create a copy of MuestraDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? fincaId = null,Object? loteId = null,Object? latitud = null,Object? longitud = null,Object? createdAt = null,}) {
  return _then(_MuestraDto(
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
