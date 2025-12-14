-- name: prompts/get_by_user
SELECT
    p.*,
    count(*) OVER() AS total_count
FROM prompts p
WHERE
    p.user_id = $1
ORDER BY p.created_at DESC
LIMIT $2
OFFSET $3;
