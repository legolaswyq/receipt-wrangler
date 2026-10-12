// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bank_statement_preview.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BankStatementPreview extends BankStatementPreview {
  @override
  final String bank;
  @override
  final String accountName;
  @override
  final String accountNumber;
  @override
  final String periodStart;
  @override
  final String periodEnd;
  @override
  final BuiltList<BankStatementTransaction> transactions;

  factory _$BankStatementPreview(
          [void Function(BankStatementPreviewBuilder)? updates]) =>
      (BankStatementPreviewBuilder()..update(updates))._build();

  _$BankStatementPreview._(
      {required this.bank,
      required this.accountName,
      required this.accountNumber,
      required this.periodStart,
      required this.periodEnd,
      required this.transactions})
      : super._();
  @override
  BankStatementPreview rebuild(
          void Function(BankStatementPreviewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BankStatementPreviewBuilder toBuilder() =>
      BankStatementPreviewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BankStatementPreview &&
        bank == other.bank &&
        accountName == other.accountName &&
        accountNumber == other.accountNumber &&
        periodStart == other.periodStart &&
        periodEnd == other.periodEnd &&
        transactions == other.transactions;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, bank.hashCode);
    _$hash = $jc(_$hash, accountName.hashCode);
    _$hash = $jc(_$hash, accountNumber.hashCode);
    _$hash = $jc(_$hash, periodStart.hashCode);
    _$hash = $jc(_$hash, periodEnd.hashCode);
    _$hash = $jc(_$hash, transactions.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BankStatementPreview')
          ..add('bank', bank)
          ..add('accountName', accountName)
          ..add('accountNumber', accountNumber)
          ..add('periodStart', periodStart)
          ..add('periodEnd', periodEnd)
          ..add('transactions', transactions))
        .toString();
  }
}

class BankStatementPreviewBuilder
    implements Builder<BankStatementPreview, BankStatementPreviewBuilder> {
  _$BankStatementPreview? _$v;

  String? _bank;
  String? get bank => _$this._bank;
  set bank(String? bank) => _$this._bank = bank;

  String? _accountName;
  String? get accountName => _$this._accountName;
  set accountName(String? accountName) => _$this._accountName = accountName;

  String? _accountNumber;
  String? get accountNumber => _$this._accountNumber;
  set accountNumber(String? accountNumber) =>
      _$this._accountNumber = accountNumber;

  String? _periodStart;
  String? get periodStart => _$this._periodStart;
  set periodStart(String? periodStart) => _$this._periodStart = periodStart;

  String? _periodEnd;
  String? get periodEnd => _$this._periodEnd;
  set periodEnd(String? periodEnd) => _$this._periodEnd = periodEnd;

  ListBuilder<BankStatementTransaction>? _transactions;
  ListBuilder<BankStatementTransaction> get transactions =>
      _$this._transactions ??= ListBuilder<BankStatementTransaction>();
  set transactions(ListBuilder<BankStatementTransaction>? transactions) =>
      _$this._transactions = transactions;

  BankStatementPreviewBuilder() {
    BankStatementPreview._defaults(this);
  }

  BankStatementPreviewBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _bank = $v.bank;
      _accountName = $v.accountName;
      _accountNumber = $v.accountNumber;
      _periodStart = $v.periodStart;
      _periodEnd = $v.periodEnd;
      _transactions = $v.transactions.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BankStatementPreview other) {
    _$v = other as _$BankStatementPreview;
  }

  @override
  void update(void Function(BankStatementPreviewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BankStatementPreview build() => _build();

  _$BankStatementPreview _build() {
    _$BankStatementPreview _$result;
    try {
      _$result = _$v ??
          _$BankStatementPreview._(
            bank: BuiltValueNullFieldError.checkNotNull(
                bank, r'BankStatementPreview', 'bank'),
            accountName: BuiltValueNullFieldError.checkNotNull(
                accountName, r'BankStatementPreview', 'accountName'),
            accountNumber: BuiltValueNullFieldError.checkNotNull(
                accountNumber, r'BankStatementPreview', 'accountNumber'),
            periodStart: BuiltValueNullFieldError.checkNotNull(
                periodStart, r'BankStatementPreview', 'periodStart'),
            periodEnd: BuiltValueNullFieldError.checkNotNull(
                periodEnd, r'BankStatementPreview', 'periodEnd'),
            transactions: transactions.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'transactions';
        transactions.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'BankStatementPreview', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
