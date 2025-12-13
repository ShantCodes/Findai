package models

// import "github.com/google/uuid"

type LoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}
type RegisterForm struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PromptForm struct {	
	Prompt     string      `json:"prompt" validate:"required"`
	UserId     string   `json:"userid" validate:"required"`
	Category   ContentType `json:"category" validate:"category"`
	RaterScore uint8       `json:"rater_score" validate:"rater_score"`
	AiModel    AiModelType `json:"ai_model" validate:"ai_model"`
}
