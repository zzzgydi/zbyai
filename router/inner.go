package router

import (
	"github.com/gin-gonic/gin"
	ctl "github.com/zzzgydi/zbyai/router/controller"
	"github.com/zzzgydi/zbyai/router/middleware"
)

func InnerRouter(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.Use(middleware.LoggerMiddleware, middleware.AuthMiddleware())

	v1.POST("/inner/auth", ctl.PostAuth)
	v1.POST("/inner/create_thread", middleware.RateLimiterMiddleware(), ctl.PostThreadCreate)
	v1.POST("/inner/append_thread", middleware.RateLimiterMiddleware(), ctl.PostThreadAppend)
	v1.POST("/inner/rewrite_thread", middleware.RateLimiterMiddleware(), ctl.PostThreadRewrite)
	v1.POST("/inner/stream_thread", ctl.PostThreadStream)
	v1.POST("/inner/detail_thread", ctl.PostThreadDetail)
	v1.POST("/inner/delete_thread", ctl.PostThreadDelete)
	v1.POST("/inner/list_thread", ctl.PostListThread)
}
