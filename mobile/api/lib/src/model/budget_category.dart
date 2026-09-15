//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'budget_category.g.dart';

/// BudgetCategory
///
/// Properties:
/// * [categoryId] - Category foreign key
/// * [name] - Category name
/// * [target] - Monthly budget target for the category
/// * [spent] - Amount spent against the category this month
/// * [over] - Whether spent exceeds target
@BuiltValue()
abstract class BudgetCategory implements Built<BudgetCategory, BudgetCategoryBuilder> {
  /// Category foreign key
  @BuiltValueField(wireName: r'categoryId')
  int? get categoryId;

  /// Category name
  @BuiltValueField(wireName: r'name')
  String? get name;

  /// Monthly budget target for the category
  @BuiltValueField(wireName: r'target')
  double? get target;

  /// Amount spent against the category this month
  @BuiltValueField(wireName: r'spent')
  double? get spent;

  /// Whether spent exceeds target
  @BuiltValueField(wireName: r'over')
  bool? get over;

  BudgetCategory._();

  factory BudgetCategory([void updates(BudgetCategoryBuilder b)]) = _$BudgetCategory;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BudgetCategoryBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BudgetCategory> get serializer => _$BudgetCategorySerializer();
}

class _$BudgetCategorySerializer implements PrimitiveSerializer<BudgetCategory> {
  @override
  final Iterable<Type> types = const [BudgetCategory, _$BudgetCategory];

  @override
  final String wireName = r'BudgetCategory';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BudgetCategory object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.categoryId != null) {
      yield r'categoryId';
      yield serializers.serialize(
        object.categoryId,
        specifiedType: const FullType(int),
      );
    }
    if (object.name != null) {
      yield r'name';
      yield serializers.serialize(
        object.name,
        specifiedType: const FullType(String),
      );
    }
    if (object.target != null) {
      yield r'target';
      yield serializers.serialize(
        object.target,
        specifiedType: const FullType(double),
      );
    }
    if (object.spent != null) {
      yield r'spent';
      yield serializers.serialize(
        object.spent,
        specifiedType: const FullType(double),
      );
    }
    if (object.over != null) {
      yield r'over';
      yield serializers.serialize(
        object.over,
        specifiedType: const FullType(bool),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    BudgetCategory object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BudgetCategoryBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'categoryId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.categoryId = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'target':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(double),
          ) as double;
          result.target = valueDes;
          break;
        case r'spent':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(double),
          ) as double;
          result.spent = valueDes;
          break;
        case r'over':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.over = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BudgetCategory deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BudgetCategoryBuilder();
    final serializedList = (serialized as Iterable<Object?>).toList();
    final unhandled = <Object?>[];
    _deserializeProperties(
      serializers,
      serialized,
      specifiedType: specifiedType,
      serializedList: serializedList,
      unhandled: unhandled,
      result: result,
    );
    return result.build();
  }
}

