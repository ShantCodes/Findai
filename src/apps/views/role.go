package views

import (
	// "findai/src/apps/auth"
	"findai/src/apps/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type RoleViews struct {
	RoleModel *models.RoleModel
}

func NewRoleViews(db *sqlx.DB) *RoleViews {
	return &RoleViews{RoleModel: models.NewRoleModel(db)}
}

func RoleGroup(router *gin.Engine, db *sqlx.DB) {
	g := router.Group("roles")
	v := NewRoleViews(db)

	g.GET("", v.GetRoles)
}

func (v *RoleViews) GetRoles(c *gin.Context) {
	roles, err := v.RoleModel.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}
