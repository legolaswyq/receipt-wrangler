import { CurrencyPipe } from "@angular/common";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, SimpleChange } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { NgxsModule } from "@ngxs/store";
import { of } from "rxjs";
import { BudgetData, BudgetService, Widget, WidgetType } from "../../open-api";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { SystemSettingsState } from "../../store/system-settings.state";

import { BudgetComponent } from "./budget.component";

describe("BudgetComponent", () => {
  let component: BudgetComponent;
  let fixture: ComponentFixture<BudgetComponent>;
  let budgetService: jest.Mocked<BudgetService>;

  const mockWidget: Widget = {
    id: 1,
    name: "My Budget",
    widgetType: WidgetType.Budget,
    configuration: {},
  };

  const mockBudgetData: BudgetData = {
    month: "2026-09",
    income: 5000,
    spent: 390,
    net: 4610,
    untracked: 90,
    categories: [
      { categoryId: 5, name: "Dining Out", target: 400, spent: 300, over: false },
    ],
  };

  beforeEach(async () => {
    const budgetServiceMock = {
      getBudgetData: jest.fn(),
      upsertBudget: jest.fn(),
      deleteBudget: jest.fn(),
    };

    await TestBed.configureTestingModule({
      imports: [BudgetComponent, SharedUiModule, NgxsModule.forRoot([SystemSettingsState])],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        { provide: BudgetService, useValue: budgetServiceMock },
        CurrencyPipe,
        CustomCurrencyPipe,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    }).compileComponents();

    budgetService = TestBed.inject(BudgetService) as jest.Mocked<BudgetService>;
    budgetService.getBudgetData.mockReturnValue(of(mockBudgetData));
    budgetService.upsertBudget.mockReturnValue(of({} as any));

    fixture = TestBed.createComponent(BudgetComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput("widget", mockWidget);
    fixture.componentRef.setInput("groupId", 1);
  });

  it("should create", () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it("loads budget data on init", () => {
    fixture.detectChanges();
    expect(budgetService.getBudgetData).toHaveBeenCalledWith(1);
    expect(component.data()?.categories.length).toBe(1);
    expect(component.data()?.net).toBe(4610);
    expect(component.isLoading()).toBe(false);
  });

  it("does not call the service when groupId is unset", () => {
    fixture.componentRef.setInput("groupId", undefined);
    fixture.detectChanges();
    expect(budgetService.getBudgetData).not.toHaveBeenCalled();
    expect(component.isLoading()).toBe(false);
  });

  it("reloads when groupId changes", () => {
    fixture.detectChanges();
    budgetService.getBudgetData.mockClear();
    component.ngOnChanges({ groupId: new SimpleChange(1, 2, false) });
    expect(budgetService.getBudgetData).toHaveBeenCalled();
  });

  describe("progressPercent", () => {
    it("caps at 100 when over target", () => {
      expect(component.progressPercent({ target: 400, spent: 500 } as any)).toBe(100);
    });

    it("computes a fractional percentage", () => {
      expect(component.progressPercent({ target: 400, spent: 300 } as any)).toBe(75);
    });

    it("returns 0 when target is missing or zero", () => {
      expect(component.progressPercent({ target: 0, spent: 100 } as any)).toBe(0);
      expect(component.progressPercent({ spent: 100 } as any)).toBe(0);
    });
  });

  describe("isOver", () => {
    it("flags over budget", () => {
      expect(component.isOver({ target: 400, spent: 500 } as any)).toBe(true);
    });

    it("is false at or under target", () => {
      expect(component.isOver({ target: 400, spent: 400 } as any)).toBe(false);
      expect(component.isOver({ target: 400, spent: 100 } as any)).toBe(false);
    });
  });

  describe("saveTarget", () => {
    it("upserts a valid target and reloads", () => {
      fixture.detectChanges();
      budgetService.getBudgetData.mockClear();
      component.saveTarget({ categoryId: 5 } as any, 450);
      expect(budgetService.upsertBudget).toHaveBeenCalledWith(1, {
        categoryId: 5,
        amount: "450",
      });
      expect(budgetService.getBudgetData).toHaveBeenCalled();
      expect(component.editingCategoryId()).toBeNull();
    });

    it("ignores an invalid amount", () => {
      component.saveTarget({ categoryId: 5 } as any, 0);
      expect(budgetService.upsertBudget).not.toHaveBeenCalled();
      expect(component.editingCategoryId()).toBeNull();
    });
  });
});
