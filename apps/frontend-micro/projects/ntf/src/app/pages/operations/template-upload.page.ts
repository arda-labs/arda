import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { ProgressBar } from 'primeng/progressbar';
import { Tag } from 'primeng/tag';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { TemplateUploadItem } from '../../models/media.models';
import { MediaService } from '../../services/media.service';

@Component({
  selector: 'app-template-upload-page',
  imports: [CommonModule, ProgressBar, Tag, Toast],
  providers: [MessageService],
  templateUrl: './template-upload.page.html',
  styleUrl: './template-upload.page.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TemplateUploadPage {
  private mediaService = inject(MediaService);
  private messageService = inject(MessageService);

  readonly ownerId = signal('template-admin');
  readonly isDragging = signal(false);
  readonly queue = signal<TemplateUploadItem[]>([]);
  readonly selectedType = signal('DOCX');

  readonly acceptedTypes = ['DOCX', 'XLSX', 'PDF', 'HTML', 'JSON'];
  readonly readyCount = computed(() => this.queue().filter(item => item.status === 'ready').length);
  readonly uploadingCount = computed(() => this.queue().filter(item => item.status === 'uploading').length);
  readonly totalBytes = computed(() => this.queue().reduce((sum, item) => sum + item.size, 0));
  readonly hasQueue = computed(() => this.queue().length > 0);

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    this.isDragging.set(true);
  }

  onDragLeave(): void {
    this.isDragging.set(false);
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    this.isDragging.set(false);
    this.addFiles(event.dataTransfer?.files);
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.addFiles(input.files);
    input.value = '';
  }

  uploadAll(): void {
    for (const item of this.queue().filter(row => row.status === 'waiting' || row.status === 'failed')) {
      this.uploadItem(item.id);
    }
  }

  removeItem(id: string): void {
    this.queue.update(items => items.filter(item => item.id !== id));
  }

  clearReady(): void {
    this.queue.update(items => items.filter(item => item.status !== 'ready'));
  }

  formatBytes(bytes: number): string {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
  }

  severity(status: TemplateUploadItem['status']): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    switch (status) {
      case 'ready': return 'success';
      case 'uploading': return 'info';
      case 'failed': return 'danger';
      default: return 'secondary';
    }
  }

  private addFiles(files: FileList | null | undefined): void {
    if (!files?.length) return;
    const items = Array.from(files).map(file => ({
      id: crypto.randomUUID(),
      name: file.name,
      type: this.detectType(file),
      size: file.size,
      module: 'template',
      status: 'waiting' as const,
      progress: 0,
      file,
    }));
    this.queue.update(current => [...items.map(({ file: _, ...item }) => item), ...current]);
    for (const item of items) {
      this.uploadFile(item.id, item.file);
    }
  }

  private uploadItem(id: string): void {
    const input = document.getElementById('template-file-input') as HTMLInputElement | null;
    input?.click();
    this.messageService.add({ severity: 'info', summary: 'Chọn lại file', detail: 'Trình duyệt không giữ binary sau khi retry. Hãy chọn lại file để upload.' });
  }

  private uploadFile(id: string, file: File): void {
    this.patchItem(id, { status: 'uploading', progress: 3, error: undefined });
    this.mediaService.uploadTemplate(file, this.ownerId()).subscribe({
      next: event => this.patchItem(id, { progress: event.progress, status: event.media ? 'ready' : 'uploading', media: event.media }),
      error: () => {
        this.patchItem(id, { status: 'failed', error: 'Không thể upload hoặc xác nhận file' });
        this.messageService.add({ severity: 'error', summary: 'Upload thất bại', detail: file.name });
      },
      complete: () => this.messageService.add({ severity: 'success', summary: 'Template sẵn sàng', detail: file.name }),
    });
  }

  private patchItem(id: string, patch: Partial<TemplateUploadItem>): void {
    this.queue.update(items => items.map(item => item.id === id ? { ...item, ...patch } : item));
  }

  private detectType(file: File): string {
    const extension = file.name.split('.').pop()?.toUpperCase() ?? '';
    return this.acceptedTypes.includes(extension) ? extension : 'FILE';
  }
}
