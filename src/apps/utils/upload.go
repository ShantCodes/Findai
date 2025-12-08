package utils

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSUploader struct {
	CDNUrl          string
	BucketName      string
	CredentialsFile string
	client          *storage.Client
}

// UploadFile uploads a file to GCS
func (u *GCSUploader) UploadFile(ctx context.Context, fileName, contentType string, file io.Reader) (string, error) {
	if u.client == nil {
		client, err := storage.NewClient(ctx, option.WithCredentialsFile(u.CredentialsFile))
		if err != nil {
			return "", err
		}
		u.client = client
	}
	bucket := u.client.Bucket(u.BucketName)
	obj := bucket.Object(fileName)

	w := obj.NewWriter(ctx)
	w.ContentType = contentType

	// Set metadata (optional)
	w.Metadata = map[string]string{
		"uploaded-by": "gin-gonic",
	}

	if _, err := io.Copy(w, file); err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	// Make the file publicly accessible (optional)
	/* if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	} */

	fileURL := fmt.Sprintf("%s/%s/%s", u.CDNUrl, u.BucketName, fileName)
	return fileURL, nil
}
