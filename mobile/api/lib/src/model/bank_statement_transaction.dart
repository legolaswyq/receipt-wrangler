//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bank_statement_transaction.g.dart';

/// BankStatementTransaction
///
/// Properties:
/// * [accountName] - Name of the account the row came from
/// * [accountNumber] - Number of the account the row came from
/// * [transfer] - True when the row moves money between the user's own accounts, so it is neither spending nor income
/// * [date] - Transaction date as YYYY-MM-DD
/// * [type] - Transaction type as printed on the statement
/// * [details] - Merchant or counterparty as printed on the statement
/// * [reference] - Extra reference line, such as a masked card number
/// * [name] - Suggested receipt name
/// * [amount] - Transaction amount, always positive
/// * [balance] - Running balance after the transaction
/// * [isIncome] - True for a deposit, false for a withdrawal
/// * [selected] - Whether the row is suggested for import by default
/// * [duplicate] - Whether a matching receipt already exists in the group
/// * [suggestedCategoryId] - Suggested category, when one could be worked out
@BuiltValue()
abstract class BankStatementTransaction implements Built<BankStatementTransaction, BankStatementTransactionBuilder> {
  /// Name of the account the row came from
  @BuiltValueField(wireName: r'accountName')
  String get accountName;

  /// Number of the account the row came from
  @BuiltValueField(wireName: r'accountNumber')
  String get accountNumber;

  /// True when the row moves money between the user's own accounts, so it is neither spending nor income
  @BuiltValueField(wireName: r'transfer')
  bool get transfer;

  /// Transaction date as YYYY-MM-DD
  @BuiltValueField(wireName: r'date')
  String get date;

  /// Transaction type as printed on the statement
  @BuiltValueField(wireName: r'type')
  String get type;

  /// Merchant or counterparty as printed on the statement
  @BuiltValueField(wireName: r'details')
  String get details;

  /// Extra reference line, such as a masked card number
  @BuiltValueField(wireName: r'reference')
  String get reference;

  /// Suggested receipt name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Transaction amount, always positive
  @BuiltValueField(wireName: r'amount')
  String get amount;

  /// Running balance after the transaction
  @BuiltValueField(wireName: r'balance')
  String get balance;

  /// True for a deposit, false for a withdrawal
  @BuiltValueField(wireName: r'isIncome')
  bool get isIncome;

  /// Whether the row is suggested for import by default
  @BuiltValueField(wireName: r'selected')
  bool get selected;

  /// Whether a matching receipt already exists in the group
  @BuiltValueField(wireName: r'duplicate')
  bool get duplicate;

  /// Suggested category, when one could be worked out
  @BuiltValueField(wireName: r'suggestedCategoryId')
  int? get suggestedCategoryId;

  BankStatementTransaction._();

  factory BankStatementTransaction([void updates(BankStatementTransactionBuilder b)]) = _$BankStatementTransaction;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BankStatementTransactionBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BankStatementTransaction> get serializer => _$BankStatementTransactionSerializer();
}

class _$BankStatementTransactionSerializer implements PrimitiveSerializer<BankStatementTransaction> {
  @override
  final Iterable<Type> types = const [BankStatementTransaction, _$BankStatementTransaction];

  @override
  final String wireName = r'BankStatementTransaction';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BankStatementTransaction object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'accountName';
    yield serializers.serialize(
      object.accountName,
      specifiedType: const FullType(String),
    );
    yield r'accountNumber';
    yield serializers.serialize(
      object.accountNumber,
      specifiedType: const FullType(String),
    );
    yield r'transfer';
    yield serializers.serialize(
      object.transfer,
      specifiedType: const FullType(bool),
    );
    yield r'date';
    yield serializers.serialize(
      object.date,
      specifiedType: const FullType(String),
    );
    yield r'type';
    yield serializers.serialize(
      object.type,
      specifiedType: const FullType(String),
    );
    yield r'details';
    yield serializers.serialize(
      object.details,
      specifiedType: const FullType(String),
    );
    yield r'reference';
    yield serializers.serialize(
      object.reference,
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
    yield r'balance';
    yield serializers.serialize(
      object.balance,
      specifiedType: const FullType(String),
    );
    yield r'isIncome';
    yield serializers.serialize(
      object.isIncome,
      specifiedType: const FullType(bool),
    );
    yield r'selected';
    yield serializers.serialize(
      object.selected,
      specifiedType: const FullType(bool),
    );
    yield r'duplicate';
    yield serializers.serialize(
      object.duplicate,
      specifiedType: const FullType(bool),
    );
    if (object.suggestedCategoryId != null) {
      yield r'suggestedCategoryId';
      yield serializers.serialize(
        object.suggestedCategoryId,
        specifiedType: const FullType(int),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    BankStatementTransaction object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BankStatementTransactionBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'accountName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.accountName = valueDes;
          break;
        case r'accountNumber':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.accountNumber = valueDes;
          break;
        case r'transfer':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.transfer = valueDes;
          break;
        case r'date':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.date = valueDes;
          break;
        case r'type':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.type = valueDes;
          break;
        case r'details':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.details = valueDes;
          break;
        case r'reference':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.reference = valueDes;
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
        case r'balance':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.balance = valueDes;
          break;
        case r'isIncome':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.isIncome = valueDes;
          break;
        case r'selected':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.selected = valueDes;
          break;
        case r'duplicate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.duplicate = valueDes;
          break;
        case r'suggestedCategoryId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.suggestedCategoryId = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BankStatementTransaction deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BankStatementTransactionBuilder();
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

