ALTER TABLE story_plans
ADD COLUMN lang text NOT NULL DEFAULT 'en';

ALTER TABLE loglines
ADD COLUMN lang text NOT NULL DEFAULT 'en';

ALTER TABLE beats_sheets
ADD COLUMN lang text NOT NULL DEFAULT 'en';

CREATE OR REPLACE VIEW story_plans_active_view AS
(
  SELECT DISTINCT ON (slug) *
  FROM story_plans
  ORDER BY slug ASC, created_at DESC
);
