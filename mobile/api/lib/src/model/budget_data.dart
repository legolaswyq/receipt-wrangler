//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:openapi/src/model/budget_category.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'budget_data.g.dart';

/// BudgetData
///
/// Properties:
/// * [month] - The month the budget data covers, formatted YYYY-MM
/// * [income] - Total income for the month
/// * [spent] - Total expense spend for the month
/// * [net] - Income minus spent
/// * [untracked] - Expense spend touching no budgeted category
/// * [categories] - Per-category budget vs. actual for the month
@BuiltValue()
abstract class BudgetData implements Built<BudgetData, BudgetDataBuilder> {
  /// The month the budget data covers, formatted YYYY-MM
  @BuiltValueField(wireName: r'month')
  String? get month;

  /// Total income for the month
  @BuiltValueField(wireName: r'income')
  double? get income;

  /// Total expense spend for the month
  @BuiltValueField(wireName: r'spent')
  double? get spent;

  /// Income minus spent
  @BuiltValueField(wireName: r'net')
  double? get net;

  /// Expense spend touching no budgeted category
  @BuiltValueField(wireName: r'untracked')
  double? get untracked;

  /// Per-category budget vs. actual for the month
  @BuiltValueField(wireName: r'categories')
  BuiltList<BudgetCategory> get categories;

  BudgetData._();

  factory BudgetData([void updates(BudgetDataBuilder b)]) = _$BudgetData;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BudgetDataBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BudgetData> get serializer => _$BudgetDataSerializer();
}

class _$BudgetDataSerializer implements PrimitiveSerializer<BudgetData> {
  @override
  final Iterable<Type> types = const [BudgetData, _$BudgetData];

  @override
  final String wireName = r'BudgetData';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BudgetData object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    if (object.month != null) {
      yield r'month';
      yield serializers.serialize(
        object.month,
        specifiedType: const FullType(String),
      );
    }
    if (object.income != null) {
      yield r'income';
      yield serializers.serialize(
        object.income,
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
    if (object.net != null) {
      yield r'net';
      yield serializers.serialize(
        object.net,
        specifiedType: const FullType(double),
      );
    }
    if (object.untracked != null) {
      yield r'untracked';
      yield serializers.serialize(
        object.untracked,
        specifiedType: const FullType(double),
      );
    }
    yield r'categories';
    yield serializers.serialize(
      object.categories,
      specifiedType: const FullType(BuiltList, [FullType(BudgetCategory)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BudgetData object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BudgetDataBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'month':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.month = valueDes;
          break;
        case r'income':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(double),
          ) as double;
          result.income = valueDes;
          break;
        case r'spent':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(double),
          ) as double;
          result.spent = valueDes;
          break;
        case r'net':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(double),
          ) as double;
          result.net = valueDes;
          break;
        case r'untracked':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(double),
          ) as double;
          result.untracked = valueDes;
          break;
        case r'categories':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(BudgetCategory)]),
          ) as BuiltList<BudgetCategory>;
          result.categories.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BudgetData deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BudgetDataBuilder();
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

