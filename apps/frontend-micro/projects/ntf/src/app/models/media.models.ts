export interface MediaMetadata {
  id: string;
  filename: string;
  contentType: string;
  sizeBytes: number;
  bucket: string;
  objectKey: string;
  ownerId: string;
  module: string;
  status: 'PENDING' | 'READY' | 'DELETED' | string;
  createdAt?: string;
  updatedAt?: string;
}

export interface InitUploadRequest {
  filename: string;
  contentType: string;
  sizeBytes: number;
  ownerId: string;
  module: string;
}

export interface InitUploadResponse {
  media: MediaMetadata;
  uploadUrl: string;
  expiresAt: string;
}

export interface TemplateUploadItem {
  id: string;
  name: string;
  type: string;
  size: number;
  module: string;
  status: 'waiting' | 'uploading' | 'ready' | 'failed';
  progress: number;
  media?: MediaMetadata;
  error?: string;
}
