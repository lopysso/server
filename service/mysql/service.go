package mysql

import (
	"database/sql"
	"sync"

	"github.com/lopysso/server/config"
	"github.com/lopysso/server/database/mysql"
)

func mysqlConfig(conf config.Config) mysql.Option {
	// conf := InjectConfig()

	confMysql := mysql.OptionDefault()

	dbConfig := conf.Db

	confMysql.Host = dbConfig.Host
	confMysql.Port = dbConfig.Port
	confMysql.Username = dbConfig.Username
	confMysql.Password = dbConfig.Password
	confMysql.DbName = dbConfig.DbName
	confMysql.Charset = dbConfig.Charset
	confMysql.MaxOpenConn = dbConfig.Pool.MaxConn
	confMysql.MaxIdleConn = dbConfig.Pool.MaxIdleConn
	return confMysql
}

func NewMysql(conf config.Config) (*sql.DB, error) {
	option := mysqlConfig(conf)
	db, err := mysql.NewMysql(option)
	if err != nil {
		return nil, err
	}
	return db, nil
}

var mysqlInstance *sql.DB
var mysqlInstanceOnce sync.Once

func MysqlInstance(conf config.Config) *sql.DB {
	mysqlInstanceOnce.Do(func() {
		theInstance, err := NewMysql(conf)
		if err != nil {
			panic(err)
		}

		mysqlInstance = theInstance

	})

	return mysqlInstance
}
