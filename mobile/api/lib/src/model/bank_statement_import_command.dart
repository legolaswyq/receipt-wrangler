//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/bank_statement_import_transaction.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bank_statement_import_command.g.dart';

/// BankStatementImportCommand
///
/// Properties:
/// * [groupId] - Group to create the receipts in
/// * [paidByUserId] - User recorded as having paid
/// * [transactions] - Reviewed transactions to import
@BuiltValue()
abstract class BankStatementImportCommand implements Built<BankStatementImportCommand, BankStatementImportCommandBuilder> {
  /// Group to create the receipts in
  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  /// User recorded as having paid
  @BuiltValueField(wireName: r'paidByUserId')
  int get paidByUserId;

  /// Reviewed transactions to import
  @BuiltValueField(wireName: r'transactions')
  BuiltList<BankStatementImportTransaction> get transactions;

  BankStatementImportCommand._();

  factory BankStatementImportCommand([void updates(BankStatementImportCommandBuilder b)]) = _$BankStatementImportCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BankStatementImportCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BankStatementImportCommand> get serializer => _$BankStatementImportCommandSerializer();
}

class _$BankStatementImportCommandSerializer implements PrimitiveSerializer<BankStatementImportCommand> {
  @override
  final Iterable<Type> types = const [BankStatementImportCommand, _$BankStatementImportCommand];

  @override
  final String wireName = r'BankStatementImportCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BankStatementImportCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(int),
    );
    yield r'paidByUserId';
    yield serializers.serialize(
      object.paidByUserId,
      specifiedType: const FullType(int),
    );
    yield r'transactions';
    yield serializers.serialize(
      object.transactions,
      specifiedType: const FullType(BuiltList, [FullType(BankStatementImportTransaction)]),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BankStatementImportCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BankStatementImportCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'paidByUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.paidByUserId = valueDes;
          break;
        case r'transactions':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(BankStatementImportTransaction)]),
          ) as BuiltList<BankStatementImportTransaction>;
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
  BankStatementImportCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BankStatementImportCommandBuilder();
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

