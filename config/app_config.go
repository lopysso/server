package config

// AppConfig app基础设置
type AppConfig struct {
	Name string

	Debug bool
}

// AppConfigLoad 加载app config
func AppConfigLoad() AppConfig {
	i := AppConfig{
		Name: "lopy_sso",
		Debug: true,
	}

	return i
}