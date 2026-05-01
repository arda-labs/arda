import { HttpClient, HttpEvent, HttpEventType, HttpHeaders, HttpRequest } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, filter, map, switchMap } from 'rxjs';
import { InitUploadRequest, InitUploadResponse, MediaMetadata } from '../models/media.models';

@Injectable({ providedIn: 'root' })
export class MediaService {
  private http = inject(HttpClient);
  private readonly baseUrl = '/api/v1/media';

  initUpload(request: InitUploadRequest): Observable<InitUploadResponse> {
    return this.http.post<any>(`${this.baseUrl}/upload/init`, {
      filename: request.filename,
      content_type: request.contentType,
      size_bytes: request.sizeBytes,
      owner_id: request.ownerId,
      module: request.module,
    }).pipe(map(resp => ({ media: this.toMedia(resp.media), uploadUrl: resp.upload_url ?? resp.uploadUrl ?? '', expiresAt: resp.expires_at ?? resp.expiresAt ?? '' })));
  }

  uploadToPresignedUrl(url: string, file: File): Observable<number> {
    const req = new HttpRequest('PUT', url, file, {
      reportProgress: true,
      headers: new HttpHeaders({ 'Content-Type': file.type || 'application/octet-stream' }),
    });
    return this.http.request(req).pipe(
      filter(event => event.type === HttpEventType.UploadProgress || event.type === HttpEventType.Response),
      map((event: HttpEvent<unknown>) => {
        if (event.type === HttpEventType.Response) return 100;
        if (event.type !== HttpEventType.UploadProgress) return 0;
        const total = event.total ?? file.size;
        return total > 0 ? Math.round((event.loaded / total) * 100) : 0;
      }),
    );
  }

  confirmUpload(id: string): Observable<MediaMetadata> {
    return this.http.post<any>(`${this.baseUrl}/${encodeURIComponent(id)}/confirm`, {}).pipe(map(resp => this.toMedia(resp)));
  }

  uploadTemplate(file: File, ownerId: string): Observable<{ progress: number; media?: MediaMetadata }> {
    return this.initUpload({ filename: file.name, contentType: file.type || 'application/octet-stream', sizeBytes: file.size, ownerId, module: 'template' }).pipe(
      switchMap(init => this.uploadToPresignedUrl(init.uploadUrl, file).pipe(
        switchMap(progress => progress === 100 ? this.confirmUpload(init.media.id).pipe(map(media => ({ progress, media }))) : [ { progress } ]),
      )),
    );
  }

  private toMedia(item: any): MediaMetadata {
    return {
      id: item?.id ?? '',
      filename: item?.filename ?? '',
      contentType: item?.content_type ?? item?.contentType ?? '',
      sizeBytes: item?.size_bytes ?? item?.sizeBytes ?? 0,
      bucket: item?.bucket ?? '',
      objectKey: item?.object_key ?? item?.objectKey ?? '',
      ownerId: item?.owner_id ?? item?.ownerId ?? '',
      module: item?.module ?? '',
      status: item?.status ?? 'PENDING',
      createdAt: item?.created_at ?? item?.createdAt,
      updatedAt: item?.updated_at ?? item?.updatedAt,
    };
  }
}
