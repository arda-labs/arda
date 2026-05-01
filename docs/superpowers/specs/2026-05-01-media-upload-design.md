# Media Upload Service Design

**Goal:** Provide a scalable, secure, and efficient file upload and management system for the Arda ecosystem using SeaweedFS S3 API and the Kratos framework.

**Architecture:** The system utilizes the **S3 Presigned URL** pattern. The backend acts as a metadata manager and authorization gatekeeper, while the actual file data flows directly from the client to the object storage. A hybrid confirmation approach is used: clients explicitly confirm uploads for fast UI response, while a background worker cleans up orphaned files.

**Tech Stack:** Go (Kratos Framework), Postgres, SeaweedFS (S3 API compatibility), AWS SDK for Go v2.

---

## 1. Data Flow

1.  **Init Upload (Client -> Media Service):**
    *   Client requests to upload a file, providing metadata (filename, size, content type, module category).
    *   Media Service verifies permissions and creates a `PENDING` record in Postgres.
    *   Media Service generates a short-lived (e.g., 15 mins) PUT Presigned URL via the S3 SDK.
    *   Returns the `media_id` and Presigned URL to the client.
2.  **Direct Upload (Client -> SeaweedFS):**
    *   Client executes an HTTP PUT request directly to the Presigned URL with the file binary payload.
3.  **Confirm Upload (Client -> Media Service):**
    *   Upon successful PUT, Client calls `POST /media/{id}/confirm`.
    *   Media Service calls S3 `HeadObject` to verify file existence and size.
    *   If valid, updates the DB record status to `READY`.
4.  **Cleanup Worker (Background Task):**
    *   A periodic worker runs (e.g., every hour) to find `PENDING` records older than 24 hours.
    *   It deletes any corresponding orphaned objects in S3 and removes the DB records.

## 2. Database Schema

Table: `media_metadata`

```sql
CREATE TABLE media_metadata (
    id UUID PRIMARY KEY,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    bucket TEXT NOT NULL,
    object_key TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    module TEXT NOT NULL,       -- e.g., 'avatar', 'invoice', 'template'
    status TEXT NOT NULL,       -- 'PENDING', 'READY', 'DELETED'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_media_status_created ON media_metadata(status, created_at);
```

## 3. Directory Structure (Kratos)

Adhering to the Arda standards for Go backend services:

*   `api/media/v1/media.proto`: gRPC and HTTP gateway definitions.
*   `internal/biz/media.go`: Core business logic (`InitUpload`, `ConfirmUpload`, `CleanupOrphans`). Defines the `MediaRepo` and `StorageRepo` interfaces.
*   `internal/data/media_repo.go`: Postgres implementation using standard database/sql or pgx.
*   `internal/data/storage_repo.go`: S3 implementation using `github.com/aws/aws-sdk-go-v2`.
*   `internal/service/media.go`: Handles incoming requests, maps DTOs, and orchestrates calls to the `biz` layer.
*   `internal/server/`: Registration of HTTP/gRPC servers and background workers.

## 4. API Definitions (Proto abstraction)

```protobuf
service Media {
  rpc InitUpload (InitUploadRequest) returns (InitUploadReply) {
    option (google.api.http) = {
      post: "/v1/media/upload/init"
      body: "*"
    };
  }
  rpc ConfirmUpload (ConfirmUploadRequest) returns (ConfirmUploadReply) {
    option (google.api.http) = {
      post: "/v1/media/{id}/confirm"
      body: "*"
    };
  }
  rpc GetMediaUrl (GetMediaUrlRequest) returns (GetMediaUrlReply) {
    option (google.api.http) = {
      get: "/v1/media/{id}/url"
    };
  }
}
```

## 5. Security & Edge Cases

*   **Path Traversal Prevention:** The `object_key` must be generated entirely by the backend (e.g., `<module>/<year>/<month>/<uuid>`). Client-provided filenames are only saved in DB metadata, never used for the S3 path.
*   **Size Limits:** The Presigned URL generated should ideally enforce content-length restrictions if the S3 backend supports it, otherwise, the `ConfirmUpload` step strictly validates the size.
*   **Access Control:** The `GetMediaUrl` endpoint must verify if the requesting user has read access to the specified `media_id` based on ownership or module policies before generating a GET Presigned URL.
