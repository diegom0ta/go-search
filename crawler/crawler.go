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
	Url string
}

func NewCrawler(url string) *Crawler {
	return &Crawler{
		Url: url,
	}
}

func (c *Crawler) Crawl(lc chan string) (*goquery.Document, chan string, string, error) {
	// Fetch the URL
	res, err := c.fetch(c.Url)
	if err != nil {
		return nil, nil, "", err
	}

	defer res.Body.Close()

	// Parse the page with goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, "", err
	}

	// Extract the page text
	text := doc.Find("body").Text()

	// Find all links in the page

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			lc <- link
		}
	})

	close(lc)

	// Return the parsed page
	return doc, lc, text, nil
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
