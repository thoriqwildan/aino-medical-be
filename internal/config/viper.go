package config

import "github.com/spf13/viper"

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic("Error reading config file: " + err.Error())
	}
	
	return config
}