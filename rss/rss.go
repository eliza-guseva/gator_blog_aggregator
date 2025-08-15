// Package rss provides RSS feed aggregation.
package rss

import (
	"fmt"
	"net/http"
	"context"
	"encoding/xml"
	"io"
	//"strings"
)

type RSSFeed struct {
	Channel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Item []RSSItem `xml:"item"`
	} `xml:"channel"`
} 

type RSSItem struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
}



func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req	, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil { return nil, fmt.Errorf("error fetching feed: %v", err) }

	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { return nil, fmt.Errorf("error fetching feed: %v", err) }
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching feed: %v", err)
	}

	feed := RSSFeed{}
	feedContent, err := io.ReadAll(resp.Body)
	if err != nil { return nil, fmt.Errorf("error fetching feed: %v", err) }

	err = xml.Unmarshal(feedContent, &feed)
	if err != nil {
    	return nil, fmt.Errorf("error parsing XML: %v", err)
	}
	fmt.Println(feed)

	return &feed, nil
}
