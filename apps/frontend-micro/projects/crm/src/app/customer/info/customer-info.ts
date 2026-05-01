import { Component, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TabsModule } from 'primeng/tabs';
import { BasicInfoComponent } from '../shared/ui/general-info/basic-info/basic-info';
import { AttachmentsComponent } from '../shared/ui/attachments/attachments';

@Component({
  selector: 'app-customer-info',
  standalone: true,
  imports: [CommonModule, TabsModule, BasicInfoComponent, AttachmentsComponent],
  templateUrl: './customer-info.html',
  styleUrl: './customer-info.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CustomerInfoComponent {}
