//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/bank_statement_account.dart';
import 'package:openapi/src/model/bank_statement_transaction.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bank_statement_preview.g.dart';

/// BankStatementPreview
///
/// Properties:
/// * [accounts] - The statements that were read, one per account
/// * [transactions] - Parsed transactions across all accounts, newest first
@BuiltValue()
abstract class BankStatementPreview implements Built<BankStatementPreview, BankStatementPreviewBuilder> {
  /// The statements that were read, one per account
  @BuiltValueField(wireName: r'accounts')
  BuiltList<BankStatementAccount> get accounts;

  /// Parsed transactions across all accounts, newest first
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
    yield r'accounts';
    yield serializers.serialize(
      object.accounts,
      specifiedType: const FullType(BuiltList, [FullType(BankStatementAccount)]),
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
        case r'accounts':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(BankStatementAccount)]),
          ) as BuiltList<BankStatementAccount>;
          result.accounts.replace(valueDes);
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

