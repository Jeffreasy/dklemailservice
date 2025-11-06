-- GECONSOLIDEERDE V10 - UPLOADED IMAGES
-- Dit is de logica van V1_30.

-- Create uploaded_images table
CREATE TABLE IF NOT EXISTS uploaded_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    public_id TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    secure_url TEXT NOT NULL,
    filename TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    folder TEXT NOT NULL,
    thumbnail_url TEXT,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_uploaded_images_user_id ON uploaded_images(user_id);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_public_id ON uploaded_images(public_id);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_folder ON uploaded_images(folder);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_deleted_at ON uploaded_images(deleted_at);
CREATE INDEX IF NOT EXISTS idx_uploaded_images_created_at ON uploaded_images(created_at DESC);