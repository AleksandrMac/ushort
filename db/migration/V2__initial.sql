CREATE TABLE IF NOT EXISTS public.url
(
    id text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
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

CREATE FUNCTION update_url_last_update() RETURNS TRIGGER AS $url_last_update$
BEGIN
    UPDATE public.url SET updated_at=now() WHERE id=NEW.id;
    RETURN NEW;
END; $url_last_update$
LANGUAGE plpgsql;

CREATE TRIGGER update_url_each_change AFTER 
UPDATE OF id, redirect_to, description ON public.url 
FOR EACH ROW EXECUTE FUNCTION update_url_last_update();