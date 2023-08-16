package crawler

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Timeout: 30 * time.Second, // Set a timeout
}

type Crawler struct {
	Url  string
	Text string
}

func NewCrawler(url string) *Crawler {
	return &Crawler{
		Url: url,
	}
}

func (c *Crawler) Crawl() (*goquery.Document, error) {
	// Fetch the URL
	res, err := c.fetch(c.Url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Parse the page with goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Extract the page text
	c.Text = doc.Find("body").Text()

	// Find all links in the page
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			lc := make(chan string)
			lc <- link
		}
	})

	// Return the parsed page
	return doc, nil
}

func (c *Crawler) fetch(url string) (*http.Response, error) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("User-Agent", "Crawler")

	// Use the client to send the request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Return the response
	return res, nil
}
