import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for BudgetApi
void main() {
  final instance = Openapi().getBudgetApi();

  group(BudgetApi, () {
    // Delete category budget
    //
    // This will delete a category's monthly budget target for a group
    //
    //Future deleteBudget(int groupId, int categoryId) async
    test('test deleteBudget', () async {
      // TODO
    });

    // Get budget data
    //
    // This will get the current month's budget data for a group
    //
    //Future<BudgetData> getBudgetData(int groupId) async
    test('test getBudgetData', () async {
      // TODO
    });

    // Upsert category budget
    //
    // This will create or update a category's monthly budget target for a group
    //
    //Future<CategoryBudget> upsertBudget(int groupId, UpsertCategoryBudgetCommand upsertCategoryBudgetCommand) async
    test('test upsertBudget', () async {
      // TODO
    });

  });
}
