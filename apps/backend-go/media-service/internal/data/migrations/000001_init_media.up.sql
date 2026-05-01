CREATE TABLE media_metadata (
    id UUID PRIMARY KEY,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    bucket TEXT NOT NULL,
    object_key TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    module TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_media_metadata_status_created ON media_metadata(status, created_at);
CREATE INDEX idx_media_metadata_owner_id ON media_metadata(owner_id);
