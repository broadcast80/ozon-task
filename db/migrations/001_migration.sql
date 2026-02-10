CREATE TABLE IF NOT EXISTS public.link (
	id serial4 NOT NULL,
	url text NOT NULL,
	alias text NOT NULL,
	created_at timestamp DEFAULT CURRENT_DATE NOT NULL,
	CONSTRAINT alias_unique UNIQUE (alias),
	CONSTRAINT user_pkey PRIMARY KEY (id)
)