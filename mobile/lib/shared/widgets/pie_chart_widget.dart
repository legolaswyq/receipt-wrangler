import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

/// A reusable pie chart widget that can display data with customizable styling.
class PieChartWidget extends StatelessWidget {
  const PieChartWidget({
    super.key,
    required this.data,
    this.height = 300,
    this.isLoading = false,
    this.noDataMessage = 'No data available',
    this.loadingMessage = 'Loading...',
  });

  /// List of data points to display in the chart
  final List<PieChartDataPoint> data;

  /// Height of the chart container
  final double height;

  /// Whether the chart is in loading state
  final bool isLoading;

  /// Message to display when there is no data
  final String noDataMessage;

  /// Message to display while loading
  final String loadingMessage;

  /// Fallback colors for slices without a stored category color (Uncategorized,
  /// tag / paid-by groupings). Mirrors the desktop pie chart's fallback palette
  /// so both clients render the same colors when the backend supplies none.
  /// A slice's own `color` (the category's stored hex) takes precedence — see
  /// [_colorFor].
  static const List<Color> defaultColors = [
    Color(0xFFFF6384),
    Color(0xFF36A2EB),
    Color(0xFFFFCE56),
    Color(0xFF4BC0C0),
    Color(0xFF9966FF),
    Color(0xFFFF9F40),
    Color(0xFFE7E9ED),
    Color(0xFF7C4DFF),
    Color(0xFFFF5252),
    Color(0xFF64FFDA),
    Color(0xFFFFD740),
    Color(0xFF448AFF),
  ];

  /// Resolves the color for slice [index]: the data point's own hex [color]
  /// (the category's stored color from the API) when present, otherwise the
  /// fallback palette by index. Matches the desktop precedence
  /// (`point.color || fallback[i % fallback.length]`).
  Color _colorFor(int index) {
    final parsed = _parseHexColor(data[index].color);
    return parsed ?? defaultColors[index % defaultColors.length];
  }

  /// Parses a `#RRGGBB` (or `RRGGBB`) hex string into an opaque [Color].
  /// Returns null for null/empty/malformed input so the caller can fall back.
  static Color? _parseHexColor(String? hex) {
    if (hex == null) return null;
    var value = hex.trim();
    if (value.startsWith('#')) value = value.substring(1);
    if (value.length != 6) return null;
    final rgb = int.tryParse(value, radix: 16);
    if (rgb == null) return null;
    return Color(0xFF000000 | rgb);
  }

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return SizedBox(
        height: height,
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const CircularProgressIndicator(),
              const SizedBox(height: 16),
              Text(loadingMessage),
            ],
          ),
        ),
      );
    }

    if (data.isEmpty) {
      return SizedBox(
        height: height,
        child: Center(
          child: Text(
            noDataMessage,
            style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                  color: Colors.grey,
                ),
          ),
        ),
      );
    }

    return SizedBox(
      height: height,
      child: Row(
        children: [
          Expanded(
            flex: 2,
            child: PieChart(
              PieChartData(
                sections: _buildSections(),
                sectionsSpace: 2,
                centerSpaceRadius: 0,
                pieTouchData: PieTouchData(enabled: false),
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            flex: 1,
            child: _buildLegend(context),
          ),
        ],
      ),
    );
  }

  List<PieChartSectionData> _buildSections() {
    final total = data.fold<double>(0, (sum, item) => sum + item.value.abs());

    return data.asMap().entries.map((entry) {
      final index = entry.key;
      final item = entry.value;
      final magnitude = item.value.abs();
      final percentage = total > 0 ? (magnitude / total * 100) : 0;
      final color = _colorFor(index);

      return PieChartSectionData(
        value: magnitude,
        // Only label slices over 5% so small slices don't clutter/overlap,
        // matching the desktop chart's datalabels threshold.
        title: percentage > 5 ? '${percentage.toStringAsFixed(1)}%' : '',
        color: color,
        radius: 80,
        titleStyle: const TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.bold,
          color: Colors.white,
        ),
      );
    }).toList();
  }

  Widget _buildLegend(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: data.asMap().entries.map((entry) {
          final index = entry.key;
          final item = entry.value;
          final color = _colorFor(index);

          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 4),
            child: Row(
              children: [
                Container(
                  width: 12,
                  height: 12,
                  decoration: BoxDecoration(
                    color: color,
                    shape: BoxShape.circle,
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    item.label,
                    style: Theme.of(context).textTheme.bodySmall,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
          );
        }).toList(),
      ),
    );
  }
}

/// A data point for the pie chart
class PieChartDataPoint {
  const PieChartDataPoint({
    required this.label,
    required this.value,
    this.color,
  });

  final String label;
  final double value;

  /// Stored category color as a hex string (e.g. `#4E79A7`) from the API, or
  /// null for buckets without one. Used to color the slice/legend; see
  /// [PieChartWidget._colorFor].
  final String? color;
}
