CREATE TABLE IF NOT EXISTS public.link (
	id serial4 NOT NULL,
	url text NOT NULL,
	alias text NOT NULL,
	created_at timestamp DEFAULT CURRENT_DATE NOT NULL,
	CONSTRAINT url_unique UNIQUE (url),
	CONSTRAINT user_pkey PRIMARY KEY (id)
)