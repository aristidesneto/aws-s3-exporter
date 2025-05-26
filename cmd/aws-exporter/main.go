package main

import (
	"aws-exporter/internal/collector"
	"aws-exporter/internal/config"
	"aws-exporter/internal/metrics"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics.InitMetrics()
	cfg := config.LoadConfigFile()

	s3Collector := collector.NewS3Collector(cfg)

	go func() {
		for {
			log.Printf("Profile de configuração carregado: %s", cfg.AwsProfile)

			if err := s3Collector.Collect(); err != nil {
				log.Printf("Erro ao coletar métricas: %v", err)
			}

			log.Printf("Aguardando %d minutos antes para a próxima coleta", cfg.Interval)
			<-time.After(time.Duration(cfg.Interval) * time.Minute)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Exporter rodando em :2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
