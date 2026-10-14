import { CurrencyPipe } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, SimpleChange } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { ChartGrouping, PieChartData, Widget, WidgetService, WidgetType } from "../../open-api";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SystemSettingsState } from "../../store/system-settings.state";

import { SpendingTableComponent } from "./spending-table.component";

describe("SpendingTableComponent", () => {
  let component: SpendingTableComponent;
  let fixture: ComponentFixture<SpendingTableComponent>;
  let widgetService: jest.Mocked<WidgetService>;

  const mockWidget: Widget = {
    id: 1,
    name: "Spending Table",
    widgetType: WidgetType.SpendingTable,
    configuration: { chartGrouping: ChartGrouping.Categories },
  };

  const mockData: PieChartData = {
    data: [
      { label: "Groceries", value: 100, color: "#4E79A7" },
      { label: "Dining Out", value: 300, color: "#F28E2B" },
      { label: "Gas", value: 100 },
    ],
  };

  beforeEach(async () => {
    const widgetServiceMock = { getPieChartData: jest.fn() };

    await TestBed.configureTestingModule({
      imports: [SpendingTableComponent, SharedUiModule, NgxsModule.forRoot([SystemSettingsState])],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        { provide: WidgetService, useValue: widgetServiceMock },
        CurrencyPipe,
        CustomCurrencyPipe,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    widgetService = TestBed.inject(WidgetService) as jest.Mocked<WidgetService>;
    widgetService.getPieChartData.mockReturnValue(of(mockData));

    fixture = TestBed.createComponent(SpendingTableComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("widget", mockWidget);
    fixture.componentRef.setInput("groupId", 1);
  });

  it("should create", () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it("loads rows sorted by absolute value desc and computes total", () => {
    fixture.detectChanges();
    const rows = component.rows();
    expect(rows.map((r) => r.label)).toEqual(["Dining Out", "Groceries", "Gas"]);
    expect(component.total()).toBe(500);
  });

  it("computes percentages against the absolute total", () => {
    fixture.detectChanges();
    const dining = component.rows().find((r) => r.label === "Dining Out")!;
    expect(Math.round(dining.percent)).toBe(60); // 300 / 500
  });

  it("passes the date range to the service", () => {
    fixture.componentRef.setInput("startDate", "2026-09-01T00:00:00Z");
    fixture.componentRef.setInput("endDate", "2026-10-01T00:00:00Z");
    fixture.detectChanges();
    expect(widgetService.getPieChartData).toHaveBeenCalledWith(1, {
      chartGrouping: ChartGrouping.Categories,
      filter: undefined,
      startDate: "2026-09-01T00:00:00Z",
      endDate: "2026-10-01T00:00:00Z",
    });
  });

  it("asks for income and words the empty state accordingly for an income widget", () => {
    widgetService.getPieChartData.mockReturnValue(of({ data: [] }));
    fixture.componentRef.setInput("widget", {
      ...mockWidget,
      configuration: { chartGrouping: ChartGrouping.Categories, incomeOnly: true },
    });
    fixture.detectChanges();

    expect(widgetService.getPieChartData.mock.calls[0][1].incomeOnly).toBe(true);
    expect(fixture.nativeElement.textContent).toContain("No income in this range.");
  });

  it("does not ask for income by default", () => {
    fixture.detectChanges();
    expect(widgetService.getPieChartData.mock.calls[0][1].incomeOnly).toBeUndefined();
  });

  it("reloads when the date range changes", () => {
    fixture.detectChanges();
    widgetService.getPieChartData.mockClear();
    component.ngOnChanges({ startDate: new SimpleChange("", "2026-09-01T00:00:00Z", false) });
    expect(widgetService.getPieChartData).toHaveBeenCalled();
  });
});
