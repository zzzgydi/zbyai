package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/common"
	"github.com/zzzgydi/zbyai/model"
	"github.com/zzzgydi/zbyai/router/utils"
	"github.com/zzzgydi/zbyai/service/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		trace := utils.GetTraceLogger(c)
		logger := trace.Logger

		var user *model.User

		// 获取cookie里的jwt
		if cookie, err := c.Cookie(common.COOKIE_SESSION); err == nil {
			// 解析jwt, 判断是否过期, 判断是否有用户信息
			jwtUser, err := auth.GetUserFromJWT(cookie)
			if err != nil {
				logger.Error("get user from jwt failed", "error", err)
			} else {
				user = jwtUser
			}
		}

		ip := utils.GetRealUserIp(c)

		if user == nil || user.AuthType == model.AUTH_NONE {
			// 如果之前是游客用户，可以刷成supabase用户
			// 判断是否有supabase的 access token
			if accessToken := c.GetHeader("Authorization"); accessToken != "" {
				accessToken = strings.TrimPrefix(accessToken, "Bearer ")
				supaUser, err := auth.GetUserFromSupabase(accessToken, ip)
				if err != nil {
					logger.Error("auth supabase error", "error", err)
				} else {
					user = supaUser
				}
			}
		}

		// 没有的话就根据ip地址生成一个游客用户
		if user == nil {
			tourUser, err := auth.GetUserFromTourist(ip)
			if err != nil {
				logger.Error("auth tourist error", "error", err)
				c.Abort()
				return
			}
			user = tourUser
		}

		if user == nil {
			logger.Error("auth middleware error")
			c.Abort()
			return
		}

		// 刷新jwt
		jwtStr, err := auth.SignNewJWT(user)
		if err != nil {
			logger.Error("sign new jwt error", "error", err)
			c.Abort()
			return
		}

		if trace != nil {
			trace.SetUid(user.Id)
		}

		// 根据请求判断是否为localhost，然后设置为localhost的domain
		if strings.Contains(c.Request.Host, "localhost") {
			c.SetCookie(common.COOKIE_SESSION, jwtStr, 60*60*24*7, "/", "localhost", false, true)
		} else {
			c.SetCookie(common.COOKIE_SESSION, jwtStr, 60*60*24*7, "/", ".zbyai.com", false, true)
		}

		c.Set(common.CTX_CURRENT_USER, user)
		c.Next()
	}
}
