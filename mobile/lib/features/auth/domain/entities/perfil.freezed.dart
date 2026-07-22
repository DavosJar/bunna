// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'perfil.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$Perfil {

 String get id; String get correo; String get nombre; String get apellido; String get telefono; String get estado; DateTime get creadoEn;
/// Create a copy of Perfil
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$PerfilCopyWith<Perfil> get copyWith => _$PerfilCopyWithImpl<Perfil>(this as Perfil, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Perfil&&(identical(other.id, id) || other.id == id)&&(identical(other.correo, correo) || other.correo == correo)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.apellido, apellido) || other.apellido == apellido)&&(identical(other.telefono, telefono) || other.telefono == telefono)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.creadoEn, creadoEn) || other.creadoEn == creadoEn));
}


@override
int get hashCode => Object.hash(runtimeType,id,correo,nombre,apellido,telefono,estado,creadoEn);

@override
String toString() {
  return 'Perfil(id: $id, correo: $correo, nombre: $nombre, apellido: $apellido, telefono: $telefono, estado: $estado, creadoEn: $creadoEn)';
}


}

/// @nodoc
abstract mixin class $PerfilCopyWith<$Res>  {
  factory $PerfilCopyWith(Perfil value, $Res Function(Perfil) _then) = _$PerfilCopyWithImpl;
@useResult
$Res call({
 String id, String correo, String nombre, String apellido, String telefono, String estado, DateTime creadoEn
});




}
/// @nodoc
class _$PerfilCopyWithImpl<$Res>
    implements $PerfilCopyWith<$Res> {
  _$PerfilCopyWithImpl(this._self, this._then);

  final Perfil _self;
  final $Res Function(Perfil) _then;

/// Create a copy of Perfil
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? correo = null,Object? nombre = null,Object? apellido = null,Object? telefono = null,Object? estado = null,Object? creadoEn = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,correo: null == correo ? _self.correo : correo // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,apellido: null == apellido ? _self.apellido : apellido // ignore: cast_nullable_to_non_nullable
as String,telefono: null == telefono ? _self.telefono : telefono // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,creadoEn: null == creadoEn ? _self.creadoEn : creadoEn // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}

}


/// Adds pattern-matching-related methods to [Perfil].
extension PerfilPatterns on Perfil {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Perfil value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Perfil() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Perfil value)  $default,){
final _that = this;
switch (_that) {
case _Perfil():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Perfil value)?  $default,){
final _that = this;
switch (_that) {
case _Perfil() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  String correo,  String nombre,  String apellido,  String telefono,  String estado,  DateTime creadoEn)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Perfil() when $default != null:
return $default(_that.id,_that.correo,_that.nombre,_that.apellido,_that.telefono,_that.estado,_that.creadoEn);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  String correo,  String nombre,  String apellido,  String telefono,  String estado,  DateTime creadoEn)  $default,) {final _that = this;
switch (_that) {
case _Perfil():
return $default(_that.id,_that.correo,_that.nombre,_that.apellido,_that.telefono,_that.estado,_that.creadoEn);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  String correo,  String nombre,  String apellido,  String telefono,  String estado,  DateTime creadoEn)?  $default,) {final _that = this;
switch (_that) {
case _Perfil() when $default != null:
return $default(_that.id,_that.correo,_that.nombre,_that.apellido,_that.telefono,_that.estado,_that.creadoEn);case _:
  return null;

}
}

}

/// @nodoc


class _Perfil implements Perfil {
  const _Perfil({required this.id, required this.correo, required this.nombre, required this.apellido, required this.telefono, required this.estado, required this.creadoEn});
  

@override final  String id;
@override final  String correo;
@override final  String nombre;
@override final  String apellido;
@override final  String telefono;
@override final  String estado;
@override final  DateTime creadoEn;

/// Create a copy of Perfil
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$PerfilCopyWith<_Perfil> get copyWith => __$PerfilCopyWithImpl<_Perfil>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Perfil&&(identical(other.id, id) || other.id == id)&&(identical(other.correo, correo) || other.correo == correo)&&(identical(other.nombre, nombre) || other.nombre == nombre)&&(identical(other.apellido, apellido) || other.apellido == apellido)&&(identical(other.telefono, telefono) || other.telefono == telefono)&&(identical(other.estado, estado) || other.estado == estado)&&(identical(other.creadoEn, creadoEn) || other.creadoEn == creadoEn));
}


@override
int get hashCode => Object.hash(runtimeType,id,correo,nombre,apellido,telefono,estado,creadoEn);

@override
String toString() {
  return 'Perfil(id: $id, correo: $correo, nombre: $nombre, apellido: $apellido, telefono: $telefono, estado: $estado, creadoEn: $creadoEn)';
}


}

/// @nodoc
abstract mixin class _$PerfilCopyWith<$Res> implements $PerfilCopyWith<$Res> {
  factory _$PerfilCopyWith(_Perfil value, $Res Function(_Perfil) _then) = __$PerfilCopyWithImpl;
@override @useResult
$Res call({
 String id, String correo, String nombre, String apellido, String telefono, String estado, DateTime creadoEn
});




}
/// @nodoc
class __$PerfilCopyWithImpl<$Res>
    implements _$PerfilCopyWith<$Res> {
  __$PerfilCopyWithImpl(this._self, this._then);

  final _Perfil _self;
  final $Res Function(_Perfil) _then;

/// Create a copy of Perfil
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? correo = null,Object? nombre = null,Object? apellido = null,Object? telefono = null,Object? estado = null,Object? creadoEn = null,}) {
  return _then(_Perfil(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,correo: null == correo ? _self.correo : correo // ignore: cast_nullable_to_non_nullable
as String,nombre: null == nombre ? _self.nombre : nombre // ignore: cast_nullable_to_non_nullable
as String,apellido: null == apellido ? _self.apellido : apellido // ignore: cast_nullable_to_non_nullable
as String,telefono: null == telefono ? _self.telefono : telefono // ignore: cast_nullable_to_non_nullable
as String,estado: null == estado ? _self.estado : estado // ignore: cast_nullable_to_non_nullable
as String,creadoEn: null == creadoEn ? _self.creadoEn : creadoEn // ignore: cast_nullable_to_non_nullable
as DateTime,
  ));
}


}

// dart format on
