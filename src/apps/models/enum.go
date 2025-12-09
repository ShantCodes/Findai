package models

import (
	"database/sql/driver"
	"fmt"
)

type (
	ContentType string
	AiModelType string
)

const (
	ContentTypeSE        ContentType = "SOFTWARE ENGINEERING"
	ContentTypeEE        ContentType = "ELECTRICAL ENGINEERING"
	ContentTypePhysics   ContentType = "PHYSICS"
	ContentTypeDesigning ContentType = "DESIGNING"

	AiModelTypeChatGPT  AiModelType = "CHATGPT"
	AiModelTypeGemini   AiModelType = "GEMINI"
	AiModelTypeGrok     AiModelType = "GROK"
	AiModelTypeDeepSeek AiModelType = "DEEPSEEK"
)

//================================================

func (v *ContentType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(v))
}
func (v ContentType) Value() (driver.Value, error) {
	return string(v), nil
}

//================================================

func (v *AiModelType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(v))
}
func (v AiModelType) Value() (driver.Value, error) {
	return string(v), nil
}

//================================================

// scanEnum is a helper function that converts an interface{} value to a string
// to support database scanning. It handles both byte slices and string values.
func scanEnum(value interface{}, target interface{}) error {
	switch v := value.(type) {
	case []byte:
		*target.(*string) = string(v) // Convert byte slice to string.
	case string:
		*target.(*string) = v // Assign string value.
	default:
		return fmt.Errorf("failed to scan type: %v", value) // Error on unsupported type.
	}
	return nil
}
