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
	cfg := config.LoadConfigFile(*configPath)

	awsConfig, err := config.LoadAWSConfig(context.TODO(), cfg.AWS.Profile, cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração AWS: %v", err)
	}
	s3Client := s3.NewFromConfig(awsConfig)
	s3Collector := collector.NewS3Collector(s3Client, cfg)

	go func() {
		for {
			err := s3Collector.Collect()
			if err != nil {
				log.Fatalf("Erro ao coletar métricas: %v", err)
			}

			log.Printf("Aguardando %d minutos antes para a próxima coleta", cfg.Interval)
			<-time.After(time.Duration(cfg.Interval) * time.Minute)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Exporter rodando em :2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
