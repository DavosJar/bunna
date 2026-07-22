// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'token_session_dto.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$TokenSessionDto {

 String get accessToken; String get refreshToken; int get expiresIn; String get tokenType; String get usuarioId; String get tenantId; String get rol;
/// Create a copy of TokenSessionDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TokenSessionDtoCopyWith<TokenSessionDto> get copyWith => _$TokenSessionDtoCopyWithImpl<TokenSessionDto>(this as TokenSessionDto, _$identity);

  /// Serializes this TokenSessionDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is TokenSessionDto&&(identical(other.accessToken, accessToken) || other.accessToken == accessToken)&&(identical(other.refreshToken, refreshToken) || other.refreshToken == refreshToken)&&(identical(other.expiresIn, expiresIn) || other.expiresIn == expiresIn)&&(identical(other.tokenType, tokenType) || other.tokenType == tokenType)&&(identical(other.usuarioId, usuarioId) || other.usuarioId == usuarioId)&&(identical(other.tenantId, tenantId) || other.tenantId == tenantId)&&(identical(other.rol, rol) || other.rol == rol));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,accessToken,refreshToken,expiresIn,tokenType,usuarioId,tenantId,rol);

@override
String toString() {
  return 'TokenSessionDto(accessToken: $accessToken, refreshToken: $refreshToken, expiresIn: $expiresIn, tokenType: $tokenType, usuarioId: $usuarioId, tenantId: $tenantId, rol: $rol)';
}


}

/// @nodoc
abstract mixin class $TokenSessionDtoCopyWith<$Res>  {
  factory $TokenSessionDtoCopyWith(TokenSessionDto value, $Res Function(TokenSessionDto) _then) = _$TokenSessionDtoCopyWithImpl;
@useResult
$Res call({
 String accessToken, String refreshToken, int expiresIn, String tokenType, String usuarioId, String tenantId, String rol
});




}
/// @nodoc
class _$TokenSessionDtoCopyWithImpl<$Res>
    implements $TokenSessionDtoCopyWith<$Res> {
  _$TokenSessionDtoCopyWithImpl(this._self, this._then);

  final TokenSessionDto _self;
  final $Res Function(TokenSessionDto) _then;

/// Create a copy of TokenSessionDto
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? accessToken = null,Object? refreshToken = null,Object? expiresIn = null,Object? tokenType = null,Object? usuarioId = null,Object? tenantId = null,Object? rol = null,}) {
  return _then(_self.copyWith(
accessToken: null == accessToken ? _self.accessToken : accessToken // ignore: cast_nullable_to_non_nullable
as String,refreshToken: null == refreshToken ? _self.refreshToken : refreshToken // ignore: cast_nullable_to_non_nullable
as String,expiresIn: null == expiresIn ? _self.expiresIn : expiresIn // ignore: cast_nullable_to_non_nullable
as int,tokenType: null == tokenType ? _self.tokenType : tokenType // ignore: cast_nullable_to_non_nullable
as String,usuarioId: null == usuarioId ? _self.usuarioId : usuarioId // ignore: cast_nullable_to_non_nullable
as String,tenantId: null == tenantId ? _self.tenantId : tenantId // ignore: cast_nullable_to_non_nullable
as String,rol: null == rol ? _self.rol : rol // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [TokenSessionDto].
extension TokenSessionDtoPatterns on TokenSessionDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _TokenSessionDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _TokenSessionDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _TokenSessionDto value)  $default,){
final _that = this;
switch (_that) {
case _TokenSessionDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _TokenSessionDto value)?  $default,){
final _that = this;
switch (_that) {
case _TokenSessionDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String accessToken,  String refreshToken,  int expiresIn,  String tokenType,  String usuarioId,  String tenantId,  String rol)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _TokenSessionDto() when $default != null:
return $default(_that.accessToken,_that.refreshToken,_that.expiresIn,_that.tokenType,_that.usuarioId,_that.tenantId,_that.rol);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String accessToken,  String refreshToken,  int expiresIn,  String tokenType,  String usuarioId,  String tenantId,  String rol)  $default,) {final _that = this;
switch (_that) {
case _TokenSessionDto():
return $default(_that.accessToken,_that.refreshToken,_that.expiresIn,_that.tokenType,_that.usuarioId,_that.tenantId,_that.rol);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String accessToken,  String refreshToken,  int expiresIn,  String tokenType,  String usuarioId,  String tenantId,  String rol)?  $default,) {final _that = this;
switch (_that) {
case _TokenSessionDto() when $default != null:
return $default(_that.accessToken,_that.refreshToken,_that.expiresIn,_that.tokenType,_that.usuarioId,_that.tenantId,_that.rol);case _:
  return null;

}
}

}

/// @nodoc

@JsonSerializable(fieldRename: FieldRename.snake)
class _TokenSessionDto implements TokenSessionDto {
  const _TokenSessionDto({required this.accessToken, required this.refreshToken, required this.expiresIn, required this.tokenType, required this.usuarioId, required this.tenantId, required this.rol});
  factory _TokenSessionDto.fromJson(Map<String, dynamic> json) => _$TokenSessionDtoFromJson(json);

@override final  String accessToken;
@override final  String refreshToken;
@override final  int expiresIn;
@override final  String tokenType;
@override final  String usuarioId;
@override final  String tenantId;
@override final  String rol;

/// Create a copy of TokenSessionDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TokenSessionDtoCopyWith<_TokenSessionDto> get copyWith => __$TokenSessionDtoCopyWithImpl<_TokenSessionDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TokenSessionDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _TokenSessionDto&&(identical(other.accessToken, accessToken) || other.accessToken == accessToken)&&(identical(other.refreshToken, refreshToken) || other.refreshToken == refreshToken)&&(identical(other.expiresIn, expiresIn) || other.expiresIn == expiresIn)&&(identical(other.tokenType, tokenType) || other.tokenType == tokenType)&&(identical(other.usuarioId, usuarioId) || other.usuarioId == usuarioId)&&(identical(other.tenantId, tenantId) || other.tenantId == tenantId)&&(identical(other.rol, rol) || other.rol == rol));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,accessToken,refreshToken,expiresIn,tokenType,usuarioId,tenantId,rol);

@override
String toString() {
  return 'TokenSessionDto(accessToken: $accessToken, refreshToken: $refreshToken, expiresIn: $expiresIn, tokenType: $tokenType, usuarioId: $usuarioId, tenantId: $tenantId, rol: $rol)';
}


}

/// @nodoc
abstract mixin class _$TokenSessionDtoCopyWith<$Res> implements $TokenSessionDtoCopyWith<$Res> {
  factory _$TokenSessionDtoCopyWith(_TokenSessionDto value, $Res Function(_TokenSessionDto) _then) = __$TokenSessionDtoCopyWithImpl;
@override @useResult
$Res call({
 String accessToken, String refreshToken, int expiresIn, String tokenType, String usuarioId, String tenantId, String rol
});




}
/// @nodoc
class __$TokenSessionDtoCopyWithImpl<$Res>
    implements _$TokenSessionDtoCopyWith<$Res> {
  __$TokenSessionDtoCopyWithImpl(this._self, this._then);

  final _TokenSessionDto _self;
  final $Res Function(_TokenSessionDto) _then;

/// Create a copy of TokenSessionDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? accessToken = null,Object? refreshToken = null,Object? expiresIn = null,Object? tokenType = null,Object? usuarioId = null,Object? tenantId = null,Object? rol = null,}) {
  return _then(_TokenSessionDto(
accessToken: null == accessToken ? _self.accessToken : accessToken // ignore: cast_nullable_to_non_nullable
as String,refreshToken: null == refreshToken ? _self.refreshToken : refreshToken // ignore: cast_nullable_to_non_nullable
as String,expiresIn: null == expiresIn ? _self.expiresIn : expiresIn // ignore: cast_nullable_to_non_nullable
as int,tokenType: null == tokenType ? _self.tokenType : tokenType // ignore: cast_nullable_to_non_nullable
as String,usuarioId: null == usuarioId ? _self.usuarioId : usuarioId // ignore: cast_nullable_to_non_nullable
as String,tenantId: null == tenantId ? _self.tenantId : tenantId // ignore: cast_nullable_to_non_nullable
as String,rol: null == rol ? _self.rol : rol // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}

// dart format on
