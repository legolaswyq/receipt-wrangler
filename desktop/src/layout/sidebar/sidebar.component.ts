import { Component, ViewEncapsulation } from "@angular/core";
import { toSignal } from "@angular/core/rxjs-interop";
import { Data, NavigationEnd, Router } from "@angular/router";
import { Store } from "@ngxs/store";
import { filter, map, startWith } from "rxjs";
import { AuthState } from "../../store";

// The app shell: a single top bar (app-header) over the routed content. The old
// left drawer rail was removed — group switching, add actions, and the admin
// menu now live in the header.
@Component({
    selector: "app-sidebar",
    templateUrl: "./sidebar.component.html",
    styleUrls: ["./sidebar.component.scss"],
    encapsulation: ViewEncapsulation.None,
    standalone: false
})
export class SidebarComponent {
  constructor(private router: Router, private store: Store) {}

  public isLoggedIn = this.store.selectSignal(AuthState.isLoggedIn);

  // True when the active route opts into a full-height, no-padding content frame
  // (route data `fullHeight`), so the routed page owns its own internal scrolling
  // instead of the shell's block-flow + p-4 padding. Used by the Report Builder.
  public readonly isContentFullHeight = toSignal(
    this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      startWith(null),
      map(() => this.deepestRouteData()["fullHeight"] === true),
    ),
    { initialValue: false },
  );

  // Walks from the router state root to the deepest activated child so the
  // fullHeight flag is picked up wherever it is declared in the route tree.
  private deepestRouteData(): Data {
    let route = this.router.routerState.snapshot.root;
    while (route.firstChild) {
      route = route.firstChild;
    }
    return route.data;
  }
}
