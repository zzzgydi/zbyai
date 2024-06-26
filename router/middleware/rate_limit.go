package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/model"
	"github.com/zzzgydi/zbyai/router/controller"
	"github.com/zzzgydi/zbyai/router/utils"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	ipLimiter := utils.NewRateLimiter("z:rl:ip", 20, 60*60)     // 1小时 最多15次
	userLimiter := utils.NewRateLimiter("z:rl:user", 60, 60*60) // 1小时 最多60次

	return func(c *gin.Context) {
		// 先判断是否为登录用户
		// 如果是登录用户，则使用用户ID作为限流key
		user := utils.GetUser(c)
		ip := utils.GetRealUserIp(c)

		var allow bool
		var err error

		// 仅正式登录的用户用用户ID作为限流key
		// 游客用ip限流
		if user != nil && user.AuthType != model.AUTH_NONE {
			allow, err = userLimiter.Allow(c, user.Id)
		} else {
			allow, err = ipLimiter.Allow(c, ip)
		}

		if err != nil {
			controller.ReturnServerError(c, err)
			return
		}

		if trace := utils.GetTraceLogger(c); trace != nil {
			trace.Trace("rl_allow", allow)
		}

		if !allow {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
