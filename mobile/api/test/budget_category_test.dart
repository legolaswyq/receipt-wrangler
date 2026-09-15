import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for BudgetCategory
void main() {
  final instance = BudgetCategoryBuilder();
  // TODO add properties to the builder and call build()

  group(BudgetCategory, () {
    // Category foreign key
    // int categoryId
    test('to test the property `categoryId`', () async {
      // TODO
    });

    // Category name
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // Monthly budget target for the category
    // double target
    test('to test the property `target`', () async {
      // TODO
    });

    // Amount spent against the category this month
    // double spent
    test('to test the property `spent`', () async {
      // TODO
    });

    // Whether spent exceeds target
    // bool over
    test('to test the property `over`', () async {
      // TODO
    });

  });
}
