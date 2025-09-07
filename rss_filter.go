package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

const (
	TITLE_INCLUDES         = 1 // タイトルに含まれる
	TITLE_OR_DESC_INCLUDES = 2 // タイトルまたは概要に含まれる
	TITLE_BEGINS           = 3 // タイトルの先頭
	TITLE_ENDS             = 4 // タイトルの終わり
	URL_INCLUDES           = 5 // URLに含まれる
	AUTHOR_INCLUDES        = 6 // 著者に含まれる
)

func GetByURL(url string) (string, error) {
	rss, err := getRSS(url, "")
	if err != nil {
		return "", fmt.Errorf("getRSS: %w", err)
	}
	res, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return "", fmt.Errorf("MarshalIndent: %w", err)
	}
	return string(res), nil
}

func GetByName(name string) (string, error) {
	ctx := context.Background()
	dbClient, err := NewDBClient(ctx)
	if err != nil {
		return "", fmt.Errorf("NewDBClient: %w", err)
	}
	defer dbClient.Close()

	cfg, err := dbClient.GetConfig(ctx, name)
	if err != nil {
		return "", fmt.Errorf("GetConfig: %w", err)
	}

	var mergedRSS feeds.Feed
	for _, r := range cfg.Rsss {
		rss, err := getRSS(r.URL, r.Title)
		if err != nil {
			return "", fmt.Errorf("getRSS: %w", err)
		}
		is := []*feeds.Item{}
		for _, v := range rss.Items {
			if isIgnore(v, cfg.Ignores) {
				continue
			}
			exists, err := dbClient.IsURLExists(ctx, v.Link.Href)
			if err != nil {
				return "", fmt.Errorf("IsURLExists: %w", err)
			}
			if exists {
				fmt.Printf("URL already exists: %s\n", v.Link.Href)
				continue
			}
			if err := dbClient.AddURL(ctx, v.Link.Href); err != nil {
				return "", fmt.Errorf("AddURL: %w", err)
			}
			is = append(is, v)
		}
		mergedRSS.Items = append(mergedRSS.Items, is...)
	}
	mergedRSS.Title = name
	mergedRSS.Description = name
	mergedRSS.Link = &feeds.Link{Href: ""}
	items := mergedRSS.Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Created.After(items[j].Created)
	})
	mergedRSS.Items = items
	res, err := mergedRSS.ToRss()
	if err != nil {
		return "", fmt.Errorf("ToRss: %w", err)
	}
	return string(res), nil
}

func isIgnore(i *feeds.Item, ignores []Ignore) bool {
	for _, ignore := range ignores {
		switch ignore.Condition {
		case TITLE_INCLUDES:
			if strings.Contains(i.Title, ignore.Text) {
				fmt.Printf("Ignored: %q in Title [%s](%s)\n", ignore.Text, i.Title, i.Link.Href)
				return true
			}
		case TITLE_OR_DESC_INCLUDES:
			if strings.Contains(i.Title, ignore.Text) || strings.Contains(i.Description, ignore.Text) {
				fmt.Printf("Ignored: %q in Title or Description [%s](%s)\n", ignore.Text, i.Title, i.Link.Href)
				return true
			}
		case TITLE_BEGINS:
			if strings.HasPrefix(i.Title, ignore.Text) {
				fmt.Printf("Ignored: %q in Title Prefix [%s](%s)\n", ignore.Text, i.Title, i.Link.Href)
				return true
			}
		case TITLE_ENDS:
			if strings.HasSuffix(i.Title, ignore.Text) {
				fmt.Printf("Ignored: %q in Title Suffix [%s](%s)\n", ignore.Text, i.Title, i.Link.Href)
				return true
			}
		case URL_INCLUDES:
			if strings.Contains(i.Link.Href, ignore.Text) {
				fmt.Printf("Ignored: %q in Link [%s](%s)\n", ignore.Text, i.Title, i.Link.Href)
				return true
			}
		case AUTHOR_INCLUDES:
			if i.Author != nil && i.Author.Name == ignore.Text {
				fmt.Printf("Ignored: %q is Author [%s](%s)\n", ignore.Text, i.Title, i.Link.Href)
				return true
			}
		}
	}
	return false
}

func getRSS(url string, title string) (*feeds.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("ParseURL(%s): %w", url, err)
	}
	now := time.Now()
	res := &feeds.Feed{
		Title:       feed.Title,
		Link:        &feeds.Link{Href: feed.Link},
		Description: feed.Description,
		Created:     now,
	}
	if title == "" {
		title = feed.Title
	}
	for _, v := range feed.Items {
		if v.PublishedParsed == nil {
			continue
		}
		item := &feeds.Item{}
		item.Title = v.Title
		item.Link = &feeds.Link{Href: v.Link}
		item.Created = *v.PublishedParsed
		item.Description = "[" + title + "] " + v.Description
		if v.Author != nil {
			item.Author = &feeds.Author{
				Name:  v.Author.Name,
				Email: v.Author.Email,
			}
		}
		res.Items = append(res.Items, item)
	}
	return res, nil
}
