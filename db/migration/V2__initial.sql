CREATE TABLE IF NOT EXISTS public.url
(
    id text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    redirect_to text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    user_id uuid NOT NULL,
    CONSTRAINT url_pkey PRIMARY KEY (id),
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES public."user" (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE public.url
    OWNER to "db-user";

COMMENT ON CONSTRAINT user_id ON public.url
    IS 'linking the short link to the owner';