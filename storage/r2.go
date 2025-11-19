package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

type FileInfo struct {
	Name         string
	LastModified time.Time
	Size         int64
}

func NewR2Client(accountID, accessKeyID, secretAccessKey, bucket, publicURL string) (*R2Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &R2Client{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

func (r *R2Client) UploadFile(ctx context.Context, filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file size for Content-Length header
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	fileName := filepath.Base(filePath)

	// Stream the file directly without loading into memory
	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.bucket),
		Key:           aws.String(fileName),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s", r.publicURL, fileName)
	return publicURL, nil
}

func (r *R2Client) DeleteFile(ctx context.Context, fileName string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from R2: %w", err)
	}
	return nil
}

func (r *R2Client) ListFiles(ctx context.Context) ([]string, error) {
	result, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var files []string
	for _, obj := range result.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}

func (r *R2Client) ListFilesWithMetadata(ctx context.Context, prefix string) ([]FileInfo, error) {
	result, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var files []FileInfo
	for _, obj := range result.Contents {
		files = append(files, FileInfo{
			Name:         *obj.Key,
			LastModified: *obj.LastModified,
			Size:         *obj.Size,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].LastModified.Before(files[j].LastModified)
	})

	return files, nil
}

func (r *R2Client) CleanupOldBackups(ctx context.Context, prefix string, retentionLimit int) ([]string, error) {
	if retentionLimit <= 0 {
		return nil, nil
	}

	files, err := r.ListFilesWithMetadata(ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list files for cleanup: %w", err)
	}

	if len(files) <= retentionLimit {
		return nil, nil
	}

	filesToDelete := files[:len(files)-retentionLimit]
	var deletedFiles []string

	for _, file := range filesToDelete {
		if err := r.DeleteFile(ctx, file.Name); err != nil {
			return deletedFiles, fmt.Errorf("failed to delete file %s: %w", file.Name, err)
		}
		deletedFiles = append(deletedFiles, file.Name)
	}

	return deletedFiles, nil
}
