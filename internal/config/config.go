package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Buckets    []string
	Interval   int
	AwsProfile string `mapstructure:"aws_profile"`
	AwsRegion  string `mapstructure:"aws_region"`
}

func LoadConfigFile(configPath string) Config {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("erro ao ler arquivo de configuração: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return cfg
}
