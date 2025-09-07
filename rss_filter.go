package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Printf("name: %s\n", name)
	if name != "" {
		res, err := GetByName(name)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}
		fmt.Fprint(w, string(res))
	} else {
		url := r.URL.Query().Get("url")
		res, err := GetByURL(url)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
			return
		}
		fmt.Fprint(w, string(res))
	}
}

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
	for _, c := range configs {
		if c.name != name {
			continue
		}
		var mergedRSS feeds.Feed
		for _, r := range c.rsss {
			rss, err := getRSS(r.url, r.title)
			if err != nil {
				return "", fmt.Errorf("getRSS: %w", err)
			}
			is := []*feeds.Item{}
			for _, v := range rss.Items {
				if isIgnore(v) {
					continue
				}
				is = append(is, v)
			}
			mergedRSS.Items = append(mergedRSS.Items, is...)
		}
		mergedRSS.Title = name
		mergedRSS.Description = name
		mergedRSS.Link = &feeds.Link{Href: apiHost + "?name=" + name}
		items := mergedRSS.Items
		sort.Slice(items, func(i, j int) bool {
			return items[i].Created.Local().String() > items[j].Created.Local().String()
		})
		mergedRSS.Items = items
		res, err := mergedRSS.ToRss()
		if err != nil {
			return "", fmt.Errorf("ToRss: %w", err)
		}
		return string(res), nil
	}
	return "", errors.New("name not found")
}

func isIgnore(i *feeds.Item) bool {
	for _, v := range filter1 {
		if strings.Contains(i.Title, v) {
			fmt.Printf("Ignored: %q in Title [%s](%s)\n", v, i.Title, i.Link.Href)
			return true
		}
		if strings.Contains(i.Description, v) {
			fmt.Printf("Ignored: %q in Description [%s](%s)\n", v, i.Title, i.Link.Href)
			return true
		}
		if strings.Contains(i.Content, v) {
			fmt.Printf("Ignored: %q in Content [%s](%s)\n", v, i.Title, i.Link.Href)
			return true
		}
	}
	for _, v := range filter2 {
		if strings.Contains(i.Title, v) {
			fmt.Printf("Ignored: %q in Title [%s](%s)\n", v, i.Title, i.Link.Href)
			return true
		}
	}
	for _, v := range filter3 {
		if strings.HasSuffix(i.Title, v) {
			fmt.Printf("Ignored: %q in Title [%s](%s)\n", v, i.Title, i.Link.Href)
			return true
		}
	}
	for _, v := range filter4 {
		if strings.Contains(i.Link.Href, v) {
			fmt.Printf("Ignored: %q in Link [%s](%s)\n", v, i.Title, i.Link.Href)
			return true
		}
	}
	for _, v := range filter5 {
		if i.Author != nil && i.Author.Name == v {
			fmt.Printf("Ignored: %q is Author [%s](%s)\n", v, i.Title, i.Link.Href)
			return true
		}
	}
	for _, v := range filter6 {
		i.Title = strings.ReplaceAll(i.Title, v, "")
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

type config struct {
	name    string
	rsss    []rss
	ignores []ignore
}

type rss struct {
	title string
	url   string
}

type condition int

const (
	TITLE_INCLUDES         condition = 1 // タイトルに含まれる
	TITLE_OR_DESC_INCLUDES condition = 2 // タイトルまたは概要に含まれる
	TITLE_BEGINS           condition = 3 // タイトルの先頭
	TITLE_ENDS             condition = 4 // タイトルの終わり
	URL_INCLUDES           condition = 5 // URLに含まれる
	AUTHOR_INCLUDES        condition = 6 // 著者に含まれる
)

type ignore struct {
	condition condition
	text      string
}
