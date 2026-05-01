import { Component, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FileUploadModule } from 'primeng/fileupload';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { MediaService } from '../../services/media.service';

@Component({
  selector: 'app-attachments',
  standalone: true,
  imports: [CommonModule, FileUploadModule, TableModule, ButtonModule, TagModule, ToastModule],
  providers: [MessageService],
  templateUrl: './attachments.html',
  styleUrl: './attachments.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AttachmentsComponent {
  private mediaService = inject(MediaService);
  private messageService = inject(MessageService);

  files = signal<any[]>([
    { id: 'm-001', name: 'CCCD_MatTruoc.jpg', type: 'Giấy tờ định danh', size: '1.2MB', uploadDate: '2026-05-01', status: 'VERIFIED' },
    { id: 'm-002', name: 'HopDongLaoDong.pdf', type: 'Hồ sơ tài chính', size: '3.5MB', uploadDate: '2026-05-01', status: 'PENDING' }
  ]);

  onUpload(event: any) {
    const uploadedFiles: File[] = event.files;

    uploadedFiles.forEach(file => {
      this.messageService.add({ severity: 'info', summary: 'Đang tải lên', detail: file.name });

      this.mediaService.uploadFile(file, 'CRM', 'temp-user').subscribe({
        next: (metadata) => {
          this.files.update(prev => [
            {
              id: metadata.id,
              name: metadata.filename,
              type: 'Tài liệu mới',
              size: (metadata.size_bytes / 1024 / 1024).toFixed(2) + 'MB',
              uploadDate: new Date().toISOString().split('T')[0],
              status: 'PENDING'
            },
            ...prev
          ]);
          this.messageService.add({ severity: 'success', summary: 'Thành công', detail: `Đã tải lên ${file.name}` });
        },
        error: (err) => {
          this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: `Không thể tải lên ${file.name}` });
        }
      });
    });
  }

  getSeverity(status: string) {
    switch (status) {
      case 'VERIFIED': return 'success';
      case 'PENDING': return 'warn';
      case 'REJECTED': return 'danger';
      default: return 'secondary';
    }
  }

  onView(file: any) {
    this.mediaService.getDownloadUrl(file.id).subscribe(res => {
      window.open(res.download_url, '_blank');
    });
  }
}
