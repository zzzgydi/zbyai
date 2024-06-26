package router

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zzzgydi/zbyai/common/config"
	L "github.com/zzzgydi/zbyai/common/logger"
)

func InitHttpServer() {
	r := gin.New()
	r.Use(gin.Recovery())

	// cors
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://zbyai.com" ||
				origin == "https://www.zbyai.com" ||
				strings.HasPrefix(origin, "http://127.0.0.1") ||
				strings.HasPrefix(origin, "http://localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	// register routers
	RootRouter(r)
	HealthRouter(r)
	InnerRouter(r)

	logger := slog.NewLogLogger(L.Handler, slog.LevelError)
	srv := &http.Server{
		Addr:     ":" + strconv.FormatInt(int64(config.AppConf.HttpPort), 10),
		Handler:  r,
		ErrorLog: logger,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
