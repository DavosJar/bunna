// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'lote_dto.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$LoteDto {

 String get id;@JsonKey(name: 'fincaID') String get fincaId; String get nombre; double get area; String get descripcion; String get estado; DateTime get createdAt;
/// Create a copy of LoteDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$LoteDtoCopyWith<LoteDto> get copyWith => _$LoteDtoCopyWithImpl<LoteDto>(this as LoteDto, _$identity);

  /// Serializes this LoteDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is LoteDto&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.area, area) || other.area == area)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,fincaId,nombre,area,descripcion,estado,createdAt);

@override
String toString() {
  return 'LoteDto(id: $id, fincaId: $fincaId, nombre: $nombre, area: $area, descripcion: $descripcion, estado: $estado, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class $LoteDtoCopyWith<$Res>  {
  factory $LoteDtoCopyWith(LoteDto value, $Res Function(LoteDto) _then) = _$LoteDtoCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'fincaID') String fincaId, String nombre, double area, String descripcion, String estado, DateTime createdAt
});




}
/// @nodoc
class _$LoteDtoCopyWithImpl<$Res>
    implements $LoteDtoCopyWith<$Res> {
  _$LoteDtoCopyWithImpl(this._self, this._then);

  final LoteDto _self;
  final $Res Function(LoteDto) _then;

/// Create a copy of LoteDto
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


/// Adds pattern-matching-related methods to [LoteDto].
extension LoteDtoPatterns on LoteDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _LoteDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _LoteDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _LoteDto value)  $default,){
final _that = this;
switch (_that) {
case _LoteDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _LoteDto value)?  $default,){
final _that = this;
switch (_that) {
case _LoteDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'fincaID')  String fincaId,  String nombre,  double area,  String descripcion,  String estado,  DateTime createdAt)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _LoteDto() when $default != null:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'fincaID')  String fincaId,  String nombre,  double area,  String descripcion,  String estado,  DateTime createdAt)  $default,) {final _that = this;
switch (_that) {
case _LoteDto():
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'fincaID')  String fincaId,  String nombre,  double area,  String descripcion,  String estado,  DateTime createdAt)?  $default,) {final _that = this;
switch (_that) {
case _LoteDto() when $default != null:
return $default(_that.id,_that.fincaId,_that.nombre,_that.area,_that.descripcion,_that.estado,_that.createdAt);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _LoteDto extends LoteDto {
  const _LoteDto({required this.id, @JsonKey(name: 'fincaID') required this.fincaId, required this.nombre, required this.area, required this.descripcion, required this.estado, required this.createdAt}): super._();
  factory _LoteDto.fromJson(Map<String, dynamic> json) => _$LoteDtoFromJson(json);

@override final  String id;
@override@JsonKey(name: 'fincaID') final  String fincaId;
@override final  String nombre;
@override final  double area;
@override final  String descripcion;
@override final  String estado;
@override final  DateTime createdAt;

/// Create a copy of LoteDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$LoteDtoCopyWith<_LoteDto> get copyWith => __$LoteDtoCopyWithImpl<_LoteDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$LoteDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _LoteDto&&(identical(other.id, id) || other.id == id)&&(identical(other.fincaId, fincaId) || other.fincaId == fincaId)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.area, area) || other.area == area)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.createdAt, createdAt) || other.createdAt == createdAt));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,fincaId,nombre,area,descripcion,estado,createdAt);

@override
String toString() {
  return 'LoteDto(id: $id, fincaId: $fincaId, nombre: $nombre, area: $area, descripcion: $descripcion, estado: $estado, createdAt: $createdAt)';
}


}

/// @nodoc
abstract mixin class _$LoteDtoCopyWith<$Res> implements $LoteDtoCopyWith<$Res> {
  factory _$LoteDtoCopyWith(_LoteDto value, $Res Function(_LoteDto) _then) = __$LoteDtoCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'fincaID') String fincaId, String nombre, double area, String descripcion, String estado, DateTime createdAt
});




}
/// @nodoc
class __$LoteDtoCopyWithImpl<$Res>
    implements _$LoteDtoCopyWith<$Res> {
  __$LoteDtoCopyWithImpl(this._self, this._then);

  final _LoteDto _self;
  final $Res Function(_LoteDto) _then;

/// Create a copy of LoteDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? fincaId = null,Object? nombre = null,Object? area = null,Object? descripcion = null,Object? estado = null,Object? createdAt = null,}) {
  return _then(_LoteDto(
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
