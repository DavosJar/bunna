// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'mis_tenants.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$MisTenants {

 List<TenantConRol> get tenants; String get propioId;
/// Create a copy of MisTenants
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$MisTenantsCopyWith<MisTenants> get copyWith => _$MisTenantsCopyWithImpl<MisTenants>(this as MisTenants, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is MisTenants&&const DeepCollectionEquality().equals(other.tenants, tenants)&&(identical(other.propioId, propioId) || other.propioId == propioId));
}


@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(tenants),propioId);

@override
String toString() {
  return 'MisTenants(tenants: $tenants, propioId: $propioId)';
}


}

/// @nodoc
abstract mixin class $MisTenantsCopyWith<$Res>  {
  factory $MisTenantsCopyWith(MisTenants value, $Res Function(MisTenants) _then) = _$MisTenantsCopyWithImpl;
@useResult
$Res call({
 List<TenantConRol> tenants, String propioId
});




}
/// @nodoc
class _$MisTenantsCopyWithImpl<$Res>
    implements $MisTenantsCopyWith<$Res> {
  _$MisTenantsCopyWithImpl(this._self, this._then);

  final MisTenants _self;
  final $Res Function(MisTenants) _then;

/// Create a copy of MisTenants
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? tenants = null,Object? propioId = null,}) {
  return _then(_self.copyWith(
tenants: null == tenants ? _self.tenants : tenants // ignore: cast_nullable_to_non_nullable
as List<TenantConRol>,propioId: null == propioId ? _self.propioId : propioId // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [MisTenants].
extension MisTenantsPatterns on MisTenants {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _MisTenants value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _MisTenants() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _MisTenants value)  $default,){
final _that = this;
switch (_that) {
case _MisTenants():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _MisTenants value)?  $default,){
final _that = this;
switch (_that) {
case _MisTenants() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( List<TenantConRol> tenants,  String propioId)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _MisTenants() when $default != null:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( List<TenantConRol> tenants,  String propioId)  $default,) {final _that = this;
switch (_that) {
case _MisTenants():
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( List<TenantConRol> tenants,  String propioId)?  $default,) {final _that = this;
switch (_that) {
case _MisTenants() when $default != null:
return $default(_that.tenants,_that.propioId);case _:
  return null;

}
}

}

/// @nodoc


class _MisTenants implements MisTenants {
  const _MisTenants({required final  List<TenantConRol> tenants, required this.propioId}): _tenants = tenants;
  

 final  List<TenantConRol> _tenants;
@override List<TenantConRol> get tenants {
  if (_tenants is EqualUnmodifiableListView) return _tenants;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_tenants);
}

@override final  String propioId;

/// Create a copy of MisTenants
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$MisTenantsCopyWith<_MisTenants> get copyWith => __$MisTenantsCopyWithImpl<_MisTenants>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _MisTenants&&const DeepCollectionEquality().equals(other._tenants, _tenants)&&(identical(other.propioId, propioId) || other.propioId == propioId));
}


@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(_tenants),propioId);

@override
String toString() {
  return 'MisTenants(tenants: $tenants, propioId: $propioId)';
}


}

/// @nodoc
abstract mixin class _$MisTenantsCopyWith<$Res> implements $MisTenantsCopyWith<$Res> {
  factory _$MisTenantsCopyWith(_MisTenants value, $Res Function(_MisTenants) _then) = __$MisTenantsCopyWithImpl;
@override @useResult
$Res call({
 List<TenantConRol> tenants, String propioId
});




}
/// @nodoc
class __$MisTenantsCopyWithImpl<$Res>
    implements _$MisTenantsCopyWith<$Res> {
  __$MisTenantsCopyWithImpl(this._self, this._then);

  final _MisTenants _self;
  final $Res Function(_MisTenants) _then;

/// Create a copy of MisTenants
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? tenants = null,Object? propioId = null,}) {
  return _then(_MisTenants(
tenants: null == tenants ? _self._tenants : tenants // ignore: cast_nullable_to_non_nullable
as List<TenantConRol>,propioId: null == propioId ? _self.propioId : propioId // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}

// dart format on
