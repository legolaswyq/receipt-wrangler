import { Component, signal, ViewEncapsulation } from "@angular/core";
import { FormControl } from "@angular/forms";
import { Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { UntilDestroy, untilDestroyed } from "@ngneat/until-destroy";
import { of, switchMap, take, tap } from "rxjs";
import { SearchResult, SearchService } from "../../open-api";
import { GroupState } from "../../store";

@UntilDestroy()
@Component({
    selector: "app-searchbar",
    templateUrl: "./searchbar.component.html",
    styleUrls: ["./searchbar.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class SearchbarComponent {
  public results = signal<SearchResult[]>([]);

  public searchFormControl = new FormControl("");

  public displayWith = (searchResult: SearchResult) => {
    return searchResult?.name;
  };

  constructor(
    private searchService: SearchService,
    private router: Router,
    private store: Store
  ) {}

  public ngOnInit(): void {
    this.searchFormControl.valueChanges
      .pipe(
        untilDestroyed(this),
        switchMap((value) =>
          value
            ? this.searchService.receiptSearch(value ?? "").pipe(take(1))
            : of([] as SearchResult[])
        ),
        tap((results) => {
          this.results.set(results);
        })
      )
      .subscribe();
  }

  public navigateToResult(result: SearchResult) {
    switch (result.type) {
      case "Receipt":
        this.router.navigateByUrl(`/receipts/${result.id}/view`);
        break;
    }
  }

  // Enter (rather than picking an autocomplete suggestion) opens the receipts
  // list filtered to the typed term, so the user gets a page of every matching
  // receipt instead of jumping to a single one.
  public submitSearch(event?: Event): void {
    // Stop the autocomplete/native form submit from also acting on Enter.
    event?.preventDefault();
    event?.stopPropagation();
    const value = this.searchFormControl.value as string | SearchResult | null;
    const term = (typeof value === "string" ? value : value?.name)?.trim();
    if (!term) {
      return;
    }
    const receiptListLink = this.store.selectSnapshot(
      GroupState.receiptListLink
    );
    this.router.navigate([receiptListLink], { queryParams: { search: term } });
  }
}
