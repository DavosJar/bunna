// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'diagnostico_flow_state.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$DiagnosticoFlowState {

 FlowFase get fase; Uint8List? get previewBytes; AnalisisYolo? get analisis; Diagnostico? get diagnostico; String? get resueltoEstado;// ACEPTADO | RECHAZADO
 String? get error;
/// Create a copy of DiagnosticoFlowState
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$DiagnosticoFlowStateCopyWith<DiagnosticoFlowState> get copyWith => _$DiagnosticoFlowStateCopyWithImpl<DiagnosticoFlowState>(this as DiagnosticoFlowState, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is DiagnosticoFlowState&&(identical(other.fase, fase) || other.fase == fase)&&const DeepCollectionEquality().equals(other.previewBytes, previewBytes)&&(identical(other.analisis, analisis) || other.analisis == analisis)&&(identical(other.diagnostico, diagnostico) || other.diagnostico == diagnostico)&&(identical(other.resueltoEstado, resueltoEstado) || other.resueltoEstado == resueltoEstado)&&(identical(other.error, error) || other.error == error));
}


@override
int get hashCode => Object.hash(runtimeType,fase,const DeepCollectionEquality().hash(previewBytes),analisis,diagnostico,resueltoEstado,error);

@override
String toString() {
  return 'DiagnosticoFlowState(fase: $fase, previewBytes: $previewBytes, analisis: $analisis, diagnostico: $diagnostico, resueltoEstado: $resueltoEstado, error: $error)';
}


}

/// @nodoc
abstract mixin class $DiagnosticoFlowStateCopyWith<$Res>  {
  factory $DiagnosticoFlowStateCopyWith(DiagnosticoFlowState value, $Res Function(DiagnosticoFlowState) _then) = _$DiagnosticoFlowStateCopyWithImpl;
@useResult
$Res call({
 FlowFase fase, Uint8List? previewBytes, AnalisisYolo? analisis, Diagnostico? diagnostico, String? resueltoEstado, String? error
});


$AnalisisYoloCopyWith<$Res>? get analisis;$DiagnosticoCopyWith<$Res>? get diagnostico;

}
/// @nodoc
class _$DiagnosticoFlowStateCopyWithImpl<$Res>
    implements $DiagnosticoFlowStateCopyWith<$Res> {
  _$DiagnosticoFlowStateCopyWithImpl(this._self, this._then);

  final DiagnosticoFlowState _self;
  final $Res Function(DiagnosticoFlowState) _then;

/// Create a copy of DiagnosticoFlowState
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? fase = null,Object? previewBytes = freezed,Object? analisis = freezed,Object? diagnostico = freezed,Object? resueltoEstado = freezed,Object? error = freezed,}) {
  return _then(_self.copyWith(
fase: null == fase ? _self.fase : fase // ignore: cast_nullable_to_non_nullable
as FlowFase,previewBytes: freezed == previewBytes ? _self.previewBytes : previewBytes // ignore: cast_nullable_to_non_nullable
as Uint8List?,analisis: freezed == analisis ? _self.analisis : analisis // ignore: cast_nullable_to_non_nullable
as AnalisisYolo?,diagnostico: freezed == diagnostico ? _self.diagnostico : diagnostico // ignore: cast_nullable_to_non_nullable
as Diagnostico?,resueltoEstado: freezed == resueltoEstado ? _self.resueltoEstado : resueltoEstado // ignore: cast_nullable_to_non_nullable
as String?,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}
/// Create a copy of DiagnosticoFlowState
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
}/// Create a copy of DiagnosticoFlowState
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$DiagnosticoCopyWith<$Res>? get diagnostico {
    if (_self.diagnostico == null) {
    return null;
  }

  return $DiagnosticoCopyWith<$Res>(_self.diagnostico!, (value) {
    return _then(_self.copyWith(diagnostico: value));
  });
}
}


/// Adds pattern-matching-related methods to [DiagnosticoFlowState].
extension DiagnosticoFlowStatePatterns on DiagnosticoFlowState {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _DiagnosticoFlowState value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _DiagnosticoFlowState() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _DiagnosticoFlowState value)  $default,){
final _that = this;
switch (_that) {
case _DiagnosticoFlowState():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _DiagnosticoFlowState value)?  $default,){
final _that = this;
switch (_that) {
case _DiagnosticoFlowState() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( FlowFase fase,  Uint8List? previewBytes,  AnalisisYolo? analisis,  Diagnostico? diagnostico,  String? resueltoEstado,  String? error)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _DiagnosticoFlowState() when $default != null:
return $default(_that.fase,_that.previewBytes,_that.analisis,_that.diagnostico,_that.resueltoEstado,_that.error);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( FlowFase fase,  Uint8List? previewBytes,  AnalisisYolo? analisis,  Diagnostico? diagnostico,  String? resueltoEstado,  String? error)  $default,) {final _that = this;
switch (_that) {
case _DiagnosticoFlowState():
return $default(_that.fase,_that.previewBytes,_that.analisis,_that.diagnostico,_that.resueltoEstado,_that.error);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( FlowFase fase,  Uint8List? previewBytes,  AnalisisYolo? analisis,  Diagnostico? diagnostico,  String? resueltoEstado,  String? error)?  $default,) {final _that = this;
switch (_that) {
case _DiagnosticoFlowState() when $default != null:
return $default(_that.fase,_that.previewBytes,_that.analisis,_that.diagnostico,_that.resueltoEstado,_that.error);case _:
  return null;

}
}

}

/// @nodoc


class _DiagnosticoFlowState extends DiagnosticoFlowState {
  const _DiagnosticoFlowState({this.fase = FlowFase.inicial, this.previewBytes, this.analisis, this.diagnostico, this.resueltoEstado, this.error}): super._();
  

@override@JsonKey() final  FlowFase fase;
@override final  Uint8List? previewBytes;
@override final  AnalisisYolo? analisis;
@override final  Diagnostico? diagnostico;
@override final  String? resueltoEstado;
// ACEPTADO | RECHAZADO
@override final  String? error;

/// Create a copy of DiagnosticoFlowState
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$DiagnosticoFlowStateCopyWith<_DiagnosticoFlowState> get copyWith => __$DiagnosticoFlowStateCopyWithImpl<_DiagnosticoFlowState>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _DiagnosticoFlowState&&(identical(other.fase, fase) || other.fase == fase)&&const DeepCollectionEquality().equals(other.previewBytes, previewBytes)&&(identical(other.analisis, analisis) || other.analisis == analisis)&&(identical(other.diagnostico, diagnostico) || other.diagnostico == diagnostico)&&(identical(other.resueltoEstado, resueltoEstado) || other.resueltoEstado == resueltoEstado)&&(identical(other.error, error) || other.error == error));
}


@override
int get hashCode => Object.hash(runtimeType,fase,const DeepCollectionEquality().hash(previewBytes),analisis,diagnostico,resueltoEstado,error);

@override
String toString() {
  return 'DiagnosticoFlowState(fase: $fase, previewBytes: $previewBytes, analisis: $analisis, diagnostico: $diagnostico, resueltoEstado: $resueltoEstado, error: $error)';
}


}

/// @nodoc
abstract mixin class _$DiagnosticoFlowStateCopyWith<$Res> implements $DiagnosticoFlowStateCopyWith<$Res> {
  factory _$DiagnosticoFlowStateCopyWith(_DiagnosticoFlowState value, $Res Function(_DiagnosticoFlowState) _then) = __$DiagnosticoFlowStateCopyWithImpl;
@override @useResult
$Res call({
 FlowFase fase, Uint8List? previewBytes, AnalisisYolo? analisis, Diagnostico? diagnostico, String? resueltoEstado, String? error
});


@override $AnalisisYoloCopyWith<$Res>? get analisis;@override $DiagnosticoCopyWith<$Res>? get diagnostico;

}
/// @nodoc
class __$DiagnosticoFlowStateCopyWithImpl<$Res>
    implements _$DiagnosticoFlowStateCopyWith<$Res> {
  __$DiagnosticoFlowStateCopyWithImpl(this._self, this._then);

  final _DiagnosticoFlowState _self;
  final $Res Function(_DiagnosticoFlowState) _then;

/// Create a copy of DiagnosticoFlowState
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? fase = null,Object? previewBytes = freezed,Object? analisis = freezed,Object? diagnostico = freezed,Object? resueltoEstado = freezed,Object? error = freezed,}) {
  return _then(_DiagnosticoFlowState(
fase: null == fase ? _self.fase : fase // ignore: cast_nullable_to_non_nullable
as FlowFase,previewBytes: freezed == previewBytes ? _self.previewBytes : previewBytes // ignore: cast_nullable_to_non_nullable
as Uint8List?,analisis: freezed == analisis ? _self.analisis : analisis // ignore: cast_nullable_to_non_nullable
as AnalisisYolo?,diagnostico: freezed == diagnostico ? _self.diagnostico : diagnostico // ignore: cast_nullable_to_non_nullable
as Diagnostico?,resueltoEstado: freezed == resueltoEstado ? _self.resueltoEstado : resueltoEstado // ignore: cast_nullable_to_non_nullable
as String?,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

/// Create a copy of DiagnosticoFlowState
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
}/// Create a copy of DiagnosticoFlowState
/// with the given fields replaced by the non-null parameter values.
@override
@pragma('vm:prefer-inline')
$DiagnosticoCopyWith<$Res>? get diagnostico {
    if (_self.diagnostico == null) {
    return null;
  }

  return $DiagnosticoCopyWith<$Res>(_self.diagnostico!, (value) {
    return _then(_self.copyWith(diagnostico: value));
  });
}
}

// dart format on
