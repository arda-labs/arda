import { Component, ChangeDetectionStrategy } from '@angular/core';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-customer',
  imports: [RouterOutlet],
  templateUrl: './customer.html',
  styleUrl: './customer.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CustomerComponent {}
