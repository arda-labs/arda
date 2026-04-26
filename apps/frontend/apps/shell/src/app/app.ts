import { Component, inject } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ToastModule } from 'primeng/toast';
import { LanguageService } from '@arda-mfe/shared-core';

@Component({
  imports: [RouterModule, ToastModule],
  selector: 'app-root',
  templateUrl: './app.html',
})
export class App {
  private langService = inject(LanguageService);
}
