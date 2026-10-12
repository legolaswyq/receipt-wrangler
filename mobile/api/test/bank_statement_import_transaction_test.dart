import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for BankStatementImportTransaction
void main() {
  final instance = BankStatementImportTransactionBuilder();
  // TODO add properties to the builder and call build()

  group(BankStatementImportTransaction, () {
    // Transaction date as YYYY-MM-DD
    // String date
    test('to test the property `date`', () async {
      // TODO
    });

    // Receipt name
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // Transaction amount, always positive
    // String amount
    test('to test the property `amount`', () async {
      // TODO
    });

    // True for a deposit, false for a withdrawal
    // bool isIncome
    test('to test the property `isIncome`', () async {
      // TODO
    });

    // Category to assign. Required, and must be an income category, when isIncome is true
    // int categoryId
    test('to test the property `categoryId`', () async {
      // TODO
    });

  });
}
