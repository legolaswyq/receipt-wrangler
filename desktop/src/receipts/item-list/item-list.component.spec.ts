import { CUSTOM_ELEMENTS_SCHEMA } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { FormArray, FormControl, FormGroup, ReactiveFormsModule } from "@angular/forms";
import { MatDialogModule } from "@angular/material/dialog";
import { MatExpansionModule } from "@angular/material/expansion";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { Store } from "@ngxs/store";
import { PipesModule } from "src/pipes/pipes.module";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { StoreModule } from "src/store/store.module";
import { Permission } from "../../open-api";
import { SetPermissions } from "../../store/auth.state.actions";
import { ItemData, ItemListComponent } from "./item-list.component";

function buildItemGroup(
  values: Partial<{
    name: string;
    amount: string;
    quantity: string;
    unitPrice: string;
    chargedToUserId: number;
  }> = {}
): FormGroup {
  return new FormGroup({
    name: new FormControl(values.name ?? "Item"),
    nameZh: new FormControl(""),
    quantity: new FormControl(values.quantity ?? ""),
    unitPrice: new FormControl(values.unitPrice ?? ""),
    amount: new FormControl(values.amount ?? "0"),
    status: new FormControl("OPEN"),
    chargedToUserId: new FormControl(values.chargedToUserId ?? undefined),
    categories: new FormArray([]),
    tags: new FormArray([]),
  });
}

describe("ItemListComponent", () => {
  let component: ItemListComponent;
  let fixture: ComponentFixture<ItemListComponent>;
  let store: Store;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ItemListComponent],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      imports: [
        ReactiveFormsModule,
        MatDialogModule,
        MatExpansionModule,
        NoopAnimationsModule,
        PipesModule,
        SharedUiModule,
        StoreModule,
      ],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: { snapshot: { data: { receipt: { id: 1, groupId: 5 } } } },
        },
      ],
    }).compileComponents();

    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(ItemListComponent);
    component = fixture.componentInstance;
  });

  function setFormWithItems(itemGroups: FormGroup[]): void {
    const form = new FormGroup({ receiptItems: new FormArray(itemGroups) });
    fixture.componentRef.setInput("form", form);
  }

  describe("getTotalAmount", () => {
    it("sums every non-share item's total price", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([buildItemGroup({ amount: "10.50" }), buildItemGroup({ amount: "4.25" })]);
      fixture.detectChanges();

      expect(component.getTotalAmount()).toBeCloseTo(14.75);
    });

    it("excludes items charged to a user (shares)", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([
        buildItemGroup({ amount: "10.00" }),
        buildItemGroup({ amount: "5.00", chargedToUserId: 2 }),
      ]);
      fixture.detectChanges();

      expect(component.getTotalAmount()).toBeCloseTo(10);
    });

    it("reflects a live edit without needing setItems() to re-run", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([buildItemGroup({ amount: "10.00" })]);
      fixture.detectChanges();

      component.receiptItems.at(0).get("amount")?.setValue("25.00");

      expect(component.getTotalAmount()).toBeCloseTo(25);
    });
  });

  describe("isAmountMismatched", () => {
    const itemData = (arrayIndex: number): ItemData => ({ item: {} as any, arrayIndex });

    it("returns false when quantity times unit price matches the total price", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([buildItemGroup({ quantity: "2", unitPrice: "3.00", amount: "6.00" })]);
      fixture.detectChanges();

      expect(component.isAmountMismatched(itemData(0))).toBe(false);
    });

    it("returns true when quantity times unit price disagrees with the total price", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([buildItemGroup({ quantity: "2", unitPrice: "3.00", amount: "9.00" })]);
      fixture.detectChanges();

      expect(component.isAmountMismatched(itemData(0))).toBe(true);
    });

    it("returns false when quantity or unit price is not set", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([buildItemGroup({ quantity: "", unitPrice: "", amount: "9.00" })]);
      fixture.detectChanges();

      expect(component.isAmountMismatched(itemData(0))).toBe(false);
    });

    it("tolerates a cent of rounding slack", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([buildItemGroup({ quantity: "3", unitPrice: "3.333", amount: "10.00" })]);
      fixture.detectChanges();

      expect(component.isAmountMismatched(itemData(0))).toBe(false);
    });
  });

  describe("addInlineItem", () => {
    it("emits a new blank item when the caller can edit", () => {
      store.dispatch(new SetPermissions([], { 5: [Permission.GroupReceiptsUpdate] }));
      setFormWithItems([]);
      fixture.detectChanges();

      const emitted: unknown[] = [];
      component.itemAdded.subscribe((item) => emitted.push(item));

      component.addInlineItem();

      expect(emitted.length).toBe(1);
    });

    it("does not emit without edit permission", () => {
      store.dispatch(new SetPermissions([], {}));
      setFormWithItems([]);
      fixture.detectChanges();

      const emitted: unknown[] = [];
      component.itemAdded.subscribe((item) => emitted.push(item));

      component.addInlineItem();

      expect(emitted.length).toBe(0);
    });
  });
});
