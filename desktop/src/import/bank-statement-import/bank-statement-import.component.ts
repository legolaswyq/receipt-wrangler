import { Component, DestroyRef, OnInit, signal, viewChild } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormArray, FormBuilder, FormControl, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { finalize, take, tap } from "rxjs";
import {
  BankStatementImportTransaction,
  BankStatementPreview,
  BankStatementTransaction,
  Category,
  ImportService,
} from "../../open-api";
import { UploadImageComponent } from "../../receipts/upload-image/upload-image.component";
import { SnackbarService } from "../../services";
import { AuthState, GroupState } from "../../store";

interface StatementRowForm {
  selected: FormControl<boolean>;
  categoryId: FormControl<number | null>;
}

export interface StatementSummary {
  selectedCount: number;
  spending: number;
  income: number;
  missingIncomeCategory: number;
}

const EMPTY_SUMMARY: StatementSummary = {
  selectedCount: 0,
  spending: 0,
  income: 0,
  missingIncomeCategory: 0,
};

@Component({
    selector: "app-bank-statement-import",
    templateUrl: "./bank-statement-import.component.html",
    styleUrl: "./bank-statement-import.component.scss",
    standalone: false
})
export class BankStatementImportComponent implements OnInit {
  public readonly uploadImageComponent = viewChild.required(UploadImageComponent);

  public readonly acceptedFileTypes = ["application/pdf"];

  public form: FormGroup = new FormGroup({});

  // One entry per statement transaction, index-aligned with preview().transactions.
  public rows = new FormArray<FormGroup<StatementRowForm>>([]);

  public readonly preview = signal<BankStatementPreview | null>(null);

  public readonly loading = signal(false);

  public readonly summary = signal<StatementSummary>(EMPTY_SUMMARY);

  public readonly expenseCategories = signal<Category[]>([]);

  public readonly incomeCategories = signal<Category[]>([]);

  // The picked statements, one per account. Read together so a transfer between
  // two of the accounts can be recognised from its two sides.
  private files: File[] = [];

  constructor(
    private destroyRef: DestroyRef,
    private dialogRef: MatDialogRef<BankStatementImportComponent>,
    private formBuilder: FormBuilder,
    private importService: ImportService,
    private snackbarService: SnackbarService,
    private store: Store,
  ) {}

  public ngOnInit(): void {
    this.form = this.formBuilder.group({
      groupId: [this.defaultGroupId(), Validators.required],
      paidByUserId: [Number(this.store.selectSnapshot(AuthState.userId)) || null, Validators.required],
    });

    // Duplicates, suggestions and the category list are all per group, so a
    // statement that is already loaded has to be read again for the new group.
    this.form.get("groupId")?.valueChanges
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => this.loadPreview());
  }

  public openFileUploadDialog(): void {
    this.uploadImageComponent().clickInput();
  }

  public filesSelected(files: File[]): void {
    this.files = files;
    this.loadPreview();
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public setAllSelected(selected: boolean): void {
    this.rows.controls.forEach((row) => row.controls.selected.setValue(selected, { emitEvent: false }));
    this.updateSummary();
  }

  public needsIncomeCategory(index: number): boolean {
    const row = this.rows.at(index).getRawValue();
    return row.selected && !!this.preview()?.transactions[index].isIncome && !row.categoryId;
  }

  public submit(): void {
    const preview = this.preview();
    if (!preview || this.loading()) {
      return;
    }

    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const summary = this.summary();
    if (summary.selectedCount === 0) {
      this.snackbarService.error("Select at least one transaction to import");
      return;
    }

    if (summary.missingIncomeCategory > 0) {
      this.snackbarService.error("Choose an income category for every selected deposit");
      return;
    }

    const transactions: BankStatementImportTransaction[] = [];
    preview.transactions.forEach((transaction, index) => {
      const row = this.rows.at(index).getRawValue();
      if (row.selected) {
        transactions.push({
          date: transaction.date,
          name: transaction.name,
          amount: transaction.amount,
          isIncome: transaction.isIncome,
          categoryId: row.categoryId ?? undefined,
        });
      }
    });

    this.loading.set(true);
    this.importService
      .importBankStatement({
        groupId: Number(this.form.value.groupId),
        paidByUserId: Number(this.form.value.paidByUserId),
        transactions,
      })
      .pipe(
        take(1),
        tap((result) => {
          this.snackbarService.success(`Imported ${result.createdCount} transactions`);
          this.dialogRef.close(true);
        }),
        finalize(() => this.loading.set(false)),
      )
      .subscribe();
  }

  private defaultGroupId(): number | null {
    const groups = this.store.selectSnapshot(GroupState.groupsWithoutAll) ?? [];
    const selectedGroupId = this.store.selectSnapshot(GroupState.selectedGroupId);
    const selected = groups.find((group) => group.id.toString() === selectedGroupId?.toString());

    return (selected ?? groups[0])?.id ?? null;
  }

  private loadPreview(): void {
    // Read the control, not form.value: this runs from the control's own
    // valueChanges, which fires before the parent form's value is updated.
    const groupId = Number(this.form.get("groupId")?.value);
    this.clearPreview();
    if (this.files.length === 0 || !groupId) {
      return;
    }

    this.loading.set(true);
    this.importService
      .previewBankStatement(this.files, groupId)
      .pipe(
        take(1),
        tap((preview) => this.setPreview(preview, groupId)),
        finalize(() => this.loading.set(false)),
      )
      .subscribe();
  }

  private clearPreview(): void {
    this.preview.set(null);
    this.rows = new FormArray<FormGroup<StatementRowForm>>([]);
    this.summary.set(EMPTY_SUMMARY);
  }

  private setPreview(preview: BankStatementPreview, groupId: number): void {
    const categories = this.store.selectSnapshot(AuthState.groupCategories(groupId));
    this.expenseCategories.set(categories.filter((category) => !category.isIncome));
    this.incomeCategories.set(categories.filter((category) => category.isIncome));

    this.rows = new FormArray(preview.transactions.map((transaction) => this.buildRow(transaction)));
    this.rows.valueChanges
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => this.updateSummary());

    this.preview.set(preview);
    this.updateSummary();
  }

  private buildRow(transaction: BankStatementTransaction): FormGroup<StatementRowForm> {
    return new FormGroup<StatementRowForm>({
      selected: new FormControl(transaction.selected, { nonNullable: true }),
      categoryId: new FormControl(transaction.suggestedCategoryId ?? null),
    });
  }

  private updateSummary(): void {
    const summary: StatementSummary = { ...EMPTY_SUMMARY };

    this.preview()?.transactions.forEach((transaction, index) => {
      const row = this.rows.at(index).getRawValue();
      if (!row.selected) {
        return;
      }

      summary.selectedCount++;
      if (transaction.isIncome) {
        summary.income += Number(transaction.amount);
        if (!row.categoryId) {
          summary.missingIncomeCategory++;
        }
      } else {
        summary.spending += Number(transaction.amount);
      }
    });

    this.summary.set(summary);
  }
}
