INSERT INTO beats_sheets (
  id,
  logline_id,
  story_plan_id,
  content,
  lang,
  created_at
) VALUES (
  ?0,
  ?1,
  ?2,
  ?3,
  ?4,
  ?5
)
RETURNING *;
