package main

import (
	"aws-s3-exporter/internal/collector"
	"aws-s3-exporter/internal/config"
	"aws-s3-exporter/internal/metrics"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configPath := flag.String("config", "", "Caminho para o arquivo de configuração")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		log.Fatal("O caminho do arquivo de configuração é obrigatório. Use a flag -config para especificar o caminho.")
	}

	metrics.InitMetrics()
	cfg := config.LoadConfigFile(*configPath)

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
