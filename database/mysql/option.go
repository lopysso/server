package mysql

// Option model
type Option struct {
	Host string
	Port uint16
	Username string
	Password string
	DbName string
	Charset string
	MaxOpenConn uint
	MaxIdleConn uint
}

// OptionDefault default
func OptionDefault() Option {
	option := Option{}
	option.Host = "localhost"
	option.Port = 3306
	option.Username = "root"
	option.Password = ""
	option.DbName = "lopy_sso"
	option.Charset = "utf8"
	option.MaxOpenConn = 10
	option.MaxIdleConn = 1

	return option
}