package config

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var minioClient *minio.Client

func Minio() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "FoPVFT3wObO61vtg"
	secretAccessKey := "ay4DjOKy6VmxTJ18HQLa0soqouik5OGF"

	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	minioClient = client
}

func UploadToMinio(ctx context.Context, objectName string, file []byte) error {
	r := bytes.NewReader(file)

	if _, err := minioClient.PutObject(ctx, "document", objectName, r, -1, minio.PutObjectOptions{}); err != nil {
		return err
	}
	return nil
}
