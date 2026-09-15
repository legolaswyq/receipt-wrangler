// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'category_budget.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CategoryBudget extends CategoryBudget {
  @override
  final int? id;
  @override
  final int? groupId;
  @override
  final int? categoryId;
  @override
  final String? amount;
  @override
  final String? createdAt;
  @override
  final String? updatedAt;

  factory _$CategoryBudget([void Function(CategoryBudgetBuilder)? updates]) =>
      (CategoryBudgetBuilder()..update(updates))._build();

  _$CategoryBudget._(
      {this.id,
      this.groupId,
      this.categoryId,
      this.amount,
      this.createdAt,
      this.updatedAt})
      : super._();
  @override
  CategoryBudget rebuild(void Function(CategoryBudgetBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CategoryBudgetBuilder toBuilder() => CategoryBudgetBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CategoryBudget &&
        id == other.id &&
        groupId == other.groupId &&
        categoryId == other.categoryId &&
        amount == other.amount &&
        createdAt == other.createdAt &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, categoryId.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'CategoryBudget')
          ..add('id', id)
          ..add('groupId', groupId)
          ..add('categoryId', categoryId)
          ..add('amount', amount)
          ..add('createdAt', createdAt)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class CategoryBudgetBuilder
    implements Builder<CategoryBudget, CategoryBudgetBuilder> {
  _$CategoryBudget? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _categoryId;
  int? get categoryId => _$this._categoryId;
  set categoryId(int? categoryId) => _$this._categoryId = categoryId;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  CategoryBudgetBuilder() {
    CategoryBudget._defaults(this);
  }

  CategoryBudgetBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _groupId = $v.groupId;
      _categoryId = $v.categoryId;
      _amount = $v.amount;
      _createdAt = $v.createdAt;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(CategoryBudget other) {
    _$v = other as _$CategoryBudget;
  }

  @override
  void update(void Function(CategoryBudgetBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CategoryBudget build() => _build();

  _$CategoryBudget _build() {
    final _$result = _$v ??
        _$CategoryBudget._(
          id: id,
          groupId: groupId,
          categoryId: categoryId,
          amount: amount,
          createdAt: createdAt,
          updatedAt: updatedAt,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
