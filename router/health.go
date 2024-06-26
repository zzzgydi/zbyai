package router

import (
	"github.com/gin-gonic/gin"
)

func HealthRouter(r *gin.Engine) {
	health := r.Group("/__internal__")

	health.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})
}
