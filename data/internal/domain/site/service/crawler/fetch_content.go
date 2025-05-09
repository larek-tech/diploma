package crawler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"

	"github.com/larek-tech/diploma/data/internal/domain/site"
)

func (s Service) fetchContent(ctx context.Context, page *site.Page) ([]string, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, page.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DataEngineCrawler/1.0)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %w", err)
	}
	rawContent := string(raw)

	// Create a goquery document from the raw HTML

	doc, err := goquery.NewDocumentFromReader(io.NopCloser(bytes.NewReader(raw)))
	if err != nil {
		return nil, fmt.Errorf("parse HTML error: %w", err)
	}
	// TODO: get back the title of the html page
	//title := doc.Find("title").Text()

	// fetchMetadata
	metadata, err := extractMetadata(doc)
	if err != nil {
		return nil, fmt.Errorf("extract metadata error: %w", err)
	}

	page.Raw = rawContent
	page.Metadata = metadata
	page.UpdatedAt = time.Now()

	links := extractLinks(doc, page.URL)
	return links, nil
}

func cleanUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				// skip invalid byte
				continue
			}
		}
		v = append(v, r)
	}
	return string(v)
}
