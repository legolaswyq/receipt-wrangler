//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/upsert_widget_command.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'upsert_dashboard_command.g.dart';

/// UpsertDashboardCommand
///
/// Properties:
/// * [name] - Dashboard name
/// * [groupId] - Group foreign key
/// * [widgets] - Widgets associated to dashboard
/// * [period] - Dashboard-level date range preset (THIS_MONTH, LAST_MONTH, LAST_3_MONTHS, THIS_YEAR, ALL_TIME, CUSTOM)
/// * [periodStartDate] - ISO start date, used only when period is CUSTOM
/// * [periodEndDate] - ISO end date, used only when period is CUSTOM
@BuiltValue()
abstract class UpsertDashboardCommand implements Built<UpsertDashboardCommand, UpsertDashboardCommandBuilder> {
  /// Dashboard name
  @BuiltValueField(wireName: r'name')
  String get name;

  /// Group foreign key
  @BuiltValueField(wireName: r'groupId')
  String get groupId;

  /// Widgets associated to dashboard
  @BuiltValueField(wireName: r'widgets')
  BuiltList<UpsertWidgetCommand>? get widgets;

  /// Dashboard-level date range preset (THIS_MONTH, LAST_MONTH, LAST_3_MONTHS, THIS_YEAR, ALL_TIME, CUSTOM)
  @BuiltValueField(wireName: r'period')
  String? get period;

  /// ISO start date, used only when period is CUSTOM
  @BuiltValueField(wireName: r'periodStartDate')
  String? get periodStartDate;

  /// ISO end date, used only when period is CUSTOM
  @BuiltValueField(wireName: r'periodEndDate')
  String? get periodEndDate;

  UpsertDashboardCommand._();

  factory UpsertDashboardCommand([void updates(UpsertDashboardCommandBuilder b)]) = _$UpsertDashboardCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(UpsertDashboardCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<UpsertDashboardCommand> get serializer => _$UpsertDashboardCommandSerializer();
}

class _$UpsertDashboardCommandSerializer implements PrimitiveSerializer<UpsertDashboardCommand> {
  @override
  final Iterable<Type> types = const [UpsertDashboardCommand, _$UpsertDashboardCommand];

  @override
  final String wireName = r'UpsertDashboardCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    UpsertDashboardCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'name';
    yield serializers.serialize(
      object.name,
      specifiedType: const FullType(String),
    );
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(String),
    );
    if (object.widgets != null) {
      yield r'widgets';
      yield serializers.serialize(
        object.widgets,
        specifiedType: const FullType(BuiltList, [FullType(UpsertWidgetCommand)]),
      );
    }
    if (object.period != null) {
      yield r'period';
      yield serializers.serialize(
        object.period,
        specifiedType: const FullType(String),
      );
    }
    if (object.periodStartDate != null) {
      yield r'periodStartDate';
      yield serializers.serialize(
        object.periodStartDate,
        specifiedType: const FullType(String),
      );
    }
    if (object.periodEndDate != null) {
      yield r'periodEndDate';
      yield serializers.serialize(
        object.periodEndDate,
        specifiedType: const FullType(String),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    UpsertDashboardCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required UpsertDashboardCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'name':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.name = valueDes;
          break;
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.groupId = valueDes;
          break;
        case r'widgets':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(UpsertWidgetCommand)]),
          ) as BuiltList<UpsertWidgetCommand>;
          result.widgets.replace(valueDes);
          break;
        case r'period':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.period = valueDes;
          break;
        case r'periodStartDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.periodStartDate = valueDes;
          break;
        case r'periodEndDate':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.periodEndDate = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  UpsertDashboardCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = UpsertDashboardCommandBuilder();
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

