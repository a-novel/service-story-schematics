SELECT
  id,
  lang,
  created_at
FROM
  beats_sheets
WHERE
  logline_id = ?0
ORDER BY
  created_at DESC
LIMIT
  ?1
OFFSET
  ?2;
