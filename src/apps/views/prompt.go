package views

import (
	"findai/src/apps/models"
	"findai/src/apps/utils"
	"findai/src/database"
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
	g.GET("", utils.Paginate(), v.GetPrompts)
}

func (v *PromptViews) CreatePrompt(c *gin.Context) {
	prompt, err := v.PromptModel.InsertPrompt(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prompt"})
		return
	}
	c.JSON(http.StatusCreated, prompt)
}

func (v *PromptViews) GetPrompts(c *gin.Context) {
	pagination := c.MustGet("paginate").(database.Paginate)

	prompts, total, err := models.GetPrompts(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"results": prompts,
		"total":   total,
		"page":    c.MustGet("page"),
		"limit":   c.MustGet("limit"),
	})
}