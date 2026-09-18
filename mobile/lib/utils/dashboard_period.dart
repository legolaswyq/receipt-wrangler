/// Dashboard-level date range presets, ported from
/// desktop/src/dashboard/utils/dashboard-period.util.ts so both clients resolve
/// the same preset to the same window and persist the same wire values on
/// Dashboard.period / periodStartDate / periodEndDate.
class DashboardPeriod {
  static const thisMonth = "THIS_MONTH";
  static const lastMonth = "LAST_MONTH";
  static const last3Months = "LAST_3_MONTHS";
  static const thisYear = "THIS_YEAR";
  static const allTime = "ALL_TIME";
  static const custom = "CUSTOM";
}

const defaultDashboardPeriod = DashboardPeriod.thisMonth;

const dashboardPeriodOptions = <(String value, String label)>[
  (DashboardPeriod.thisMonth, "This month"),
  (DashboardPeriod.lastMonth, "Last month"),
  (DashboardPeriod.last3Months, "Last 3 months"),
  (DashboardPeriod.thisYear, "This year"),
  (DashboardPeriod.allTime, "All time"),
  (DashboardPeriod.custom, "Custom"),
];

class DashboardRange {
  const DashboardRange({required this.startDate, required this.endDate});

  final String startDate;
  final String endDate;
}

/// Turns a preset (plus optional custom dates, `yyyy-MM-dd`) into an
/// inclusive-start / exclusive-end RFC3339 window, matching
/// `resolveDashboardRange` on desktop. All Time (and anything unknown) yields
/// empty bounds, which the backend treats as no filter.
DashboardRange resolveDashboardRange(
  String? period, {
  String? startDate,
  String? endDate,
}) {
  final now = DateTime.now();
  final y = now.year;
  final m = now.month;
  String iso(DateTime d) => d.toUtc().toIso8601String();

  switch (period) {
    case DashboardPeriod.thisMonth:
      return DashboardRange(
        startDate: iso(DateTime(y, m, 1)),
        endDate: iso(DateTime(y, m + 1, 1)),
      );
    case DashboardPeriod.lastMonth:
      return DashboardRange(
        startDate: iso(DateTime(y, m - 1, 1)),
        endDate: iso(DateTime(y, m, 1)),
      );
    case DashboardPeriod.last3Months:
      return DashboardRange(
        startDate: iso(DateTime(y, m - 2, 1)),
        endDate: iso(DateTime(y, m + 1, 1)),
      );
    case DashboardPeriod.thisYear:
      return DashboardRange(
        startDate: iso(DateTime(y, 1, 1)),
        endDate: iso(DateTime(y + 1, 1, 1)),
      );
    case DashboardPeriod.custom:
      final start = (startDate != null && startDate.isNotEmpty)
          ? iso(DateTime.parse(startDate))
          : "";
      // The custom end date is inclusive to the user; make it exclusive by
      // advancing one day so the whole end day is included.
      final end = (endDate != null && endDate.isNotEmpty)
          ? iso(DateTime.parse(endDate).add(const Duration(days: 1)))
          : "";
      return DashboardRange(startDate: start, endDate: end);
    case DashboardPeriod.allTime:
    default:
      return const DashboardRange(startDate: "", endDate: "");
  }
}
