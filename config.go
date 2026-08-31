package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Antenna AntennaConfig  `mapstructure:"antenna"`
	Signal  SignalConfig   `mapstructure:"signal"`
	Targets []TargetConfig `mapstructure:"targets"`
}

type AntennaConfig struct {
	NumElements    int     `mapstructure:"numElements"`
	NumBeams       int     `mapstructure:"numBeams"`
	ElementSpacing float64 `mapstructure:"elementSpacing"`
	ElementAzHPBW  float64 `mapstructure:"elementAzHPBW"`
	ElementElHPBW  float64 `mapstructure:"elementElHPBW"`
	AntennaTilt    float64 `mapstructure:"antennaTilt"`
	Frequency      float64 `mapstructure:"frequency"`
}

type SignalConfig struct {
	SampleRate float64 `mapstructure:"sample_rate"`
	PulseWidth float64 `mapstructure:"pulse_width"`
	Bandwidth  float64 `mapstructure:"bandwidth"`
}

type TargetConfig struct {
	Position [3]float64 `mapstructure:"location"`
	Velocity [3]float64 `mapstructure:"velocity"`
	RCS      float64    `mapstructure:"rcs"`
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
