WITH active_story_plans AS (
  SELECT DISTINCT ON (slug) *
  FROM story_plans
  ORDER BY slug ASC, created_at DESC
)
SELECT *
FROM active_story_plans
WHERE slug = ?0;
