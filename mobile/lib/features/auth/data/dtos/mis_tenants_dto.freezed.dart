// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'mis_tenants_dto.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$TenantConRolDto {

 String get id; String get nombre; String get slug; String get rol; bool get esPropio;
/// Create a copy of TenantConRolDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$TenantConRolDtoCopyWith<TenantConRolDto> get copyWith => _$TenantConRolDtoCopyWithImpl<TenantConRolDto>(this as TenantConRolDto, _$identity);

  /// Serializes this TenantConRolDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is TenantConRolDto&&(identical(other.id, id) || other.id == id)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.slug, slug) || other.slug == slug)&&(identical(other.rol, rol) || other.rol == rol)&&(identical(other.esPropio, esPropio) || other.esPropio == esPropio));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,nombre,slug,rol,esPropio);

@override
String toString() {
  return 'TenantConRolDto(id: $id, nombre: $nombre, slug: $slug, rol: $rol, esPropio: $esPropio)';
}


}

/// @nodoc
abstract mixin class $TenantConRolDtoCopyWith<$Res>  {
  factory $TenantConRolDtoCopyWith(TenantConRolDto value, $Res Function(TenantConRolDto) _then) = _$TenantConRolDtoCopyWithImpl;
@useResult
$Res call({
 String id, String nombre, String slug, String rol, bool esPropio
});




}
/// @nodoc
class _$TenantConRolDtoCopyWithImpl<$Res>
    implements $TenantConRolDtoCopyWith<$Res> {
  _$TenantConRolDtoCopyWithImpl(this._self, this._then);

  final TenantConRolDto _self;
  final $Res Function(TenantConRolDto) _then;

/// Create a copy of TenantConRolDto
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? nombre = null,Object? slug = null,Object? rol = null,Object? esPropio = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,slug: null == slug ? _self.slug : slug // ignore: cast_nullable_to_non_nullable
as String,rol: null == rol ? _self.rol : rol // ignore: cast_nullable_to_non_nullable
as String,esPropio: null == esPropio ? _self.esPropio : esPropio // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [TenantConRolDto].
extension TenantConRolDtoPatterns on TenantConRolDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _TenantConRolDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _TenantConRolDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _TenantConRolDto value)  $default,){
final _that = this;
switch (_that) {
case _TenantConRolDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _TenantConRolDto value)?  $default,){
final _that = this;
switch (_that) {
case _TenantConRolDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String nombre,  String slug,  String rol,  bool esPropio)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _TenantConRolDto() when $default != null:
return $default(_that.id,_that.nombre,_that.slug,_that.rol,_that.esPropio);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String nombre,  String slug,  String rol,  bool esPropio)  $default,) {final _that = this;
switch (_that) {
case _TenantConRolDto():
return $default(_that.id,_that.nombre,_that.slug,_that.rol,_that.esPropio);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String nombre,  String slug,  String rol,  bool esPropio)?  $default,) {final _that = this;
switch (_that) {
case _TenantConRolDto() when $default != null:
return $default(_that.id,_that.nombre,_that.slug,_that.rol,_that.esPropio);case _:
  return null;

}
}

}

/// @nodoc

@JsonSerializable(fieldRename: FieldRename.snake)
class _TenantConRolDto implements TenantConRolDto {
  const _TenantConRolDto({required this.id, required this.nombre, required this.slug, required this.rol, required this.esPropio});
  factory _TenantConRolDto.fromJson(Map<String, dynamic> json) => _$TenantConRolDtoFromJson(json);

@override final  String id;
@override final  String nombre;
@override final  String slug;
@override final  String rol;
@override final  bool esPropio;

/// Create a copy of TenantConRolDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$TenantConRolDtoCopyWith<_TenantConRolDto> get copyWith => __$TenantConRolDtoCopyWithImpl<_TenantConRolDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$TenantConRolDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _TenantConRolDto&&(identical(other.id, id) || other.id == id)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.slug, slug) || other.slug == slug)&&(identical(other.rol, rol) || other.rol == rol)&&(identical(other.esPropio, esPropio) || other.esPropio == esPropio));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,nombre,slug,rol,esPropio);

@override
String toString() {
  return 'TenantConRolDto(id: $id, nombre: $nombre, slug: $slug, rol: $rol, esPropio: $esPropio)';
}


}

/// @nodoc
abstract mixin class _$TenantConRolDtoCopyWith<$Res> implements $TenantConRolDtoCopyWith<$Res> {
  factory _$TenantConRolDtoCopyWith(_TenantConRolDto value, $Res Function(_TenantConRolDto) _then) = __$TenantConRolDtoCopyWithImpl;
@override @useResult
$Res call({
 String id, String nombre, String slug, String rol, bool esPropio
});




}
/// @nodoc
class __$TenantConRolDtoCopyWithImpl<$Res>
    implements _$TenantConRolDtoCopyWith<$Res> {
  __$TenantConRolDtoCopyWithImpl(this._self, this._then);

  final _TenantConRolDto _self;
  final $Res Function(_TenantConRolDto) _then;

/// Create a copy of TenantConRolDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? nombre = null,Object? slug = null,Object? rol = null,Object? esPropio = null,}) {
  return _then(_TenantConRolDto(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,slug: null == slug ? _self.slug : slug // ignore: cast_nullable_to_non_nullable
as String,rol: null == rol ? _self.rol : rol // ignore: cast_nullable_to_non_nullable
as String,esPropio: null == esPropio ? _self.esPropio : esPropio // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}


/// @nodoc
mixin _$MisTenantsDto {

 List<TenantConRolDto> get tenants; String get propioId;
/// Create a copy of MisTenantsDto
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MisTenantsDtoCopyWith<MisTenantsDto> get copyWith => _$MisTenantsDtoCopyWithImpl<MisTenantsDto>(this as MisTenantsDto, _$identity);

  /// Serializes this MisTenantsDto to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MisTenantsDto&&const DeepCollectionEquality().equals(other.tenants, tenants)&&(identical(other.propioId, propioId) || other.propioId == propioId));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(tenants),propioId);

@override
String toString() {
  return 'MisTenantsDto(tenants: $tenants, propioId: $propioId)';
}


}

/// @nodoc
abstract mixin class $MisTenantsDtoCopyWith<$Res>  {
  factory $MisTenantsDtoCopyWith(MisTenantsDto value, $Res Function(MisTenantsDto) _then) = _$MisTenantsDtoCopyWithImpl;
@useResult
$Res call({
 List<TenantConRolDto> tenants, String propioId
});




}
/// @nodoc
class _$MisTenantsDtoCopyWithImpl<$Res>
    implements $MisTenantsDtoCopyWith<$Res> {
  _$MisTenantsDtoCopyWithImpl(this._self, this._then);

  final MisTenantsDto _self;
  final $Res Function(MisTenantsDto) _then;

/// Create a copy of MisTenantsDto
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? tenants = null,Object? propioId = null,}) {
  return _then(_self.copyWith(
tenants: null == tenants ? _self.tenants : tenants // ignore: cast_nullable_to_non_nullable
as List<TenantConRolDto>,propioId: null == propioId ? _self.propioId : propioId // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [MisTenantsDto].
extension MisTenantsDtoPatterns on MisTenantsDto {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MisTenantsDto value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MisTenantsDto() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MisTenantsDto value)  $default,){
final _that = this;
switch (_that) {
case _MisTenantsDto():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MisTenantsDto value)?  $default,){
final _that = this;
switch (_that) {
case _MisTenantsDto() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( List<TenantConRolDto> tenants,  String propioId)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MisTenantsDto() when $default != null:
return $default(_that.tenants,_that.propioId);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( List<TenantConRolDto> tenants,  String propioId)  $default,) {final _that = this;
switch (_that) {
case _MisTenantsDto():
return $default(_that.tenants,_that.propioId);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( List<TenantConRolDto> tenants,  String propioId)?  $default,) {final _that = this;
switch (_that) {
case _MisTenantsDto() when $default != null:
return $default(_that.tenants,_that.propioId);case _:
  return null;

}
}

}

/// @nodoc

@JsonSerializable(fieldRename: FieldRename.snake)
class _MisTenantsDto implements MisTenantsDto {
  const _MisTenantsDto({required final  List<TenantConRolDto> tenants, required this.propioId}): _tenants = tenants;
  factory _MisTenantsDto.fromJson(Map<String, dynamic> json) => _$MisTenantsDtoFromJson(json);

 final  List<TenantConRolDto> _tenants;
@override List<TenantConRolDto> get tenants {
  if (_tenants is EqualUnmodifiableListView) return _tenants;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_tenants);
}

@override final  String propioId;

/// Create a copy of MisTenantsDto
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MisTenantsDtoCopyWith<_MisTenantsDto> get copyWith => __$MisTenantsDtoCopyWithImpl<_MisTenantsDto>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$MisTenantsDtoToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _MisTenantsDto&&const DeepCollectionEquality().equals(other._tenants, _tenants)&&(identical(other.propioId, propioId) || other.propioId == propioId));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(_tenants),propioId);

@override
String toString() {
  return 'MisTenantsDto(tenants: $tenants, propioId: $propioId)';
}


}

/// @nodoc
abstract mixin class _$MisTenantsDtoCopyWith<$Res> implements $MisTenantsDtoCopyWith<$Res> {
  factory _$MisTenantsDtoCopyWith(_MisTenantsDto value, $Res Function(_MisTenantsDto) _then) = __$MisTenantsDtoCopyWithImpl;
@override @useResult
$Res call({
 List<TenantConRolDto> tenants, String propioId
});




}
/// @nodoc
class __$MisTenantsDtoCopyWithImpl<$Res>
    implements _$MisTenantsDtoCopyWith<$Res> {
  __$MisTenantsDtoCopyWithImpl(this._self, this._then);

  final _MisTenantsDto _self;
  final $Res Function(_MisTenantsDto) _then;

/// Create a copy of MisTenantsDto
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? tenants = null,Object? propioId = null,}) {
  return _then(_MisTenantsDto(
tenants: null == tenants ? _self._tenants : tenants // ignore: cast_nullable_to_non_nullable
as List<TenantConRolDto>,propioId: null == propioId ? _self.propioId : propioId // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}

// dart format on
