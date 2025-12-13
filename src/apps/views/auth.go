package views

import (
	"findai/src/apps/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthViews struct {
	AuthService *auth.AuthService
}

func NewAuthViews(db *sqlx.DB) *AuthViews {
	authService := &auth.AuthService{Db: db}
	return &AuthViews{AuthService: authService}
}

func AuthGroup(router *gin.Engine, db *sqlx.DB) {
	g := router.Group("auth")
	v := NewAuthViews(db)

	g.POST("/register", v.Register)
	g.POST("/login", v.Login)
}

func (v *AuthViews) Register(c *gin.Context) {
	user, err := v.AuthService.RegisterUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (v *AuthViews) Login(c *gin.Context) {
	user, err := v.AuthService.LoginUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokens, err := auth.GenerateFullTokens(user.Id.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
