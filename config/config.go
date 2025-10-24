package config

import (
	"context"

	"github.com/spf13/viper"
)

func Configure(ctx context.Context) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
