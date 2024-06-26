package auth

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/zzzgydi/zbyai/common"
	"github.com/zzzgydi/zbyai/common/logger"
	"github.com/zzzgydi/zbyai/model"
	"gorm.io/gorm"
)

// 可以直接通过jwt解析得到一些用户信息，不需要调用supabase的api
func VerifySupabaseJwt(accessToken string) (*SupaUserClaims, error) {
	claims := &SupaUserClaims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return supaJwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims = token.Claims.(*SupaUserClaims)
	return claims, nil
}

// get or register
func GetUserFromSupabase(accessToken, ip string) (*model.User, error) {
	var id string
	var email string
	var name string
	var meta map[string]interface{}

	// 先尝试这个
	supaClaims, err := VerifySupabaseJwt(accessToken)
	if err != nil {
		logger.Logger.Error("verify supabase jwt error", "error", err)

		supaUser, err := common.VerifyToken(accessToken)
		if err != nil {
			return nil, err
		}

		id = supaUser.ID
		email = supaUser.Email
		name = supaUser.Email
		meta = supaUser.UserMetadata
	} else {
		id = supaClaims.Sub
		email = supaClaims.Email
		name = supaClaims.Email
		meta = supaClaims.UserMetadata
	}

	if metaName, ok := meta["name"].(string); ok {
		name = metaName
	}

	user, err := model.GetUserById(id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// register new user
		user = model.NewSupabaseUser(id, name, email, ip)
		if err := user.Create(); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func GetUserFromTourist(ip string) (*model.User, error) {
	user := model.NewTouristUser(ip)
	if err := user.Create(); err != nil {
		return nil, err
	}
	return user, nil
}
