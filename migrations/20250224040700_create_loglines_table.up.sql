CREATE TABLE loglines
(
  id uuid PRIMARY KEY NOT NULL,
  user_id uuid NOT NULL,
  slug text NOT NULL,

  name text NOT NULL,
  content text NOT NULL,

  CONSTRAINT unique_slug_per_user UNIQUE (user_id, slug),

  created_at timestamp(6) with time zone NOT NULL
);

CREATE INDEX loglines_user_id_idx ON loglines (user_id);

CREATE INDEX loglines_slug_idx ON loglines (slug);
