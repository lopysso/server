package config

import (
	"github.com/spf13/viper"
)

// DbConfig config of database
type DbConfig struct {
	// 暂时不用管，只支持mysql
	Adapter string

	Host string
	Port uint16

	Username string

	Password string

	DbName string

	Charset string

	Pool Pool
}

// Pool 连接池
type Pool struct {
	MaxConn uint

	MaxIdleConn uint
}

// DbConfigLoad load config from env and db config file
func DbConfigLoad() DbConfig {
	v := viper.New()

	// default
	v.SetDefault("adapter", "mysql")
	v.SetDefault("host", "localhost")
	v.SetDefault("port", "3306")
	v.SetDefault("username", "root")
	v.SetDefault("password", "")
	v.SetDefault("charset", "utf8")
	v.SetDefault("dbName", "lopy_sso")
	v.SetDefault("pool.maxConn", "2")
	v.SetDefault("pool.maxConn", "10")

	// config
	v.AddConfigPath("config/")
	v.SetConfigName("db")
	v.SetConfigType("yml")

	v.ReadInConfig()

	// env
	v.SetEnvPrefix("DB")

	v.BindEnv("adapter")
	v.BindEnv("host")
	v.BindEnv("port")
	v.BindEnv("password")
	v.BindEnv("charset")
	v.BindEnv("dbName")

	// prefix not working ??
	v.BindEnv("pool.MaxIdleConn", "DB_POOL_IDLE")
	v.BindEnv("pool.MaxConn", "DB_POOL_CONN")

	i := DbConfig{}
	v.Unmarshal(&i)

	return i
}
