package auth

import (
	"findai/src/apps/models"
	"findai/src/apps/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthService struct {
	Db *sqlx.DB
}

func (s *AuthService) RegisterUser(c *gin.Context) (*models.User, error) {
	form := new(models.RegisterForm)
	if err := c.ShouldBindJSON(&form); err != nil {
		return nil, err
	}

	var user models.User
	if err := utils.Copy(&form, &user); err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(form.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	row, err := utils.QuerySelect(c.Request.Context(), "register",
		user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) LoginUser(c *gin.Context) (*models.User, error) {
	var form models.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		return nil, err
	}

	var user models.User
	row, err := utils.QuerySelect(c.Request.Context(), "login", form.Email)
	if err != nil {
		return nil, err
	}
	err = row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	err = CheckPasswordHash(form.Password, user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
