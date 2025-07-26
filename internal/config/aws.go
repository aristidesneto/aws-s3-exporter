package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadAWSConfig(awsConfigs []Aws) ([]aws.Config, error) {
	var configs []aws.Config

	for _, awsCfg := range awsConfigs {
		var opts []func(*config.LoadOptions) error

		if awsCfg.Region != "" {
			opts = append(opts, config.WithRegion(awsCfg.Region))
		}

		if awsCfg.Profile != "" {
			opts = append(opts, config.WithSharedConfigProfile(awsCfg.Profile))
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
		if err != nil {
			return nil, err
		}

		configs = append(configs, cfg)
	}

	log.Println("Configuração AWS carregada")
	return configs, nil
}
