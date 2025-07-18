SELECT slug
FROM story_plans
WHERE slug ~ ?0
ORDER BY created_at DESC;
