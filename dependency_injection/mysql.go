//+build wireinject

package dependency_injection

import (
	"database/sql"

	"github.com/google/wire"
	// "github.com/lopysso/server/config"
	// "github.com/lopysso/server/database/mysql"
	mysql2 "github.com/lopysso/server/service/mysql"
)

// // var mysqlInstance *sql.DB = mysql.NewMysql(option mysql.Option)

// func mysqlConfig(conf config.Config) mysql.Option {
// 	// conf := InjectConfig()

// 	confMysql := mysql.OptionDefault()

// 	dbConfig := conf.Db

// 	confMysql.Host = dbConfig.Host
// 	confMysql.Port = dbConfig.Port
// 	confMysql.Username = dbConfig.Username
// 	confMysql.Password = dbConfig.Password
// 	confMysql.DbName = dbConfig.DbName
// 	confMysql.Charset = dbConfig.Charset
// 	confMysql.MaxOpenConn = dbConfig.Pool.MaxConn
// 	confMysql.MaxIdleConn = dbConfig.Pool.MaxIdleConn
// 	return confMysql
// }

// func NewMysql() (*sql.DB, error) {
// 	panic(wire.Build(InjectConfig,mysqlConfig, mysql.NewMysql))
// }

// func instance() *sql.DB {
// 	mysqlInstance, err := NewMysql()

// 	if err != nil {
// 		panic("mysql init error")
// 	}

// 	return mysqlInstance
// }

// var _mysqlInstance = instance()

func InjectMysql() *sql.DB {
	panic(wire.Build(InjectConfig, mysql2.MysqlInstance))
}
