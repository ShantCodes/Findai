-- name: prompts/get_filtered
SELECT
    p.*,
    count(*) OVER() AS total_count
FROM prompts p
WHERE
    (CASE WHEN $1 = '' THEN TRUE ELSE p.category = CAST($1 AS content_type) END)
AND
    (CASE WHEN $2 = '' THEN TRUE ELSE p.prompt ILIKE '%' || $2 || '%' END)
ORDER BY p.created_at DESC
LIMIT $3
OFFSET $4;
