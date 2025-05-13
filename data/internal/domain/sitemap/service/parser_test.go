package service

import (
	"fmt"
	"log/slog"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testURL = "https://grafana.com/sitemap.xml"

func TestParseSitemapContent(t *testing.T) {
	URL, err := url.Parse(testURL)
	assert.NoError(t, err)

	res, err := New().GetAndParseSitemap(*URL)
	slog.Info("results:", "res", len(res))
	assert.NoError(t, err)
	fmt.Println(res)
}
