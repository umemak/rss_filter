package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const urlCollectionName = "rss_filter_urls"
const configCollectionName = "rss_configurations"

type DBClient struct {
	fs *firestore.Client
}

// Config represents the structure of the configuration in Firestore.
// Note that the field names must be exported (start with a capital letter)
// to be used by the Firestore client.
type Config struct {
	Name    string   `firestore:"name"`
	Rsss    []Rss    `firestore:"rsss"`
	Ignores []Ignore `firestore:"ignores"`
}

type Rss struct {
	Title string `firestore:"title"`
	URL   string `firestore:"url"`
}

type Ignore struct {
	Condition int    `firestore:"condition"`
	Text      string `firestore:"text"`
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

func (c *DBClient) GetConfig(ctx context.Context, name string) (*Config, error) {
	dsnap, err := c.fs.Collection(configCollectionName).Doc(name).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("config %q not found: %w", name, err)
		}
		return nil, fmt.Errorf("Get: %w", err)
	}
	var cfg Config
	if err := dsnap.DataTo(&cfg); err != nil {
		return nil, fmt.Errorf("DataTo: %w", err)
	}
	return &cfg, nil
}

func (c *DBClient) ListConfigs(ctx context.Context) ([]string, error) {
	var names []string
	iter := c.fs.Collection(configCollectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Next: %w", err)
		}
		names = append(names, doc.Ref.ID)
	}
	return names, nil
}

func (c *DBClient) SaveConfig(ctx context.Context, config *Config) error {
	_, err := c.fs.Collection(configCollectionName).Doc(config.Name).Set(ctx, config)
	if err != nil {
		return fmt.Errorf("Set: %w", err)
	}
	return nil
}

func (c *DBClient) DeleteConfig(ctx context.Context, name string) error {
	_, err := c.fs.Collection(configCollectionName).Doc(name).Delete(ctx)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}

func (c *DBClient) IsURLExists(ctx context.Context, url string) (bool, error) {
	docID := fmt.Sprintf("%x", sha1.Sum([]byte(url)))
	doc, err := c.fs.Collection(urlCollectionName).Doc(docID).Get(ctx)
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
	_, err := c.fs.Collection(urlCollectionName).Doc(docID).Set(ctx, map[string]interface{}{
		"url":       url,
		"createdAt": firestore.ServerTimestamp,
	})
	if err != nil {
		return fmt.Errorf("Set: %w", err)
	}
	return nil
}
