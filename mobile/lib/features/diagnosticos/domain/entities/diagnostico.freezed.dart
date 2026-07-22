// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'diagnostico.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$Diagnostico {

 String get id; String get muestraId; String get estado; bool get tieneClorosis; double get confianza; String? get imageUrl; String? get imageBase64;
/// Create a copy of Diagnostico
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$DiagnosticoCopyWith<Diagnostico> get copyWith => _$DiagnosticoCopyWithImpl<Diagnostico>(this as Diagnostico, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Diagnostico&&(identical(other.id, id) || other.id == id)&&(identical(other.muestraId, muestraId) || other.muestraId == muestraId)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.tieneClorosis, tieneClorosis) || other.tieneClorosis == tieneClorosis)&&(identical(other.confianza, confianza) || other.confianza == confianza)&&(identical(other.imageUrl, imageUrl) || other.imageUrl == imageUrl)&&(identical(other.imageBase64, imageBase64) || other.imageBase64 == imageBase64));
}


@override
int get hashCode => Object.hash(runtimeType,id,muestraId,estado,tieneClorosis,confianza,imageUrl,imageBase64);

@override
String toString() {
  return 'Diagnostico(id: $id, muestraId: $muestraId, estado: $estado, tieneClorosis: $tieneClorosis, confianza: $confianza, imageUrl: $imageUrl, imageBase64: $imageBase64)';
}


}

/// @nodoc
abstract mixin class $DiagnosticoCopyWith<$Res>  {
  factory $DiagnosticoCopyWith(Diagnostico value, $Res Function(Diagnostico) _then) = _$DiagnosticoCopyWithImpl;
@useResult
$Res call({
 String id, String muestraId, String estado, bool tieneClorosis, double confianza, String? imageUrl, String? imageBase64
});




}
/// @nodoc
class _$DiagnosticoCopyWithImpl<$Res>
    implements $DiagnosticoCopyWith<$Res> {
  _$DiagnosticoCopyWithImpl(this._self, this._then);

  final Diagnostico _self;
  final $Res Function(Diagnostico) _then;

/// Create a copy of Diagnostico
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


/// Adds pattern-matching-related methods to [Diagnostico].
extension DiagnosticoPatterns on Diagnostico {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Diagnostico value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Diagnostico() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Diagnostico value)  $default,){
final _that = this;
switch (_that) {
case _Diagnostico():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Diagnostico value)?  $default,){
final _that = this;
switch (_that) {
case _Diagnostico() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String muestraId,  String estado,  bool tieneClorosis,  double confianza,  String? imageUrl,  String? imageBase64)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Diagnostico() when $default != null:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String muestraId,  String estado,  bool tieneClorosis,  double confianza,  String? imageUrl,  String? imageBase64)  $default,) {final _that = this;
switch (_that) {
case _Diagnostico():
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String muestraId,  String estado,  bool tieneClorosis,  double confianza,  String? imageUrl,  String? imageBase64)?  $default,) {final _that = this;
switch (_that) {
case _Diagnostico() when $default != null:
return $default(_that.id,_that.muestraId,_that.estado,_that.tieneClorosis,_that.confianza,_that.imageUrl,_that.imageBase64);case _:
  return null;

}
}

}

/// @nodoc


class _Diagnostico implements Diagnostico {
  const _Diagnostico({required this.id, required this.muestraId, required this.estado, required this.tieneClorosis, required this.confianza, this.imageUrl, this.imageBase64});
  

@override final  String id;
@override final  String muestraId;
@override final  String estado;
@override final  bool tieneClorosis;
@override final  double confianza;
@override final  String? imageUrl;
@override final  String? imageBase64;

/// Create a copy of Diagnostico
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$DiagnosticoCopyWith<_Diagnostico> get copyWith => __$DiagnosticoCopyWithImpl<_Diagnostico>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Diagnostico&&(identical(other.id, id) || other.id == id)&&(identical(other.muestraId, muestraId) || other.muestraId == muestraId)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.tieneClorosis, tieneClorosis) || other.tieneClorosis == tieneClorosis)&&(identical(other.confianza, confianza) || other.confianza == confianza)&&(identical(other.imageUrl, imageUrl) || other.imageUrl == imageUrl)&&(identical(other.imageBase64, imageBase64) || other.imageBase64 == imageBase64));
}


@override
int get hashCode => Object.hash(runtimeType,id,muestraId,estado,tieneClorosis,confianza,imageUrl,imageBase64);

@override
String toString() {
  return 'Diagnostico(id: $id, muestraId: $muestraId, estado: $estado, tieneClorosis: $tieneClorosis, confianza: $confianza, imageUrl: $imageUrl, imageBase64: $imageBase64)';
}


}

/// @nodoc
abstract mixin class _$DiagnosticoCopyWith<$Res> implements $DiagnosticoCopyWith<$Res> {
  factory _$DiagnosticoCopyWith(_Diagnostico value, $Res Function(_Diagnostico) _then) = __$DiagnosticoCopyWithImpl;
@override @useResult
$Res call({
 String id, String muestraId, String estado, bool tieneClorosis, double confianza, String? imageUrl, String? imageBase64
});




}
/// @nodoc
class __$DiagnosticoCopyWithImpl<$Res>
    implements _$DiagnosticoCopyWith<$Res> {
  __$DiagnosticoCopyWithImpl(this._self, this._then);

  final _Diagnostico _self;
  final $Res Function(_Diagnostico) _then;

/// Create a copy of Diagnostico
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? muestraId = null,Object? estado = null,Object? tieneClorosis = null,Object? confianza = null,Object? imageUrl = freezed,Object? imageBase64 = freezed,}) {
  return _then(_Diagnostico(
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
