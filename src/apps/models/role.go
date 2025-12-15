package models

import (
	"context"
	"findai/src/apps/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Role struct {
	Id     uuid.UUID `db:"id" json:"id"`
	Name   RoleType  `db:"name" json:"name"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
}

type RoleModel struct {
	Db *sqlx.DB
}

func NewRoleModel(db *sqlx.DB) *RoleModel{
	return &RoleModel{Db: db}
}

func GetRolesByUserID(db *sqlx.DB, userID uuid.UUID) ([]Role, error) {
	var roles []Role
	rows, err := utils.QuerySelectRows(context.Background(), db, "get_roles_by_user_id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role Role
		if err := rows.StructScan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (m *RoleModel) GetRoles() ([]Role, error) {
	var roles []Role
	rows, err := utils.QuerySelectRows(context.Background(), m.Db, "get_roles")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var role Role
		if err := rows.StructScan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}
