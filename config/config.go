package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Debug       bool
	URL         string

	Rate             uint // in seconds
	LoginCapacity    uint
	PasswordCapacity uint
	IPCapacity       uint

	ContextTimeout uint // in ms
	Db             map[string]string
}

func (c *Config) IsProd() bool {
	return c.Environment == "prod"
}

func (c *Config) isDev() bool {
	return c.Environment == "dev"
}

func GetConfig(configPath string) (*Config, error) {
	var config Config
	splits := strings.Split(filepath.Base(configPath), ".")
	viper.SetConfigName(filepath.Base(splits[0]))
	viper.AddConfigPath(filepath.Dir(configPath))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read config %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to unmarshal config: %v", err)
	}
	if config.Environment == "" {
		log.Fatal("Unable to find environment parameter in config")
	}
	return &config, nil
}
