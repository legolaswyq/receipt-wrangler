import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { QuickScanProgressJob, QuickScanProgressService } from "../../services";
import { QuickScanProgressBannerComponent } from "./quick-scan-progress-banner.component";

describe("QuickScanProgressBannerComponent", () => {
  let component: QuickScanProgressBannerComponent;
  let fixture: ComponentFixture<QuickScanProgressBannerComponent>;
  let jobs: QuickScanProgressJob[];
  let quickScanProgressService: { jobs: jest.Mock; dismissJob: jest.Mock };

  beforeEach(() => {
    jobs = [];
    quickScanProgressService = {
      jobs: jest.fn(() => jobs),
      dismissJob: jest.fn(),
    };

    TestBed.configureTestingModule({
      declarations: [QuickScanProgressBannerComponent],
      imports: [MatButtonModule, MatIconModule, MatProgressBarModule, NoopAnimationsModule],
      providers: [{ provide: QuickScanProgressService, useValue: quickScanProgressService }],
    });

    fixture = TestBed.createComponent(QuickScanProgressBannerComponent);
    component = fixture.componentInstance;
  });

  it("creates", () => {
    fixture.componentRef.setInput("groupId", 1);
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it("shows no banner when there are no jobs for this group", () => {
    jobs = [{ id: 1, groupIds: [2], total: 3, completed: 1, stillProcessing: false }];
    fixture.componentRef.setInput("groupId", 1);
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('[data-testid="quick-scan-progress-banner"]')).toBeNull();
  });

  it("shows a banner for a job targeting this group, with in-progress text", () => {
    jobs = [{ id: 1, groupIds: [1], total: 3, completed: 1, stillProcessing: false }];
    fixture.componentRef.setInput("groupId", "1");
    fixture.detectChanges();

    const banner = fixture.nativeElement.querySelector('[data-testid="quick-scan-progress-banner"]');
    expect(banner).not.toBeNull();
    expect(banner.querySelector('[data-testid="quick-scan-progress-text"]').textContent).toContain(
      "Extracting text and structuring data… (1/3 images)"
    );
    // No fractional progress is trackable per image (only two discrete backend
    // stages exist server-side, with no queryable in-between state), so an
    // in-progress job animates indeterminately rather than sitting static.
    expect(banner.querySelector("mat-progress-bar").getAttribute("mode")).toBe("indeterminate");
  });

  it("shows completed text and a full determinate bar once the job is done", () => {
    jobs = [{ id: 1, groupIds: [1], total: 2, completed: 2, stillProcessing: false }];
    fixture.componentRef.setInput("groupId", 1);
    fixture.detectChanges();

    const banner = fixture.nativeElement.querySelector('[data-testid="quick-scan-progress-banner"]');
    expect(banner.querySelector('[data-testid="quick-scan-progress-text"]').textContent).toContain(
      "All 2 images processed"
    );
    expect(banner.querySelector("mat-progress-bar").getAttribute("mode")).toBe("determinate");
  });

  it("shows still-processing text once polling has given up", () => {
    jobs = [{ id: 1, groupIds: [1], total: 2, completed: 0, stillProcessing: true }];
    fixture.componentRef.setInput("groupId", 1);
    fixture.detectChanges();

    const text = fixture.nativeElement.querySelector('[data-testid="quick-scan-progress-text"]').textContent;
    expect(text).toContain("Still processing in the background");
  });

  it("dismisses the job when the dismiss button is clicked", () => {
    jobs = [{ id: 42, groupIds: [1], total: 2, completed: 1, stillProcessing: false }];
    fixture.componentRef.setInput("groupId", 1);
    fixture.detectChanges();

    fixture.nativeElement.querySelector('[data-testid="quick-scan-progress-dismiss"]').click();

    expect(quickScanProgressService.dismissJob).toHaveBeenCalledWith(42);
  });

  it("shows every job on the virtual All group page, regardless of which real group it targets", () => {
    jobs = [
      { id: 1, groupIds: [1], total: 3, completed: 1, stillProcessing: false },
      { id: 2, groupIds: [5], total: 1, completed: 0, stillProcessing: false },
    ];
    // The All group's own id (e.g. 2) matches neither job's groupIds.
    fixture.componentRef.setInput("groupId", 2);
    fixture.componentRef.setInput("isAllGroup", true);
    fixture.detectChanges();

    const banners = fixture.nativeElement.querySelectorAll('[data-testid="quick-scan-progress-banner"]');
    expect(banners.length).toBe(2);
  });
});
