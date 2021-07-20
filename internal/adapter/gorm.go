package adapter

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/raulcarlin/go-backend/internal/config"
)

func NewGORM(conf *config.Config) (*gorm.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.DBConfig.Host, conf.DBConfig.Port),
		DBName:               conf.DBConfig.DbName,
		User:                 conf.DBConfig.Username,
		Passwd:               conf.DBConfig.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	return gorm.Open("mysql", cfg.FormatDSN())
}
