package config

import (
	"log"

	"github.com/spf13/viper"
)

type Aws struct {
	Profile        string   `mapstructure:"profile"`
	Region         string   `mapstructure:"region"`
	Buckets        []string `mapstructure:"buckets"`
	ScrapeInterval int      `mapstructure:"scrape_interval"`
}

type AwsConfig struct {
	Aws []Aws `mapstructure:"aws"`
}

func LoadConfigFile(configPath string) AwsConfig {
	viper.AutomaticEnv()

	var cfg AwsConfig

	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("erro não foi possível ler o arquivo de configuração: %v", err)
		}
		log.Printf("Arquivo de configuração carregado: %s", viper.ConfigFileUsed())

		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatalf("error unmarshalling config: %v", err)
		}

		for i := range cfg.Aws {
			if cfg.Aws[i].ScrapeInterval == 0 {
				cfg.Aws[i].ScrapeInterval = 5
			}
		}
	}

	return cfg
}
