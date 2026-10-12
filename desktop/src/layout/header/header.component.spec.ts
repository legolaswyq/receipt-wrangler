import { OverlayContainer } from "@angular/cdk/overlay";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { provideHttpClientTesting } from "@angular/common/http/testing";
import { NO_ERRORS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";
import { MatIconModule } from "@angular/material/icon";
import { MatMenuModule, MatMenuTrigger } from "@angular/material/menu";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { By } from "@angular/platform-browser";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { Router, RouteReuseStrategy } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { Store } from "@ngxs/store";
import { NgbPopoverModule } from "@ng-bootstrap/ng-bootstrap";
import { of } from "rxjs";
import { DirectivesModule } from "../../directives/directives.module";
import { ApiModule, AuthService, Permission } from "../../open-api";
import { SetAuthState, SetPermissions } from "../../store/auth.state.actions";
import { SetFeatureConfig } from "../../store/feature-config.state.actions";
import { SetSelectedGroupId } from "../../store/group.state.actions";
import { Logout } from "../../store";
import { StoreModule } from "../../store/store.module";
import { HeaderComponent } from "./header.component";

describe("HeaderComponent", () => {
  let component: HeaderComponent;
  let fixture: ComponentFixture<HeaderComponent>;
  let store: Store;
  let overlayContainer: OverlayContainer;

  const logIn = () =>
    store.dispatch(
      new SetAuthState({ exp: Math.floor(Date.now() / 1000) + 3600 } as any)
    );

  const searchbarRendered = (): boolean =>
    !!fixture.nativeElement.querySelector("app-searchbar");

  const bellRendered = (): boolean =>
    !!fixture.nativeElement.querySelector('app-button[icon="notifications"]');

  const queryTestId = (id: string) =>
    fixture.nativeElement.querySelector(`[data-testid="${id}"]`);

  // The settings items live inside a MatMenu, so they render into the overlay
  // only once the gear trigger is opened.
  const openSettingsMenu = async () => {
    await fixture.whenStable();
    const trigger = fixture.debugElement
      .queryAll(By.directive(MatMenuTrigger))
      .find(
        (el) =>
          el.nativeElement.getAttribute("data-testid") ===
          "header-settings-menu"
      )!
      .injector.get(MatMenuTrigger);
    trigger.openMenu();
    await fixture.whenStable();
  };

  const overlayText = (): string =>
    overlayContainer.getContainerElement().textContent ?? "";

  const overlayTestId = (id: string) =>
    overlayContainer
      .getContainerElement()
      .querySelector(`[data-testid="${id}"]`);

  const grantReceiptsSearch = () =>
    store.dispatch(new SetPermissions(["app.receipts.search"], {}));

  const selectGroupOne = () => store.dispatch(new SetSelectedGroupId("1"));

  const grantGroupPerms = (groupId: number, perms: string[]) =>
    store.dispatch(new SetPermissions([], { [groupId]: perms }));

  beforeEach(async () => {
    // StoreModule persists auth (incl. permissions) to localStorage; clear it so
    // granted permissions don't leak across tests in this file.
    localStorage.clear();

    await TestBed.configureTestingModule({
    declarations: [HeaderComponent],
    imports: [ApiModule,
        DirectivesModule,
        MatDialogModule,
        MatIconModule,
        MatMenuModule,
        MatSnackBarModule,
        NgbPopoverModule,
        NoopAnimationsModule,
        RouterTestingModule,
        StoreModule,
    ],
    providers: [provideZonelessChangeDetection(), provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()],
    schemas: [NO_ERRORS_SCHEMA],
}).compileComponents();

    store = TestBed.inject(Store);
    overlayContainer = TestBed.inject(OverlayContainer);
    fixture = TestBed.createComponent(HeaderComponent);
    component = fixture.componentInstance;
    fixture.autoDetectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should not render a notifications bell", async () => {
    logIn();
    await fixture.whenStable();

    expect(bellRendered()).toBe(false);
  });

  it("should hide the search bar without app.receipts.search", async () => {
    logIn();
    await fixture.whenStable();

    expect(searchbarRendered()).toBe(false);
  });

  it("should show the search bar with app.receipts.search", async () => {
    logIn();
    grantReceiptsSearch();
    await fixture.whenStable();

    expect(searchbarRendered()).toBe(true);
  });

  it("hides the mobile search toggle without app.receipts.search", async () => {
    logIn();
    await fixture.whenStable();

    expect(queryTestId("header-mobile-search-toggle")).toBeFalsy();
  });

  it("shows the mobile search toggle with app.receipts.search", async () => {
    logIn();
    grantReceiptsSearch();
    await fixture.whenStable();

    expect(queryTestId("header-mobile-search-toggle")).toBeTruthy();
  });

  it("toggles mobileSearchOpen when the mobile search button is clicked", async () => {
    logIn();
    grantReceiptsSearch();
    await fixture.whenStable();

    expect(component.mobileSearchOpen()).toBe(false);

    queryTestId("header-mobile-search-toggle").click();
    await fixture.whenStable();
    expect(component.mobileSearchOpen()).toBe(true);
    expect(
      fixture.nativeElement.querySelector("#header-container").classList
    ).toContain("mobile-search-active");

    queryTestId("header-mobile-search-toggle").click();
    await fixture.whenStable();
    expect(component.mobileSearchOpen()).toBe(false);
  });

  it("should hide the Dashboard tab without group.dashboards.read for the selected group", async () => {
    logIn();
    selectGroupOne();
    await fixture.whenStable();

    expect(queryTestId("header-dashboard-tab")).toBeFalsy();
  });

  it("should show the Dashboard tab with group.dashboards.read for the selected group", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, ["group.dashboards.read"]);
    await fixture.whenStable();

    expect(queryTestId("header-dashboard-tab")).toBeTruthy();
  });

  it("should always render the Receipts tab", async () => {
    logIn();
    await fixture.whenStable();

    expect(queryTestId("header-receipts-tab")).toBeTruthy();
  });

  it("gates the primary Add Receipt (quick scan) button on the group's quick-scan permission", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, ["group.receipts.quick-scan"]);
    await fixture.whenStable();

    expect(queryTestId("header-add-receipt")).toBeTruthy();
    expect(queryTestId("header-add-manual")).toBeFalsy();
  });

  it("gates the manual add button on the group's create permission", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, ["group.receipts.create"]);
    await fixture.whenStable();

    expect(queryTestId("header-add-manual")).toBeTruthy();
    expect(queryTestId("header-add-receipt")).toBeFalsy();
  });

  it("hides the admin manage menu items without the matching read permissions", async () => {
    logIn();
    await openSettingsMenu();

    const text = overlayText();
    expect(text).not.toContain("Manage Categories");
    expect(text).not.toContain("Manage Tags");
    expect(text).not.toContain("Manage Groups");
    expect(text).not.toContain("Manage Custom Fields");
  });

  it("shows the admin manage menu items with the matching read permissions", async () => {
    logIn();
    store.dispatch(
      new SetPermissions(
        [
          "app.categories.read",
          "app.tags.read",
          "app.groups.read",
          "app.custom-fields.read",
        ],
        {}
      )
    );
    await openSettingsMenu();

    const text = overlayText();
    expect(text).toContain("Manage Categories");
    expect(text).toContain("Manage Tags");
    expect(text).toContain("Manage Groups");
    expect(text).toContain("Manage Custom Fields");
  });

  it("shows the Add Group item with the app group-create permission", async () => {
    logIn();
    store.dispatch(new SetPermissions([Permission.AppGroupsCreate], {}));
    await openSettingsMenu();

    expect(overlayTestId("header-add-group")).toBeTruthy();
  });

  it("hides the Add Group item without the app group-create permission", async () => {
    logIn();
    await openSettingsMenu();

    expect(overlayTestId("header-add-group")).toBeFalsy();
  });

  it("shows Import Bank Statement with receipt-create on the selected group", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, [Permission.GroupReceiptsCreate]);
    await openSettingsMenu();

    expect(overlayTestId("header-import-bank-statement")).toBeTruthy();
  });

  it("hides Import Bank Statement without receipt-create on the selected group", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, [Permission.GroupReceiptsRead]);
    await openSettingsMenu();

    expect(overlayTestId("header-import-bank-statement")).toBeFalsy();
  });

  it("rebuilds the current page after a bank statement import, and not after a cancel", async () => {
    const dialog = TestBed.inject(MatDialog);
    const router = TestBed.inject(Router);
    const reuseStrategy = TestBed.inject(RouteReuseStrategy);
    const originalShouldReuse = reuseStrategy.shouldReuseRoute;
    const navigate = jest.spyOn(router, "navigateByUrl").mockResolvedValue(true);
    const open = jest.spyOn(dialog, "open");

    open.mockReturnValue({ afterClosed: () => of(false) } as any);
    component.openBankStatementImportDialog();
    expect(navigate).not.toHaveBeenCalled();

    open.mockReturnValue({ afterClosed: () => of(true) } as any);
    component.openBankStatementImportDialog();
    expect(navigate).toHaveBeenCalledWith(router.url, { onSameUrlNavigation: "reload" });

    // Route reuse is only switched off for that one navigation.
    await fixture.whenStable();
    expect(reuseStrategy.shouldReuseRoute).toBe(originalShouldReuse);
  });

  it("hides User Settings without any settings read permission", async () => {
    logIn();
    await openSettingsMenu();

    expect(overlayText()).not.toContain("User Settings");
  });

  it("shows User Settings with a settings read permission", async () => {
    logIn();
    store.dispatch(new SetPermissions(["app.account.read"], {}));
    await openSettingsMenu();

    expect(overlayText()).toContain("User Settings");
  });

  it("renders the Dashboard/Receipts links in the bottom tab bar", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, ["group.dashboards.read"]);
    await fixture.whenStable();

    expect(queryTestId("header-dashboard-tab-mobile")).toBeTruthy();
    expect(queryTestId("header-receipts-tab-mobile")).toBeTruthy();
  });

  it("hides the Dashboard link in the bottom tab bar without group.dashboards.read", async () => {
    logIn();
    selectGroupOne();
    await fixture.whenStable();

    expect(queryTestId("header-dashboard-tab-mobile")).toBeFalsy();
    expect(queryTestId("header-receipts-tab-mobile")).toBeTruthy();
  });

  it("does not render the bottom tab bar when logged out", async () => {
    await fixture.whenStable();

    expect(queryTestId("header-bottom-nav")).toBeFalsy();
  });

  it("uses Quick Scan as the bottom tab bar add action when it is available", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, ["group.receipts.quick-scan", "group.receipts.create"]);
    store.dispatch(new SetFeatureConfig({ aiPoweredReceipts: true } as any));
    await fixture.whenStable();

    expect(queryTestId("header-add-receipt-mobile")).toBeTruthy();
    expect(queryTestId("header-add-manual-mobile")).toBeFalsy();
  });

  it("falls back to manual add in the bottom tab bar when the AI feature is off", async () => {
    logIn();
    selectGroupOne();
    grantGroupPerms(1, ["group.receipts.quick-scan", "group.receipts.create"]);
    store.dispatch(new SetFeatureConfig({ aiPoweredReceipts: false } as any));
    await fixture.whenStable();

    expect(queryTestId("header-add-receipt-mobile")).toBeFalsy();
    expect(queryTestId("header-add-manual-mobile")).toBeTruthy();
  });

  it("shows no bottom tab bar add action without quick-scan or create permission", async () => {
    logIn();
    selectGroupOne();
    await fixture.whenStable();

    expect(queryTestId("header-add-receipt-mobile")).toBeFalsy();
    expect(queryTestId("header-add-manual-mobile")).toBeFalsy();
  });

  it("dispatches Logout on logout", async () => {
    const authService = TestBed.inject(AuthService);
    jest.spyOn(authService, "logout").mockReturnValue(of(undefined) as any);
    jest.spyOn(TestBed.inject(Router), "navigate").mockResolvedValue(true);
    const dispatch = jest.spyOn(store, "dispatch");

    component.logout();
    await fixture.whenStable();

    expect(dispatch).toHaveBeenCalledWith(new Logout());
  });
});
