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

func LoadConfigFile() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return cfg
}
