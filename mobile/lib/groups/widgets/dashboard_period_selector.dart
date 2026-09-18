import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../utils/dashboard_period.dart';

/// Dashboard-level date range selector — the mobile counterpart of the
/// desktop dashboard's "Date range" dropdown (+ custom From/To date inputs).
/// Purely a controlled input: [period]/[customStart]/[customEnd] come from
/// the parent (GroupDashboard), which owns persistence and range resolution.
class DashboardPeriodSelector extends StatelessWidget {
  const DashboardPeriodSelector({
    super.key,
    required this.period,
    required this.customStart,
    required this.customEnd,
    required this.onPeriodChanged,
    required this.onCustomStartChanged,
    required this.onCustomEndChanged,
  });

  final String period;
  final String customStart;
  final String customEnd;
  final ValueChanged<String> onPeriodChanged;
  final ValueChanged<String> onCustomStartChanged;
  final ValueChanged<String> onCustomEndChanged;

  static final _isoFormat = DateFormat('yyyy-MM-dd');

  Future<void> _pickDate(
    BuildContext context,
    String current,
    ValueChanged<String> onChanged,
  ) async {
    final initial =
        current.isNotEmpty ? DateTime.tryParse(current) : null;
    final picked = await showDatePicker(
      context: context,
      initialDate: initial ?? DateTime.now(),
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
    );
    if (picked != null) {
      onChanged(_isoFormat.format(picked));
    }
  }

  Widget _buildDateField(
    BuildContext context,
    String label,
    String value,
    ValueChanged<String> onChanged,
  ) {
    return Expanded(
      child: InkWell(
        onTap: () => _pickDate(context, value, onChanged),
        child: InputDecorator(
          decoration: InputDecoration(labelText: label, isDense: true),
          child: Text(value.isEmpty ? 'Select' : value),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        DropdownButtonFormField<String>(
          // `initialValue` (not `value`) is intentional: this is a FormField,
          // which only reads its initial-value param once per State instance.
          // The caller gives this widget a fresh key when the selected
          // dashboard changes (see GroupDashboard.build), so a dashboard with
          // a different stored period gets a fresh FormField state seeded
          // correctly rather than stale from the previous dashboard.
          initialValue: period,
          decoration: const InputDecoration(labelText: 'Date range', isDense: true),
          items: dashboardPeriodOptions
              .map((option) => DropdownMenuItem(
                    value: option.$1,
                    child: Text(option.$2),
                  ))
              .toList(),
          onChanged: (value) {
            if (value != null) onPeriodChanged(value);
          },
        ),
        if (period == DashboardPeriod.custom) ...[
          const SizedBox(height: 8),
          Row(
            children: [
              _buildDateField(context, 'From', customStart, onCustomStartChanged),
              const SizedBox(width: 8),
              _buildDateField(context, 'To', customEnd, onCustomEndChanged),
            ],
          ),
        ],
      ],
    );
  }
}
