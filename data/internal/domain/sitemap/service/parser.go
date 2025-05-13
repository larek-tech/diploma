package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/larek-tech/diploma/data/internal/domain/sitemap"
)

type SitemapParser struct{}

func New() *SitemapParser {
	return &SitemapParser{}
}

func fetchSitemapContent(sitemapURL string) (string, error) {
	resp, err := http.Get(sitemapURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}

func (sp *SitemapParser) ParseSitemapContentRecursive(content string, visited map[string]struct{}) ([]sitemap.URLResult, error) {

	var urlset sitemap.URLSet
	if err := xml.Unmarshal([]byte(content), &urlset); err == nil && len(urlset.URLs) > 0 {
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

	var sitemapIndex sitemap.SitemapIndex
	if err := xml.Unmarshal([]byte(content), &sitemapIndex); err == nil && len(sitemapIndex.Sitemaps) > 0 {
		var allResults []sitemap.URLResult
		for _, sm := range sitemapIndex.Sitemaps {
			loc := strings.TrimSpace(sm.Loc)
			if loc == "" {
				continue
			}
			if _, seen := visited[loc]; seen {
				continue // avoid cycles
			}
			visited[loc] = struct{}{}
			subContent, err := fetchSitemapContent(loc)
			if err != nil {
				continue // skip broken links
			}
			subResults, err := sp.ParseSitemapContentRecursive(subContent, visited)
			if err == nil {
				allResults = append(allResults, subResults...)
			}
		}
		return allResults, nil
	}

	return nil, fmt.Errorf("failed to parse sitemap as urlset or sitemapindex")
}

func (sp *SitemapParser) GetAndParseSitemap(siteURL url.URL) ([]sitemap.URLResult, error) {
	siteURL.Path = "/sitemap.xml"
	siteURL.RawQuery = ""
	content, err := fetchSitemapContent(siteURL.String())
	if err != nil {
		return nil, err
	}
	return sp.ParseSitemapContentRecursive(content, map[string]struct{}{siteURL.String(): {}})
}
