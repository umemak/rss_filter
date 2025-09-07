package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collectionName = "rss_filter_urls"

type DBClient struct {
	fs *firestore.Client
}

func NewDBClient(ctx context.Context) (*DBClient, error) {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT_ID environment variable must be set")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}
	return &DBClient{fs: client}, nil
}

func (c *DBClient) Close() {
	c.fs.Close()
}

func (c *DBClient) IsURLExists(ctx context.Context, url string) (bool, error) {
	docID := fmt.Sprintf("%x", sha1.Sum([]byte(url)))
	doc, err := c.fs.Collection(collectionName).Doc(docID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, fmt.Errorf("Get: %w", err)
	}
	return doc.Exists(), nil
}

func (c *DBClient) AddURL(ctx context.Context, url string) error {
	docID := fmt.Sprintf("%x", sha1.Sum([]byte(url)))
	_, err := c.fs.Collection(collectionName).Doc(docID).Set(ctx, map[string]interface{}{
		"url":       url,
		"createdAt": firestore.ServerTimestamp,
	})
	if err != nil {
		return fmt.Errorf("Set: %w", err)
	}
	return nil
}
