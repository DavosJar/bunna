// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'quick_scan_controller.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$QuickScanState {

 bool get analizando; Uint8List? get previewBytes; AnalisisYolo? get analisis; String? get error;
/// Create a copy of QuickScanState
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$QuickScanStateCopyWith<QuickScanState> get copyWith => _$QuickScanStateCopyWithImpl<QuickScanState>(this as QuickScanState, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is QuickScanState&&(identical(other.analizando, analizando) || other.analizando == analizando)&&const DeepCollectionEquality().equals(other.previewBytes, previewBytes)&&(identical(other.analisis, analisis) || other.analisis == analisis)&&(identical(other.error, error) || other.error == error));
}


@override
int get hashCode => Object.hash(runtimeType,analizando,const DeepCollectionEquality().hash(previewBytes),analisis,error);

@override
String toString() {
  return 'QuickScanState(analizando: $analizando, previewBytes: $previewBytes, analisis: $analisis, error: $error)';
}


}

/// @nodoc
abstract mixin class $QuickScanStateCopyWith<$Res>  {
  factory $QuickScanStateCopyWith(QuickScanState value, $Res Function(QuickScanState) _then) = _$QuickScanStateCopyWithImpl;
@useResult
$Res call({
 bool analizando, Uint8List? previewBytes, AnalisisYolo? analisis, String? error
});


$AnalisisYoloCopyWith<$Res>? get analisis;

}
/// @nodoc
class _$QuickScanStateCopyWithImpl<$Res>
    implements $QuickScanStateCopyWith<$Res> {
  _$QuickScanStateCopyWithImpl(this._self, this._then);

  final QuickScanState _self;
  final $Res Function(QuickScanState) _then;

/// Create a copy of QuickScanState
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? analizando = null,Object? previewBytes = freezed,Object? analisis = freezed,Object? error = freezed,}) {
  return _then(_self.copyWith(
analizando: null == analizando ? _self.analizando : analizando // ignore: cast_nullable_to_non_nullable
as bool,previewBytes: freezed == previewBytes ? _self.previewBytes : previewBytes // ignore: cast_nullable_to_non_nullable
as Uint8List?,analisis: freezed == analisis ? _self.analisis : analisis // ignore: cast_nullable_to_non_nullable
as AnalisisYolo?,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}
/// Create a copy of QuickScanState
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$AnalisisYoloCopyWith<$Res>? get analisis {
    if (_self.analisis == null) {
    return null;
  }

  return $AnalisisYoloCopyWith<$Res>(_self.analisis!, (value) {
    return _then(_self.copyWith(analisis: value));
  });
}
}


/// Adds pattern-matching-related methods to [QuickScanState].
extension QuickScanStatePatterns on QuickScanState {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _QuickScanState value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _QuickScanState() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _QuickScanState value)  $default,){
final _that = this;
switch (_that) {
case _QuickScanState():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _QuickScanState value)?  $default,){
final _that = this;
switch (_that) {
case _QuickScanState() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( bool analizando,  Uint8List? previewBytes,  AnalisisYolo? analisis,  String? error)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _QuickScanState() when $default != null:
return $default(_that.analizando,_that.previewBytes,_that.analisis,_that.error);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( bool analizando,  Uint8List? previewBytes,  AnalisisYolo? analisis,  String? error)  $default,) {final _that = this;
switch (_that) {
case _QuickScanState():
return $default(_that.analizando,_that.previewBytes,_that.analisis,_that.error);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( bool analizando,  Uint8List? previewBytes,  AnalisisYolo? analisis,  String? error)?  $default,) {final _that = this;
switch (_that) {
case _QuickScanState() when $default != null:
return $default(_that.analizando,_that.previewBytes,_that.analisis,_that.error);case _:
  return null;

}
}

}

/// @nodoc


class _QuickScanState implements QuickScanState {
  const _QuickScanState({this.analizando = false, this.previewBytes, this.analisis, this.error});
  

@override@JsonKey() final  bool analizando;
@override final  Uint8List? previewBytes;
@override final  AnalisisYolo? analisis;
@override final  String? error;

/// Create a copy of QuickScanState
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$QuickScanStateCopyWith<_QuickScanState> get copyWith => __$QuickScanStateCopyWithImpl<_QuickScanState>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _QuickScanState&&(identical(other.analizando, analizando) || other.analizando == analizando)&&const DeepCollectionEquality().equals(other.previewBytes, previewBytes)&&(identical(other.analisis, analisis) || other.analisis == analisis)&&(identical(other.error, error) || other.error == error));
}


@override
int get hashCode => Object.hash(runtimeType,analizando,const DeepCollectionEquality().hash(previewBytes),analisis,error);

@override
String toString() {
  return 'QuickScanState(analizando: $analizando, previewBytes: $previewBytes, analisis: $analisis, error: $error)';
}


}

/// @nodoc
abstract mixin class _$QuickScanStateCopyWith<$Res> implements $QuickScanStateCopyWith<$Res> {
  factory _$QuickScanStateCopyWith(_QuickScanState value, $Res Function(_QuickScanState) _then) = __$QuickScanStateCopyWithImpl;
@override @useResult
$Res call({
 bool analizando, Uint8List? previewBytes, AnalisisYolo? analisis, String? error
});


@override $AnalisisYoloCopyWith<$Res>? get analisis;

}
/// @nodoc
class __$QuickScanStateCopyWithImpl<$Res>
    implements _$QuickScanStateCopyWith<$Res> {
  __$QuickScanStateCopyWithImpl(this._self, this._then);

  final _QuickScanState _self;
  final $Res Function(_QuickScanState) _then;

/// Create a copy of QuickScanState
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? analizando = null,Object? previewBytes = freezed,Object? analisis = freezed,Object? error = freezed,}) {
  return _then(_QuickScanState(
analizando: null == analizando ? _self.analizando : analizando // ignore: cast_nullable_to_non_nullable
as bool,previewBytes: freezed == previewBytes ? _self.previewBytes : previewBytes // ignore: cast_nullable_to_non_nullable
as Uint8List?,analisis: freezed == analisis ? _self.analisis : analisis // ignore: cast_nullable_to_non_nullable
as AnalisisYolo?,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

/// Create a copy of QuickScanState
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$AnalisisYoloCopyWith<$Res>? get analisis {
    if (_self.analisis == null) {
    return null;
  }

  return $AnalisisYoloCopyWith<$Res>(_self.analisis!, (value) {
    return _then(_self.copyWith(analisis: value));
  });
}
}

// dart format on
