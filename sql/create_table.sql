CREATE TABLE IF NOT EXISTS diary (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    embedding vector(3072) NOT NULL -- openai text embedding (large)
);