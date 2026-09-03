import { Component, computed, inject, input } from "@angular/core";
import { QuickScanProgressJob, QuickScanProgressService } from "../../services";

@Component({
  selector: "app-quick-scan-progress-banner",
  templateUrl: "./quick-scan-progress-banner.component.html",
  styleUrls: ["./quick-scan-progress-banner.component.scss"],
  standalone: false,
})
export class QuickScanProgressBannerComponent {
  public readonly groupId = input<string | number | undefined>(undefined);

  // The virtual "All" group aggregates every real group's receipts, but a quick-scanned
  // image is always assigned to one specific real group (picked in the dialog) -- never to
  // the All group itself. So on the All page, surface every in-flight job rather than
  // filtering by groupId, or nothing would ever show there.
  public readonly isAllGroup = input<boolean>(false);

  private readonly quickScanProgressService = inject(QuickScanProgressService);

  // Only the jobs targeting this page's group -- a single quick scan can span multiple
  // groups (each image picks its own), so the banner only surfaces on a group it's part of.
  public readonly jobs = computed(() => {
    if (this.isAllGroup()) {
      return this.quickScanProgressService.jobs();
    }

    const groupId = Number(this.groupId());
    if (!groupId) {
      return [];
    }

    return this.quickScanProgressService.jobs().filter((job) => job.groupIds.includes(groupId));
  });

  public isDone(job: QuickScanProgressJob): boolean {
    return job.total > 0 && job.completed >= job.total;
  }

  public dismiss(jobId: number): void {
    this.quickScanProgressService.dismissJob(jobId);
  }
}
