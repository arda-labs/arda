import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, switchMap } from 'rxjs';

export interface MediaMetadata {
  id: string;
  filename: string;
  content_type: string;
  size_bytes: number;
  owner_id: string;
  module: string;
  status: string;
}

export interface InitUploadResponse {
  media: MediaMetadata;
  upload_url: string;
  expires_at: string;
}

@Injectable({
  providedIn: 'root'
})
export class MediaService {
  private http = inject(HttpClient);
  private readonly API_URL = '/api/v1/media';

  /**
   * 1. Khởi tạo upload: Lấy Pre-signed URL từ Media Service
   */
  initUpload(file: File, module: string, ownerId: string): Observable<InitUploadResponse> {
    const body = {
      filename: file.name,
      content_type: file.type,
      size_bytes: file.size,
      owner_id: ownerId,
      module: module
    };
    return this.http.post<InitUploadResponse>(`${this.API_URL}/upload/init`, body);
  }

  /**
   * 2. Upload trực tiếp lên S3/SeaweedFS bằng Pre-signed URL
   */
  uploadToS3(uploadUrl: string, file: File): Observable<any> {
    return this.http.put(uploadUrl, file, {
      headers: { 'Content-Type': file.type }
    });
  }

  /**
   * 3. Xác nhận upload thành công với Media Service
   */
  confirmUpload(mediaId: string): Observable<MediaMetadata> {
    return this.http.post<MediaMetadata>(`${this.API_URL}/${mediaId}/confirm`, {});
  }

  /**
   * Quy trình hoàn chỉnh: Init -> Upload -> Confirm
   */
  uploadFile(file: File, module: string, ownerId: string): Observable<MediaMetadata> {
    return this.initUpload(file, module, ownerId).pipe(
      switchMap(resp => this.uploadToS3(resp.upload_url, file).pipe(
        switchMap(() => this.confirmUpload(resp.media.id))
      ))
    );
  }

  getDownloadUrl(mediaId: string): Observable<{ download_url: string }> {
    return this.http.get<{ download_url: string }>(`${this.API_URL}/${mediaId}/url`);
  }
}
