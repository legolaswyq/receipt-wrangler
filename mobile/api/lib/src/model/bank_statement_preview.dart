//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/bank_statement_transaction.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bank_statement_preview.g.dart';

/// BankStatementPreview
///
/// Properties:
/// * [bank] - Bank the statement was recognised as
/// * [accountName] - Account name
/// * [accountNumber] - Account number
/// * [periodStart] - Statement period start as YYYY-MM-DD
/// * [periodEnd] - Statement period end as YYYY-MM-DD
/// * [transactions] - Parsed transactions, newest first
@BuiltValue()
abstract class BankStatementPreview implements Built<BankStatementPreview, BankStatementPreviewBuilder> {
  /// Bank the statement was recognised as
  @BuiltValueField(wireName: r'bank')
  String get bank;

  /// Account name
  @BuiltValueField(wireName: r'accountName')
  String get accountName;

  /// Account number
  @BuiltValueField(wireName: r'accountNumber')
  String get accountNumber;

  /// Statement period start as YYYY-MM-DD
  @BuiltValueField(wireName: r'periodStart')
  String get periodStart;

  /// Statement period end as YYYY-MM-DD
  @BuiltValueField(wireName: r'periodEnd')
  String get periodEnd;

  /// Parsed transactions, newest first
  @BuiltValueField(wireName: r'transactions')
  BuiltList<BankStatementTransaction> get transactions;

  BankStatementPreview._();

  factory BankStatementPreview([void updates(BankStatementPreviewBuilder b)]) = _$BankStatementPreview;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BankStatementPreviewBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BankStatementPreview> get serializer => _$BankStatementPreviewSerializer();
}

class _$BankStatementPreviewSerializer implements PrimitiveSerializer<BankStatementPreview> {
  @override
  final Iterable<Type> types = const [BankStatementPreview, _$BankStatementPreview];

  @override
  final String wireName = r'BankStatementPreview';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BankStatementPreview object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'bank';
    yield serializers.serialize(
      object.bank,
      specifiedType: const FullType(String),
    );
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
    yield r'periodStart';
    yield serializers.serialize(
      object.periodStart,
      specifiedType: const FullType(String),
    );
    yield r'periodEnd';
    yield serializers.serialize(
      object.periodEnd,
      specifiedType: const FullType(String),
    );
    yield r'transactions';
    yield serializers.serialize(
      object.transactions,
      specifiedType: const FullType(BuiltList, [FullType(BankStatementTransaction)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BankStatementPreview object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BankStatementPreviewBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'bank':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.bank = valueDes;
          break;
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
        case r'periodStart':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.periodStart = valueDes;
          break;
        case r'periodEnd':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.periodEnd = valueDes;
          break;
        case r'transactions':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(BankStatementTransaction)]),
          ) as BuiltList<BankStatementTransaction>;
          result.transactions.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BankStatementPreview deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BankStatementPreviewBuilder();
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

