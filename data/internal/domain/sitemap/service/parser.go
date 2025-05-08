package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/larek-tech/diploma/data/internal/domain/sitemap"
)

// SitemapParser handles sitemap parsing operations
type SitemapParser struct{}

func New() *SitemapParser {
	return &SitemapParser{}
}

// GetAndParseSitemap fetches a sitemap from a URL and parses it
func (sp *SitemapParser) GetAndParseSitemap(siteURL url.URL) ([]sitemap.URLResult, error) {
	siteURL.Path = "/sitemap.xml"

	resp, err := http.Get(siteURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return sp.ParseSitemapContent(string(body))
}

// ParseSitemapContent parses a sitemap XML string and returns URLs with changefreq
// ParseSitemapContent parses the XML content of a sitemap and converts it into a slice of URLResult objects.
// It unmarshals the provided XML string into a sitemap.URLSet structure and extracts relevant information
// from each URL entry such as the location, change frequency, and last modification time.
//
// Parameters:
//   - content: A string containing the XML sitemap content to be parsed
//
// Returns:
//   - []sitemap.URLResult: A slice of URLResult objects containing the extracted information
//   - error: An error if the XML parsing fails, nil otherwise
//
// Example usage:
//
//	content := `<?xml version="1.0" encoding="UTF-8"?>
//	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
//	  <url>
//	    <loc>https://example.com/</loc>
//	    <lastmod>2023-01-01</lastmod>
//	    <changefreq>daily</changefreq>
//	  </url>
//	</urlset>`
//	results, err := sitemapParser.ParseSitemapContent(content)
//	if err != nil {
//	    log.Fatalf("Failed to parse sitemap: %v", err)
//	}
//	for _, result := range results {
//	    fmt.Printf("URL: %s, LastMod: %s, ChangeFreq: %s\n", result.URL, result.LastMod, result.ChangeFreq)
//	}
func (sp *SitemapParser) ParseSitemapContent(content string) ([]sitemap.URLResult, error) {
	var urlset sitemap.URLSet
	if err := xml.Unmarshal([]byte(content), &urlset); err != nil {
		return nil, fmt.Errorf("failed to parse sitemap XML: %w", err)
	}

	var results []sitemap.URLResult
	for _, url := range urlset.URLs {
		results = append(results, sitemap.URLResult{
			URL:        url.Loc,
			ChangeFreq: url.ChangeFreq,
			LastMod:    url.LastMod,
		})
	}

	return results, nil
}
