package views

import (
	"findai/src/apps/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type PromptViews struct {
	PromptModel *models.PromptModel
}

func NewPromptViews(db *sqlx.DB) *PromptViews {
	return &PromptViews{PromptModel: models.NewPromptModel(db)}
}

func PromptGroup(router *gin.Engine, db *sqlx.DB) {
	g := router.Group("prompts")
	v := NewPromptViews(db)

	g.POST("", v.CreatePrompt)
}

func (v *PromptViews) CreatePrompt(c *gin.Context) {
	prompt, err := v.PromptModel.InsertPrompt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prompt"})
		return
	}
	c.JSON(http.StatusCreated, prompt)
}