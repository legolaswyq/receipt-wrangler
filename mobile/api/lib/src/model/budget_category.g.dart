// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'budget_category.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BudgetCategory extends BudgetCategory {
  @override
  final int? categoryId;
  @override
  final String? name;
  @override
  final double? target;
  @override
  final double? spent;
  @override
  final bool? over;

  factory _$BudgetCategory([void Function(BudgetCategoryBuilder)? updates]) =>
      (BudgetCategoryBuilder()..update(updates))._build();

  _$BudgetCategory._(
      {this.categoryId, this.name, this.target, this.spent, this.over})
      : super._();
  @override
  BudgetCategory rebuild(void Function(BudgetCategoryBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BudgetCategoryBuilder toBuilder() => BudgetCategoryBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BudgetCategory &&
        categoryId == other.categoryId &&
        name == other.name &&
        target == other.target &&
        spent == other.spent &&
        over == other.over;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, categoryId.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, target.hashCode);
    _$hash = $jc(_$hash, spent.hashCode);
    _$hash = $jc(_$hash, over.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BudgetCategory')
          ..add('categoryId', categoryId)
          ..add('name', name)
          ..add('target', target)
          ..add('spent', spent)
          ..add('over', over))
        .toString();
  }
}

class BudgetCategoryBuilder
    implements Builder<BudgetCategory, BudgetCategoryBuilder> {
  _$BudgetCategory? _$v;

  int? _categoryId;
  int? get categoryId => _$this._categoryId;
  set categoryId(int? categoryId) => _$this._categoryId = categoryId;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  double? _target;
  double? get target => _$this._target;
  set target(double? target) => _$this._target = target;

  double? _spent;
  double? get spent => _$this._spent;
  set spent(double? spent) => _$this._spent = spent;

  bool? _over;
  bool? get over => _$this._over;
  set over(bool? over) => _$this._over = over;

  BudgetCategoryBuilder() {
    BudgetCategory._defaults(this);
  }

  BudgetCategoryBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _categoryId = $v.categoryId;
      _name = $v.name;
      _target = $v.target;
      _spent = $v.spent;
      _over = $v.over;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BudgetCategory other) {
    _$v = other as _$BudgetCategory;
  }

  @override
  void update(void Function(BudgetCategoryBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BudgetCategory build() => _build();

  _$BudgetCategory _build() {
    final _$result = _$v ??
        _$BudgetCategory._(
          categoryId: categoryId,
          name: name,
          target: target,
          spent: spent,
          over: over,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
