import { FormOption } from "../../interfaces/form-option.interface";

// Dashboard-level date range presets. The value is what's persisted on
// Dashboard.period; the client resolves it to concrete dates at query time so a
// preset like "This month" always means the current month.
export const DashboardPeriod = {
  ThisMonth: "THIS_MONTH",
  LastMonth: "LAST_MONTH",
  Last3Months: "LAST_3_MONTHS",
  ThisYear: "THIS_YEAR",
  AllTime: "ALL_TIME",
  Custom: "CUSTOM",
} as const;

export const DEFAULT_DASHBOARD_PERIOD = DashboardPeriod.ThisMonth;

export const dashboardPeriodOptions: FormOption[] = [
  { value: DashboardPeriod.ThisMonth, displayValue: "This month" },
  { value: DashboardPeriod.LastMonth, displayValue: "Last month" },
  { value: DashboardPeriod.Last3Months, displayValue: "Last 3 months" },
  { value: DashboardPeriod.ThisYear, displayValue: "This year" },
  { value: DashboardPeriod.AllTime, displayValue: "All time" },
  { value: DashboardPeriod.Custom, displayValue: "Custom" },
];

export interface DashboardRange {
  startDate: string;
  endDate: string;
}

// resolveDashboardRange turns a preset (plus optional custom dates) into an
// inclusive-start / exclusive-end RFC3339 window. All Time (and anything
// unknown) yields empty bounds, which the backend treats as no filter.
export function resolveDashboardRange(
  period: string | undefined,
  startDate?: string,
  endDate?: string
): DashboardRange {
  const now = new Date();
  const y = now.getFullYear();
  const m = now.getMonth();
  const iso = (d: Date) => d.toISOString();

  switch (period) {
    case DashboardPeriod.ThisMonth:
      return { startDate: iso(new Date(y, m, 1)), endDate: iso(new Date(y, m + 1, 1)) };
    case DashboardPeriod.LastMonth:
      return { startDate: iso(new Date(y, m - 1, 1)), endDate: iso(new Date(y, m, 1)) };
    case DashboardPeriod.Last3Months:
      return { startDate: iso(new Date(y, m - 2, 1)), endDate: iso(new Date(y, m + 1, 1)) };
    case DashboardPeriod.ThisYear:
      return { startDate: iso(new Date(y, 0, 1)), endDate: iso(new Date(y + 1, 0, 1)) };
    case DashboardPeriod.Custom: {
      const start = startDate ? iso(new Date(startDate)) : "";
      // The custom end date is inclusive to the user; make it exclusive by
      // advancing one day so the whole end day is included.
      const end = endDate
        ? iso(new Date(new Date(endDate).getTime() + 24 * 60 * 60 * 1000))
        : "";
      return { startDate: start, endDate: end };
    }
    case DashboardPeriod.AllTime:
    default:
      return { startDate: "", endDate: "" };
  }
}
