// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'analisis_yolo.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$YoloFeedback {

 String? get label; String? get level;// low | medium | high
 String? get recommendation;
/// Create a copy of YoloFeedback
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$YoloFeedbackCopyWith<YoloFeedback> get copyWith => _$YoloFeedbackCopyWithImpl<YoloFeedback>(this as YoloFeedback, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is YoloFeedback&&(identical(other.label, label) || other.label == label)&&(identical(other.level, level) || other.level == level)&&(identical(other.recommendation, recommendation) || other.recommendation == recommendation));
}


@override
int get hashCode => Object.hash(runtimeType,label,level,recommendation);

@override
String toString() {
  return 'YoloFeedback(label: $label, level: $level, recommendation: $recommendation)';
}


}

/// @nodoc
abstract mixin class $YoloFeedbackCopyWith<$Res>  {
  factory $YoloFeedbackCopyWith(YoloFeedback value, $Res Function(YoloFeedback) _then) = _$YoloFeedbackCopyWithImpl;
@useResult
$Res call({
 String? label, String? level, String? recommendation
});




}
/// @nodoc
class _$YoloFeedbackCopyWithImpl<$Res>
    implements $YoloFeedbackCopyWith<$Res> {
  _$YoloFeedbackCopyWithImpl(this._self, this._then);

  final YoloFeedback _self;
  final $Res Function(YoloFeedback) _then;

/// Create a copy of YoloFeedback
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? label = freezed,Object? level = freezed,Object? recommendation = freezed,}) {
  return _then(_self.copyWith(
label: freezed == label ? _self.label : label // ignore: cast_nullable_to_non_nullable
as String?,level: freezed == level ? _self.level : level // ignore: cast_nullable_to_non_nullable
as String?,recommendation: freezed == recommendation ? _self.recommendation : recommendation // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

}


/// Adds pattern-matching-related methods to [YoloFeedback].
extension YoloFeedbackPatterns on YoloFeedback {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _YoloFeedback value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _YoloFeedback() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _YoloFeedback value)  $default,){
final _that = this;
switch (_that) {
case _YoloFeedback():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _YoloFeedback value)?  $default,){
final _that = this;
switch (_that) {
case _YoloFeedback() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String? label,  String? level,  String? recommendation)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _YoloFeedback() when $default != null:
return $default(_that.label,_that.level,_that.recommendation);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String? label,  String? level,  String? recommendation)  $default,) {final _that = this;
switch (_that) {
case _YoloFeedback():
return $default(_that.label,_that.level,_that.recommendation);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String? label,  String? level,  String? recommendation)?  $default,) {final _that = this;
switch (_that) {
case _YoloFeedback() when $default != null:
return $default(_that.label,_that.level,_that.recommendation);case _:
  return null;

}
}

}

/// @nodoc


class _YoloFeedback implements YoloFeedback {
  const _YoloFeedback({this.label, this.level, this.recommendation});
  

@override final  String? label;
@override final  String? level;
// low | medium | high
@override final  String? recommendation;

/// Create a copy of YoloFeedback
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$YoloFeedbackCopyWith<_YoloFeedback> get copyWith => __$YoloFeedbackCopyWithImpl<_YoloFeedback>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _YoloFeedback&&(identical(other.label, label) || other.label == label)&&(identical(other.level, level) || other.level == level)&&(identical(other.recommendation, recommendation) || other.recommendation == recommendation));
}


@override
int get hashCode => Object.hash(runtimeType,label,level,recommendation);

@override
String toString() {
  return 'YoloFeedback(label: $label, level: $level, recommendation: $recommendation)';
}


}

/// @nodoc
abstract mixin class _$YoloFeedbackCopyWith<$Res> implements $YoloFeedbackCopyWith<$Res> {
  factory _$YoloFeedbackCopyWith(_YoloFeedback value, $Res Function(_YoloFeedback) _then) = __$YoloFeedbackCopyWithImpl;
@override @useResult
$Res call({
 String? label, String? level, String? recommendation
});




}
/// @nodoc
class __$YoloFeedbackCopyWithImpl<$Res>
    implements _$YoloFeedbackCopyWith<$Res> {
  __$YoloFeedbackCopyWithImpl(this._self, this._then);

  final _YoloFeedback _self;
  final $Res Function(_YoloFeedback) _then;

/// Create a copy of YoloFeedback
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? label = freezed,Object? level = freezed,Object? recommendation = freezed,}) {
  return _then(_YoloFeedback(
label: freezed == label ? _self.label : label // ignore: cast_nullable_to_non_nullable
as String?,level: freezed == level ? _self.level : level // ignore: cast_nullable_to_non_nullable
as String?,recommendation: freezed == recommendation ? _self.recommendation : recommendation // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}


}

/// @nodoc
mixin _$AnalisisYolo {

 YoloFeedback? get feedback; int get numDetections; double get avgConfidence; List<String> get clasesDetectadas; String? get imageBase64;
/// Create a copy of AnalisisYolo
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$AnalisisYoloCopyWith<AnalisisYolo> get copyWith => _$AnalisisYoloCopyWithImpl<AnalisisYolo>(this as AnalisisYolo, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is AnalisisYolo&&(identical(other.feedback, feedback) || other.feedback == feedback)&&(identical(other.numDetections, numDetections) || other.numDetections == numDetections)&&(identical(other.avgConfidence, avgConfidence) || other.avgConfidence == avgConfidence)&&const DeepCollectionEquality().equals(other.clasesDetectadas, clasesDetectadas)&&(identical(other.imageBase64, imageBase64) || other.imageBase64 == imageBase64));
}


@override
int get hashCode => Object.hash(runtimeType,feedback,numDetections,avgConfidence,const DeepCollectionEquality().hash(clasesDetectadas),imageBase64);

@override
String toString() {
  return 'AnalisisYolo(feedback: $feedback, numDetections: $numDetections, avgConfidence: $avgConfidence, clasesDetectadas: $clasesDetectadas, imageBase64: $imageBase64)';
}


}

/// @nodoc
abstract mixin class $AnalisisYoloCopyWith<$Res>  {
  factory $AnalisisYoloCopyWith(AnalisisYolo value, $Res Function(AnalisisYolo) _then) = _$AnalisisYoloCopyWithImpl;
@useResult
$Res call({
 YoloFeedback? feedback, int numDetections, double avgConfidence, List<String> clasesDetectadas, String? imageBase64
});


$YoloFeedbackCopyWith<$Res>? get feedback;

}
/// @nodoc
class _$AnalisisYoloCopyWithImpl<$Res>
    implements $AnalisisYoloCopyWith<$Res> {
  _$AnalisisYoloCopyWithImpl(this._self, this._then);

  final AnalisisYolo _self;
  final $Res Function(AnalisisYolo) _then;

/// Create a copy of AnalisisYolo
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? feedback = freezed,Object? numDetections = null,Object? avgConfidence = null,Object? clasesDetectadas = null,Object? imageBase64 = freezed,}) {
  return _then(_self.copyWith(
feedback: freezed == feedback ? _self.feedback : feedback // ignore: cast_nullable_to_non_nullable
as YoloFeedback?,numDetections: null == numDetections ? _self.numDetections : numDetections // ignore: cast_nullable_to_non_nullable
as int,avgConfidence: null == avgConfidence ? _self.avgConfidence : avgConfidence // ignore: cast_nullable_to_non_nullable
as double,clasesDetectadas: null == clasesDetectadas ? _self.clasesDetectadas : clasesDetectadas // ignore: cast_nullable_to_non_nullable
as List<String>,imageBase64: freezed == imageBase64 ? _self.imageBase64 : imageBase64 // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}
/// Create a copy of AnalisisYolo
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$YoloFeedbackCopyWith<$Res>? get feedback {
    if (_self.feedback == null) {
    return null;
  }

  return $YoloFeedbackCopyWith<$Res>(_self.feedback!, (value) {
    return _then(_self.copyWith(feedback: value));
  });
}
}


/// Adds pattern-matching-related methods to [AnalisisYolo].
extension AnalisisYoloPatterns on AnalisisYolo {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _AnalisisYolo value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _AnalisisYolo() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _AnalisisYolo value)  $default,){
final _that = this;
switch (_that) {
case _AnalisisYolo():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _AnalisisYolo value)?  $default,){
final _that = this;
switch (_that) {
case _AnalisisYolo() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( YoloFeedback? feedback,  int numDetections,  double avgConfidence,  List<String> clasesDetectadas,  String? imageBase64)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _AnalisisYolo() when $default != null:
return $default(_that.feedback,_that.numDetections,_that.avgConfidence,_that.clasesDetectadas,_that.imageBase64);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( YoloFeedback? feedback,  int numDetections,  double avgConfidence,  List<String> clasesDetectadas,  String? imageBase64)  $default,) {final _that = this;
switch (_that) {
case _AnalisisYolo():
return $default(_that.feedback,_that.numDetections,_that.avgConfidence,_that.clasesDetectadas,_that.imageBase64);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( YoloFeedback? feedback,  int numDetections,  double avgConfidence,  List<String> clasesDetectadas,  String? imageBase64)?  $default,) {final _that = this;
switch (_that) {
case _AnalisisYolo() when $default != null:
return $default(_that.feedback,_that.numDetections,_that.avgConfidence,_that.clasesDetectadas,_that.imageBase64);case _:
  return null;

}
}

}

/// @nodoc


class _AnalisisYolo extends AnalisisYolo {
  const _AnalisisYolo({this.feedback, required this.numDetections, required this.avgConfidence, required final  List<String> clasesDetectadas, this.imageBase64}): _clasesDetectadas = clasesDetectadas,super._();
  

@override final  YoloFeedback? feedback;
@override final  int numDetections;
@override final  double avgConfidence;
 final  List<String> _clasesDetectadas;
@override List<String> get clasesDetectadas {
  if (_clasesDetectadas is EqualUnmodifiableListView) return _clasesDetectadas;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_clasesDetectadas);
}

@override final  String? imageBase64;

/// Create a copy of AnalisisYolo
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$AnalisisYoloCopyWith<_AnalisisYolo> get copyWith => __$AnalisisYoloCopyWithImpl<_AnalisisYolo>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _AnalisisYolo&&(identical(other.feedback, feedback) || other.feedback == feedback)&&(identical(other.numDetections, numDetections) || other.numDetections == numDetections)&&(identical(other.avgConfidence, avgConfidence) || other.avgConfidence == avgConfidence)&&const DeepCollectionEquality().equals(other._clasesDetectadas, _clasesDetectadas)&&(identical(other.imageBase64, imageBase64) || other.imageBase64 == imageBase64));
}


@override
int get hashCode => Object.hash(runtimeType,feedback,numDetections,avgConfidence,const DeepCollectionEquality().hash(_clasesDetectadas),imageBase64);

@override
String toString() {
  return 'AnalisisYolo(feedback: $feedback, numDetections: $numDetections, avgConfidence: $avgConfidence, clasesDetectadas: $clasesDetectadas, imageBase64: $imageBase64)';
}


}

/// @nodoc
abstract mixin class _$AnalisisYoloCopyWith<$Res> implements $AnalisisYoloCopyWith<$Res> {
  factory _$AnalisisYoloCopyWith(_AnalisisYolo value, $Res Function(_AnalisisYolo) _then) = __$AnalisisYoloCopyWithImpl;
@override @useResult
$Res call({
 YoloFeedback? feedback, int numDetections, double avgConfidence, List<String> clasesDetectadas, String? imageBase64
});


@override $YoloFeedbackCopyWith<$Res>? get feedback;

}
/// @nodoc
class __$AnalisisYoloCopyWithImpl<$Res>
    implements _$AnalisisYoloCopyWith<$Res> {
  __$AnalisisYoloCopyWithImpl(this._self, this._then);

  final _AnalisisYolo _self;
  final $Res Function(_AnalisisYolo) _then;

/// Create a copy of AnalisisYolo
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? feedback = freezed,Object? numDetections = null,Object? avgConfidence = null,Object? clasesDetectadas = null,Object? imageBase64 = freezed,}) {
  return _then(_AnalisisYolo(
feedback: freezed == feedback ? _self.feedback : feedback // ignore: cast_nullable_to_non_nullable
as YoloFeedback?,numDetections: null == numDetections ? _self.numDetections : numDetections // ignore: cast_nullable_to_non_nullable
as int,avgConfidence: null == avgConfidence ? _self.avgConfidence : avgConfidence // ignore: cast_nullable_to_non_nullable
as double,clasesDetectadas: null == clasesDetectadas ? _self._clasesDetectadas : clasesDetectadas // ignore: cast_nullable_to_non_nullable
as List<String>,imageBase64: freezed == imageBase64 ? _self.imageBase64 : imageBase64 // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

/// Create a copy of AnalisisYolo
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$YoloFeedbackCopyWith<$Res>? get feedback {
    if (_self.feedback == null) {
    return null;
  }

  return $YoloFeedbackCopyWith<$Res>(_self.feedback!, (value) {
    return _then(_self.copyWith(feedback: value));
  });
}
}

// dart format on
