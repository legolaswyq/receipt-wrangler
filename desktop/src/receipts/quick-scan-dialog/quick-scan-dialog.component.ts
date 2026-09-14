import { Component, DestroyRef, OnInit, ViewEncapsulation, inject, viewChild } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validators } from "@angular/forms";
import { MatDialogRef } from "@angular/material/dialog";
import { Store } from "@ngxs/store";
import { take, tap } from "rxjs";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { setRequired } from "../../form";
import { Category, GroupReceiptSettings, Permission, ReceiptService, ReceiptStatus, Tag } from "../../open-api";
import { QuickScanProgressService, SnackbarService } from "../../services";
import { AuthState, GroupState } from "../../store";
import { codePointMaxLengthValidator, trimmedRequiredValidator } from "../../validators";
import { UploadImageComponent } from "../upload-image/upload-image.component";

// Mirrors the backend's models.MaxCommentLength (the Comment column is varchar(500), which
// MySQL/Postgres measure in characters) and mobile's FormBuilderTextField maxLength, so an
// over-length comment is caught here instead of coming back as a 400 after submit. Counted in
// code points, like the backend's rune count -- Validators.maxLength counts UTF-16 code units,
// which would reject 500 emoji the API accepts.
export const QUICK_SCAN_COMMENT_MAX_LENGTH = 500;

// Built once so setRequired can remove them by reference on a later recompute.
const commentMaxLengthValidator = codePointMaxLengthValidator(QUICK_SCAN_COMMENT_MAX_LENGTH);
// QuickScanCommand trims each comment on parse, so a whitespace-only one arrives empty and fails
// the group's required check -- Validators.required would have called it valid and eaten a 400.
const commentRequiredValidator = trimmedRequiredValidator();

@Component({
    selector: "app-quick-scan-dialog",
    templateUrl: "./quick-scan-dialog.component.html",
    styleUrls: ["./quick-scan-dialog.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class QuickScanDialogComponent implements OnInit {
  public readonly uploadImageComponent = viewChild.required(UploadImageComponent);

  public form: FormGroup = new FormGroup({});

  public images: ReceiptFileUploadCommand[] = [];

  public currentlySelectedIndex: number = 0;

  // When on (and 2+ images), all images are one long receipt: sent with combineImages=true so the
  // backend transcribes each, combines the text, and runs a single extraction into one receipt with
  // every image attached. Kept OUT of `form` so form.value stays the per-image arrays. Defaults on
  // because a multi-image upload is most often one receipt too long to fit in a single photo.
  public combineImages = new FormControl(true);

  // base-input has no built-in message for the maxlength error, so an unmapped one would render an
  // empty mat-error. Held as a field rather than an inline object literal so the binding is stable.
  public readonly commentErrorMessages: { [key: string]: string } = {
    maxlength: `Comment must be ${QUICK_SCAN_COMMENT_MAX_LENGTH} characters or fewer.`,
  };

  private readonly commentPermissionByGroup = new Map<number, boolean>();

  private readonly destroyRef = inject(DestroyRef);

  constructor(
    private dialogRef: MatDialogRef<QuickScanDialogComponent>,
    private formBuilder: FormBuilder,
    private quickScanProgressService: QuickScanProgressService,
    private receiptService: ReceiptService,
    private snackbarService: SnackbarService,
    private store: Store
  ) {}

  public get paidByUserIds(): FormArray {
    return this.form.get("paidByUserIds") as FormArray;
  }

  public get statuses(): FormArray {
    return this.form.get("statuses") as FormArray;
  }

  public get groupIds(): FormArray {
    return this.form.get("groupIds") as FormArray;
  }

  public get categories(): FormArray {
    return this.form.get("categories") as FormArray;
  }

  public get tags(): FormArray {
    return this.form.get("tags") as FormArray;
  }

  public get comments(): FormArray {
    return this.form.get("comments") as FormArray;
  }

  public ngOnInit(): void {
    this.initForm();
  }

  private initForm(): void {
    this.form = this.formBuilder.group({
      paidByUserIds: this.formBuilder.array<number>([]),
      statuses: this.formBuilder.array<ReceiptStatus>([]),
      groupIds: this.formBuilder.array<number>([]),
      categories: this.formBuilder.array<Category[]>([]),
      tags: this.formBuilder.array<Tag[]>([]),
      comments: this.formBuilder.array<string>([]),
    });

    // Re-resolve each image's field config whenever anything changes (a group selection can flip
    // which fields are shown/required for that image).
    this.form.valueChanges
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe(() => this.configureImages());
  }

  public fileLoaded(fileData: ReceiptFileUploadCommand): void {
    if (fileData.file && !fileData.encodedImage) {
      fileData.url = URL.createObjectURL(fileData.file);
    }
    this.images.push(fileData);
    const userPreferences = this.store.selectSnapshot(AuthState.userPreferences);

    // Fall back to sensible defaults when the user hasn't configured quick-scan preferences:
    // the first real group, the uploader (current user) as payer, and OPEN status -- so a fresh
    // upload isn't blocked on a "Group is required" error with empty paid-by/status.
    const firstGroupId = this.store.selectSnapshot(GroupState.groupsWithoutAll)?.[0]?.id ?? "";
    const currentUserId = Number(this.store.selectSnapshot(AuthState.userId)) || "";

    // Group is always required; paid-by/status validators are applied per-image by configureImages.
    this.paidByUserIds.push(new FormControl(userPreferences?.quickScanDefaultPaidById || currentUserId));
    this.statuses.push(new FormControl(userPreferences?.quickScanDefaultStatus || ReceiptStatus.Open));
    this.groupIds.push(new FormControl(userPreferences?.quickScanDefaultGroupId || firstGroupId, Validators.required));
    // Categories/tags are multi-selects backed by app-category/tag-autocomplete,
    // which push selections onto a FormArray (see the receipt form). A plain
    // FormControl has no push(), so selecting one would throw — use a FormArray.
    this.categories.push(this.formBuilder.array([]));
    this.tags.push(this.formBuilder.array([]));
    // The comment is a scalar text field, so a plain FormControl - unlike categories/tags above.
    // maxLength is applied here rather than in configureImages because setRequired composes
    // validators additively, so it survives every show/require recompute and is never re-added.
    this.comments.push(new FormControl("", commentMaxLengthValidator));

    this.configureImages();
  }

  // Returns the receipt settings for the group currently selected on the given image, if any.
  public settingsForIndex(index: number): GroupReceiptSettings | undefined {
    const groupId = this.groupIds.at(index)?.value;
    if (!groupId) {
      return undefined;
    }

    return this.store.selectSnapshot(GroupState.getGroupById(groupId.toString()))?.groupReceiptSettings;
  }

  public showPaidBy(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanPaidByEnabled ?? true;
  }

  public showStatus(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanStatusEnabled ?? true;
  }

  public showCategories(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanCategoriesEnabled ?? false;
  }

  public showTags(index: number): boolean {
    return this.settingsForIndex(index)?.quickScanTagsEnabled ?? false;
  }

  // The comment field additionally requires group.comments.create: without it the field is hidden,
  // never required, and a comment sent anyway is dropped server-side. hideComments hides the whole
  // group's comments, so it hides this too. Mirrors the backend's IsQuickScanCommentShown.
  public showComment(index: number): boolean {
    const settings = this.settingsForIndex(index);
    if (!(settings?.quickScanCommentEnabled ?? false) || (settings?.hideComments ?? false)) {
      return false;
    }

    return this.canCommentForIndex(index);
  }

  // Cached per group: AuthState.hasGroupPermission allocates a new selector on each call, and
  // showComment is read from the template on every change-detection pass.
  private canCommentForIndex(index: number): boolean {
    const groupId = Number(this.groupIds.at(index)?.value);
    if (!groupId) {
      return false;
    }

    if (!this.commentPermissionByGroup.has(groupId)) {
      this.commentPermissionByGroup.set(
        groupId,
        this.store.selectSnapshot(AuthState.hasGroupPermission(groupId, Permission.GroupCommentsCreate))
      );
    }

    return this.commentPermissionByGroup.get(groupId) ?? false;
  }

  public categoriesForIndex(index: number): Category[] {
    const groupId = this.groupIds.at(index)?.value;
    return groupId ? this.store.selectSnapshot(AuthState.groupCategories(Number(groupId))) : [];
  }

  public tagsForIndex(index: number): Tag[] {
    const groupId = this.groupIds.at(index)?.value;
    return groupId ? this.store.selectSnapshot(AuthState.groupTags(Number(groupId))) : [];
  }

  // Applies the per-image required validators and clears hidden fields so the server backfills their
  // group-configured defaults (a hidden paid-by/status must be sent empty, not with a stale value).
  private configureImages(): void {
    for (let i = 0; i < this.images.length; i++) {
      const settings = this.settingsForIndex(i);
      const paidByShown = settings?.quickScanPaidByEnabled ?? true;
      const statusShown = settings?.quickScanStatusEnabled ?? true;
      const categoriesShown = settings?.quickScanCategoriesEnabled ?? false;
      const tagsShown = settings?.quickScanTagsEnabled ?? false;

      setRequired(this.form.get("paidByUserIds." + i), paidByShown && (settings?.quickScanPaidByRequired ?? true));
      setRequired(this.form.get("statuses." + i), statusShown && (settings?.quickScanStatusRequired ?? true));
      setRequired(this.form.get("categories." + i), categoriesShown && (settings?.quickScanCategoriesRequired ?? false));
      setRequired(this.form.get("tags." + i), tagsShown && (settings?.quickScanTagsRequired ?? false));

      if (!paidByShown) {
        this.form.get("paidByUserIds." + i)?.setValue("", { emitEvent: false });
      }
      if (!statusShown) {
        this.form.get("statuses." + i)?.setValue("", { emitEvent: false });
      }
      if (!categoriesShown) {
        (this.form.get("categories." + i) as FormArray | null)?.clear({ emitEvent: false });
      }
      if (!tagsShown) {
        (this.form.get("tags." + i) as FormArray | null)?.clear({ emitEvent: false });
      }

      const commentShown = this.showComment(i);
      setRequired(
        this.form.get("comments." + i),
        commentShown && (settings?.quickScanCommentRequired ?? false),
        commentRequiredValidator
      );
      if (!commentShown) {
        this.form.get("comments." + i)?.setValue("", { emitEvent: false });
      }
    }
  }

  public openImageUploadComponent(): void {
    this.uploadImageComponent().clickInput();
  }

  public removeImage(index: number): void {
    this.paidByUserIds.removeAt(index);
    this.statuses.removeAt(index);
    this.groupIds.removeAt(index);
    this.categories.removeAt(index);
    this.tags.removeAt(index);
    this.comments.removeAt(index);
    this.images.splice(index, 1);
  }

  // Combine is only meaningful with 2+ images; a single image is always its own receipt.
  public isCombineActive(): boolean {
    return !!this.combineImages.value && this.images.length > 1;
  }

  // In combine mode only the shared (index 0) field set matters; the other per-image controls are
  // ignored, so validity is judged on index 0 alone.
  private isCombineValid(): boolean {
    return [this.groupIds, this.paidByUserIds, this.statuses, this.categories, this.tags, this.comments]
      .every((array) => array.at(0)?.valid ?? true);
  }

  public submitButtonClicked(): void {
    if (this.images.length === 0) {
      this.snackbarService.error("Please select images to upload");
      return;
    }

    const combine = this.isCombineActive();
    if (!(combine ? this.isCombineValid() : this.form.valid)) {
      this.snackbarService.error("Please fill in all required fields. Some images are missing required fields.");
      return;
    }

    const count = this.images.length;
    // In combine mode the per-file arrays must still carry one entry per file (the API validates
    // len == files), so the shared index-0 value is repeated across every file; the backend reads
    // index 0. Otherwise each file keeps its own value.
    const repeat = <T>(value: T): T[] => Array(count).fill(value);
    const groupIds = combine ? repeat(this.groupIds.at(0)?.value) : this.groupIds.value;
    const paidByUserIds = combine ? repeat(this.paidByUserIds.at(0)?.value) : this.paidByUserIds.value;
    const statuses = combine ? repeat(this.statuses.at(0)?.value) : this.statuses.value;
    const categoryIds = combine ? repeat(this.joinIdsForControl(this.categories.at(0))) : this.joinIds(this.categories);
    const tagIds = combine ? repeat(this.joinIdsForControl(this.tags.at(0))) : this.joinIds(this.tags);
    const comments = combine ? repeat(this.comments.at(0)?.value) : this.comments.value;

    this.receiptService
      .quickScanReceipt(
        this.images.map((i) => i.file),
        groupIds,
        paidByUserIds,
        statuses,
        categoryIds,
        tagIds,
        comments,
        combine
      )
      .pipe(
        take(1),
        tap(() => {
          // Quick scan is fire-and-forget (no task/receipt id comes back), so progress is
          // shown via a persistent snackbar that polls for the new receipts to appear
          // rather than a one-shot "queued" toast. Combine produces ONE receipt.
          this.quickScanProgressService.trackQuickScan(
            combine ? [this.groupIds.at(0)?.value] : this.groupIds.value,
            combine ? 1 : this.images.length
          );
          this.dialogRef.close();
        }),
      )
      .subscribe();
  }

  // Serializes each image's selected category/tag objects into a comma-joined id string (one entry
  // per image), matching the multipart shape the API expects.
  private joinIds(array: FormArray): string[] {
    return array.controls.map((control) => this.joinIdsForControl(control));
  }

  // Comma-joins one control's selected category/tag ids into the string the multipart form expects.
  private joinIdsForControl(control: AbstractControl | null): string {
    return ((control?.value ?? []) as { id: number }[]).map((entity) => entity.id).join(",");
  }

  public cancelButtonClicked(): void {
    this.dialogRef.close();
  }

  public navigateImages(delta: number): void {
    const newValue = this.currentlySelectedIndex + delta;
    if (newValue === -1) {
      this.currentlySelectedIndex = this.images.length - 1;
      return;
    }

    if (newValue > this.images.length - 1) {
      this.currentlySelectedIndex = 0;
      return;
    }

    this.currentlySelectedIndex = newValue;
  }
}
