import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { Button } from 'primeng/button';

@Component({
  selector: 'app-not-found',
  standalone: true,
  imports: [RouterLink, Button],
  template: `
    <div class="min-h-screen flex flex-col items-center justify-center px-4 bg-surface-50 dark:bg-surface-950 text-center">
      <div class="w-24 h-24 rounded-full bg-surface-100 dark:bg-surface-800 flex items-center justify-center mb-6">
        <i class="pi pi-map-marker text-4xl text-surface-400"></i>
      </div>
      <h1 class="text-6xl font-bold text-primary-500 mb-2">404</h1>
      <h2 class="text-2xl font-semibold text-surface-900 dark:text-surface-0 mb-4">Không tìm thấy trang</h2>
      <p class="text-surface-500 dark:text-surface-400 mb-8 max-w-md">
        Đường dẫn bạn truy cập không tồn tại hoặc đã bị di chuyển.
      </p>
      <div class="flex gap-4">
        <p-button label="Về trang chủ" icon="pi pi-home" routerLink="/" />
        <p-button label="Quay lại" icon="pi pi-arrow-left" [variant]="'text'" (click)="goBack()" />
      </div>
    </div>
  `,
})
export class NotFoundPage {
  goBack(): void {
    window.history.back();
  }
}
