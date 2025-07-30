package main

import (
	"aws-s3-exporter/internal/collector"
	"aws-s3-exporter/internal/config"
	"aws-s3-exporter/internal/metrics"
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configPath := flag.String("config", "", "Caminho para o arquivo de configuração")
	flag.Parse()

	metrics.InitMetrics()
	configFile := config.LoadConfigFile(*configPath)

	awsConfigs, err := config.LoadAWSConfig(configFile.Aws)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}

	// Configs de contas AWS
	for i := range awsConfigs {
		cfg := configFile.Aws[i]
		profile := cfg.Profile
		interval := cfg.ScrapeInterval

		log.Printf("[ %s ] 🚀 Iniciando loop de coleta a cada %d minutos", profile, interval)

		s3Client := s3.NewFromConfig(awsConfigs[i])
		s3Collector := collector.NewS3Collector(s3Client, cfg)

		go func(profile string, interval int, s3Collector *collector.S3Collector) {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
				err := s3Collector.Collect(ctx)
				if err != nil {
					log.Printf("[ %s ] ❌ Erro na  coleta de métricas: %v", profile, err)
				}

				cancel()
				log.Printf("[ %s ] ⏳ Nova coleta em %d minutos", profile, interval)
				<-time.After(time.Duration(interval) * time.Minute)
			}
		}(profile, interval, s3Collector)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Println("🌐 Exporter rodando em :2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
