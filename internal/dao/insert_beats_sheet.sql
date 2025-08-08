INSERT INTO beats_sheets (
  id,
  logline_id,
  content,
  lang,
  created_at
) VALUES (
  ?0,
  ?1,
  ?2,
  ?3,
  ?4
)
RETURNING *;
