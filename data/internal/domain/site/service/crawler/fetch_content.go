package crawler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/larek-tech/diploma/data/internal/domain/site"
)

func (s Service) fetchContent(ctx context.Context, page *site.Page) ([]string, error) {
	ctx, span := s.tracer.Start(ctx, "crawlerService.fetchContent", trace.WithAttributes(
		attribute.String("url", page.URL),
		attribute.String("pageID", page.ID),
	))
	defer span.End()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, page.URL, nil)
	if err != nil {
		err = fmt.Errorf("create request error: %w", err)
		span.RecordError(err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DataEngineCrawler/1.0)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("http request error: %w", err)
		span.RecordError(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("read response body error: %w", err)
		span.RecordError(err)
		return nil, err
	}
	rawContent := string(raw)

	// Create a goquery document from the raw HTML
	doc, err := goquery.NewDocumentFromReader(io.NopCloser(bytes.NewReader(raw)))
	if err != nil {
		err = fmt.Errorf("create goquery document error: %w", err)
		span.RecordError(err)
		return nil, err
	}
	// TODO: get back the title of the html page
	//title := doc.Find("title").Text()

	// fetchMetadata
	metadata, err := extractMetadata(doc)
	if err != nil {
		err = fmt.Errorf("extract metadata error: %w", err)
		span.RecordError(err)
		return nil, err
	}

	page.Raw = rawContent
	page.Metadata = metadata
	page.UpdatedAt = time.Now()

	links := extractLinks(doc, page.URL)
	return links, nil
}
