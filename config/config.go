package config

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Configure(ctx context.Context) error {
	godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
