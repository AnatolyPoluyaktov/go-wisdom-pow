package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Client ClientConfig `mapstructure:"client"`
}

type ServerConfig struct {
	Port int       `mapstructure:"port"`
	POW  POWConfig `mapstructure:"pow"`
}

type POWConfig struct {
	Difficulty  int `mapstructure:"difficulty"`
	TTL         int `mapstructure:"ttl"`
	MaxAttempts int `mapstructure:"max_attempts"`
}

type ClientConfig struct {
	ServerAddr string `mapstructure:"server_addr"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
