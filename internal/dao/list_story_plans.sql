WITH active_story_plans AS (
  SELECT DISTINCT ON (slug) *
  FROM story_plans
  ORDER BY slug ASC, created_at DESC
)

SELECT
  id,
  slug,
  name,
  description,
  lang,
  created_at
FROM active_story_plans
ORDER BY created_at DESC
LIMIT ?0 OFFSET ?1;
