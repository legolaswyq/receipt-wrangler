import 'package:built_collection/built_collection.dart';
import 'package:built_value/json_object.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/income_only.dart';

api.Widget widgetWith(Map<String, JsonObject?>? configuration) {
  return api.Widget((b) => b
    ..id = 1
    ..dashboardId = 1
    ..widgetType = api.WidgetType.PIE_CHART
    ..configuration =
        configuration == null ? null : MapBuilder(configuration));
}

void main() {
  test('is income-only when the configuration flag is true', () {
    expect(isIncomeOnlyWidget(widgetWith({'incomeOnly': JsonObject(true)})),
        isTrue);
  });

  test('is a spending widget when the flag is false, absent or malformed', () {
    expect(isIncomeOnlyWidget(widgetWith({'incomeOnly': JsonObject(false)})),
        isFalse);
    expect(isIncomeOnlyWidget(widgetWith({'incomeOnly': JsonObject('yes')})),
        isFalse);
    expect(isIncomeOnlyWidget(widgetWith({'incomeOnly': null})), isFalse);
    expect(isIncomeOnlyWidget(widgetWith({})), isFalse);
    expect(isIncomeOnlyWidget(widgetWith(null)), isFalse);
  });
}
