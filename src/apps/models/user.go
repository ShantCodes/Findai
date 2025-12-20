package models

import (
	"context"
	"findai/src/apps/utils"
	"findai/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserWithCount struct {
	User
	TotalCount int `db:"total_count"`
}

type UserModel struct {
	Db *sqlx.DB
}

func NewUserModel(db *sqlx.DB) *RoleModel {
	return &RoleModel{Db: db}
}

func GetAllUsers(p database.Paginate) ([]*User, int, error) {
	db := database.DB()

	res, err := utils.QuerySelectRows(context.Background(), db, "get_all_users")
	if err != nil {
		return nil, 0, err
	}
	defer res.Close()
	users := []*User{}
	totalCount := 0

	for res.Next() {
		var result UserWithCount
		if err := res.StructScan(&result); err != nil {
			return nil, 0, err
		}
		users = append(users, &result.User)
		totalCount = result.TotalCount
	}

	return users, totalCount, nil
}
