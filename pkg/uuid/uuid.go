package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

// UUIDv7を生成する関数
func GenerateUUID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID v7: %w", err)
	}
	return id.String(), nil
}
