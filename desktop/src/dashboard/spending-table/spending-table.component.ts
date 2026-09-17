import { CommonModule, CurrencyPipe } from "@angular/common";
import { Component, OnChanges, OnInit, SimpleChanges, computed, input, signal } from "@angular/core";
import { take, tap } from "rxjs";
import { ChartGrouping, PieChartData, PieChartDataCommand, PieChartDataPoint, Widget, WidgetService } from "../../open-api";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { PipesModule } from "../../pipes/pipes.module";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

interface SpendingRow {
  label: string;
  value: number;
  color: string;
  percent: number;
}

@Component({
  selector: "app-spending-table",
  templateUrl: "./spending-table.component.html",
  styleUrls: ["./spending-table.component.scss"],
  standalone: true,
  imports: [CommonModule, SharedUiModule, PipesModule],
  providers: [CurrencyPipe, CustomCurrencyPipe],
})
export class SpendingTableComponent implements OnInit, OnChanges {
  public readonly widget = input.required<Widget>();
  public readonly groupId = input<number>();
  public readonly startDate = input<string>("");
  public readonly endDate = input<string>("");

  public readonly rows = signal<SpendingRow[]>([]);
  public readonly isLoading = signal(true);

  public readonly total = computed(() =>
    this.rows().reduce((sum, row) => sum + row.value, 0)
  );

  constructor(private widgetService: WidgetService) {}

  public ngOnInit(): void {
    this.loadData();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (
      (changes["groupId"] && !changes["groupId"].firstChange) ||
      (changes["startDate"] && !changes["startDate"].firstChange) ||
      (changes["endDate"] && !changes["endDate"].firstChange)
    ) {
      this.loadData();
    }
  }

  public getGroupingLabel(): string {
    const config = this.widget()?.configuration as { chartGrouping?: ChartGrouping };
    switch (config?.chartGrouping) {
      case ChartGrouping.Categories:
        return "Categories";
      case ChartGrouping.Tags:
        return "Tags";
      case ChartGrouping.Paidby:
        return "Paid By";
      default:
        return "";
    }
  }

  private loadData(): void {
    const groupId = this.groupId();
    const widget = this.widget();
    if (!groupId || !widget?.configuration) {
      this.isLoading.set(false);
      return;
    }

    const config = widget.configuration as { chartGrouping?: ChartGrouping; filter?: any };
    if (!config.chartGrouping) {
      this.isLoading.set(false);
      return;
    }

    const command: PieChartDataCommand = {
      chartGrouping: config.chartGrouping,
      filter: config.filter,
      startDate: this.startDate() || undefined,
      endDate: this.endDate() || undefined,
    };

    this.isLoading.set(true);
    this.widgetService
      .getPieChartData(groupId, command)
      .pipe(
        take(1),
        tap((response: PieChartData) => {
          this.rows.set(this.toRows(response.data ?? []));
          this.isLoading.set(false);
        })
      )
      .subscribe();
  }

  private toRows(points: PieChartDataPoint[]): SpendingRow[] {
    const totalAbs = points.reduce((sum, p) => sum + Math.abs(p.value ?? 0), 0);
    return points
      .map((p) => ({
        label: p.label || "Unknown",
        value: p.value ?? 0,
        color: p.color || "",
        percent: totalAbs > 0 ? (Math.abs(p.value ?? 0) / totalAbs) * 100 : 0,
      }))
      .sort((a, b) => Math.abs(b.value) - Math.abs(a.value));
  }
}
