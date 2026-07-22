// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'diagnostico_dto.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$DiagnosticoDto {

 String get id;@JsonKey(name: 'muestraID') String get muestraId; String get estado; bool get tieneClorosis; double get confianza;@JsonKey(name: 'imageURL') String? get imageUrl; String? get imageBase64;
/// Create a copy of DiagnosticoDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$DiagnosticoDtoCopyWith<DiagnosticoDto> get copyWith => _$DiagnosticoDtoCopyWithImpl<DiagnosticoDto>(this as DiagnosticoDto, _$identity);

  /// Serializes this DiagnosticoDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is DiagnosticoDto&&(identical(other.id, id) || other.id == id)&&(identical(other.muestraId, muestraId) || other.muestraId == muestraId)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.tieneClorosis, tieneClorosis) || other.tieneClorosis == tieneClorosis)&&(identical(other.confianza, confianza) || other.confianza == confianza)&&(identical(other.imageUrl, imageUrl) || other.imageUrl == imageUrl)&&(identical(other.imageBase64, imageBase64) || other.imageBase64 == imageBase64));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,muestraId,estado,tieneClorosis,confianza,imageUrl,imageBase64);

@override
String toString() {
  return 'DiagnosticoDto(id: $id, muestraId: $muestraId, estado: $estado, tieneClorosis: $tieneClorosis, confianza: $confianza, imageUrl: $imageUrl, imageBase64: $imageBase64)';
}


}

/// @nodoc
abstract mixin class $DiagnosticoDtoCopyWith<$Res>  {
  factory $DiagnosticoDtoCopyWith(DiagnosticoDto value, $Res Function(DiagnosticoDto) _then) = _$DiagnosticoDtoCopyWithImpl;
@useResult
$Res call({
 String id,@JsonKey(name: 'muestraID') String muestraId, String estado, bool tieneClorosis, double confianza,@JsonKey(name: 'imageURL') String? imageUrl, String? imageBase64
});




}
/// @nodoc
class _$DiagnosticoDtoCopyWithImpl<$Res>
    implements $DiagnosticoDtoCopyWith<$Res> {
  _$DiagnosticoDtoCopyWithImpl(this._self, this._then);

  final DiagnosticoDto _self;
  final $Res Function(DiagnosticoDto) _then;

/// Create a copy of DiagnosticoDto
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? muestraId = null,Object? estado = null,Object? tieneClorosis = null,Object? confianza = null,Object? imageUrl = freezed,Object? imageBase64 = freezed,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,muestraId: null == muestraId ? _self.muestraId : muestraId // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,tieneClorosis: null == tieneClorosis ? _self.tieneClorosis : tieneClorosis // ignore: cast_nullable_to_non_nullable
as bool,confianza: null == confianza ? _self.confianza : confianza // ignore: cast_nullable_to_non_nullable
as double,imageUrl: freezed == imageUrl ? _self.imageUrl : imageUrl // ignore: cast_nullable_to_non_nullable
as String?,imageBase64: freezed == imageBase64 ? _self.imageBase64 : imageBase64 // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

}


/// Adds pattern-matching-related methods to [DiagnosticoDto].
extension DiagnosticoDtoPatterns on DiagnosticoDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _DiagnosticoDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _DiagnosticoDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _DiagnosticoDto value)  $default,){
final _that = this;
switch (_that) {
case _DiagnosticoDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _DiagnosticoDto value)?  $default,){
final _that = this;
switch (_that) {
case _DiagnosticoDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'muestraID')  String muestraId,  String estado,  bool tieneClorosis,  double confianza, @JsonKey(name: 'imageURL')  String? imageUrl,  String? imageBase64)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _DiagnosticoDto() when $default != null:
return $default(_that.id,_that.muestraId,_that.estado,_that.tieneClorosis,_that.confianza,_that.imageUrl,_that.imageBase64);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id, @JsonKey(name: 'muestraID')  String muestraId,  String estado,  bool tieneClorosis,  double confianza, @JsonKey(name: 'imageURL')  String? imageUrl,  String? imageBase64)  $default,) {final _that = this;
switch (_that) {
case _DiagnosticoDto():
return $default(_that.id,_that.muestraId,_that.estado,_that.tieneClorosis,_that.confianza,_that.imageUrl,_that.imageBase64);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id, @JsonKey(name: 'muestraID')  String muestraId,  String estado,  bool tieneClorosis,  double confianza, @JsonKey(name: 'imageURL')  String? imageUrl,  String? imageBase64)?  $default,) {final _that = this;
switch (_that) {
case _DiagnosticoDto() when $default != null:
return $default(_that.id,_that.muestraId,_that.estado,_that.tieneClorosis,_that.confianza,_that.imageUrl,_that.imageBase64);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _DiagnosticoDto extends DiagnosticoDto {
  const _DiagnosticoDto({required this.id, @JsonKey(name: 'muestraID') required this.muestraId, required this.estado, required this.tieneClorosis, required this.confianza, @JsonKey(name: 'imageURL') this.imageUrl, this.imageBase64}): super._();
  factory _DiagnosticoDto.fromJson(Map<String, dynamic> json) => _$DiagnosticoDtoFromJson(json);

@override final  String id;
@override@JsonKey(name: 'muestraID') final  String muestraId;
@override final  String estado;
@override final  bool tieneClorosis;
@override final  double confianza;
@override@JsonKey(name: 'imageURL') final  String? imageUrl;
@override final  String? imageBase64;

/// Create a copy of DiagnosticoDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$DiagnosticoDtoCopyWith<_DiagnosticoDto> get copyWith => __$DiagnosticoDtoCopyWithImpl<_DiagnosticoDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$DiagnosticoDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _DiagnosticoDto&&(identical(other.id, id) || other.id == id)&&(identical(other.muestraId, muestraId) || other.muestraId == muestraId)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.tieneClorosis, tieneClorosis) || other.tieneClorosis == tieneClorosis)&&(identical(other.confianza, confianza) || other.confianza == confianza)&&(identical(other.imageUrl, imageUrl) || other.imageUrl == imageUrl)&&(identical(other.imageBase64, imageBase64) || other.imageBase64 == imageBase64));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,muestraId,estado,tieneClorosis,confianza,imageUrl,imageBase64);

@override
String toString() {
  return 'DiagnosticoDto(id: $id, muestraId: $muestraId, estado: $estado, tieneClorosis: $tieneClorosis, confianza: $confianza, imageUrl: $imageUrl, imageBase64: $imageBase64)';
}


}

/// @nodoc
abstract mixin class _$DiagnosticoDtoCopyWith<$Res> implements $DiagnosticoDtoCopyWith<$Res> {
  factory _$DiagnosticoDtoCopyWith(_DiagnosticoDto value, $Res Function(_DiagnosticoDto) _then) = __$DiagnosticoDtoCopyWithImpl;
@override @useResult
$Res call({
 String id,@JsonKey(name: 'muestraID') String muestraId, String estado, bool tieneClorosis, double confianza,@JsonKey(name: 'imageURL') String? imageUrl, String? imageBase64
});




}
/// @nodoc
class __$DiagnosticoDtoCopyWithImpl<$Res>
    implements _$DiagnosticoDtoCopyWith<$Res> {
  __$DiagnosticoDtoCopyWithImpl(this._self, this._then);

  final _DiagnosticoDto _self;
  final $Res Function(_DiagnosticoDto) _then;

/// Create a copy of DiagnosticoDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? muestraId = null,Object? estado = null,Object? tieneClorosis = null,Object? confianza = null,Object? imageUrl = freezed,Object? imageBase64 = freezed,}) {
  return _then(_DiagnosticoDto(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,muestraId: null == muestraId ? _self.muestraId : muestraId // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,tieneClorosis: null == tieneClorosis ? _self.tieneClorosis : tieneClorosis // ignore: cast_nullable_to_non_nullable
as bool,confianza: null == confianza ? _self.confianza : confianza // ignore: cast_nullable_to_non_nullable
as double,imageUrl: freezed == imageUrl ? _self.imageUrl : imageUrl // ignore: cast_nullable_to_non_nullable
as String?,imageBase64: freezed == imageBase64 ? _self.imageBase64 : imageBase64 // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}


}

// dart format on
