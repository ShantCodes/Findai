package models

import (
	"context"
	"findai/src/apps/utils"
	"findai/src/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Prompt struct {
	ID         uuid.UUID   `db:"id" json:"id"`
	Prompt     string      `db:"prompt" json:"prompt"`
	UserId     uuid.UUID   `db:"user_id" json:"userid"`
	Category   ContentType `db:"category" json:"category"`
	RaterScore uint8       `db:"rater_score" json:"rater_score"`
	AiModel    AiModelType `db:"ai_model" json:"ai_model"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type PromptModel struct {
	Db *sqlx.DB
}

func NewPromptModel(db *sqlx.DB) *PromptModel {
	return &PromptModel{Db: db}
}

func (m *PromptModel) InsertPrompt(ctx *gin.Context) (*Prompt, error) {
	form := new(PromptForm)
	if err := ctx.ShouldBindJSON(&form); err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(form.UserId)
	if err != nil {
		return nil, err
	}
	p := new(Prompt)
	if err := utils.Copy(form, p); err != nil {
		return nil, err
	}
	p.UserId = userID

	if p.RaterScore > 5 {
		p.RaterScore = 5
	}
	row := utils.QuerySelect(ctx.Request.Context(), m.Db, "insert_prompt",
		p.Prompt,
		p.UserId,
		p.Category,
		p.RaterScore,
		p.AiModel)

	err = row.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPrompts(p database.Paginate) ([]*Prompt, int, error) {
	var userID, category, search string
	for _, filter := range p.Filters {
		switch filter.Key {
		case "user_id":
			userID = filter.Value
		case "category":
			category = filter.Value
		case "q":
			search = filter.Value
		}
	}

	var rows *sqlx.Rows
	var err error
	db := database.DB()

	if userID != "" {
		rows, err = utils.QuerySelectRows(context.Background(), db, "get_prompts_by_user", userID, p.Limit, p.Offset)
	} else {
		rows, err = utils.QuerySelectRows(context.Background(), db, "get_prompts", category, search, p.Limit, p.Offset)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	prompts := []*Prompt{}
	totalCount := 0

	type PromptWithCount struct {
		Prompt
		TotalCount int `db:"total_count"`
	}

	for rows.Next() {
		var result PromptWithCount
		if err := rows.StructScan(&result); err != nil {
			return nil, 0, err
		}
		prompts = append(prompts, &result.Prompt)
		totalCount = result.TotalCount
	}

	return prompts, totalCount, nil
}
