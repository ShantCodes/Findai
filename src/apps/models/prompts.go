package models

import (
	"time"

	"github.com/google/uuid"
)

type Prompt struct {
	ID         uuid.UUID   `db:"primaryKey" json:"id"`
	Prompt     string      `db:"type:text;not null" json:"prompt"`
	UserId     uuid.UUID   `db:"userid" json:"user_id"`
	Category   ContentType `db:"category" json:"category"`
	RaterScore uint8       `db:"rater_score" json:"rater_score"`
	AiModel    ContentType `db:"ai_model" json:"ai_model"`

	CreatedAt time.Time `db:"created-at" json:"created_at"`
	UpdatedAt time.Time `db:"updated-at" json:"updated_at"`
}
