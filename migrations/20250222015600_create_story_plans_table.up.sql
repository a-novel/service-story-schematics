CREATE TABLE story_plans
(
  id uuid PRIMARY KEY NOT NULL,
  slug text NOT NULL,

  name text NOT NULL,
  description text,

  beats jsonb NOT NULL,

  created_at timestamp(0) with time zone NOT NULL
);

CREATE INDEX story_plans_slug_timestamp_idx ON story_plans (
  slug ASC, created_at DESC
);

CREATE VIEW story_plans_active_view AS
(
  SELECT DISTINCT ON (slug) *
  FROM story_plans
  ORDER BY slug ASC, created_at DESC
);
