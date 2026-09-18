import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { CUSTOM_ELEMENTS_SCHEMA, provideZonelessChangeDetection } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { Router } from "@angular/router";
import { RouterTestingModule } from "@angular/router/testing";
import { NgxsModule } from "@ngxs/store";
import { LayoutState } from "src/store/layout.state";
import { ApiModule } from "../../open-api";
import { AuthState, FeatureConfigState, GroupState } from "../../store";
import { SidebarComponent } from "./sidebar.component";

describe("SidebarComponent (app shell)", () => {
  let component: SidebarComponent;
  let fixture: ComponentFixture<SidebarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
    declarations: [SidebarComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [NgxsModule.forRoot([AuthState, LayoutState, GroupState, FeatureConfigState]),
        RouterTestingModule,
        ApiModule],
    providers: [provideZonelessChangeDetection(), provideHttpClient(withInterceptorsFromDi())]
}).compileComponents();

    fixture = TestBed.createComponent(SidebarComponent);
    component = fixture.componentInstance;
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("reflects the active route's fullHeight data in isContentFullHeight", async () => {
    const router = TestBed.inject(Router);
    router.resetConfig([
      { path: "plain", children: [] },
      { path: "full", data: { fullHeight: true }, children: [] },
    ]);

    await router.navigateByUrl("/plain");
    await fixture.whenStable();
    expect(component.isContentFullHeight()).toBe(false);

    await router.navigateByUrl("/full");
    await fixture.whenStable();
    expect(component.isContentFullHeight()).toBe(true);
  });
});
