import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CUSTOMER_FIELD_NAMES } from '../../../constants/customer-field-names.constant';

@Component({
  selector: 'app-basic-info',
  imports: [],
  templateUrl: './basic-info.html',
  styleUrl: './basic-info.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class BasicInfoComponent {
  fieldNames = signal(CUSTOMER_FIELD_NAMES.basicInfo);
}
