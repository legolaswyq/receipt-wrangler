import { Component, computed, signal } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { Router, RouteReuseStrategy } from "@angular/router";
import { Store } from "@ngxs/store";
import { switchMap, take } from "rxjs";
import { LayoutState } from "src/store/layout.state";
import { AboutComponent } from "../../about/about/about.component";
import { DEFAULT_DIALOG_CONFIG } from "../../constants";
import { BankStatementImportComponent } from "../../import/bank-statement-import/bank-statement-import.component";
import { ImportFormComponent } from "../../import/import-form/import-form.component";
import { AuthService, Permission } from "../../open-api";
import { AuthState, GroupState, Logout } from "../../store";
import { FeatureConfigState } from "../../store/feature-config.state";
import { hasAll } from "../../utils/permission.utils";

// The review list has a row per transaction, so it gets more room than the
// default half-width dialog.
const BANK_STATEMENT_DIALOG_CONFIG = {
  ...DEFAULT_DIALOG_CONFIG,
  width: "70%",
  maxWidth: "960px",
};

@Component({
    selector: "app-header",
    templateUrl: "./header.component.html",
    styleUrls: ["./header.component.scss"],
    standalone: false
})
export class HeaderComponent {
  public isLoggedIn = this.store.selectSignal(AuthState.isLoggedIn);

  public selectedGroupId = this.store.selectSignal(GroupState.selectedGroupId);

  public selectedGroupIdNumber = computed(() => {
    const id = Number.parseInt(this.selectedGroupId() ?? "");
    return Number.isNaN(id) ? 0 : id;
  });

  public showProgressBar = this.store.selectSignal(LayoutState.showProgressBar);

  public userPreferences = this.store.selectSignal(AuthState.userPreferences);

  // Phone-width viewports have no room for a permanent search input, so below
  // $breakpoint-sm it is toggled open in place of the logo (header.component.scss).
  public mobileSearchOpen = signal(false);

  public receiptHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.receiptListLink)];
  });

  public dashboardHeaderLink = computed(() => {
    this.selectedGroupId();
    return [this.store.selectSnapshot(GroupState.dashboardLink)];
  });

  protected readonly Permission = Permission;

  private aiPoweredReceipts = this.store.selectSignal(
    FeatureConfigState.aiPoweredReceipts
  );

  private groupPermissions = this.store.selectSignal(AuthState.groupPermissions);

  // Whether app-quick-scan-button will actually render (it self-gates on the
  // feature flag), so the bottom tab bar can fall back to manual add instead
  // of showing no add action at all.
  protected readonly canQuickScan = computed(
    () =>
      !!this.aiPoweredReceipts() &&
      hasAll(
        this.groupPermissions()?.[this.selectedGroupIdNumber()] ?? [],
        Permission.GroupReceiptsQuickScan
      )
  );

  // Reports are reachable with read OR the readAll bypass (the *hasAppPermission
  // directive is single-key AND-only, so the OR is resolved through the selector).
  protected readonly canViewReports = this.store.selectSignal(
    AuthState.hasAnyAppPermission([
      Permission.AppReportsRead,
      Permission.AppReportsReadAll,
    ])
  );

  protected readonly canViewSystemSettings = this.store.selectSignal(
    AuthState.hasAnyAppPermission([
      Permission.AppSystemSettingsRead,
      Permission.AppPromptsRead,
      Permission.AppReceiptProcessingSettingsRead,
      Permission.AppSystemEmailsRead,
      Permission.AppSystemTasksRead,
    ])
  );

  protected readonly canViewUserSettings = this.store.selectSignal(
    AuthState.hasAnyAppPermission([
      Permission.AppAccountRead,
      Permission.AppUserPreferencesRead,
      Permission.AppApiKeysRead,
    ])
  );

  constructor(
    private authService: AuthService,
    private matDialog: MatDialog,
    private router: Router,
    private routeReuseStrategy: RouteReuseStrategy,
    private store: Store
  ) {}

  public logout(): void {
    this.authService
      .logout()
      .pipe(
        take(1),
        switchMap(() => this.store.dispatch(new Logout())),
        switchMap(() => this.router.navigate(["/"])),
      )
      .subscribe();
  }

  public openImportDialog(): void {
    this.matDialog.open(ImportFormComponent, DEFAULT_DIALOG_CONFIG);
  }

  public openBankStatementImportDialog(): void {
    this.matDialog
      .open(BankStatementImportComponent, BANK_STATEMENT_DIALOG_CONFIG)
      .afterClosed()
      .pipe(take(1))
      .subscribe((imported) => {
        if (imported) {
          this.reloadCurrentRoute();
        }
      });
  }

  // Rebuilds the routed page in place so its list or dashboard picks up newly
  // imported receipts. A same-URL navigation alone keeps the existing component,
  // so route reuse is switched off for just this one navigation.
  private reloadCurrentRoute(): void {
    const shouldReuseRoute = this.routeReuseStrategy.shouldReuseRoute;
    this.routeReuseStrategy.shouldReuseRoute = () => false;

    this.router
      .navigateByUrl(this.router.url, { onSameUrlNavigation: "reload" })
      .finally(() => (this.routeReuseStrategy.shouldReuseRoute = shouldReuseRoute));
  }

  public openAboutDialog(): void {
    this.matDialog.open(AboutComponent, DEFAULT_DIALOG_CONFIG);
  }
}
