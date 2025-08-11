SELECT
  slug
FROM
  loglines
WHERE
  slug ~ ?0
  AND user_id = ?1
ORDER BY
  created_at DESC;
