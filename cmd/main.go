package main

import (
	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
	"github.com/zzzgydi/zbyai/common/logger"
	"github.com/zzzgydi/zbyai/router"
)

func main() {
	env := config.GetEnv()
	rootDir := config.GetRootDir()
	logger.InitLogger(rootDir)
	config.InitConfig(rootDir, env)
	initializer.InitInitializer()
	router.InitHttpServer()
}
