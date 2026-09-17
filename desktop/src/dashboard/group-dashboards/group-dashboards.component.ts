import { Component, computed, effect } from "@angular/core";
import { Router } from "@angular/router";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { DashboardState } from "src/store/dashboard.state";
import { AddDashboardToGroup } from "src/store/dashboard.state.actions";
import {
  ChartGrouping,
  Dashboard,
  DashboardService,
  Permission,
  UpsertDashboardCommand,
  WidgetType,
} from "../../open-api";
import { AuthState, GroupState, SetSelectedDashboardId } from "../../store";
import { DEFAULT_DASHBOARD_PERIOD } from "../utils/dashboard-period.util";

// Stable empty reference so the dashboards computed doesn't emit a fresh [] on
// every recompute (which would needlessly re-run the auto-select effect).
const EMPTY_DASHBOARDS: Dashboard[] = [];

@UntilDestroy()
@Component({
    selector: "app-group-dashboards",
    templateUrl: "./group-dashboards.component.html",
    styleUrls: ["./group-dashboards.component.scss"],
    standalone: false
})
export class GroupDashboardsComponent {
  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  protected readonly selectedGroupIdNum = computed(
    () => +(this.selectedGroupId() ?? 0)
  );

  public selectedDashboardId = this.store.selectSignal(
    GroupState.selectedDashboardId
  );

  private dashboardsByGroup = this.store.selectSignal(DashboardState.dashboards);

  public dashboards = computed<Dashboard[]>(
    () => this.dashboardsByGroup()[this.selectedGroupId()] ?? EMPTY_DASHBOARDS
  );

  // Guards the auto-create so the effect firing several times before the store
  // updates doesn't POST multiple default dashboards.
  private creatingDefault = false;

  constructor(
    private dashboardService: DashboardService,
    private router: Router,
    private store: Store
  ) {
    // The route resolver has already loaded this group's dashboards, so once
    // they're available: select the single dashboard and show its outlet, or —
    // if the group has none — create the default one. Depends only on
    // dashboards(); selectedDashboardId is read untracked.
    effect(() => {
      const dashboards = this.dashboards();
      const selectedDashboardId = this.store.selectSnapshot(
        GroupState.selectedDashboardId
      );

      if (selectedDashboardId) {
        this.navigateToDashboard(+selectedDashboardId);
      } else if (dashboards.length > 0) {
        this.setSelectedDashboardId(dashboards[0].id);
        this.navigateToDashboard(dashboards[0].id);
      } else {
        this.ensureDefaultDashboard();
      }
    });
  }

  public navigateToDashboard(dashboardId: number): void {
    const selectedGroupId = this.store.selectSnapshot(
      GroupState.selectedGroupId
    );

    setTimeout(() => {
      this.router.navigateByUrl(
        `/dashboard/group/${selectedGroupId}/${dashboardId}`
      );
    }, 0);
  }

  public setSelectedDashboardId(dashboardId: number): void {
    this.store.dispatch(new SetSelectedDashboardId(dashboardId?.toString()));
  }

  // Creates the single default dashboard for a group that has none, so every
  // group shows a working dashboard without any setup. Only runs for callers who
  // can create dashboards (otherwise the API would 403).
  private ensureDefaultDashboard(): void {
    const groupId = this.store.selectSnapshot(GroupState.selectedGroupId);
    if (
      this.creatingDefault ||
      !groupId ||
      !this.store.selectSnapshot(
        AuthState.hasGroupPermission(+groupId, Permission.GroupDashboardsCreate)
      )
    ) {
      return;
    }

    this.creatingDefault = true;

    const command: UpsertDashboardCommand = {
      name: "home",
      groupId: groupId,
      period: DEFAULT_DASHBOARD_PERIOD,
      widgets: [
        {
          name: "Spending by Category",
          widgetType: WidgetType.PieChart,
          configuration: { chartGrouping: ChartGrouping.Categories },
        },
        {
          name: "Category Breakdown",
          widgetType: WidgetType.SpendingTable,
          configuration: { chartGrouping: ChartGrouping.Categories },
        },
        {
          name: "Monthly Budget",
          widgetType: WidgetType.Budget,
          configuration: {},
        },
      ],
    };

    this.dashboardService
      .createDashboard(command)
      .pipe(
        take(1),
        tap((dashboard) => {
          this.store.dispatch(new AddDashboardToGroup(groupId, dashboard));
          this.creatingDefault = false;
        })
      )
      .subscribe({
        error: () => {
          this.creatingDefault = false;
        },
      });
  }
}
