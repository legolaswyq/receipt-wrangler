import { TestBed } from "@angular/core/testing";
import { of } from "rxjs";
import { ReceiptService } from "../open-api";
import { QuickScanProgressService } from "./quick-scan-progress.service";

describe("QuickScanProgressService", () => {
  let service: QuickScanProgressService;
  let receiptService: { getReceiptsForGroup: jest.Mock };

  beforeEach(() => {
    jest.useFakeTimers();
    receiptService = { getReceiptsForGroup: jest.fn() };

    TestBed.configureTestingModule({
      providers: [{ provide: ReceiptService, useValue: receiptService }],
    });

    service = TestBed.inject(QuickScanProgressService);
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it("does nothing when there are no files", () => {
    service.trackQuickScan([1], 0);

    expect(service.jobs()).toEqual([]);
    expect(receiptService.getReceiptsForGroup).not.toHaveBeenCalled();
  });

  it("creates a job seeded with the total file count and its unique groups", () => {
    receiptService.getReceiptsForGroup.mockReturnValue(of({ data: [], totalCount: 5 }));

    service.trackQuickScan([1, 1, 2], 3);

    expect(service.jobs()).toEqual([
      { id: expect.any(Number), groupIds: [1, 2], total: 3, completed: 0, stillProcessing: false },
    ]);
    // Baseline fetch: one call per unique group id.
    expect(receiptService.getReceiptsForGroup).toHaveBeenCalledTimes(2);
  });

  it("marks files complete as each group's receipt count grows, then auto-removes the job", () => {
    let callCount = 0;
    receiptService.getReceiptsForGroup.mockImplementation(() => {
      callCount++;
      // First call is the baseline (5); every call after simulates one more receipt landing.
      return of({ data: [], totalCount: callCount === 1 ? 5 : 6 });
    });

    service.trackQuickScan([1], 1);

    jest.advanceTimersByTime(3000);

    expect(service.jobs()[0].completed).toBe(1);
    expect(service.jobs().length).toBe(1);

    jest.advanceTimersByTime(2000);
    expect(service.jobs()).toEqual([]);
  });

  it("caps a group's contribution at the number of files sent to it", () => {
    // A single group receiving many new receipts should still only count as 1 completed
    // file when only 1 file targeted that group.
    let callCount = 0;
    receiptService.getReceiptsForGroup.mockImplementation(() => {
      callCount++;
      return of({ data: [], totalCount: callCount === 1 ? 0 : 10 });
    });

    service.trackQuickScan([1], 1);

    jest.advanceTimersByTime(3000);

    expect(service.jobs()[0].completed).toBe(1);
  });

  it("gives up after the max poll window and flags the job as still processing", () => {
    // Never increases -- nothing ever completes.
    receiptService.getReceiptsForGroup.mockReturnValue(of({ data: [], totalCount: 0 }));

    service.trackQuickScan([1], 1);

    // 40 attempts * 3000ms
    jest.advanceTimersByTime(40 * 3000);

    expect(service.jobs()[0].stillProcessing).toBe(true);
    expect(service.jobs().length).toBe(1);

    jest.advanceTimersByTime(4000);
    expect(service.jobs()).toEqual([]);
  });

  it("stops polling and removes the job on manual dismiss", () => {
    receiptService.getReceiptsForGroup.mockReturnValue(of({ data: [], totalCount: 0 }));

    service.trackQuickScan([1], 1);
    const jobId = service.jobs()[0].id;
    const callsBeforeDismiss = receiptService.getReceiptsForGroup.mock.calls.length;

    service.dismissJob(jobId);
    expect(service.jobs()).toEqual([]);

    jest.advanceTimersByTime(30000);
    expect(receiptService.getReceiptsForGroup.mock.calls.length).toBe(callsBeforeDismiss);
  });

  it("tracks concurrent jobs independently", () => {
    let callCount = 0;
    receiptService.getReceiptsForGroup.mockImplementation(() => {
      callCount++;
      // Group 1 baseline 0 then jumps to 1; group 2 baseline 0 then stays 0.
      return of({ data: [], totalCount: callCount <= 2 ? 0 : 1 });
    });

    service.trackQuickScan([1], 1);
    service.trackQuickScan([2], 1);

    expect(service.jobs().length).toBe(2);
  });
});
