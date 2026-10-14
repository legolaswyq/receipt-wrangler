// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bank_statement_account.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BankStatementAccount extends BankStatementAccount {
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

  factory _$BankStatementAccount(
          [void Function(BankStatementAccountBuilder)? updates]) =>
      (BankStatementAccountBuilder()..update(updates))._build();

  _$BankStatementAccount._(
      {required this.bank,
      required this.accountName,
      required this.accountNumber,
      required this.periodStart,
      required this.periodEnd})
      : super._();
  @override
  BankStatementAccount rebuild(
          void Function(BankStatementAccountBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BankStatementAccountBuilder toBuilder() =>
      BankStatementAccountBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BankStatementAccount &&
        bank == other.bank &&
        accountName == other.accountName &&
        accountNumber == other.accountNumber &&
        periodStart == other.periodStart &&
        periodEnd == other.periodEnd;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, bank.hashCode);
    _$hash = $jc(_$hash, accountName.hashCode);
    _$hash = $jc(_$hash, accountNumber.hashCode);
    _$hash = $jc(_$hash, periodStart.hashCode);
    _$hash = $jc(_$hash, periodEnd.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BankStatementAccount')
          ..add('bank', bank)
          ..add('accountName', accountName)
          ..add('accountNumber', accountNumber)
          ..add('periodStart', periodStart)
          ..add('periodEnd', periodEnd))
        .toString();
  }
}

class BankStatementAccountBuilder
    implements Builder<BankStatementAccount, BankStatementAccountBuilder> {
  _$BankStatementAccount? _$v;

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

  BankStatementAccountBuilder() {
    BankStatementAccount._defaults(this);
  }

  BankStatementAccountBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _bank = $v.bank;
      _accountName = $v.accountName;
      _accountNumber = $v.accountNumber;
      _periodStart = $v.periodStart;
      _periodEnd = $v.periodEnd;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BankStatementAccount other) {
    _$v = other as _$BankStatementAccount;
  }

  @override
  void update(void Function(BankStatementAccountBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BankStatementAccount build() => _build();

  _$BankStatementAccount _build() {
    final _$result = _$v ??
        _$BankStatementAccount._(
          bank: BuiltValueNullFieldError.checkNotNull(
              bank, r'BankStatementAccount', 'bank'),
          accountName: BuiltValueNullFieldError.checkNotNull(
              accountName, r'BankStatementAccount', 'accountName'),
          accountNumber: BuiltValueNullFieldError.checkNotNull(
              accountNumber, r'BankStatementAccount', 'accountNumber'),
          periodStart: BuiltValueNullFieldError.checkNotNull(
              periodStart, r'BankStatementAccount', 'periodStart'),
          periodEnd: BuiltValueNullFieldError.checkNotNull(
              periodEnd, r'BankStatementAccount', 'periodEnd'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
