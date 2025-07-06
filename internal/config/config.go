package config

import (
	"log"

	"github.com/spf13/viper"
)

type AWSConfig struct {
	Profile string `mapstructure:"profile"`
	Region  string `mapstructure:"region"`
}

type S3Config struct {
	Buckets []string `mapstructure:"buckets"`
}

type Config struct {
	AWS      AWSConfig `mapstructure:"aws"`
	S3       S3Config  `mapstructure:"s3"`
	Interval int       `mapstructure:"interval"`
}

func LoadConfigFile(configPath string) Config {
	// Definir valores default
	viper.SetDefault("interval", 5)
	viper.SetDefault("s3_buckets", []string{})

	viper.AutomaticEnv()

	viper.BindEnv("s3.buckets", "S3_BUCKETS")
	viper.BindEnv("aws.profile", "AWS_PROFILE")
	viper.BindEnv("aws.region", "AWS_DEFAULT_REGION")

	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("aviso: não foi possível ler o arquivo de configuração: %v", err)
		}
		log.Printf("Arquivo de configuração carregado: %s", viper.ConfigFileUsed())
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	return cfg
}
