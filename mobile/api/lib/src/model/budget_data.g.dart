// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'budget_data.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BudgetData extends BudgetData {
  @override
  final String? month;
  @override
  final double? income;
  @override
  final double? spent;
  @override
  final double? net;
  @override
  final double? untracked;
  @override
  final BuiltList<BudgetCategory> categories;

  factory _$BudgetData([void Function(BudgetDataBuilder)? updates]) =>
      (BudgetDataBuilder()..update(updates))._build();

  _$BudgetData._(
      {this.month,
      this.income,
      this.spent,
      this.net,
      this.untracked,
      required this.categories})
      : super._();
  @override
  BudgetData rebuild(void Function(BudgetDataBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BudgetDataBuilder toBuilder() => BudgetDataBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BudgetData &&
        month == other.month &&
        income == other.income &&
        spent == other.spent &&
        net == other.net &&
        untracked == other.untracked &&
        categories == other.categories;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, month.hashCode);
    _$hash = $jc(_$hash, income.hashCode);
    _$hash = $jc(_$hash, spent.hashCode);
    _$hash = $jc(_$hash, net.hashCode);
    _$hash = $jc(_$hash, untracked.hashCode);
    _$hash = $jc(_$hash, categories.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BudgetData')
          ..add('month', month)
          ..add('income', income)
          ..add('spent', spent)
          ..add('net', net)
          ..add('untracked', untracked)
          ..add('categories', categories))
        .toString();
  }
}

class BudgetDataBuilder implements Builder<BudgetData, BudgetDataBuilder> {
  _$BudgetData? _$v;

  String? _month;
  String? get month => _$this._month;
  set month(String? month) => _$this._month = month;

  double? _income;
  double? get income => _$this._income;
  set income(double? income) => _$this._income = income;

  double? _spent;
  double? get spent => _$this._spent;
  set spent(double? spent) => _$this._spent = spent;

  double? _net;
  double? get net => _$this._net;
  set net(double? net) => _$this._net = net;

  double? _untracked;
  double? get untracked => _$this._untracked;
  set untracked(double? untracked) => _$this._untracked = untracked;

  ListBuilder<BudgetCategory>? _categories;
  ListBuilder<BudgetCategory> get categories =>
      _$this._categories ??= ListBuilder<BudgetCategory>();
  set categories(ListBuilder<BudgetCategory>? categories) =>
      _$this._categories = categories;

  BudgetDataBuilder() {
    BudgetData._defaults(this);
  }

  BudgetDataBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _month = $v.month;
      _income = $v.income;
      _spent = $v.spent;
      _net = $v.net;
      _untracked = $v.untracked;
      _categories = $v.categories.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BudgetData other) {
    _$v = other as _$BudgetData;
  }

  @override
  void update(void Function(BudgetDataBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BudgetData build() => _build();

  _$BudgetData _build() {
    _$BudgetData _$result;
    try {
      _$result = _$v ??
          _$BudgetData._(
            month: month,
            income: income,
            spent: spent,
            net: net,
            untracked: untracked,
            categories: categories.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'categories';
        categories.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'BudgetData', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
