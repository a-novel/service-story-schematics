CREATE TABLE story_plans
(
  id uuid PRIMARY KEY NOT NULL,
  slug text NOT NULL,

  name text NOT NULL,
  description text,

  beats jsonb NOT NULL,

  created_at timestamp(0) with time zone NOT NULL
);

CREATE INDEX story_plans_slug_idx ON story_plans (slug);
