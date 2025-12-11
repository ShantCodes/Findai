package models

type LoginForm struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password" validate:"required"`	
}
type RegisterForm struct {
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`	
}
