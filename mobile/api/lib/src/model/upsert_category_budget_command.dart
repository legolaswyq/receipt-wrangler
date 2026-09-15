//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_category_budget_command.g.dart';

/// UpsertCategoryBudgetCommand
///
/// Properties:
/// * [categoryId] - Category foreign key
/// * [amount] - Monthly budget target amount for the category
@BuiltValue()
abstract class UpsertCategoryBudgetCommand implements Built<UpsertCategoryBudgetCommand, UpsertCategoryBudgetCommandBuilder> {
  /// Category foreign key
  @BuiltValueField(wireName: r'categoryId')
  int get categoryId;

  /// Monthly budget target amount for the category
  @BuiltValueField(wireName: r'amount')
  String get amount;

  UpsertCategoryBudgetCommand._();

  factory UpsertCategoryBudgetCommand([void updates(UpsertCategoryBudgetCommandBuilder b)]) = _$UpsertCategoryBudgetCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertCategoryBudgetCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertCategoryBudgetCommand> get serializer => _$UpsertCategoryBudgetCommandSerializer();
}

class _$UpsertCategoryBudgetCommandSerializer implements PrimitiveSerializer<UpsertCategoryBudgetCommand> {
  @override
  final Iterable<Type> types = const [UpsertCategoryBudgetCommand, _$UpsertCategoryBudgetCommand];

  @override
  final String wireName = r'UpsertCategoryBudgetCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertCategoryBudgetCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'categoryId';
    yield serializers.serialize(
      object.categoryId,
      specifiedType: const FullType(int),
    );
    yield r'amount';
    yield serializers.serialize(
      object.amount,
      specifiedType: const FullType(String),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertCategoryBudgetCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertCategoryBudgetCommandBuilder result,
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
        case r'amount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.amount = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertCategoryBudgetCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertCategoryBudgetCommandBuilder();
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

