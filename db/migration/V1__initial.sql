CREATE TABLE IF NOT EXISTS public.user(
    id uuid,
    email text,
    password text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    PRIMARY KEY (id),
    CONSTRAINT email UNIQUE (email)
);

ALTER TABLE public.user
    OWNER to "db-user";