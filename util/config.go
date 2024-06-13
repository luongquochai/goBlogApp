package util

import "github.com/spf13/viper"

type Config struct {
	SECRET_KEY       string `mapstructure:"SECRET_KEY"`
	DB_SOURCE        string `mapstructure:"DB_SOURCE"`
	DB_DRIVER        string `mapstructure:"DB_DRIVER"`
	SUPER_SECRET_KEY string `mapstructure:"SUPER_SECRET_KEY"`
	MY_SECRET_KEY    string `mapstructure:"MY_SECRET_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
