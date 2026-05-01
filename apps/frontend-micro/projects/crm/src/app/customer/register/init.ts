import { Component, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { StepperModule } from 'primeng/stepper';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { BasicInfoComponent } from '../shared/ui/general-info/basic-info/basic-info';
import { AttachmentsComponent } from '../shared/ui/attachments/attachments';
import { CustomerService } from '../shared/services/customer.service';

@Component({
  selector: 'app-init',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    StepperModule,
    ButtonModule,
    CardModule,
    ToastModule,
    BasicInfoComponent,
    AttachmentsComponent
  ],
  providers: [MessageService],
  templateUrl: './init.html',
  styleUrl: './init.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InitComponent {
  private customerService = inject(CustomerService);
  private messageService = inject(MessageService);

  activeStep = signal(1);
  formData = signal({
    name: '',
    idNumber: '',
    address: '',
    email: '',
    phone: ''
  });

  nextStep() {
    this.activeStep.update(s => s + 1);
  }

  prevStep() {
    this.activeStep.update(s => s - 1);
  }

  onSubmit() {
    this.messageService.add({ severity: 'info', summary: 'Đang gửi hồ sơ', detail: 'Vui lòng đợi trong giây lát...' });

    this.customerService.createRequest(this.formData()).subscribe({
      next: (res) => {
        this.messageService.add({
          severity: 'success',
          summary: 'Thành công',
          detail: 'Hồ sơ đăng ký đã được gửi và đang chờ phê duyệt tại BPM.'
        });
      },
      error: (err) => {
        this.messageService.add({
          severity: 'error',
          summary: 'Lỗi',
          detail: 'Không thể kết nối đến máy chủ. Vui lòng thử lại sau.'
        });
      }
    });
  }
}
