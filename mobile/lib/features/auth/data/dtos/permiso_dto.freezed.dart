// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'permiso_dto.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$PermisoDto {

 String get codigo; String get nombre; String get descripcion; String get modulo;
/// Create a copy of PermisoDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PermisoDtoCopyWith<PermisoDto> get copyWith => _$PermisoDtoCopyWithImpl<PermisoDto>(this as PermisoDto, _$identity);

  /// Serializes this PermisoDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is PermisoDto&&(identical(other.codigo, codigo) || other.codigo == codigo)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.modulo, modulo) || other.modulo == modulo));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,codigo,nombre,descripcion,modulo);

@override
String toString() {
  return 'PermisoDto(codigo: $codigo, nombre: $nombre, descripcion: $descripcion, modulo: $modulo)';
}


}

/// @nodoc
abstract mixin class $PermisoDtoCopyWith<$Res>  {
  factory $PermisoDtoCopyWith(PermisoDto value, $Res Function(PermisoDto) _then) = _$PermisoDtoCopyWithImpl;
@useResult
$Res call({
 String codigo, String nombre, String descripcion, String modulo
});




}
/// @nodoc
class _$PermisoDtoCopyWithImpl<$Res>
    implements $PermisoDtoCopyWith<$Res> {
  _$PermisoDtoCopyWithImpl(this._self, this._then);

  final PermisoDto _self;
  final $Res Function(PermisoDto) _then;

/// Create a copy of PermisoDto
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? codigo = null,Object? nombre = null,Object? descripcion = null,Object? modulo = null,}) {
  return _then(_self.copyWith(
codigo: null == codigo ? _self.codigo : codigo // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,modulo: null == modulo ? _self.modulo : modulo // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [PermisoDto].
extension PermisoDtoPatterns on PermisoDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _PermisoDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _PermisoDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _PermisoDto value)  $default,){
final _that = this;
switch (_that) {
case _PermisoDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _PermisoDto value)?  $default,){
final _that = this;
switch (_that) {
case _PermisoDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String codigo,  String nombre,  String descripcion,  String modulo)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _PermisoDto() when $default != null:
return $default(_that.codigo,_that.nombre,_that.descripcion,_that.modulo);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String codigo,  String nombre,  String descripcion,  String modulo)  $default,) {final _that = this;
switch (_that) {
case _PermisoDto():
return $default(_that.codigo,_that.nombre,_that.descripcion,_that.modulo);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String codigo,  String nombre,  String descripcion,  String modulo)?  $default,) {final _that = this;
switch (_that) {
case _PermisoDto() when $default != null:
return $default(_that.codigo,_that.nombre,_that.descripcion,_that.modulo);case _:
  return null;

}
}

}

/// @nodoc

@JsonSerializable(fieldRename: FieldRename.snake)
class _PermisoDto implements PermisoDto {
  const _PermisoDto({required this.codigo, required this.nombre, required this.descripcion, required this.modulo});
  factory _PermisoDto.fromJson(Map<String, dynamic> json) => _$PermisoDtoFromJson(json);

@override final  String codigo;
@override final  String nombre;
@override final  String descripcion;
@override final  String modulo;

/// Create a copy of PermisoDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PermisoDtoCopyWith<_PermisoDto> get copyWith => __$PermisoDtoCopyWithImpl<_PermisoDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$PermisoDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _PermisoDto&&(identical(other.codigo, codigo) || other.codigo == codigo)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.descripcion, descripcion) || other.descripcion == descripcion)&&(identical(other.modulo, modulo) || other.modulo == modulo));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,codigo,nombre,descripcion,modulo);

@override
String toString() {
  return 'PermisoDto(codigo: $codigo, nombre: $nombre, descripcion: $descripcion, modulo: $modulo)';
}


}

/// @nodoc
abstract mixin class _$PermisoDtoCopyWith<$Res> implements $PermisoDtoCopyWith<$Res> {
  factory _$PermisoDtoCopyWith(_PermisoDto value, $Res Function(_PermisoDto) _then) = __$PermisoDtoCopyWithImpl;
@override @useResult
$Res call({
 String codigo, String nombre, String descripcion, String modulo
});




}
/// @nodoc
class __$PermisoDtoCopyWithImpl<$Res>
    implements _$PermisoDtoCopyWith<$Res> {
  __$PermisoDtoCopyWithImpl(this._self, this._then);

  final _PermisoDto _self;
  final $Res Function(_PermisoDto) _then;

/// Create a copy of PermisoDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? codigo = null,Object? nombre = null,Object? descripcion = null,Object? modulo = null,}) {
  return _then(_PermisoDto(
codigo: null == codigo ? _self.codigo : codigo // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,descripcion: null == descripcion ? _self.descripcion : descripcion // ignore: cast_nullable_to_non_nullable
as String,modulo: null == modulo ? _self.modulo : modulo // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}


/// @nodoc
mixin _$MisPermisosDto {

 List<PermisoDto> get permisos;
/// Create a copy of MisPermisosDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MisPermisosDtoCopyWith<MisPermisosDto> get copyWith => _$MisPermisosDtoCopyWithImpl<MisPermisosDto>(this as MisPermisosDto, _$identity);

  /// Serializes this MisPermisosDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MisPermisosDto&&const DeepCollectionEquality().equals(other.permisos, permisos));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(permisos));

@override
String toString() {
  return 'MisPermisosDto(permisos: $permisos)';
}


}

/// @nodoc
abstract mixin class $MisPermisosDtoCopyWith<$Res>  {
  factory $MisPermisosDtoCopyWith(MisPermisosDto value, $Res Function(MisPermisosDto) _then) = _$MisPermisosDtoCopyWithImpl;
@useResult
$Res call({
 List<PermisoDto> permisos
});




}
/// @nodoc
class _$MisPermisosDtoCopyWithImpl<$Res>
    implements $MisPermisosDtoCopyWith<$Res> {
  _$MisPermisosDtoCopyWithImpl(this._self, this._then);

  final MisPermisosDto _self;
  final $Res Function(MisPermisosDto) _then;

/// Create a copy of MisPermisosDto
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? permisos = null,}) {
  return _then(_self.copyWith(
permisos: null == permisos ? _self.permisos : permisos // ignore: cast_nullable_to_non_nullable
as List<PermisoDto>,
  ));
}

}


/// Adds pattern-matching-related methods to [MisPermisosDto].
extension MisPermisosDtoPatterns on MisPermisosDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MisPermisosDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MisPermisosDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MisPermisosDto value)  $default,){
final _that = this;
switch (_that) {
case _MisPermisosDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MisPermisosDto value)?  $default,){
final _that = this;
switch (_that) {
case _MisPermisosDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( List<PermisoDto> permisos)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MisPermisosDto() when $default != null:
return $default(_that.permisos);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( List<PermisoDto> permisos)  $default,) {final _that = this;
switch (_that) {
case _MisPermisosDto():
return $default(_that.permisos);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( List<PermisoDto> permisos)?  $default,) {final _that = this;
switch (_that) {
case _MisPermisosDto() when $default != null:
return $default(_that.permisos);case _:
  return null;

}
}

}

/// @nodoc

@JsonSerializable(fieldRename: FieldRename.snake)
class _MisPermisosDto implements MisPermisosDto {
  const _MisPermisosDto({required final  List<PermisoDto> permisos}): _permisos = permisos;
  factory _MisPermisosDto.fromJson(Map<String, dynamic> json) => _$MisPermisosDtoFromJson(json);

 final  List<PermisoDto> _permisos;
@override List<PermisoDto> get permisos {
  if (_permisos is EqualUnmodifiableListView) return _permisos;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_permisos);
}


/// Create a copy of MisPermisosDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MisPermisosDtoCopyWith<_MisPermisosDto> get copyWith => __$MisPermisosDtoCopyWithImpl<_MisPermisosDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MisPermisosDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _MisPermisosDto&&const DeepCollectionEquality().equals(other._permisos, _permisos));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(_permisos));

@override
String toString() {
  return 'MisPermisosDto(permisos: $permisos)';
}


}

/// @nodoc
abstract mixin class _$MisPermisosDtoCopyWith<$Res> implements $MisPermisosDtoCopyWith<$Res> {
  factory _$MisPermisosDtoCopyWith(_MisPermisosDto value, $Res Function(_MisPermisosDto) _then) = __$MisPermisosDtoCopyWithImpl;
@override @useResult
$Res call({
 List<PermisoDto> permisos
});




}
/// @nodoc
class __$MisPermisosDtoCopyWithImpl<$Res>
    implements _$MisPermisosDtoCopyWith<$Res> {
  __$MisPermisosDtoCopyWithImpl(this._self, this._then);

  final _MisPermisosDto _self;
  final $Res Function(_MisPermisosDto) _then;

/// Create a copy of MisPermisosDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? permisos = null,}) {
  return _then(_MisPermisosDto(
permisos: null == permisos ? _self._permisos : permisos // ignore: cast_nullable_to_non_nullable
as List<PermisoDto>,
  ));
}


}

// dart format on
