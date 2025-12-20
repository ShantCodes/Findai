package views

import (
	"findai/src/apps/auth"
	"findai/src/apps/models"
	"findai/src/apps/utils"
	"findai/src/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UserViews struct {
	AuthService *auth.AuthService
}

func NewUserViews(db *sqlx.DB) *UserViews {
	return &UserViews{AuthService: (*auth.AuthService)(models.NewUserModel(db))}
}

func UserGroup(router *gin.Engine, db *sqlx.DB) {
	g := router.Group("users")
	v := NewUserViews(db)

	g.GET("/all", auth.LoginRequired(), auth.AdminOnly(), utils.Paginate(), v.GetAllUsers)
}

func (v *UserViews) GetAllUsers(c *gin.Context) {
	pagination := c.MustGet("paginate").(database.Paginate)

	users, total, err := models.GetAllUsers(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"results": users,
		"total":   total,
		"page":    c.MustGet("page"),
		"limit":   c.MustGet("limit"),
	})
}
