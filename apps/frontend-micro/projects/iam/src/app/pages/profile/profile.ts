import { ChangeDetectionStrategy, Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule, FormControl, FormGroup, Validators } from '@angular/forms';
import { rxResource } from '@angular/core/rxjs-interop';
import { Avatar } from 'primeng/avatar';
import { Button } from 'primeng/button';
import { InputText } from 'primeng/inputtext';
import { TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    Avatar,
    Button,
    InputText,
    TableModule,
    Tag,
    Toast,
  ],
  providers: [MessageService],
  templateUrl: './profile.html',
  styleUrl: './profile.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Profile implements OnInit {
  private authService = inject(AuthService);
  private messageService = inject(MessageService);

  user = this.authService.currentUser;
  initials = this.authService.userInitials;

  profileForm = new FormGroup({
    email: new FormControl({ value: '', disabled: true }, { nonNullable: true }),
    displayName: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(3)] })
  });

  auditLogsResource = rxResource({
    stream: () => this.authService.getMyAuditLogs()
  });

  isSaving = signal(false);

  ngOnInit(): void {
    const currentUser = this.user();
    if (currentUser) {
      this.profileForm.patchValue({
        email: currentUser.email,
        displayName: currentUser.name
      });
    }
  }

  saveProfile() {
    if (this.profileForm.invalid) return;

    this.isSaving.set(true);
    const { displayName } = this.profileForm.getRawValue();
    this.authService.updateProfile(displayName).subscribe({
      next: () => {
        this.messageService.add({
          severity: 'success',
          summary: 'Thành công',
          detail: 'Thông tin cá nhân đã được cập nhật',
        });
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({
          severity: 'error',
          summary: 'Lỗi',
          detail: 'Không thể cập nhật thông tin',
        });
        this.isSaving.set(false);
      }
    });
  }

  getActionSeverity(action: string): 'success' | 'info' | 'warn' | 'danger' | 'secondary' | 'contrast' | undefined {
    if (action.includes('create') || action.includes('grant')) return 'success';
    if (action.includes('update') || action.includes('edit')) return 'warn';
    if (action.includes('delete') || action.includes('revoke')) return 'danger';
    return 'info';
  }
}
