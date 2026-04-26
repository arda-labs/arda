import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AvatarModule } from 'primeng/avatar';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { AuthService, AuditLog } from '../../services/auth.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    ReactiveFormsModule,
    AvatarModule,
    ButtonModule,
    InputTextModule,
    TableModule,
    TagModule,
    ToastModule
  ],
  providers: [MessageService],
  templateUrl: './profile.html',
  styleUrl: './profile.css',
})
export class Profile implements OnInit {
  private authService = inject(AuthService);
  private fb = inject(FormBuilder);
  private messageService = inject(MessageService);

  user = this.authService.currentUser;
  initials = this.authService.userInitials;

  profileForm: FormGroup;
  auditLogs = signal<AuditLog[]>([]);
  isSaving = signal(false);
  isLoadingLogs = signal(false);

  constructor() {
    this.profileForm = this.fb.group({
      email: [{ value: '', disabled: true }],
      displayName: ['', [Validators.required, Validators.minLength(3)]],
    });
  }

  ngOnInit(): void {
    const currentUser = this.user();
    if (currentUser) {
      this.profileForm.patchValue({
        email: currentUser.email,
        displayName: currentUser.name
      });
    }
    this.loadAuditLogs();
  }

  async saveProfile() {
    if (this.profileForm.invalid) return;

    this.isSaving.set(true);
    try {
      await this.authService.updateProfile(this.profileForm.value.displayName);
      this.messageService.add({
        severity: 'success',
        summary: 'Thành công',
        detail: 'Thông tin cá nhân đã được cập nhật'
      });
      this.profileForm.markAsPristine();
    } catch (err) {
      this.messageService.add({
        severity: 'error',
        summary: 'Lỗi',
        detail: 'Không thể cập nhật thông tin'
      });
    } finally {
      this.isSaving.set(false);
    }
  }

  async loadAuditLogs() {
    this.isLoadingLogs.set(true);
    try {
      const logs = await this.authService.getMyAuditLogs();
      this.auditLogs.set(logs);
    } catch (err) {
      console.error('Failed to load audit logs', err);
    } finally {
      this.isLoadingLogs.set(false);
    }
  }

  getActionSeverity(action: string): "success" | "info" | "warn" | "danger" | "secondary" | "contrast" | undefined {
    if (action.includes('create') || action.includes('grant')) return 'success';
    if (action.includes('update') || action.includes('edit')) return 'warn';
    if (action.includes('delete') || action.includes('revoke')) return 'danger';
    return 'info';
  }
}
