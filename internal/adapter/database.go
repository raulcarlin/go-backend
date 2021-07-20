package adapter

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/raulcarlin/go-backend/internal/config"
)

func NewDB(conf *config.Config) (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.DBConfig.Host, conf.DBConfig.Port),
		DBName:               conf.DBConfig.DbName,
		User:                 conf.DBConfig.Username,
		Passwd:               conf.DBConfig.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	return sql.Open("mysql", cfg.FormatDSN())
}
