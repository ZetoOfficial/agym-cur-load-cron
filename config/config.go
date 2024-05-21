package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type MobiFitness struct {
	AccessToken string `validate:"required"`
	ApiURL      string `validate:"required,url"`
}

type Config struct {
	MobiFitness *MobiFitness
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{
		MobiFitness: &MobiFitness{
			AccessToken: viper.GetString("mobifitness.access_token"),
			ApiURL:      viper.GetString("mobifitness.api_url"),
		},
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	return config, nil
}
