// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bank_statement_preview.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BankStatementPreview extends BankStatementPreview {
  @override
  final BuiltList<BankStatementAccount> accounts;
  @override
  final BuiltList<BankStatementTransaction> transactions;

  factory _$BankStatementPreview(
          [void Function(BankStatementPreviewBuilder)? updates]) =>
      (BankStatementPreviewBuilder()..update(updates))._build();

  _$BankStatementPreview._({required this.accounts, required this.transactions})
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
        accounts == other.accounts &&
        transactions == other.transactions;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, accounts.hashCode);
    _$hash = $jc(_$hash, transactions.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BankStatementPreview')
          ..add('accounts', accounts)
          ..add('transactions', transactions))
        .toString();
  }
}

class BankStatementPreviewBuilder
    implements Builder<BankStatementPreview, BankStatementPreviewBuilder> {
  _$BankStatementPreview? _$v;

  ListBuilder<BankStatementAccount>? _accounts;
  ListBuilder<BankStatementAccount> get accounts =>
      _$this._accounts ??= ListBuilder<BankStatementAccount>();
  set accounts(ListBuilder<BankStatementAccount>? accounts) =>
      _$this._accounts = accounts;

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
      _accounts = $v.accounts.toBuilder();
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
            accounts: accounts.build(),
            transactions: transactions.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'accounts';
        accounts.build();
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
