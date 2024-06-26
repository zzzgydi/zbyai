package auth

import (
	"fmt"

	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
)

var (
	jwtSecret     []byte
	supaJwtSecret []byte
)

func InitAuth() error {
	jwtSecret = []byte(config.AppConf.JwtSecret)

	supaConf := &config.AppConf.Supabase
	if supaConf.JwtSecret == "" {
		return fmt.Errorf("supabase jwt secret error")
	}
	supaJwtSecret = []byte(supaConf.JwtSecret)

	return nil
}

func init() {
	initializer.Register("auth", InitAuth)
}
