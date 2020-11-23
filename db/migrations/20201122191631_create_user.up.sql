CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL UNIQUE,
    online BOOLEAN NOT NULL,
    seen timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_ref on public.users(id);
