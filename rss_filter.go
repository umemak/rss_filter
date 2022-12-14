package rssfilter

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Language    string    `xml:"language"`
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Creator     string `xml:"dc:creator"`
}

func Exec(url string) string {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	return fmt.Sprintf("%#v\n", feed)
}

func Fetch(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "RSSFilter/0.1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func Parse(buf []byte) (*RSS, error) {
	rss := RSS{}
	err := xml.Unmarshal(buf, &rss)
	if err != nil {
		return nil, fmt.Errorf("xml.Unmarshal: %w", err)
	}
	for i, v := range rss.Channel.Items {
		rss.Channel.Items[i].Description = strings.ReplaceAll(v.Description, "\n", "")
	}
	return &rss, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
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
		var mergedRSS RSS
		for i, r := range c.rsss {
			rss, err := getRSS(r.url, r.title)
			if err != nil {
				return "", fmt.Errorf("getRSS: %w", err)
			}
			if i == 0 {
				mergedRSS = *rss
			} else {
				mergedRSS.Channel.Items = append(mergedRSS.Channel.Items, rss.Channel.Items...)
			}
		}
		mergedRSS.Channel.Title = name
		mergedRSS.Channel.Description = name
		mergedRSS.Channel.Link = apiHost + "?name=" + name
		items := mergedRSS.Channel.Items
		sort.Slice(items, func(i, j int) bool {
			ti, _ := time.Parse(time.RFC1123Z, items[i].PubDate)
			tj, _ := time.Parse(time.RFC1123Z, items[j].PubDate)
			return ti.Unix() > tj.Unix()
		})
		mergedRSS.Channel.Items = items
		res, err := xml.MarshalIndent(mergedRSS, "", "  ")
		if err != nil {
			return "", fmt.Errorf("MarshalIndent: %w", err)
		}
		return string(res), nil
	}
	return "", errors.New("name not found")
}

func getRSS(url string, title string) (*RSS, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("ParseURL(%s): %w", url, err)
	}
	// res, err := Parse(buf)
	// if err != nil {
	// 	return nil, fmt.Errorf("Parse: %w", err)
	// }
	res := &RSS{Version: "2.0"}
	res.Channel.Language = feed.Language
	res.Channel.Title = feed.Title
	res.Channel.Link = feed.Link
	res.Channel.Description = feed.Description
	if title == "" {
		title = feed.Title
	}
	for _, v := range feed.Items {
		item := RSSItem{}
		item.Title = v.Title
		item.Link = v.Link
		item.PubDate = v.Published
		item.Description = "[" + title + "] " + v.Description
		if v.Author != nil {
			item.Creator = v.Author.Name
		}
		res.Channel.Items = append(res.Channel.Items, item)
	}
	return res, nil
}

type config struct {
	name string
	rsss []rss
}

type rss struct {
	title string
	url   string
}
