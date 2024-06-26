package model

import (
	"github.com/zzzgydi/zbyai/common"
	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
)

func initMigrate() error {
	if config.AppConf.Mysql.AutoMigrate {
		// 设置引擎
		err := common.MDB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
		).AutoMigrate(
			&Thread{},
			&ThreadRun{},
			&ThreadAnswer{},
			&ThreadSearch{},
			&User{},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	initializer.Register("migrate", initMigrate)
}
