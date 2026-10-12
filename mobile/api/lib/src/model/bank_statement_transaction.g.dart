// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bank_statement_transaction.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BankStatementTransaction extends BankStatementTransaction {
  @override
  final String date;
  @override
  final String type;
  @override
  final String details;
  @override
  final String reference;
  @override
  final String name;
  @override
  final String amount;
  @override
  final String balance;
  @override
  final bool isIncome;
  @override
  final bool selected;
  @override
  final bool duplicate;
  @override
  final int? suggestedCategoryId;

  factory _$BankStatementTransaction(
          [void Function(BankStatementTransactionBuilder)? updates]) =>
      (BankStatementTransactionBuilder()..update(updates))._build();

  _$BankStatementTransaction._(
      {required this.date,
      required this.type,
      required this.details,
      required this.reference,
      required this.name,
      required this.amount,
      required this.balance,
      required this.isIncome,
      required this.selected,
      required this.duplicate,
      this.suggestedCategoryId})
      : super._();
  @override
  BankStatementTransaction rebuild(
          void Function(BankStatementTransactionBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BankStatementTransactionBuilder toBuilder() =>
      BankStatementTransactionBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BankStatementTransaction &&
        date == other.date &&
        type == other.type &&
        details == other.details &&
        reference == other.reference &&
        name == other.name &&
        amount == other.amount &&
        balance == other.balance &&
        isIncome == other.isIncome &&
        selected == other.selected &&
        duplicate == other.duplicate &&
        suggestedCategoryId == other.suggestedCategoryId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, details.hashCode);
    _$hash = $jc(_$hash, reference.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, balance.hashCode);
    _$hash = $jc(_$hash, isIncome.hashCode);
    _$hash = $jc(_$hash, selected.hashCode);
    _$hash = $jc(_$hash, duplicate.hashCode);
    _$hash = $jc(_$hash, suggestedCategoryId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BankStatementTransaction')
          ..add('date', date)
          ..add('type', type)
          ..add('details', details)
          ..add('reference', reference)
          ..add('name', name)
          ..add('amount', amount)
          ..add('balance', balance)
          ..add('isIncome', isIncome)
          ..add('selected', selected)
          ..add('duplicate', duplicate)
          ..add('suggestedCategoryId', suggestedCategoryId))
        .toString();
  }
}

class BankStatementTransactionBuilder
    implements
        Builder<BankStatementTransaction, BankStatementTransactionBuilder> {
  _$BankStatementTransaction? _$v;

  String? _date;
  String? get date => _$this._date;
  set date(String? date) => _$this._date = date;

  String? _type;
  String? get type => _$this._type;
  set type(String? type) => _$this._type = type;

  String? _details;
  String? get details => _$this._details;
  set details(String? details) => _$this._details = details;

  String? _reference;
  String? get reference => _$this._reference;
  set reference(String? reference) => _$this._reference = reference;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  String? _balance;
  String? get balance => _$this._balance;
  set balance(String? balance) => _$this._balance = balance;

  bool? _isIncome;
  bool? get isIncome => _$this._isIncome;
  set isIncome(bool? isIncome) => _$this._isIncome = isIncome;

  bool? _selected;
  bool? get selected => _$this._selected;
  set selected(bool? selected) => _$this._selected = selected;

  bool? _duplicate;
  bool? get duplicate => _$this._duplicate;
  set duplicate(bool? duplicate) => _$this._duplicate = duplicate;

  int? _suggestedCategoryId;
  int? get suggestedCategoryId => _$this._suggestedCategoryId;
  set suggestedCategoryId(int? suggestedCategoryId) =>
      _$this._suggestedCategoryId = suggestedCategoryId;

  BankStatementTransactionBuilder() {
    BankStatementTransaction._defaults(this);
  }

  BankStatementTransactionBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _date = $v.date;
      _type = $v.type;
      _details = $v.details;
      _reference = $v.reference;
      _name = $v.name;
      _amount = $v.amount;
      _balance = $v.balance;
      _isIncome = $v.isIncome;
      _selected = $v.selected;
      _duplicate = $v.duplicate;
      _suggestedCategoryId = $v.suggestedCategoryId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BankStatementTransaction other) {
    _$v = other as _$BankStatementTransaction;
  }

  @override
  void update(void Function(BankStatementTransactionBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BankStatementTransaction build() => _build();

  _$BankStatementTransaction _build() {
    final _$result = _$v ??
        _$BankStatementTransaction._(
          date: BuiltValueNullFieldError.checkNotNull(
              date, r'BankStatementTransaction', 'date'),
          type: BuiltValueNullFieldError.checkNotNull(
              type, r'BankStatementTransaction', 'type'),
          details: BuiltValueNullFieldError.checkNotNull(
              details, r'BankStatementTransaction', 'details'),
          reference: BuiltValueNullFieldError.checkNotNull(
              reference, r'BankStatementTransaction', 'reference'),
          name: BuiltValueNullFieldError.checkNotNull(
              name, r'BankStatementTransaction', 'name'),
          amount: BuiltValueNullFieldError.checkNotNull(
              amount, r'BankStatementTransaction', 'amount'),
          balance: BuiltValueNullFieldError.checkNotNull(
              balance, r'BankStatementTransaction', 'balance'),
          isIncome: BuiltValueNullFieldError.checkNotNull(
              isIncome, r'BankStatementTransaction', 'isIncome'),
          selected: BuiltValueNullFieldError.checkNotNull(
              selected, r'BankStatementTransaction', 'selected'),
          duplicate: BuiltValueNullFieldError.checkNotNull(
              duplicate, r'BankStatementTransaction', 'duplicate'),
          suggestedCategoryId: suggestedCategoryId,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
