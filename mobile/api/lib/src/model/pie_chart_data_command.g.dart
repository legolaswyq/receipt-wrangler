// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'pie_chart_data_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PieChartDataCommand extends PieChartDataCommand {
  @override
  final ChartGrouping chartGrouping;
  @override
  final ReceiptPagedRequestFilter? filter;
  @override
  final String? startDate;
  @override
  final String? endDate;
  @override
  final bool? incomeOnly;

  factory _$PieChartDataCommand(
          [void Function(PieChartDataCommandBuilder)? updates]) =>
      (PieChartDataCommandBuilder()..update(updates))._build();

  _$PieChartDataCommand._(
      {required this.chartGrouping,
      this.filter,
      this.startDate,
      this.endDate,
      this.incomeOnly})
      : super._();
  @override
  PieChartDataCommand rebuild(
          void Function(PieChartDataCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PieChartDataCommandBuilder toBuilder() =>
      PieChartDataCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PieChartDataCommand &&
        chartGrouping == other.chartGrouping &&
        filter == other.filter &&
        startDate == other.startDate &&
        endDate == other.endDate &&
        incomeOnly == other.incomeOnly;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, chartGrouping.hashCode);
    _$hash = $jc(_$hash, filter.hashCode);
    _$hash = $jc(_$hash, startDate.hashCode);
    _$hash = $jc(_$hash, endDate.hashCode);
    _$hash = $jc(_$hash, incomeOnly.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PieChartDataCommand')
          ..add('chartGrouping', chartGrouping)
          ..add('filter', filter)
          ..add('startDate', startDate)
          ..add('endDate', endDate)
          ..add('incomeOnly', incomeOnly))
        .toString();
  }
}

class PieChartDataCommandBuilder
    implements Builder<PieChartDataCommand, PieChartDataCommandBuilder> {
  _$PieChartDataCommand? _$v;

  ChartGrouping? _chartGrouping;
  ChartGrouping? get chartGrouping => _$this._chartGrouping;
  set chartGrouping(ChartGrouping? chartGrouping) =>
      _$this._chartGrouping = chartGrouping;

  ReceiptPagedRequestFilterBuilder? _filter;
  ReceiptPagedRequestFilterBuilder get filter =>
      _$this._filter ??= ReceiptPagedRequestFilterBuilder();
  set filter(ReceiptPagedRequestFilterBuilder? filter) =>
      _$this._filter = filter;

  String? _startDate;
  String? get startDate => _$this._startDate;
  set startDate(String? startDate) => _$this._startDate = startDate;

  String? _endDate;
  String? get endDate => _$this._endDate;
  set endDate(String? endDate) => _$this._endDate = endDate;

  bool? _incomeOnly;
  bool? get incomeOnly => _$this._incomeOnly;
  set incomeOnly(bool? incomeOnly) => _$this._incomeOnly = incomeOnly;

  PieChartDataCommandBuilder() {
    PieChartDataCommand._defaults(this);
  }

  PieChartDataCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _chartGrouping = $v.chartGrouping;
      _filter = $v.filter?.toBuilder();
      _startDate = $v.startDate;
      _endDate = $v.endDate;
      _incomeOnly = $v.incomeOnly;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PieChartDataCommand other) {
    _$v = other as _$PieChartDataCommand;
  }

  @override
  void update(void Function(PieChartDataCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PieChartDataCommand build() => _build();

  _$PieChartDataCommand _build() {
    _$PieChartDataCommand _$result;
    try {
      _$result = _$v ??
          _$PieChartDataCommand._(
            chartGrouping: BuiltValueNullFieldError.checkNotNull(
                chartGrouping, r'PieChartDataCommand', 'chartGrouping'),
            filter: _filter?.build(),
            startDate: startDate,
            endDate: endDate,
            incomeOnly: incomeOnly,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'filter';
        _filter?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'PieChartDataCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
