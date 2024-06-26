package common

import (
	"fmt"
	"time"

	slogGorm "github.com/zzzgydi/slog-gorm"
	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
	L "github.com/zzzgydi/zbyai/common/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MDB *gorm.DB

func initMysql() error {
	conf := &config.AppConf.Mysql
	if conf.Host == "" || conf.Port == "" || conf.User == "" || conf.Password == "" || conf.DbName == "" {
		return fmt.Errorf("mysql conf error")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DbName)

	slowThreshold := time.Duration(conf.SlowThreshold) * time.Millisecond
	if conf.SlowThreshold <= 0 {
		slowThreshold = 200 * time.Millisecond
	}

	gormLogger := slogGorm.New(
		slogGorm.WithHandler(L.Logger.Handler()),
		slogGorm.WithParameterizedQueries(true),
		slogGorm.WithSlowThreshold(slowThreshold),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return err
	}

	MDB = db
	return nil
}

func init() {
	initializer.Register("mysql", initMysql)
}
