import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_period_selector.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/group_activities.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/group_summary.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/pie_chart.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/report_widget.dart';
import 'package:receipt_wrangler_mobile/utils/dashboard_period.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../client/client.dart';
import 'dashboard_widgets/budget_widget.dart';
import 'dashboard_widgets/filtered_receipts.dart';
import 'dashboard_widgets/spending_table.dart';

class GroupDashboard extends StatefulWidget {
  GroupDashboard({super.key, required this.dashboards});

  @override
  State<GroupDashboard> createState() => _GroupDashboard();

  List<api.Dashboard> dashboards = [];
}

class _GroupDashboard extends State<GroupDashboard> {
  int? selectedDashboardIndex;

  // Local mutable copy so persisting a period change (via updateDashboard) can
  // replace the affected entry without mutating the widget-owned list.
  late List<api.Dashboard> _dashboards;

  String _period = defaultDashboardPeriod;
  String _customStart = "";
  String _customEnd = "";

  @override
  void initState() {
    super.initState();
    _dashboards = List.of(widget.dashboards);
    _syncPeriodFromDashboard(getSelectedDashboard(_dashboards));
  }

  void onGroupTap(api.Group group) {
    context.go("/groups/${group.id}");
  }

  void setSelectedDashboardIndex(int index) {
    setState(() {
      selectedDashboardIndex = index;
      _syncPeriodFromDashboard(getSelectedDashboard(_dashboards));
    });
  }

  /// Reflects the selected dashboard's stored period/custom dates in local
  /// state without persisting (mirrors desktop's `syncPeriodFromDashboard`).
  void _syncPeriodFromDashboard(api.Dashboard? dashboard) {
    _period = dashboard?.period ?? defaultDashboardPeriod;
    _customStart = dashboard?.periodStartDate ?? "";
    _customEnd = dashboard?.periodEndDate ?? "";
  }

  void _onPeriodChanged(String period) {
    setState(() => _period = period);
    _persistPeriod();
  }

  void _onCustomStartChanged(String value) {
    setState(() => _customStart = value);
    if (_period == DashboardPeriod.custom) _persistPeriod();
  }

  void _onCustomEndChanged(String value) {
    setState(() => _customEnd = value);
    if (_period == DashboardPeriod.custom) _persistPeriod();
  }

  /// Persists the current period selection to the selected dashboard (mirrors
  /// desktop's `persistPeriod`), then swaps the local copy for the server's
  /// response so `getSelectedDashboard` reflects it on rebuild.
  Future<void> _persistPeriod() async {
    final dashboard = getSelectedDashboard(_dashboards);
    if (dashboard == null) return;

    final widgetCommands = (dashboard.widgets ?? const <api.Widget>[])
        .map((w) => api.UpsertWidgetCommand((wb) => wb
              ..name = w.name
              ..widgetType = w.widgetType!
              ..configuration.replace(w.configuration ?? {})))
        .toList();

    var command = api.UpsertDashboardCommand((b) => b
      ..name = dashboard.name
      ..groupId = getGroupId(context)
      ..period = _period
      ..periodStartDate = _customStart
      ..periodEndDate = _customEnd
      ..widgets.replace(widgetCommands));

    var response = await OpenApiClient.client
        .getDashboardApi()
        .updateDashboard(dashboardId: dashboard.id, upsertDashboardCommand: command);

    var updated = response.data;
    if (!mounted || updated == null) return;

    setState(() {
      var index = _dashboards.indexWhere((d) => d.id == dashboard.id);
      if (index != -1) _dashboards[index] = updated;
    });
  }

  DashboardRange get _range => resolveDashboardRange(
        _period,
        startDate: _customStart,
        endDate: _customEnd,
      );

  Widget buildChoiceChipList(List<api.Dashboard> dashboards) {
    var widgets = <Widget>[];
    var effectiveIndex = selectedDashboardIndex ?? 0;

    for (int i = 0; i < dashboards.length; i++) {
      var dashboard = dashboards[i];
      var selected = i == effectiveIndex;
      var theme = Theme.of(context);

      widgets.add(ChoiceChip(
        key: Key(dashboard.id.toString()),
        label: Text(dashboards[i].name),
        selected: selected,
        selectedColor: theme.primaryColor,
        onSelected: (value) => setSelectedDashboardIndex(i),
      ));
      widgets.add(const SizedBox(width: 10));
    }

    return SizedBox(
        height: 50,
        child: ListView(
          scrollDirection: Axis.horizontal,
          children: widgets,
        ));
  }

  List<Widget> buildDashboardWidgets(
      api.Dashboard? dashboard, double widgetHeight) {
    var widgets = <Widget>[];
    var range = _range;

    if (dashboard != null) {
      for (var widget in (dashboard.widgets)?.toList() ?? []) {
        switch (widget.widgetType) {
          case api.WidgetType.FILTERED_RECEIPTS:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: FilteredReceipts(
                dashboardWidget: widget,
              ),
            ));
            break;
          case api.WidgetType.GROUP_SUMMARY:
            widgets.add(GroupSummary(
              dashboardWidget: widget,
            ));
            break;
          case api.WidgetType.GROUP_ACTIVITY:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: GroupActivities(
                dashboardWidget: widget,
              ),
            ));
            break;
          case api.WidgetType.PIE_CHART:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: DashboardPieChart(
                dashboardWidget: widget,
                startDate: range.startDate,
                endDate: range.endDate,
              ),
            ));
            break;
          case api.WidgetType.REPORT:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: ReportWidget(
                dashboardWidget: widget,
              ),
            ));
            break;
          case api.WidgetType.SPENDING_TABLE:
            // Sizes to its content (like GROUP_SUMMARY) so the table flows in
            // the outer dashboard ListView rather than scrolling in a fixed box.
            widgets.add(SpendingTable(
              dashboardWidget: widget,
              startDate: range.startDate,
              endDate: range.endDate,
            ));
            break;
          case api.WidgetType.BUDGET:
            // Always current-month (no date-range dependency), matching
            // desktop's BudgetComponent (no startDate/endDate inputs).
            widgets.add(BudgetWidget(
              dashboardWidget: widget,
            ));
            break;
        }
      }
    }

    return widgets;
  }

  api.Dashboard? getSelectedDashboard(List<api.Dashboard>? dashboards) {
    if (dashboards == null || dashboards.isEmpty) {
      return null;
    }
    var index = (selectedDashboardIndex ?? 0).clamp(0, dashboards.length - 1);
    return dashboards[index];
  }

  @override
  Widget build(BuildContext context) {
    if (_dashboards.isEmpty) {
      return const Center(child: Text("No dashboards found"));
    }

    api.Dashboard? selectedDashboard = getSelectedDashboard(_dashboards);
    var widgetHeight = MediaQuery.of(context).size.height * 0.6;
    List<Widget> children =
        buildDashboardWidgets(selectedDashboard, widgetHeight);

    // Only surface the dashboard switcher when there's more than one to pick
    // from; with a single dashboard the lone chip is just noise (and the group
    // concept is hidden, so a single unified view is the norm).
    return Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (_dashboards.length > 1) buildChoiceChipList(_dashboards),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: DashboardPeriodSelector(
              // Fresh key per selected dashboard so its internal FormField
              // state re-seeds from `period` when you switch dashboards,
              // instead of keeping the previous dashboard's stale selection
              // (DropdownButtonFormField only reads `initialValue` once).
              key: ValueKey('period-selector-${selectedDashboardIndex ?? 0}'),
              period: _period,
              customStart: _customStart,
              customEnd: _customEnd,
              onPeriodChanged: _onPeriodChanged,
              onCustomStartChanged: _onCustomStartChanged,
              onCustomEndChanged: _onCustomEndChanged,
            ),
          ),
          Expanded(child: ListView(children: children))
        ]);
  }
}
