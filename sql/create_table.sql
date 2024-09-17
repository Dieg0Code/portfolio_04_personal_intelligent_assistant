CREATE TABLE IF NOT EXISTS diary (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    embedding VECTOR(3072) NOT NULL -- openai text embedding (large)
);