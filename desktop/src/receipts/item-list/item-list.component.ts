import {
  Component,
  ElementRef,
  HostListener,
  Input,
  OnChanges,
  OnDestroy,
  OnInit,
  Signal,
  SimpleChanges,
  ViewEncapsulation,
  input,
  output,
  signal,
  viewChild,
  viewChildren
} from "@angular/core";
import { FormArray, FormGroup } from "@angular/forms";
import { MatDialog } from "@angular/material/dialog";
import { MatExpansionPanel } from "@angular/material/expansion";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { InputComponent } from "../../input";
import { Category, Group, Item, Permission, Receipt, Tag } from "../../open-api";
import { AuthState } from "../../store";
import { KeyboardShortcutService } from "../../services/keyboard-shortcut.service";
import { Subject, takeUntil } from "rxjs";
import { KEYBOARD_SHORTCUT_ACTIONS } from "../../constants/keyboard-shortcuts.constant";
import { QuickActionsDialogComponent } from "../quick-actions-dialog/quick-actions-dialog.component";
import { DEFAULT_DIALOG_CONFIG } from "src/constants";

export interface ItemData {
  item: Item;
  arrayIndex: number;
}

@Component({
  selector: "app-item-list",
  templateUrl: "./item-list.component.html",
  styleUrls: ["./item-list.component.scss"],
  encapsulation: ViewEncapsulation.None,
  standalone: false
})
export class ItemListComponent implements OnInit, OnChanges, OnDestroy {
  public readonly itemsExpansionPanel = viewChild.required<MatExpansionPanel>("itemsExpansionPanel");

  public readonly nameFields = viewChildren<InputComponent>("nameField");


  public readonly form = input.required<FormGroup>();

  @Input() public originalReceipt?: Receipt;

  public readonly categories = input<Category[]>([]);

  public readonly tags = input<Tag[]>([]);

  public readonly selectedGroup = input<Group>();

  public readonly triggerAddMode = input<boolean>(false);

  public readonly itemAdded = output<Item>();

  public readonly itemRemoved = output<{
    item: Item;
    arrayIndex: number;
}>();

  public readonly itemSplit = output<{
    items: Item[];
    itemIndex: number;
}>();


  public items = signal<ItemData[]>([]);

  protected readonly Permission = Permission;

  /**
   * Reactive gate for editing items: holds `group.receipts.update` for the
   * receipt's group. Reassigned in `ngOnInit` once the receipt is resolved;
   * deny-by-default until then.
   */
  public canEdit: Signal<boolean> = signal(false);


  private destroy$ = new Subject<void>();

  public get receiptItems(): FormArray {
    return this.form().get("receiptItems") as FormArray;
  }

  constructor(
    private activatedRoute: ActivatedRoute,
    private keyboardShortcutService: KeyboardShortcutService,
    private matDialog: MatDialog,
    private store: Store
  ) {}

  public ngOnInit(): void {
    this.originalReceipt = this.activatedRoute.snapshot.data["receipt"];
    this.canEdit = this.store.selectSignal(
      AuthState.hasGroupPermission(
        this.originalReceipt?.groupId ?? 0,
        Permission.GroupReceiptsUpdate
      )
    );
    this.setItems();
    this.setupKeyboardShortcuts();
  }

  public ngOnChanges(changes: SimpleChanges): void {
    if (changes["triggerAddMode"] && changes["triggerAddMode"].currentValue) {
      this.addInlineItem();
    }
    if (changes["form"]) {
      this.setItems();
    }
  }

  @HostListener("document:keydown", ["$event"])
  public handleKeyboardShortcut(event: KeyboardEvent): void {
    // Items are editable whenever the caller holds edit permission, regardless of the
    // page's view/edit mode -- there is no separate "unlock to edit" step for items.
    if (!this.canEdit()) {
      return;
    }

    // Let the service handle the keyboard event
    this.keyboardShortcutService.handleKeyboardEvent(event);
  }

  private setupKeyboardShortcuts(): void {
    // Subscribe to keyboard shortcut events
    this.keyboardShortcutService.shortcutTriggered
      .pipe(takeUntil(this.destroy$))
      .subscribe(shortcutEvent => {
        this.handleShortcutAction(shortcutEvent.action);
      });
  }

  private handleShortcutAction(action: string): void {
    switch (action) {
      case KEYBOARD_SHORTCUT_ACTIONS.ADD_ITEM:
        this.addInlineItem();
        break;
    }
  }

  public setItems(): void {
    const receiptItems = this.form().get("receiptItems");
    if (receiptItems) {
      const items = this.form().get("receiptItems")?.value as Item[];
      const itemDataArray: ItemData[] = [];

      if (items?.length > 0) {
        items.forEach((item, index) => {
          if (!item?.chargedToUserId) {
            const itemData: ItemData = {
              item: item,
              arrayIndex: index,
            };
            itemDataArray.push(itemData);
          }
        });
      }
      this.items.set(itemDataArray);
    }
  }

  public removeItem(itemData: ItemData): void {
    this.itemRemoved.emit({ item: itemData.item, arrayIndex: itemData.arrayIndex });
  }

  public splitItem(itemData: ItemData): void {
    const dialogRef = this.matDialog.open(QuickActionsDialogComponent, {
      ...DEFAULT_DIALOG_CONFIG,
      data: {
        originalReceipt: this.originalReceipt,
        amountToSplit: parseFloat(itemData.item.amount) || 0,
        itemIndex: itemData.arrayIndex,
        usersToOmit: []
      }
    });

    // Pass data directly to component since it uses @Input decorators
    dialogRef.componentInstance.originalReceipt = this.originalReceipt;
    dialogRef.componentInstance.amountToSplit = parseFloat(itemData.item.amount) || 0;
    dialogRef.componentInstance.itemIndex = itemData.arrayIndex;
    dialogRef.componentInstance.usersToOmit = [];

    // Subscribe to the component's itemsToAdd output
    dialogRef.componentInstance.itemsToAdd.subscribe((data: { items: Item[], itemIndex?: number }) => {
      if (data.itemIndex !== undefined) {
        this.itemSplit.emit({ items: data.items, itemIndex: data.itemIndex });
      }
    });

    dialogRef.afterClosed().subscribe((result: boolean) => {
      // Dialog closed
    });
  }

  public addInlineItem(event?: MouseEvent): void {
    if (event) {
      event?.stopImmediatePropagation();
    }

    if (this.canEdit()) {
      const newItem = {
        name: "",
        chargedToUserId: undefined,
      } as Item;
      this.itemAdded.emit(newItem);

      const itemsExpansionPanel = this.itemsExpansionPanel();
      if (itemsExpansionPanel && !itemsExpansionPanel.expanded) {
        itemsExpansionPanel.open();
      }
    }
  }

  public addInlineItemOnBlur(index: number): void {
    const items = this.items();
    if (items && items.length - 1 === index) {
      const item = items.at(index) as ItemData;
      const itemInput = this.receiptItems.at(item?.arrayIndex);
      if (itemInput.valid) {
        this.addInlineItem();
      }
    }
  }

  public checkLastInlineItem(): void {
    if (this.canEdit()) {
      const items = this.items();
      if (items && items.length > 1) {
        const lastItem = items[items.length - 1];
        const formGroup = this.receiptItems.at(lastItem.arrayIndex);
        const nameValue = formGroup.get("name")?.value;
        const amountValue = formGroup.get("amount")?.value;

        if (formGroup.pristine && (!nameValue || nameValue.trim() === "") && (!amountValue || amountValue === 0)) {
          this.itemRemoved.emit({ item: lastItem.item, arrayIndex: lastItem.arrayIndex });
        }
      }
    }
  }

  // Reads live FormArray control values (not the items() snapshot, which only refreshes on
  // add/remove) so the total reflects a Total Price edit immediately, without needing to blur
  // out to the next row first. Excludes shares (chargedToUserId set), matching setItems().
  public getTotalAmount(): number {
    return this.receiptItems.controls.reduce((total, control) => {
      if (control.get("chargedToUserId")?.value) {
        return total;
      }
      const amount = parseFloat(control.get("amount")?.value) || 0;
      return total + amount;
    }, 0);
  }

  // An item's total price should equal quantity * unitPrice when the receipt actually printed
  // both. When the AI extraction's amount disagrees with that calculation (beyond a cent of
  // rounding slack), flag it rather than silently trusting either value.
  public isAmountMismatched(itemData: ItemData): boolean {
    const control = this.receiptItems.at(itemData.arrayIndex);
    const quantity = parseFloat(control.get("quantity")?.value);
    const unitPrice = parseFloat(control.get("unitPrice")?.value);
    const amount = parseFloat(control.get("amount")?.value);

    if (Number.isNaN(quantity) || Number.isNaN(unitPrice) || Number.isNaN(amount)) {
      return false;
    }

    return Math.abs(quantity * unitPrice - amount) > 0.01;
  }

  // Keyboard event handlers

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
    this.keyboardShortcutService.clearHintTimeout();
  }
}
