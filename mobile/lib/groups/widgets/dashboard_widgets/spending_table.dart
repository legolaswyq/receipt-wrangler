import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/utils/color.dart';
import 'package:receipt_wrangler_mobile/utils/currency.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../client/client.dart';
import '../constants/text_styles.dart';

/// Category Breakdown table — the mobile counterpart of the desktop
/// SpendingTable widget. Reuses the same `getPieChartData` endpoint as the pie
/// chart, then renders the buckets as a Category / Amount / % table sorted by
/// magnitude, with a Total footer. Kept in sync with
/// desktop/src/dashboard/spending-table/spending-table.component.*.
class SpendingTable extends StatefulWidget {
  const SpendingTable({super.key, required this.dashboardWidget});

  final api.Widget dashboardWidget;

  @override
  State<SpendingTable> createState() => _SpendingTableState();
}

class _SpendingTableState extends State<SpendingTable> {
  late Future _future;
  bool _isInitialized = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (!_isInitialized) {
      _loadData();
      _isInitialized = true;
    }
  }

  api.ChartGrouping _getChartGrouping() {
    var config = widget.dashboardWidget.configuration;
    if (config != null && config.containsKey('chartGrouping')) {
      var value = config['chartGrouping'];
      if (value != null) {
        try {
          return api.ChartGrouping.valueOf(value.asString);
        } catch (_) {
          // Fall through to the default.
        }
      }
    }
    return api.ChartGrouping.CATEGORIES;
  }

  String _getGroupingLabel() {
    switch (_getChartGrouping()) {
      case api.ChartGrouping.TAGS:
        return 'Tags';
      case api.ChartGrouping.PAIDBY:
        return 'Paid By';
      default:
        return 'Categories';
    }
  }

  void _loadData() {
    var groupId = int.tryParse(getGroupId(context)) ?? 0;
    var command = api.PieChartDataCommand((b) => b
      ..chartGrouping = _getChartGrouping());

    _future = OpenApiClient.client.getWidgetApi().getPieChartData(
          groupId: groupId,
          pieChartDataCommand: command,
        );
  }

  Widget _buildHeader() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          widget.dashboardWidget.name ?? 'Spending Table',
          style: dashboardWidgetNameStyle,
        ),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
          decoration: BoxDecoration(
            color: Colors.grey[200],
            borderRadius: BorderRadius.circular(12),
          ),
          child: Text(
            _getGroupingLabel(),
            style: const TextStyle(fontSize: 12),
          ),
        ),
      ],
    );
  }

  Widget _buildRow(String label, String amount, String percent,
      {Color? swatch, bool isTotal = false}) {
    var style = TextStyle(
      fontWeight: isTotal ? FontWeight.bold : FontWeight.normal,
    );
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          if (swatch != null) ...[
            Container(
              width: 12,
              height: 12,
              decoration:
                  BoxDecoration(color: swatch, shape: BoxShape.circle),
            ),
            const SizedBox(width: 8),
          ] else if (!isTotal) ...[
            const SizedBox(width: 20),
          ],
          Expanded(child: Text(label, style: style)),
          const SizedBox(width: 8),
          Text(amount, style: style),
          SizedBox(
            width: 56,
            child: Text(percent, style: style, textAlign: TextAlign.right),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      future: _future,
      builder: (context, snapshot) {
        var isLoading = snapshot.connectionState != ConnectionState.done;

        var points = <api.PieChartDataPoint>[];
        if (snapshot.hasData && snapshot.data?.data != null) {
          points = (snapshot.data!.data as api.PieChartData).data.toList();
        }

        // Sort by magnitude, largest first (matches the desktop table).
        points.sort((a, b) => b.value.abs().compareTo(a.value.abs()));
        var totalAbs =
            points.fold<double>(0, (sum, p) => sum + p.value.abs());
        var total = points.fold<double>(0, (sum, p) => sum + p.value);

        Widget body;
        if (isLoading) {
          body = const Padding(
            padding: EdgeInsets.all(16),
            child: Center(child: CircularProgressIndicator()),
          );
        } else if (points.isEmpty) {
          body = Padding(
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: Text('No spending in this range.',
                style: TextStyle(color: Colors.grey[600])),
          );
        } else {
          var rows = <Widget>[];
          for (var point in points) {
            var percent =
                totalAbs > 0 ? (point.value.abs() / totalAbs * 100) : 0;
            rows.add(_buildRow(
              point.label,
              formatCurrency(context, point.value.toString()) ?? '',
              '${percent.toStringAsFixed(1)}%',
              swatch: hexToColor(point.color),
            ));
            rows.add(const Divider(height: 1));
          }
          rows.add(_buildRow(
            'Total',
            formatCurrency(context, total.toString()) ?? '',
            '',
            isTotal: true,
          ));
          body = Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: rows,
          );
        }

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 10),
            _buildHeader(),
            const SizedBox(height: 4),
            body,
          ],
        );
      },
    );
  }
}
