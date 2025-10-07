package services

import (
	"context"
	"dklautomationgo/models"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type NewsletterFetcher struct {
	sources []string
}

func NewNewsletterFetcher() *NewsletterFetcher {
	sourcesEnv := os.Getenv("NEWSLETTER_SOURCES")
	var sources []string
	if sourcesEnv != "" {
		parts := strings.Split(sourcesEnv, ",")
		for _, p := range parts {
			s := strings.TrimSpace(p)
			if s != "" {
				sources = append(sources, s)
			}
		}
	}
	return &NewsletterFetcher{sources: sources}
}

// minimal RSS parsing structures
type rss struct {
	Channel rssChannel `xml:"channel"`
}
type rssChannel struct {
	Items []rssItem `xml:"item"`
}
type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
}

func (f *NewsletterFetcher) Fetch(ctx context.Context) ([]models.NewsItem, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	var items []models.NewsItem
	for _, src := range f.sources {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		func() {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				io.Copy(io.Discard, resp.Body)
				return
			}
			var doc rss
			dec := xml.NewDecoder(resp.Body)
			if err := dec.Decode(&doc); err != nil {
				return
			}
			for _, it := range doc.Channel.Items {
				if strings.TrimSpace(it.Title) == "" || strings.TrimSpace(it.Link) == "" {
					continue
				}
				var ts time.Time
				if t, err := time.Parse(time.RFC1123Z, it.PubDate); err == nil {
					ts = t
				}
				items = append(items, models.NewsItem{
					Title:       it.Title,
					Description: it.Description,
					Link:        it.Link,
					PubDate:     ts,
					Category:    it.Category,
				})
			}
		}()
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].PubDate.After(items[j].PubDate) })
	return items, nil
}
