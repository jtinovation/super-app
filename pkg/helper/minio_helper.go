package helper

import (
	"context"
	"jti-super-app-go/config"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

func UploadFile(bucketName, objectName string, file multipart.File, size int64, contentType string) error {
	_, err := config.MinioClient.PutObject(context.Background(), bucketName, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func DeleteFile(bucketName, objectName string) error {
	return config.MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
}
