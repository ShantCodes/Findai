package views

import (
	"net/http"
	"findai/src/apps/auth"
	"findai/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthViews struct {
	Db *sqlx.DB
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (v *AuthViews) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	query := `INSERT INTO users (username, email, password, is_active) VALUES ($1, $2, $3, $4) RETURNING id`
	err = v.Db.QueryRow(query, user.Username, user.Email, user.Password, true).Scan(&user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (v *AuthViews) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	query := `SELECT id, email, password FROM users WHERE email=$1`
	err := v.Db.Get(&user, query, req.Email)
	if err != nil {
	 c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = auth.CheckPasswordHash(req.Password, user.Password)
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
