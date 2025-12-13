INSERT INTO prompts (prompt, user_id, category, rater_score, ai_model)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at;
