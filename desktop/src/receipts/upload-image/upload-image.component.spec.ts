import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { ActivatedRoute } from "@angular/router";
import { of } from "rxjs";
import { ApiModule, ReceiptImageService } from "../../open-api";
import { UploadImageComponent } from "./upload-image.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";

describe("UploadImageComponent", () => {
  let component: UploadImageComponent;
  let fixture: ComponentFixture<UploadImageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [UploadImageComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule, MatSnackBarModule],
    providers: [
        { provide: ActivatedRoute, useValue: { snapshot: { data: {} } } },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
}).compileComponents();

    fixture = TestBed.createComponent(UploadImageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  describe("PDF uploads", () => {
    const upload = async (convertToPreview: boolean) => {
      const receiptImageService = TestBed.inject(ReceiptImageService);
      const convertToJpg = jest
        .spyOn(receiptImageService, "convertToJpg")
        .mockReturnValue(of({ encodedImage: "encoded" }) as any);
      const emitted: any[] = [];
      component.fileLoaded.subscribe((command) => emitted.push(command));
      fixture.componentRef.setInput("convertToPreview", convertToPreview);

      const file = new File(["%PDF"], "statement.pdf", { type: "application/pdf" });
      await component.onFileChange({ target: { files: [file] } });
      await new Promise((resolve) => setTimeout(resolve, 50));

      return { convertToJpg, emitted, file };
    };

    it("renders a JPG preview by default", async () => {
      const { convertToJpg, emitted } = await upload(true);

      expect(convertToJpg).toHaveBeenCalled();
      expect(emitted[0].encodedImage).toBe("encoded");
    });

    it("emits the raw file without a preview when convertToPreview is off", async () => {
      const { convertToJpg, emitted, file } = await upload(false);

      expect(convertToJpg).not.toHaveBeenCalled();
      expect(emitted).toEqual([{ file, receiptId: 0 }]);
    });

    it("emits the accepted files of one pick together, leaving out other types", async () => {
      const picked: File[][] = [];
      component.filesSelected.subscribe((files) => picked.push(files));
      fixture.componentRef.setInput("acceptFileTypes", ["application/pdf"]);
      fixture.componentRef.setInput("convertToPreview", false);

      const first = new File(["%PDF"], "a.pdf", { type: "application/pdf" });
      const second = new File(["%PDF"], "b.pdf", { type: "application/pdf" });
      const other = new File(["x"], "notes.txt", { type: "text/plain" });
      await component.onFileChange({ target: { files: [first, other, second] } });

      expect(picked).toEqual([[first, second]]);
    });
  });
});
