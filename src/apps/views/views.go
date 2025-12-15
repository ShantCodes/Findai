package views

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Init(r *gin.Engine, db *sqlx.DB) {
	AuthGroup(r, db)
	PromptGroup(r, db)
	RoleGroup(r, db)
}
