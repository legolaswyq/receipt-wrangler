// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'bank_statement_import_result.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$BankStatementImportResult extends BankStatementImportResult {
  @override
  final int createdCount;

  factory _$BankStatementImportResult(
          [void Function(BankStatementImportResultBuilder)? updates]) =>
      (BankStatementImportResultBuilder()..update(updates))._build();

  _$BankStatementImportResult._({required this.createdCount}) : super._();
  @override
  BankStatementImportResult rebuild(
          void Function(BankStatementImportResultBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  BankStatementImportResultBuilder toBuilder() =>
      BankStatementImportResultBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is BankStatementImportResult &&
        createdCount == other.createdCount;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdCount.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'BankStatementImportResult')
          ..add('createdCount', createdCount))
        .toString();
  }
}

class BankStatementImportResultBuilder
    implements
        Builder<BankStatementImportResult, BankStatementImportResultBuilder> {
  _$BankStatementImportResult? _$v;

  int? _createdCount;
  int? get createdCount => _$this._createdCount;
  set createdCount(int? createdCount) => _$this._createdCount = createdCount;

  BankStatementImportResultBuilder() {
    BankStatementImportResult._defaults(this);
  }

  BankStatementImportResultBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdCount = $v.createdCount;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(BankStatementImportResult other) {
    _$v = other as _$BankStatementImportResult;
  }

  @override
  void update(void Function(BankStatementImportResultBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  BankStatementImportResult build() => _build();

  _$BankStatementImportResult _build() {
    final _$result = _$v ??
        _$BankStatementImportResult._(
          createdCount: BuiltValueNullFieldError.checkNotNull(
              createdCount, r'BankStatementImportResult', 'createdCount'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
