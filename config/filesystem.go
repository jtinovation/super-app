package config

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinio() {
	var err error
	MinioClient, err = minio.New(AppConfig.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AppConfig.Minio.AccessKeyID, AppConfig.Minio.SecretAccessKey, ""),
		Secure: AppConfig.Minio.UseSSL,
	})

	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}
}
