SELECT
  slug,
  name,
  content,
  lang,
  created_at
FROM loglines
WHERE user_id = ?0
ORDER BY created_at DESC, name DESC, slug DESC
LIMIT ?1 OFFSET ?2;
