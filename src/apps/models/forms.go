package models

type LoginForm struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required,min=6,max=30"`
}
type RegisterForm struct {
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`	
}
