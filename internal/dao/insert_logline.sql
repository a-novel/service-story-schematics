INSERT INTO
  loglines (
    id,
    user_id,
    slug,
    name,
    content,
    lang,
    created_at
  )
VALUES
  (?0, ?1, ?2, ?3, ?4, ?5, ?6)
RETURNING
  *;
