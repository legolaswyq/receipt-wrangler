import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { ButtonModule } from "../button";
import { CheckboxModule } from "../checkbox/checkbox.module";
import { PipesModule } from "../pipes";
import { ReceiptsModule } from "../receipts/receipts.module";
import { SelectModule } from "../select/select.module";
import { SharedUiModule } from "../shared-ui/shared-ui.module";
import { UserAutocompleteModule } from "../user-autocomplete/user-autocomplete.module";
import { BankStatementImportComponent } from "./bank-statement-import/bank-statement-import.component";
import { ImportFormComponent } from "./import-form/import-form.component";


@NgModule({
  declarations: [
    BankStatementImportComponent,
    ImportFormComponent
  ],
  imports: [
    ButtonModule,
    CheckboxModule,
    CommonModule,
    MatProgressBarModule,
    PipesModule,
    ReactiveFormsModule,
    ReceiptsModule,
    SelectModule,
    SharedUiModule,
    UserAutocompleteModule,
  ],
  exports: [
    BankStatementImportComponent,
    ImportFormComponent
  ]
})
export class ImportModule {}
