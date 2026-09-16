import { Component, OnInit, computed, signal } from "@angular/core";
import { FormControl } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { DEFAULT_HOST_CLASS } from "src/constants";
import { DashboardState } from "src/store/dashboard.state";
import { UpdateDashBoardForGroup } from "src/store/dashboard.state.actions";
import { Dashboard, DashboardService, UpsertDashboardCommand, WidgetType } from "../../open-api";
import { GroupState } from "../../store";
import {
  DEFAULT_DASHBOARD_PERIOD,
  DashboardPeriod,
  dashboardPeriodOptions,
  resolveDashboardRange,
} from "../utils/dashboard-period.util";

@UntilDestroy()
@Component({
    selector: "app-dashboard",
    templateUrl: "./dashboard.component.html",
    styleUrls: ["./dashboard.component.scss"],
    host: DEFAULT_HOST_CLASS,
    standalone: false
})
export class DashboardComponent implements OnInit {
  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  public dashboards = signal<Dashboard[]>([]);

  public selectedDashboard = signal<Dashboard | undefined>(undefined);

  public widgetType = WidgetType;

  public readonly dashboardPeriodOptions = dashboardPeriodOptions;
  public readonly customPeriod = DashboardPeriod.Custom;

  public readonly periodControl = new FormControl<string>(DEFAULT_DASHBOARD_PERIOD);
  public readonly customStartControl = new FormControl<string>("");
  public readonly customEndControl = new FormControl<string>("");

  public readonly selectedPeriod = signal<string>(DEFAULT_DASHBOARD_PERIOD);
  private readonly customStart = signal<string>("");
  private readonly customEnd = signal<string>("");

  // The resolved [startDate, endDate) window passed to spend-over-time widgets.
  public readonly range = computed(() =>
    resolveDashboardRange(this.selectedPeriod(), this.customStart(), this.customEnd())
  );

  constructor(
    private activatedRoute: ActivatedRoute,
    private store: Store,
    private dashboardService: DashboardService
  ) {}

  public ngOnInit(): void {
    this.listenForDashboardChanges();
    this.listenForParamChanges();
    this.listenForPeriodChanges();
  }

  public onPeriodChange(period: string): void {
    this.selectedPeriod.set(period);
    this.persistPeriod();
  }

  public onCustomDateChange(): void {
    this.customStart.set(this.customStartControl.value ?? "");
    this.customEnd.set(this.customEndControl.value ?? "");
    if (this.selectedPeriod() === DashboardPeriod.Custom) {
      this.persistPeriod();
    }
  }

  private listenForPeriodChanges(): void {
    this.periodControl.valueChanges
      .pipe(
        untilDestroyed(this),
        tap((period) => this.onPeriodChange(period ?? DEFAULT_DASHBOARD_PERIOD))
      )
      .subscribe();

    this.customStartControl.valueChanges
      .pipe(untilDestroyed(this), tap(() => this.onCustomDateChange()))
      .subscribe();
    this.customEndControl.valueChanges
      .pipe(untilDestroyed(this), tap(() => this.onCustomDateChange()))
      .subscribe();
  }

  private listenForParamChanges(): void {
    this.activatedRoute.params
      .pipe(
        tap(() => {
          this.setSelectedDashboard();
        })
      )
      .subscribe();
  }

  private listenForDashboardChanges(): void {
    this.store
      .select(DashboardState.dashboards)
      .pipe(
        untilDestroyed(this),
        tap(() => {
          const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);
          this.dashboards.set(this.store.selectSnapshot(
            DashboardState.getDashboardsByGroupId(groupId)
          ));
          this.setSelectedDashboard();
        })
      )
      .subscribe();
  }

  private setSelectedDashboard(): void {
    const selectedDashboardId = this.store.selectSnapshot(
      GroupState.selectedDashboardId
    );
    const dashboard = this.dashboards().find(
      (dashboard) => dashboard?.id?.toString() === selectedDashboardId
    );
    this.selectedDashboard.set(dashboard);
    this.syncPeriodFromDashboard(dashboard);
  }

  private syncPeriodFromDashboard(dashboard: Dashboard | undefined): void {
    const period = dashboard?.period || DEFAULT_DASHBOARD_PERIOD;
    this.selectedPeriod.set(period);
    this.customStart.set(dashboard?.periodStartDate ?? "");
    this.customEnd.set(dashboard?.periodEndDate ?? "");
    // Reflect the stored values in the controls without re-triggering a persist.
    this.periodControl.setValue(period, { emitEvent: false });
    this.customStartControl.setValue(dashboard?.periodStartDate ?? "", { emitEvent: false });
    this.customEndControl.setValue(dashboard?.periodEndDate ?? "", { emitEvent: false });
  }

  private persistPeriod(): void {
    const dashboard = this.selectedDashboard();
    if (!dashboard?.id) {
      return;
    }

    const command: UpsertDashboardCommand = {
      name: dashboard.name ?? "",
      groupId: (dashboard.groupId ?? "").toString(),
      widgets: (dashboard.widgets ?? []).map((widget) => ({
        name: widget.name,
        widgetType: widget.widgetType!,
        configuration: widget.configuration,
      })),
      period: this.selectedPeriod(),
      periodStartDate: this.customStart(),
      periodEndDate: this.customEnd(),
    };

    this.dashboardService
      .updateDashboard(dashboard.id, command)
      .pipe(
        take(1),
        tap((updated) => {
          const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);
          this.store.dispatch(
            new UpdateDashBoardForGroup(groupId, dashboard.id!, updated)
          );
        })
      )
      .subscribe();
  }
}
