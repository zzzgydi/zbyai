package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/zzzgydi/zbyai/model"
)

func GetUserFromJWT(jwtStr string) (*model.User, error) {
	claims := &AuthClaims{}

	token, err := jwt.ParseWithClaims(jwtStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims = token.Claims.(*AuthClaims)
	user := &model.User{
		Id:       claims.Id,
		Name:     claims.Name,
		AuthType: claims.AuthType,
	}
	return user, nil
}

func SignNewJWT(user *model.User) (string, error) {
	claims := &AuthClaims{
		Id:        user.Id,
		Name:      user.Name,
		AuthType:  user.AuthType,
		ExpiresAt: time.Now().Add(7 * time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
