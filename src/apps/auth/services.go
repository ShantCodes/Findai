package auth

import (
	"findai/src/apps/models"

	"github.com/jmoiron/sqlx"
)

type AuthService struct {
	Db *sqlx.DB
}

func (s *AuthService) RegisterUser(username, email, password string) (*models.User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	query := `INSERT INTO users (username, email, password, is_active) VALUES ($1, $2, $3, $4) RETURNING id`
	err = s.Db.QueryRow(query, user.Username, user.Email, user.Password, true).Scan(&user.Id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) LoginUser(email, password string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, password FROM users WHERE email=$1`
	err := s.Db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	err = CheckPasswordHash(password, user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
