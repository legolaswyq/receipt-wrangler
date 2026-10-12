// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bank_statement_import_transaction.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BankStatementImportTransaction extends BankStatementImportTransaction {
  @override
  final String date;
  @override
  final String name;
  @override
  final String amount;
  @override
  final bool isIncome;
  @override
  final int? categoryId;

  factory _$BankStatementImportTransaction(
          [void Function(BankStatementImportTransactionBuilder)? updates]) =>
      (BankStatementImportTransactionBuilder()..update(updates))._build();

  _$BankStatementImportTransaction._(
      {required this.date,
      required this.name,
      required this.amount,
      required this.isIncome,
      this.categoryId})
      : super._();
  @override
  BankStatementImportTransaction rebuild(
          void Function(BankStatementImportTransactionBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BankStatementImportTransactionBuilder toBuilder() =>
      BankStatementImportTransactionBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BankStatementImportTransaction &&
        date == other.date &&
        name == other.name &&
        amount == other.amount &&
        isIncome == other.isIncome &&
        categoryId == other.categoryId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, isIncome.hashCode);
    _$hash = $jc(_$hash, categoryId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BankStatementImportTransaction')
          ..add('date', date)
          ..add('name', name)
          ..add('amount', amount)
          ..add('isIncome', isIncome)
          ..add('categoryId', categoryId))
        .toString();
  }
}

class BankStatementImportTransactionBuilder
    implements
        Builder<BankStatementImportTransaction,
            BankStatementImportTransactionBuilder> {
  _$BankStatementImportTransaction? _$v;

  String? _date;
  String? get date => _$this._date;
  set date(String? date) => _$this._date = date;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  bool? _isIncome;
  bool? get isIncome => _$this._isIncome;
  set isIncome(bool? isIncome) => _$this._isIncome = isIncome;

  int? _categoryId;
  int? get categoryId => _$this._categoryId;
  set categoryId(int? categoryId) => _$this._categoryId = categoryId;

  BankStatementImportTransactionBuilder() {
    BankStatementImportTransaction._defaults(this);
  }

  BankStatementImportTransactionBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _date = $v.date;
      _name = $v.name;
      _amount = $v.amount;
      _isIncome = $v.isIncome;
      _categoryId = $v.categoryId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BankStatementImportTransaction other) {
    _$v = other as _$BankStatementImportTransaction;
  }

  @override
  void update(void Function(BankStatementImportTransactionBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BankStatementImportTransaction build() => _build();

  _$BankStatementImportTransaction _build() {
    final _$result = _$v ??
        _$BankStatementImportTransaction._(
          date: BuiltValueNullFieldError.checkNotNull(
              date, r'BankStatementImportTransaction', 'date'),
          name: BuiltValueNullFieldError.checkNotNull(
              name, r'BankStatementImportTransaction', 'name'),
          amount: BuiltValueNullFieldError.checkNotNull(
              amount, r'BankStatementImportTransaction', 'amount'),
          isIncome: BuiltValueNullFieldError.checkNotNull(
              isIncome, r'BankStatementImportTransaction', 'isIncome'),
          categoryId: categoryId,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
