import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { Button } from 'primeng/button';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-access-denied',
  standalone: true,
  imports: [RouterLink, Button],
  template: `
    <div class="min-h-screen flex flex-col items-center justify-center px-4 bg-surface-50 dark:bg-surface-950 text-center">
      <div class="w-24 h-24 rounded-full bg-red-100 dark:bg-red-900/20 flex items-center justify-center mb-6">
        <i class="pi pi-lock text-4xl text-red-500"></i>
      </div>
      <h1 class="text-6xl font-bold text-red-500 mb-2">403</h1>
      <h2 class="text-2xl font-semibold text-surface-900 dark:text-surface-0 mb-4">Truy cập bị từ chối</h2>
      <p class="text-surface-500 dark:text-surface-400 mb-8 max-w-md">
        Bạn không có quyền truy cập vào tài nguyên này. Vui lòng liên hệ quản trị viên.
      </p>
      <div class="flex gap-4">
        <p-button label="Về trang chủ" icon="pi pi-home" routerLink="/" severity="danger" />
        <p-button label="Đăng xuất" icon="pi pi-sign-out" [variant]="'text'" severity="secondary" (click)="logout()" />
      </div>
    </div>
  `,
})
export class AccessDeniedPage {
  private authService = inject(AuthService);

  logout(): void {
    this.authService.logout();
  }
}
