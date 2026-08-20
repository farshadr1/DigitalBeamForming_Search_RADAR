package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Signal SignalConfig `mapstructure:"signal"`
}

type SignalConfig struct {
	SampleRate float64 `mapstructure:"sample_rate"`
	PulseWidth float64 `mapstructure:"pulse_width"`
	Bandwidth  float64 `mapstructure:"bandwidth"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	// Creat instant from Config
	var cfg Config

	// Parsing config.yaml
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
