package sitemap

import "encoding/xml"

// URLSet represents the root element of a sitemap XML
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URL represents an entry in the sitemap
type URL struct {
	Loc        string `xml:"loc"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	LastMod    string `xml:"lastmod,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// URLResult represents the JSON output structure
type URLResult struct {
	URL        string `json:"url"`
	ChangeFreq string `json:"changefreq,omitempty"`
	LastMod    string `json:"lastmod,omitempty"`
}

type SitemapIndex struct {
	XMLName  xml.Name      `xml:"sitemapindex"`
	Sitemaps []SitemapInfo `xml:"sitemap"`
}

type SitemapInfo struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}
