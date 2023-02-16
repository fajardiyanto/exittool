package config

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"log"
)

const (
	bucketName = "golekno.com"
)

var gcsClient *storage.Client

func GoogleCloudStorage() {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile("key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}

	gcsClient = client
}

func UploadFile(ctx context.Context, object string, file []byte) (string, error) {
	reader := bytes.NewReader(file)

	wc := gcsClient.Bucket(bucketName).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, reader); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	u := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, wc.Attrs().Name)
	return u, nil
}
