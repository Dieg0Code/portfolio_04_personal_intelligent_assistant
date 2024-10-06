create table if not exists user_message (
    id uuid primary key default uuid_generate_v4(),
    message_content text not null,
    sender_location text not null,
    created_at timestamp not null default now(),
    embedding vector(3072) not null -- openai text embedding (large)
);