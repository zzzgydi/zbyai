package auth

import (
	"errors"
	"time"

	"github.com/zzzgydi/zbyai/model"
)

type AuthClaims struct {
	Id        string         `json:"id"`
	Name      string         `json:"name,omitempty"`
	AuthType  model.AuthType `json:"auth"`
	ExpiresAt int64          `json:"exp,omitempty"`
}

func (a AuthClaims) Valid() error {
	if a.Id == "" {
		return errors.New("invalid jwt")
	}
	if a.ExpiresAt < time.Now().Unix() {
		return errors.New("jwt expired")
	}
	return nil
}

// 用于解析supabase的jwt token
// 忽略其他
type SupaUserClaims struct {
	Sub          string                 `json:"sub"`
	Email        string                 `json:"email"`
	Exp          int64                  `json:"exp"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

func (s SupaUserClaims) Valid() error {
	if s.Sub == "" || s.Email == "" {
		return errors.New("invalid jwt")
	}
	if s.Exp < time.Now().Unix() {
		return errors.New("jwt expired")
	}
	return nil
}
