package common

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// tokenをUUIDとして生成する
func GenerateToken() string {
	token := uuid.New()
	return token.String()
}

// contextにtokenを保存するための型
type TokenType struct{}

func SetToken(r *http.Request, token string) *http.Request {
	// contextにtokenを保存する
	ctx := r.Context()
	ctx = context.WithValue(ctx, TokenType{}, token)
	return r.WithContext(ctx)
}

func GetToken(r *http.Request) string {
	// contextからtokenを取得する
	ctx := r.Context()
	token, _ := ctx.Value(TokenType{}).(string)
	return token
}
