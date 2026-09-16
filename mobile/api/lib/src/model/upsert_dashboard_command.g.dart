// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_dashboard_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertDashboardCommand extends UpsertDashboardCommand {
  @override
  final String name;
  @override
  final String groupId;
  @override
  final BuiltList<UpsertWidgetCommand>? widgets;
  @override
  final String? period;
  @override
  final String? periodStartDate;
  @override
  final String? periodEndDate;

  factory _$UpsertDashboardCommand(
          [void Function(UpsertDashboardCommandBuilder)? updates]) =>
      (UpsertDashboardCommandBuilder()..update(updates))._build();

  _$UpsertDashboardCommand._(
      {required this.name,
      required this.groupId,
      this.widgets,
      this.period,
      this.periodStartDate,
      this.periodEndDate})
      : super._();
  @override
  UpsertDashboardCommand rebuild(
          void Function(UpsertDashboardCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertDashboardCommandBuilder toBuilder() =>
      UpsertDashboardCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertDashboardCommand &&
        name == other.name &&
        groupId == other.groupId &&
        widgets == other.widgets &&
        period == other.period &&
        periodStartDate == other.periodStartDate &&
        periodEndDate == other.periodEndDate;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, widgets.hashCode);
    _$hash = $jc(_$hash, period.hashCode);
    _$hash = $jc(_$hash, periodStartDate.hashCode);
    _$hash = $jc(_$hash, periodEndDate.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertDashboardCommand')
          ..add('name', name)
          ..add('groupId', groupId)
          ..add('widgets', widgets)
          ..add('period', period)
          ..add('periodStartDate', periodStartDate)
          ..add('periodEndDate', periodEndDate))
        .toString();
  }
}

class UpsertDashboardCommandBuilder
    implements Builder<UpsertDashboardCommand, UpsertDashboardCommandBuilder> {
  _$UpsertDashboardCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _groupId;
  String? get groupId => _$this._groupId;
  set groupId(String? groupId) => _$this._groupId = groupId;

  ListBuilder<UpsertWidgetCommand>? _widgets;
  ListBuilder<UpsertWidgetCommand> get widgets =>
      _$this._widgets ??= ListBuilder<UpsertWidgetCommand>();
  set widgets(ListBuilder<UpsertWidgetCommand>? widgets) =>
      _$this._widgets = widgets;

  String? _period;
  String? get period => _$this._period;
  set period(String? period) => _$this._period = period;

  String? _periodStartDate;
  String? get periodStartDate => _$this._periodStartDate;
  set periodStartDate(String? periodStartDate) =>
      _$this._periodStartDate = periodStartDate;

  String? _periodEndDate;
  String? get periodEndDate => _$this._periodEndDate;
  set periodEndDate(String? periodEndDate) =>
      _$this._periodEndDate = periodEndDate;

  UpsertDashboardCommandBuilder() {
    UpsertDashboardCommand._defaults(this);
  }

  UpsertDashboardCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _groupId = $v.groupId;
      _widgets = $v.widgets?.toBuilder();
      _period = $v.period;
      _periodStartDate = $v.periodStartDate;
      _periodEndDate = $v.periodEndDate;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertDashboardCommand other) {
    _$v = other as _$UpsertDashboardCommand;
  }

  @override
  void update(void Function(UpsertDashboardCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertDashboardCommand build() => _build();

  _$UpsertDashboardCommand _build() {
    _$UpsertDashboardCommand _$result;
    try {
      _$result = _$v ??
          _$UpsertDashboardCommand._(
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UpsertDashboardCommand', 'name'),
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'UpsertDashboardCommand', 'groupId'),
            widgets: _widgets?.build(),
            period: period,
            periodStartDate: periodStartDate,
            periodEndDate: periodEndDate,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'widgets';
        _widgets?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'UpsertDashboardCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
