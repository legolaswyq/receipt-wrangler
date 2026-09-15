import { CommonModule, CurrencyPipe } from "@angular/common";
import { Component, OnChanges, OnInit, SimpleChanges, input, signal } from "@angular/core";
import { take, tap } from "rxjs";
import { BudgetCategory, BudgetData, BudgetService, Widget } from "../../open-api";
import { CustomCurrencyPipe } from "../../pipes/custom-currency.pipe";
import { PipesModule } from "../../pipes/pipes.module";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";

@Component({
  selector: "app-budget",
  templateUrl: "./budget.component.html",
  styleUrls: ["./budget.component.scss"],
  standalone: true,
  imports: [CommonModule, SharedUiModule, PipesModule],
  providers: [CurrencyPipe, CustomCurrencyPipe],
})
export class BudgetComponent implements OnInit, OnChanges {
  public readonly widget = input.required<Widget>();
  public readonly groupId = input<number>();

  public readonly data = signal<BudgetData | undefined>(undefined);
  public readonly isLoading = signal(true);
  public readonly editingCategoryId = signal<number | null>(null);

  constructor(private budgetService: BudgetService) {}

  public ngOnInit(): void {
    this.loadData();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["groupId"] && !changes["groupId"].firstChange) {
      this.loadData();
    }
  }

  public progressPercent(category: BudgetCategory): number {
    const target = category.target ?? 0;
    if (target <= 0) {
      return 0;
    }
    return Math.min(100, Math.round(((category.spent ?? 0) / target) * 100));
  }

  public isOver(category: BudgetCategory): boolean {
    const target = category.target ?? 0;
    return target > 0 && (category.spent ?? 0) > target;
  }

  public startEdit(category: BudgetCategory): void {
    this.editingCategoryId.set(category.categoryId ?? null);
  }

  public saveTarget(category: BudgetCategory, amount: number): void {
    const groupId = this.groupId();
    if (!groupId || category.categoryId == null || isNaN(amount) || amount <= 0) {
      this.editingCategoryId.set(null);
      return;
    }

    this.budgetService
      .upsertBudget(groupId, { categoryId: category.categoryId, amount: String(amount) })
      .pipe(
        take(1),
        tap(() => {
          this.editingCategoryId.set(null);
          this.loadData();
        }),
      )
      .subscribe();
  }

  private loadData(): void {
    const groupId = this.groupId();
    if (!groupId) {
      this.isLoading.set(false);
      return;
    }

    this.isLoading.set(true);
    this.budgetService
      .getBudgetData(groupId)
      .pipe(
        take(1),
        tap((response: BudgetData) => {
          this.data.set(response);
          this.isLoading.set(false);
        }),
      )
      .subscribe();
  }
}
