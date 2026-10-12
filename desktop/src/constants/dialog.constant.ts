import { MatDialogConfig } from '@angular/material/dialog';

export const DEFAULT_DIALOG_CONFIG: MatDialogConfig<any> = {
  width: '50%',
  // The fixed 50% width is unusable below ~900px (e.g. a 390px phone renders a
  // ~195px dialog); app-responsive.scss overrides the overlay pane's width via
  // this panel class on narrow viewports, since MatDialogConfig has no media
  // query concept of its own.
  panelClass: 'app-dialog-responsive',
};
