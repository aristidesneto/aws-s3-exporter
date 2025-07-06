package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadAWSConfig(ctx context.Context, profile, region string) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	if region != "" {
		opts = append(opts, config.WithRegion(region))
		log.Printf("Usando região AWS da configuração: %s", region)
	}

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
		log.Printf("Usando perfil AWS da configuração: %s", profile)
	}

	log.Println("Carregando configuração AWS")
	return config.LoadDefaultConfig(ctx, opts...)
}
