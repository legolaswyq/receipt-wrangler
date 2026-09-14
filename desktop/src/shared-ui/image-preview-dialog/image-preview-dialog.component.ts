import { Component, Inject } from "@angular/core";
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from "@angular/material/dialog";
import { MatIconModule } from "@angular/material/icon";

export interface ImagePreviewDialogData {
  src: string;
  alt?: string;
}

// A full-height, vertically-scrollable preview of a single receipt image. Unlike the carousel's
// zoom controls, this shows the image at the dialog's full width and lets a tall receipt scroll
// top-to-bottom -- the natural way to read a long receipt.
@Component({
  selector: "app-image-preview-dialog",
  standalone: true,
  imports: [MatDialogModule, MatIconModule],
  templateUrl: "./image-preview-dialog.component.html",
  styleUrl: "./image-preview-dialog.component.scss",
})
export class ImagePreviewDialogComponent {
  constructor(
    private dialogRef: MatDialogRef<ImagePreviewDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ImagePreviewDialogData
  ) {}

  public close(): void {
    this.dialogRef.close();
  }
}
