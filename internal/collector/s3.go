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
	s3Client  *s3.Client
	awsConfig config.Aws
}

func NewS3Collector(client *s3.Client, awsConfig config.Aws) *S3Collector {
	return &S3Collector{
		s3Client:  client,
		awsConfig: awsConfig,
	}
}

func (c *S3Collector) Collect() error {
	ctx := context.TODO()

	metrics.FileCount.Reset()
	metrics.TotalSize.Reset()
	metrics.LastUpload.Reset()

	if len(c.awsConfig.Buckets) == 0 {
		return fmt.Errorf("nenhum bucket foi especificado")
	}

	// Em vez de listar todos os buckets, processa apenas os buckets configurados
	for _, bucketName := range c.awsConfig.Buckets {
		// Verifica se temos acesso ao bucket
		_, err := c.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucketName),
		})

		profile := c.awsConfig.Profile
		if err != nil {
			log.Printf("[ %s] Erro ao verificar acesso ao bucket %s: %v", profile, bucketName, err)
			continue
		}

		log.Printf("[ %s] Processando bucket: %s", profile, bucketName)

		if err := c.collectBucketMetrics(ctx, c.s3Client, bucketName); err != nil {
			log.Printf("[ %s] Erro ao coletar métricas do bucket %s: %v", profile, bucketName, err)
		}

		log.Printf("[ %s] Métricas do bucket %s coletadas com sucesso", profile, bucketName)
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
