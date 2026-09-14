import { Component, OnChanges, SimpleChanges, ViewEncapsulation, inject, input, output } from "@angular/core";
import { MatDialog } from "@angular/material/dialog";
import { UntilDestroy } from "@ngneat/until-destroy";
import { FormMode } from "src/enums/form-mode.enum";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { FileDataView } from "../../open-api";
import { ImagePreviewDialogComponent } from "../../shared-ui/image-preview-dialog/image-preview-dialog.component";

@UntilDestroy()
@Component({
    selector: "app-carousel",
    templateUrl: "./carousel.component.html",
    styleUrls: ["./carousel.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class CarouselComponent implements OnChanges {
  public readonly images = input<FileDataView[]>([]);

  public readonly imagePreviews = input<ReceiptFileUploadCommand[]>([]);

  public readonly disabled = input<boolean>(false);

  public readonly mode = input.required<FormMode>();

  public readonly hideButtonControls = input<boolean>(false);

  public readonly initialIndex = input<number>(-1);

  public readonly removeButtonClicked = output<number>();

  public scale: number = 1;

  public currentlyShownImageIndex: number = 0;

  private readonly matDialog = inject(MatDialog);

  // Opens a full-width, vertically-scrollable preview of the clicked image -- reading a long
  // receipt top-to-bottom, rather than the inline zoom the +/- controls provide.
  public openImagePreview(src: string): void {
    if (!src) {
      return;
    }
    this.matDialog.open(ImagePreviewDialogComponent, {
      maxWidth: "95vw",
      width: "600px",
      panelClass: "image-preview-dialog-panel",
      data: { src },
    });
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["initialIndex"]) {
      this.currentlyShownImageIndex = this.initialIndex();
    }
  }

  public emitRemoveButtonClicked(index: number): void {
    this.removeButtonClicked.emit(index);
  }

  public zoomOut() {
    this.adjustScale(-0.1);
  }

  public zoomIn() {
    this.adjustScale(0.1);
  }

  public updateCurrentlyShownImage(index: number): void {
    this.currentlyShownImageIndex = index;
  }

  public adjustScale(amount: number): void {
    const newScale = this.scale + amount;
    this.scale = Math.max(newScale, 0.1);
  }
}
