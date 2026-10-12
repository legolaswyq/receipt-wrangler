//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bank_statement_import_transaction.g.dart';

/// BankStatementImportTransaction
///
/// Properties:
/// * [date] - Transaction date as YYYY-MM-DD
/// * [name] - Receipt name
/// * [amount] - Transaction amount, always positive
/// * [isIncome] - True for a deposit, false for a withdrawal
/// * [categoryId] - Category to assign. Required, and must be an income category, when isIncome is true
@BuiltValue()
abstract class BankStatementImportTransaction implements Built<BankStatementImportTransaction, BankStatementImportTransactionBuilder> {
  /// Transaction date as YYYY-MM-DD
  @BuiltValueField(wireName: r'date')
  String get date;

  /// Receipt name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Transaction amount, always positive
  @BuiltValueField(wireName: r'amount')
  String get amount;

  /// True for a deposit, false for a withdrawal
  @BuiltValueField(wireName: r'isIncome')
  bool get isIncome;

  /// Category to assign. Required, and must be an income category, when isIncome is true
  @BuiltValueField(wireName: r'categoryId')
  int? get categoryId;

  BankStatementImportTransaction._();

  factory BankStatementImportTransaction([void updates(BankStatementImportTransactionBuilder b)]) = _$BankStatementImportTransaction;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BankStatementImportTransactionBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BankStatementImportTransaction> get serializer => _$BankStatementImportTransactionSerializer();
}

class _$BankStatementImportTransactionSerializer implements PrimitiveSerializer<BankStatementImportTransaction> {
  @override
  final Iterable<Type> types = const [BankStatementImportTransaction, _$BankStatementImportTransaction];

  @override
  final String wireName = r'BankStatementImportTransaction';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BankStatementImportTransaction object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'date';
    yield serializers.serialize(
      object.date,
      specifiedType: const FullType(String),
    );
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'amount';
    yield serializers.serialize(
      object.amount,
      specifiedType: const FullType(String),
    );
    yield r'isIncome';
    yield serializers.serialize(
      object.isIncome,
      specifiedType: const FullType(bool),
    );
    if (object.categoryId != null) {
      yield r'categoryId';
      yield serializers.serialize(
        object.categoryId,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    BankStatementImportTransaction object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BankStatementImportTransactionBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'date':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.date = valueDes;
          break;
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'amount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.amount = valueDes;
          break;
        case r'isIncome':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isIncome = valueDes;
          break;
        case r'categoryId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.categoryId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BankStatementImportTransaction deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BankStatementImportTransactionBuilder();
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

