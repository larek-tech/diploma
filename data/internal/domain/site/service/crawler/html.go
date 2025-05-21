package crawler

import (
	"log/slog"

	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractMetadata(doc *goquery.Document) (map[string]string, error) {
	metaTags := make(map[string]string)
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {

		name, _ := s.Attr("name")
		if name == "" {
			name, _ = s.Attr("property")
		}
		if name == "" {
			name, _ = s.Attr("http-equiv")
		}

		content, _ := s.Attr("content")
		if name != "" && content != "" {
			metaTags[name] = content
		}
	})

	return metaTags, nil
}

func extractLinks(doc *goquery.Document, pageUrl string) []string {
	links := make([]string, 0)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// Skip links with hash/fragment
			if strings.Contains(href, "#") {
				slog.Debug("skipping link with fragment", "href", href)
				return
			}

			if !strings.HasPrefix(href, "http://") && !strings.HasPrefix(href, "https://") {
				// Handle any relative URL, not just those starting with "/"
				u, err := url.Parse(href)
				if err != nil {
					slog.Error("failed to parse URL", "href", href, "err", err)
					return
				}
				baseURL, err := url.Parse(pageUrl)
				if err != nil {
					slog.Error("failed to parse base URL", "pageUrl", pageUrl, "err", err)
					return
				}
				href = baseURL.ResolveReference(u).String()
			}
			links = append(links, href)
		}
	})
	return links
}
