package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	FileCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "s3_backup_file_count",
			Help: "Número de arquivos por diretório (ano/mes/dia)",
		},
		[]string{"bucket", "prefix"},
	)

	TotalSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "s3_backup_total_bytes",
			Help: "Total de bytes por diretório (ano/mes/dia)",
		},
		[]string{"bucket", "prefix"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(FileCount)
	prometheus.MustRegister(TotalSize)
}
