import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for BankStatementTransaction
void main() {
  final instance = BankStatementTransactionBuilder();
  // TODO add properties to the builder and call build()

  group(BankStatementTransaction, () {
    // Transaction date as YYYY-MM-DD
    // String date
    test('to test the property `date`', () async {
      // TODO
    });

    // Transaction type as printed on the statement
    // String type
    test('to test the property `type`', () async {
      // TODO
    });

    // Merchant or counterparty as printed on the statement
    // String details
    test('to test the property `details`', () async {
      // TODO
    });

    // Extra reference line, such as a masked card number
    // String reference
    test('to test the property `reference`', () async {
      // TODO
    });

    // Suggested receipt name
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // Transaction amount, always positive
    // String amount
    test('to test the property `amount`', () async {
      // TODO
    });

    // Running balance after the transaction
    // String balance
    test('to test the property `balance`', () async {
      // TODO
    });

    // True for a deposit, false for a withdrawal
    // bool isIncome
    test('to test the property `isIncome`', () async {
      // TODO
    });

    // Whether the row is suggested for import by default
    // bool selected
    test('to test the property `selected`', () async {
      // TODO
    });

    // Whether a matching receipt already exists in the group
    // bool duplicate
    test('to test the property `duplicate`', () async {
      // TODO
    });

    // Suggested category, when one could be worked out
    // int suggestedCategoryId
    test('to test the property `suggestedCategoryId`', () async {
      // TODO
    });

  });
}
