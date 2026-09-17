import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed, } from "@angular/core/testing";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute, Params, Router } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { BehaviorSubject, of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { DashboardState } from "src/store/dashboard.state";
import { AddDashboardToGroup } from "src/store/dashboard.state.actions";
import { ButtonModule } from "../../button";
import { Dashboard, DashboardService } from "../../open-api";
import { GroupState, SetSelectedDashboardId } from "../../store";
import { AuthState } from "../../store/auth.state";
import { SetPermissions } from "../../store/auth.state.actions";
import { DirectivesModule } from "../../directives/directives.module";
import { GroupDashboardsComponent } from "./group-dashboards.component";

describe("GroupDashboardsComponent", () => {
  let component: GroupDashboardsComponent;
  let fixture: ComponentFixture<GroupDashboardsComponent>;
  let store: Store;

  const dashboard = (id: number, groupId = 1): Dashboard => ({
    id,
    name: `dashboard-${id}`,
    groupId,
    userId: 1,
  });

  // navigateToDashboard defers the actual navigation via setTimeout(0); spying
  // navigateByUrl keeps a real (routeless) Router from producing rejected
  // navigations, and lets the navigation tests assert the target url.
  const spyNavigate = () =>
    jest
      .spyOn(TestBed.inject(Router), "navigateByUrl")
      .mockResolvedValue(true as any);

  const seed = (
    groups: Partial<{ selectedGroupId: string; selectedDashboardId: string }>,
    dashboards: { [groupId: string]: Dashboard[] } = {}
  ) => {
    store.reset({
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
        ...groups,
      },
      dashboards: { dashboards },
    });
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [GroupDashboardsComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [PipesModule,
        MatDialogModule,
        NgxsModule.forRoot([GroupState, DashboardState, AuthState]),
        PipesModule,
        ButtonModule,
        DirectivesModule,
        MatSnackBarModule],
      providers: [
        DashboardService,
        {
          provide: ActivatedRoute,
          useValue: {
            params: new BehaviorSubject<Params>({}),
            snapshot: {
              data: {
                dashboards: [],
              },
            },
          },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ]
    });
    store = TestBed.inject(Store);
    store.reset({
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
        selectedGroupId: "1",
        groups: [],
      },
    });
    fixture = TestBed.createComponent(GroupDashboardsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("derives the current group's dashboards from the store", () => {
    spyNavigate();
    const dashboards = [dashboard(1)];
    seed({ selectedGroupId: "1", selectedDashboardId: "" }, { "1": dashboards });

    TestBed.flushEffects();

    expect(component.dashboards()).toEqual(dashboards);
  });

  it("reacts to the selected group id changing", () => {
    spyNavigate();
    const g1 = [dashboard(1, 1)];
    const g2 = [dashboard(2, 2)];
    seed({ selectedGroupId: "1", selectedDashboardId: "" }, { "1": g1, "2": g2 });
    TestBed.flushEffects();

    expect(component.dashboards()).toEqual(g1);

    store.reset({
      ...store.snapshot(),
      groups: {
        ...store.snapshot().groups,
        selectedGroupId: "2",
        selectedDashboardId: "",
      },
    });
    TestBed.flushEffects();

    expect(component.dashboards()).toEqual(g2);
  });

  it("auto-selects and navigates to the first dashboard once a cold group's data lands", () => {
    jest.useFakeTimers();
    const navSpy = spyNavigate();
    seed({ selectedGroupId: "1", selectedDashboardId: "" }, { "1": [dashboard(7)] });

    TestBed.flushEffects();

    expect(store.selectSnapshot(GroupState.selectedDashboardId)).toBe("7");

    jest.runOnlyPendingTimers();
    expect(navSpy).toHaveBeenCalledWith("/dashboard/group/1/7");
    jest.useRealTimers();
  });

  it("navigates to the already-selected dashboard when data lands (warm/deep-link)", () => {
    jest.useFakeTimers();
    const navSpy = spyNavigate();
    seed(
      { selectedGroupId: "1", selectedDashboardId: "5" },
      { "1": [dashboard(5), dashboard(8)] }
    );

    TestBed.flushEffects();
    jest.runOnlyPendingTimers();

    expect(navSpy).toHaveBeenCalledWith("/dashboard/group/1/5");
    expect(store.selectSnapshot(GroupState.selectedDashboardId)).toBe("5");
    jest.useRealTimers();
  });

  it("creates a default dashboard when the group has none and the user can create", () => {
    spyNavigate();
    const created = dashboard(9);
    const createSpy = jest
      .spyOn(TestBed.inject(DashboardService), "createDashboard")
      .mockReturnValue(of(created) as any);

    store.dispatch(new SetPermissions([], { 1: ["group.dashboards.create"] }));
    seed({ selectedGroupId: "1", selectedDashboardId: "" }, { "1": [] });

    // A fresh component runs its constructor effect against the seeded state
    // (an empty->empty dashboards signal wouldn't re-fire the shared one).
    const dispatchSpy = jest.spyOn(store, "dispatch");
    const freshFixture = TestBed.createComponent(GroupDashboardsComponent);
    freshFixture.detectChanges();
    TestBed.flushEffects();

    expect(createSpy).toHaveBeenCalledTimes(1);
    const command = createSpy.mock.calls[0][0] as any;
    expect(command.widgets.map((w: any) => w.widgetType)).toEqual([
      "PIE_CHART",
      "SPENDING_TABLE",
      "BUDGET",
    ]);
    expect(dispatchSpy).toHaveBeenCalledWith(
      new AddDashboardToGroup("1", created)
    );
  });

  it("does not create a default dashboard without create permission", () => {
    spyNavigate();
    const createSpy = jest.spyOn(
      TestBed.inject(DashboardService),
      "createDashboard"
    );

    store.dispatch(new SetPermissions([], { 1: ["group.dashboards.read"] }));
    seed({ selectedGroupId: "1", selectedDashboardId: "" }, { "1": [] });

    const freshFixture = TestBed.createComponent(GroupDashboardsComponent);
    freshFixture.detectChanges();
    TestBed.flushEffects();

    expect(createSpy).not.toHaveBeenCalled();
  });

  it("should set selected dashboard id", () => {
    const store = TestBed.inject(Store);
    const storeSpy = jest.spyOn(store, "dispatch");

    component.setSelectedDashboardId(1);

    expect(storeSpy).toHaveBeenCalledWith(new SetSelectedDashboardId("1"));
  });
});
