import { Injectable, signal } from "@angular/core";
import { forkJoin, interval, map, Observable, Subscription } from "rxjs";
import { switchMap, take } from "rxjs/operators";
import { ReceiptService } from "../open-api";

const POLL_INTERVAL_MS = 3000;
// ~2 minutes of polling before giving up and just leaving the job's last known state.
const MAX_POLL_ATTEMPTS = 40;

export interface QuickScanProgressJob {
  id: number;
  groupIds: number[];
  total: number;
  completed: number;
  // Set once polling gave up waiting for every receipt to appear -- the receipts may still be
  // processing on the server, this only means the client stopped watching.
  stillProcessing: boolean;
}

// Quick scan is fire-and-forget: the backend enqueues background OCR/AI jobs and returns
// immediately, with no task or receipt id to correlate against. This service approximates
// progress by polling each target group's receipt count and treating any increase (up to the
// number of files sent to that group) as "done" -- an inference, not a guarantee. A receipt
// created concurrently by someone else in the same group would also count.
@Injectable({
  providedIn: "root",
})
export class QuickScanProgressService {
  private readonly jobsSignal = signal<QuickScanProgressJob[]>([]);

  public readonly jobs = this.jobsSignal.asReadonly();

  private nextJobId = 0;

  private readonly pollSubscriptions = new Map<number, Subscription>();

  constructor(private receiptService: ReceiptService) {}

  public trackQuickScan(groupIds: number[], totalFiles: number): void {
    if (totalFiles <= 0) {
      return;
    }

    const targetCountByGroup = new Map<number, number>();
    for (const groupId of groupIds) {
      targetCountByGroup.set(groupId, (targetCountByGroup.get(groupId) ?? 0) + 1);
    }
    const uniqueGroupIds = [...targetCountByGroup.keys()];

    const jobId = ++this.nextJobId;
    this.jobsSignal.update((jobs) => [
      ...jobs,
      { id: jobId, groupIds: uniqueGroupIds, total: totalFiles, completed: 0, stillProcessing: false },
    ]);

    this.fetchReceiptCounts(uniqueGroupIds)
      .pipe(take(1))
      .subscribe((baselineCounts) => {
        let attempts = 0;
        const subscription: Subscription = interval(POLL_INTERVAL_MS)
          .pipe(switchMap(() => this.fetchReceiptCounts(uniqueGroupIds)))
          .subscribe((currentCounts) => {
            attempts++;

            let completed = 0;
            for (const [groupId, target] of targetCountByGroup) {
              const gained = (currentCounts.get(groupId) ?? 0) - (baselineCounts.get(groupId) ?? 0);
              completed += Math.min(Math.max(gained, 0), target);
            }
            this.patchJob(jobId, { completed });

            if (completed >= totalFiles) {
              this.stopPolling(jobId);
              setTimeout(() => this.removeJob(jobId), 2000);
            } else if (attempts >= MAX_POLL_ATTEMPTS) {
              this.stopPolling(jobId);
              this.patchJob(jobId, { stillProcessing: true });
              setTimeout(() => this.removeJob(jobId), 4000);
            }
          });

        this.pollSubscriptions.set(jobId, subscription);
      });
  }

  public dismissJob(jobId: number): void {
    this.stopPolling(jobId);
    this.removeJob(jobId);
  }

  private stopPolling(jobId: number): void {
    this.pollSubscriptions.get(jobId)?.unsubscribe();
    this.pollSubscriptions.delete(jobId);
  }

  private patchJob(jobId: number, patch: Partial<QuickScanProgressJob>): void {
    this.jobsSignal.update((jobs) => jobs.map((job) => (job.id === jobId ? { ...job, ...patch } : job)));
  }

  private removeJob(jobId: number): void {
    this.jobsSignal.update((jobs) => jobs.filter((job) => job.id !== jobId));
  }

  private fetchReceiptCounts(groupIds: number[]): Observable<Map<number, number>> {
    const requests = groupIds.map((groupId) =>
      this.receiptService
        .getReceiptsForGroup(groupId, { page: 1, pageSize: 1 })
        .pipe(map((paged) => [groupId, paged.totalCount ?? 0] as [number, number]))
    );

    return forkJoin(requests).pipe(map((entries) => new Map(entries)));
  }
}
