package router

import "github.com/gin-gonic/gin"

func RootRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "https://www.zbyai.com")
	})
}
