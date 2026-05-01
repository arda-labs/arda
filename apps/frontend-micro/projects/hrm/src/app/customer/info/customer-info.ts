import { Component, ChangeDetectionStrategy } from '@angular/core';
import { BasicInfoComponent } from '../shared/ui/general-info/basic-info/basic-info';

@Component({
  selector: 'app-customer-info',
  imports: [BasicInfoComponent],
  templateUrl: './customer-info.html',
  styleUrl: './customer-info.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CustomerInfoComponent {}
