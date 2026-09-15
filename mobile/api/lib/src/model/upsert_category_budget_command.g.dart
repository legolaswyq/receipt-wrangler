// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_category_budget_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertCategoryBudgetCommand extends UpsertCategoryBudgetCommand {
  @override
  final int categoryId;
  @override
  final String amount;

  factory _$UpsertCategoryBudgetCommand(
          [void Function(UpsertCategoryBudgetCommandBuilder)? updates]) =>
      (UpsertCategoryBudgetCommandBuilder()..update(updates))._build();

  _$UpsertCategoryBudgetCommand._(
      {required this.categoryId, required this.amount})
      : super._();
  @override
  UpsertCategoryBudgetCommand rebuild(
          void Function(UpsertCategoryBudgetCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertCategoryBudgetCommandBuilder toBuilder() =>
      UpsertCategoryBudgetCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertCategoryBudgetCommand &&
        categoryId == other.categoryId &&
        amount == other.amount;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, categoryId.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertCategoryBudgetCommand')
          ..add('categoryId', categoryId)
          ..add('amount', amount))
        .toString();
  }
}

class UpsertCategoryBudgetCommandBuilder
    implements
        Builder<UpsertCategoryBudgetCommand,
            UpsertCategoryBudgetCommandBuilder> {
  _$UpsertCategoryBudgetCommand? _$v;

  int? _categoryId;
  int? get categoryId => _$this._categoryId;
  set categoryId(int? categoryId) => _$this._categoryId = categoryId;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  UpsertCategoryBudgetCommandBuilder() {
    UpsertCategoryBudgetCommand._defaults(this);
  }

  UpsertCategoryBudgetCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _categoryId = $v.categoryId;
      _amount = $v.amount;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertCategoryBudgetCommand other) {
    _$v = other as _$UpsertCategoryBudgetCommand;
  }

  @override
  void update(void Function(UpsertCategoryBudgetCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertCategoryBudgetCommand build() => _build();

  _$UpsertCategoryBudgetCommand _build() {
    final _$result = _$v ??
        _$UpsertCategoryBudgetCommand._(
          categoryId: BuiltValueNullFieldError.checkNotNull(
              categoryId, r'UpsertCategoryBudgetCommand', 'categoryId'),
          amount: BuiltValueNullFieldError.checkNotNull(
              amount, r'UpsertCategoryBudgetCommand', 'amount'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
