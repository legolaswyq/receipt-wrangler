import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { of } from "rxjs";
import { PipesModule } from "src/pipes/pipes.module";
import { ReceiptTableState } from "src/store/receipt-table.state";
import { ApiModule, Permission, Receipt } from "../../open-api";
import { QuickScanProgressJob, QuickScanProgressService } from "../../services";
import { AuthState, GroupState } from "../../store";
import { SetPermissions } from "../../store/auth.state.actions";
import { ReceiptsTableComponent } from "./receipts-table.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("ReceiptsTableComponent", () => {
  let component: ReceiptsTableComponent;
  let fixture: ComponentFixture<ReceiptsTableComponent>;
  let store: Store;
  let quickScanJobs: ReturnType<typeof signal<QuickScanProgressJob[]>>;

  beforeEach(async () => {
    quickScanJobs = signal<QuickScanProgressJob[]>([]);

    await TestBed.configureTestingModule({
    declarations: [ReceiptsTableComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        NgxsModule.forRoot([ReceiptTableState, AuthState, GroupState]),
        ReactiveFormsModule,
        MatSnackBarModule,
        MatTooltipModule,
        MatDialogModule,
        PipesModule],
    providers: [
        {
            provide: ActivatedRoute,
            useValue: {
                snapshot: {
                    data: {
                        categories: [],
                        tags: [],
                    },
                },
            },
        },
        { provide: QuickScanProgressService, useValue: { jobs: quickScanJobs, dismissJob: jest.fn() } },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(ReceiptsTableComponent);
    component = fixture.componentInstance;
    Object.defineProperty(component, 'table', {
      value: () => ({
        selection: {},
        changed: of(undefined),
      }),
    });
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("gates each header action on its own group permission", () => {
    store.dispatch(
      new SetPermissions([], {
        5: [
          Permission.GroupReceiptsCreate,
          Permission.GroupReceiptsQuickScan,
          Permission.GroupEmailPoll,
        ],
      })
    );
    component.groupId = "5";

    (component as any).setCanEdit();

    expect(component.canCreate).toEqual(true);
    expect(component.canQuickScan).toEqual(true);
    expect(component.canPollEmail).toEqual(true);
    // None of those imply update, which gates Edit / Bulk Status Update.
    expect(component.canEdit).toEqual(false);
  });

  it("keeps canEdit tied to group.receipts.update only", () => {
    store.dispatch(
      new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] })
    );
    component.groupId = "5";

    (component as any).setCanEdit();

    expect(component.canEdit).toEqual(true);
    expect(component.canCreate).toEqual(false);
    expect(component.canQuickScan).toEqual(false);
    expect(component.canPollEmail).toEqual(false);
  });

  it("should map selected ids from selecton", () => {
    const selectedReceipts: Receipt[] = [
      {
        id: 1,
      } as Receipt,
      {
        id: 2,
      } as Receipt,
    ];
    Object.defineProperty(component, 'table', {
      value: () => ({
        selection: {
          changed: of({
            source: {
              selected: selectedReceipts,
            },
          }),
        },
      }),
    });
    component.ngAfterViewInit();

    expect(component.selectedReceiptIds()).toEqual([1, 2]);
  });

  describe("quick-scan completion reload (handleQuickScanJobsChanged)", () => {
    // Quick scan creates receipts entirely server-side with no push notification, so without
    // this the new row only ever appears after a manual page refresh. Calls the handler the
    // constructor's effect() delegates to directly -- this TestBed setup never renders a real
    // template, so forcing a full change-detection tick (as flushing the effect itself would
    // require) trips up unrelated lifecycle code (ngOnInit/ngAfterViewInit) that this spec
    // otherwise never exercises.
    it("reloads receipts once a job targeting this group completes an image", () => {
      component.groupId = "5";
      (component as any).group = { isAllGroup: false };
      const reloadSpy = jest.spyOn(component, "getFilteredReceipts").mockImplementation(() => {});

      (component as any).handleQuickScanJobsChanged([
        { id: 1, groupIds: [5], total: 1, completed: 0, stillProcessing: false },
      ]);
      expect(reloadSpy).not.toHaveBeenCalled();

      (component as any).handleQuickScanJobsChanged([
        { id: 1, groupIds: [5], total: 1, completed: 1, stillProcessing: false },
      ]);
      expect(reloadSpy).toHaveBeenCalledTimes(1);
    });

    it("ignores a job's completion for a group other than the one being viewed", () => {
      component.groupId = "5";
      (component as any).group = { isAllGroup: false };
      const reloadSpy = jest.spyOn(component, "getFilteredReceipts").mockImplementation(() => {});

      (component as any).handleQuickScanJobsChanged([
        { id: 1, groupIds: [9], total: 1, completed: 1, stillProcessing: false },
      ]);

      expect(reloadSpy).not.toHaveBeenCalled();
    });

    it("reloads for any group's completion when viewing the virtual All group", () => {
      component.groupId = "2";
      (component as any).group = { isAllGroup: true };
      const reloadSpy = jest.spyOn(component, "getFilteredReceipts").mockImplementation(() => {});

      (component as any).handleQuickScanJobsChanged([
        { id: 1, groupIds: [9], total: 1, completed: 1, stillProcessing: false },
      ]);

      expect(reloadSpy).toHaveBeenCalledTimes(1);
    });
  });
});
