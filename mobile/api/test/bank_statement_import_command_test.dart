import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for BankStatementImportCommand
void main() {
  final instance = BankStatementImportCommandBuilder();
  // TODO add properties to the builder and call build()

  group(BankStatementImportCommand, () {
    // Group to create the receipts in
    // int groupId
    test('to test the property `groupId`', () async {
      // TODO
    });

    // User recorded as having paid
    // int paidByUserId
    test('to test the property `paidByUserId`', () async {
      // TODO
    });

    // Reviewed transactions to import
    // BuiltList<BankStatementImportTransaction> transactions
    test('to test the property `transactions`', () async {
      // TODO
    });

  });
}
