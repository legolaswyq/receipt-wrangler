//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'bank_statement_import_result.g.dart';

/// BankStatementImportResult
///
/// Properties:
/// * [createdCount] - Number of receipts created
@BuiltValue()
abstract class BankStatementImportResult implements Built<BankStatementImportResult, BankStatementImportResultBuilder> {
  /// Number of receipts created
  @BuiltValueField(wireName: r'createdCount')
  int get createdCount;

  BankStatementImportResult._();

  factory BankStatementImportResult([void updates(BankStatementImportResultBuilder b)]) = _$BankStatementImportResult;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(BankStatementImportResultBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<BankStatementImportResult> get serializer => _$BankStatementImportResultSerializer();
}

class _$BankStatementImportResultSerializer implements PrimitiveSerializer<BankStatementImportResult> {
  @override
  final Iterable<Type> types = const [BankStatementImportResult, _$BankStatementImportResult];

  @override
  final String wireName = r'BankStatementImportResult';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    BankStatementImportResult object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'createdCount';
    yield serializers.serialize(
      object.createdCount,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    BankStatementImportResult object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required BankStatementImportResultBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'createdCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.createdCount = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  BankStatementImportResult deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = BankStatementImportResultBuilder();
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

