import 'package:openapi/openapi.dart' as api;

/// Whether a pie chart / spending table widget is configured to show income
/// (receipts in income-flagged categories) instead of spending. Mirrors the
/// desktop widgets' `configuration.incomeOnly`.
bool isIncomeOnlyWidget(api.Widget dashboardWidget) {
  var value = dashboardWidget.configuration?['incomeOnly'];
  return value != null && value.isBool && value.asBool;
}
