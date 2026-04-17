CREATE TABLE images (
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    original_path TEXT,
    processed_path TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
