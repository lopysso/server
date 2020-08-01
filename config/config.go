package config

// Config config
type Config struct {
	
	App AppConfig

	Db DbConfig

	NumberGenerator NumberGenerator
}

// Load 加载配置
func Load() Config {

	conf := Config{}

	conf.App = AppConfigLoad()

	conf.Db = DbConfigLoad()

	conf.NumberGenerator = NumberGeneratorConfigLoad()

	return conf
}
