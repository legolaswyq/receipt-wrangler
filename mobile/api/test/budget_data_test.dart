import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for BudgetData
void main() {
  final instance = BudgetDataBuilder();
  // TODO add properties to the builder and call build()

  group(BudgetData, () {
    // The month the budget data covers, formatted YYYY-MM
    // String month
    test('to test the property `month`', () async {
      // TODO
    });

    // Total income for the month
    // double income
    test('to test the property `income`', () async {
      // TODO
    });

    // Total expense spend for the month
    // double spent
    test('to test the property `spent`', () async {
      // TODO
    });

    // Income minus spent
    // double net
    test('to test the property `net`', () async {
      // TODO
    });

    // Expense spend touching no budgeted category
    // double untracked
    test('to test the property `untracked`', () async {
      // TODO
    });

    // Per-category budget vs. actual for the month
    // BuiltList<BudgetCategory> categories
    test('to test the property `categories`', () async {
      // TODO
    });

  });
}
