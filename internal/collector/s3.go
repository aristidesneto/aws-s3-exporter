package collector

import (
	"context"
	"fmt"
	"log"
	"time"

	"aws-s3-exporter/internal/config"
	"aws-s3-exporter/internal/helper"
	"aws-s3-exporter/internal/metrics"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Collector struct {
	client *s3.Client
	cfg    config.Config
}

func NewS3Collector(client *s3.Client, cfg config.Config) *S3Collector {
	return &S3Collector{
		client: client,
		cfg:    cfg,
	}
}

func (c *S3Collector) Collect() error {
	ctx := context.TODO()

	metrics.FileCount.Reset()
	metrics.TotalSize.Reset()
	metrics.LastUpload.Reset()

	if len(c.cfg.S3.Buckets) == 0 {
		return fmt.Errorf("nenhum bucket foi especificado")
	}

	// Em vez de listar todos os buckets, processa apenas os buckets configurados
	for _, bucketName := range c.cfg.S3.Buckets {
		// Verifica se temos acesso ao bucket
		_, err := c.client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			log.Printf("Erro ao verificar acesso ao bucket %s: %v", bucketName, err)
			continue
		}

		log.Printf("Processando bucket: %s", bucketName)

		if err := c.collectBucketMetrics(ctx, c.client, bucketName); err != nil {
			log.Printf("Erro ao coletar métricas do bucket %s: %v", bucketName, err)
		}

		log.Printf("Métricas do bucket %s coletadas com sucesso", bucketName)
	}

	return nil
}

func (c *S3Collector) collectBucketMetrics(ctx context.Context, client *s3.Client, bucketName string) error {
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})

	const retentionDays = 60

	countMap := make(map[string]int)
	sizeMap := make(map[string]int64)
	lastUploadMap := make(map[string]time.Time)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			size := aws.ToInt64(obj.Size)
			lastMod := aws.ToTime(obj.LastModified)

			// prefix := helper.ExtractDatePrefix(key)
			prefix, valid := helper.ExtractDatePrefixAndCheck(key, lastMod, retentionDays)
			if !valid {
				continue
			}

			countMap[prefix]++
			sizeMap[prefix] += size
			// Salvar maior timestamp por prefixo
			if current, ok := lastUploadMap[prefix]; !ok || lastMod.After(current) {
				lastUploadMap[prefix] = lastMod
			}
		}
	}

	for prefix, count := range countMap {
		metrics.FileCount.WithLabelValues(bucketName, prefix).Set(float64(count))
		metrics.TotalSize.WithLabelValues(bucketName, prefix).Set(float64(sizeMap[prefix]))
		metrics.LastUpload.WithLabelValues(bucketName, prefix).Set(float64(lastUploadMap[prefix].Unix()))
	}

	return nil
}
