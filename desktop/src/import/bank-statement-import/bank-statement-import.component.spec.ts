import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { CurrencyPipe } from "@angular/common";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { BankStatementPreview, BankStatementTransaction, ImportService } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services";
import { AuthState, GroupState } from "../../store";
import { SystemSettingsState } from "../../store/system-settings.state";
import { BankStatementImportComponent } from "./bank-statement-import.component";

const transaction = (overrides: Partial<BankStatementTransaction>): BankStatementTransaction => ({
  date: "2026-01-20",
  type: "Visa Purchase",
  details: "Corner Groc",
  reference: "",
  name: "Corner Groc",
  amount: "42.19",
  balance: "100.00",
  isIncome: false,
  selected: true,
  duplicate: false,
  ...overrides,
});

const buildPreview = (): BankStatementPreview => ({
  bank: "ANZ",
  accountName: "Everyday Account",
  accountNumber: "01-0001-0000001-00",
  periodStart: "2026-01-01",
  periodEnd: "2026-01-31",
  transactions: [
    transaction({ suggestedCategoryId: 10 }),
    transaction({ name: "Employer Limited", amount: "1500.00", isIncome: true, type: "Direct Credit" }),
    transaction({ name: "Transfer From: 06", amount: "1000.00", isIncome: true, type: "Transfer", selected: false }),
    transaction({ name: "Acme Trading", amount: "40.00", duplicate: true, selected: false }),
  ],
});

describe("BankStatementImportComponent", () => {
  let component: BankStatementImportComponent;
  let fixture: ComponentFixture<BankStatementImportComponent>;
  let importService: { previewBankStatement: jest.Mock; importBankStatement: jest.Mock };
  let snackbarService: { error: jest.Mock; success: jest.Mock };
  let dialogRef: { close: jest.Mock };
  let store: Store;

  const queryAll = (id: string): HTMLElement[] =>
    Array.from(fixture.nativeElement.querySelectorAll(`[data-testid="${id}"]`));

  const loadFile = async () => {
    component.fileLoaded({ file: new File(["%PDF"], "statement.pdf", { type: "application/pdf" }) } as any);
    await fixture.whenStable();
  };

  beforeEach(async () => {
    importService = {
      previewBankStatement: jest.fn(() => of(buildPreview())),
      importBankStatement: jest.fn(() => of({ createdCount: 2 })),
    };
    snackbarService = { error: jest.fn(), success: jest.fn() };
    dialogRef = { close: jest.fn() };

    await TestBed.configureTestingModule({
      declarations: [BankStatementImportComponent],
      imports: [ReactiveFormsModule, PipesModule, NgxsModule.forRoot([AuthState, GroupState, SystemSettingsState])],
      providers: [
        provideZonelessChangeDetection(),
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
        CurrencyPipe,
        { provide: ImportService, useValue: importService },
        { provide: SnackbarService, useValue: snackbarService },
        { provide: MatDialogRef, useValue: dialogRef },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();

    store = TestBed.inject(Store);
    store.reset({
      ...store.snapshot(),
      auth: {
        userId: "7",
        groupCategories: {
          1: [
            { id: 10, name: "Groceries", isIncome: false },
            { id: 20, name: "Salary", isIncome: true },
          ],
          3: [{ id: 30, name: "Other", isIncome: false }],
        },
      },
      groups: {
        groups: [
          { id: 2, name: "All", isAllGroup: true, groupMembers: [] },
          { id: 1, name: "My Receipts", isAllGroup: false, groupMembers: [] },
          { id: 3, name: "Shared", isAllGroup: false, groupMembers: [] },
        ],
        selectedGroupId: "2",
      },
    });

    fixture = TestBed.createComponent(BankStatementImportComponent);
    component = fixture.componentInstance;
    fixture.autoDetectChanges();
    await fixture.whenStable();
  });

  it("defaults to the first real group when the All group is selected, and to the current user", () => {
    expect(component.form.value).toEqual({ groupId: 1, paidByUserId: 7 });
    expect(queryAll("bank-statement-empty").length).toBe(1);
  });

  it("defaults to the selected group when it is a real group", () => {
    store.reset({ ...store.snapshot(), groups: { ...store.snapshot().groups, selectedGroupId: "3" } });

    component.ngOnInit();

    expect(component.form.value.groupId).toBe(3);
  });

  it("previews the uploaded statement and applies the server's defaults", async () => {
    await loadFile();

    expect(importService.previewBankStatement).toHaveBeenCalledWith(expect.any(File), 1);
    expect(queryAll("bank-statement-row").length).toBe(4);
    expect(queryAll("bank-statement-duplicate").length).toBe(1);
    expect(component.rows.getRawValue()).toEqual([
      { selected: true, categoryId: 10 },
      { selected: true, categoryId: null },
      { selected: false, categoryId: null },
      { selected: false, categoryId: null },
    ]);
    expect(component.expenseCategories().map((c) => c.id)).toEqual([10]);
    expect(component.incomeCategories().map((c) => c.id)).toEqual([20]);
    expect(component.summary()).toEqual({
      selectedCount: 2,
      spending: 42.19,
      income: 1500,
      missingIncomeCategory: 1,
    });
  });

  it("refuses to import a selected deposit without an income category", async () => {
    await loadFile();

    expect(queryAll("bank-statement-income-required").length).toBe(1);

    component.submit();

    expect(importService.importBankStatement).not.toHaveBeenCalled();
    expect(snackbarService.error).toHaveBeenCalledWith("Choose an income category for every selected deposit");
  });

  it("refuses to import when nothing is selected", async () => {
    await loadFile();
    component.setAllSelected(false);

    component.submit();

    expect(component.summary().selectedCount).toBe(0);
    expect(importService.importBankStatement).not.toHaveBeenCalled();
    expect(snackbarService.error).toHaveBeenCalledWith("Select at least one transaction to import");
  });

  it("imports only the selected rows with their categories and closes with true", async () => {
    await loadFile();
    component.rows.at(1).controls.categoryId.setValue(20);
    await fixture.whenStable();

    expect(queryAll("bank-statement-income-required").length).toBe(0);

    component.submit();

    expect(importService.importBankStatement).toHaveBeenCalledWith({
      groupId: 1,
      paidByUserId: 7,
      transactions: [
        { date: "2026-01-20", name: "Corner Groc", amount: "42.19", isIncome: false, categoryId: 10 },
        { date: "2026-01-20", name: "Employer Limited", amount: "1500.00", isIncome: true, categoryId: 20 },
      ],
    });
    expect(snackbarService.success).toHaveBeenCalledWith("Imported 2 transactions");
    expect(dialogRef.close).toHaveBeenCalledWith(true);
  });

  it("select all includes transfers and duplicates", async () => {
    await loadFile();

    component.setAllSelected(true);

    expect(component.summary().selectedCount).toBe(4);
    expect(component.summary().missingIncomeCategory).toBe(2);
  });

  it("reads the statement again when the group changes", async () => {
    await loadFile();

    component.form.get("groupId")?.setValue(3);
    await fixture.whenStable();

    expect(importService.previewBankStatement).toHaveBeenLastCalledWith(expect.any(File), 3);
    expect(component.expenseCategories().map((c) => c.id)).toEqual([30]);
    expect(component.incomeCategories()).toEqual([]);
  });

  it("does nothing on submit before a statement is loaded, and cancel closes with false", () => {
    component.submit();
    component.closeDialog();

    expect(importService.importBankStatement).not.toHaveBeenCalled();
    expect(dialogRef.close).toHaveBeenCalledWith(false);
  });
});
