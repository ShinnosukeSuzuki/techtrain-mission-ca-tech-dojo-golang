package common

import "github.com/google/uuid"

// tokenをUUIDとして生成する
func GenerateToken() string {
	token := uuid.New()
	return token.String()
}
