package config

import "github.com/spf13/viper"

// NumberGenerator 号码生成器
type NumberGenerator struct {
	Snowflake Snowflake
}

// Snowflake 雪花算法
type Snowflake struct {
	Node int64
}

// NumberGeneratorConfigLoad load config from env and db config file
func NumberGeneratorConfigLoad() NumberGenerator {
	v := viper.New()

	// default
	v.SetDefault("snowflake.node", "1")

	// config
	v.AddConfigPath("config/")
	v.SetConfigName("number_generator")
	v.SetConfigType("yml")

	v.ReadInConfig()

	v.BindEnv("snowflake.node", "NUMBER_GENERATOR_SNOWFLAKE_NODE")

	i := NumberGenerator{}
	v.Unmarshal(&i)

	return i
}
