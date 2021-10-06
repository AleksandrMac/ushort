CREATE TABLE IF NOT EXISTS public.user(
    id uuid,
    email text,
    password text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT email UNIQUE (email)
);
