create index on public.diary
using ivfflat (embedding vector_cosine_ops)
with (lists = 100);