// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bank_statement_import_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BankStatementImportCommand extends BankStatementImportCommand {
  @override
  final int groupId;
  @override
  final int paidByUserId;
  @override
  final BuiltList<BankStatementImportTransaction> transactions;

  factory _$BankStatementImportCommand(
          [void Function(BankStatementImportCommandBuilder)? updates]) =>
      (BankStatementImportCommandBuilder()..update(updates))._build();

  _$BankStatementImportCommand._(
      {required this.groupId,
      required this.paidByUserId,
      required this.transactions})
      : super._();
  @override
  BankStatementImportCommand rebuild(
          void Function(BankStatementImportCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BankStatementImportCommandBuilder toBuilder() =>
      BankStatementImportCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BankStatementImportCommand &&
        groupId == other.groupId &&
        paidByUserId == other.paidByUserId &&
        transactions == other.transactions;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, paidByUserId.hashCode);
    _$hash = $jc(_$hash, transactions.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BankStatementImportCommand')
          ..add('groupId', groupId)
          ..add('paidByUserId', paidByUserId)
          ..add('transactions', transactions))
        .toString();
  }
}

class BankStatementImportCommandBuilder
    implements
        Builder<BankStatementImportCommand, BankStatementImportCommandBuilder> {
  _$BankStatementImportCommand? _$v;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _paidByUserId;
  int? get paidByUserId => _$this._paidByUserId;
  set paidByUserId(int? paidByUserId) => _$this._paidByUserId = paidByUserId;

  ListBuilder<BankStatementImportTransaction>? _transactions;
  ListBuilder<BankStatementImportTransaction> get transactions =>
      _$this._transactions ??= ListBuilder<BankStatementImportTransaction>();
  set transactions(ListBuilder<BankStatementImportTransaction>? transactions) =>
      _$this._transactions = transactions;

  BankStatementImportCommandBuilder() {
    BankStatementImportCommand._defaults(this);
  }

  BankStatementImportCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _groupId = $v.groupId;
      _paidByUserId = $v.paidByUserId;
      _transactions = $v.transactions.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BankStatementImportCommand other) {
    _$v = other as _$BankStatementImportCommand;
  }

  @override
  void update(void Function(BankStatementImportCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BankStatementImportCommand build() => _build();

  _$BankStatementImportCommand _build() {
    _$BankStatementImportCommand _$result;
    try {
      _$result = _$v ??
          _$BankStatementImportCommand._(
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'BankStatementImportCommand', 'groupId'),
            paidByUserId: BuiltValueNullFieldError.checkNotNull(
                paidByUserId, r'BankStatementImportCommand', 'paidByUserId'),
            transactions: transactions.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'transactions';
        transactions.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'BankStatementImportCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
