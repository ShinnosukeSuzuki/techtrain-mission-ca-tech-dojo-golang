package uuid

import "github.com/google/uuid"

// UUIDを生成する関数
func GenerateUUID() string {
	token := uuid.New()
	return token.String()
}
