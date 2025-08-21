package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")

	err := config.ReadInConfig()
	if err != nil {
		fmt.Println("⚠️  Config file not found, falling back to environment variables:", err)
		config.AutomaticEnv()
	} else {
		fmt.Println("✅ Config loaded from file:", config.ConfigFileUsed())
	}

	return config
}
