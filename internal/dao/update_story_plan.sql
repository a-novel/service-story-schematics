INSERT INTO story_plans (
  id,
  slug,
  name,
  description,
  beats,
  lang,
  created_at
) VALUES (
  ?0,
  ?1,
  ?2,
  ?3,
  ?4,
  ?5,
  ?6
)
RETURNING *;
