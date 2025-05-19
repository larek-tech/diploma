package site

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidSiteID = errors.New("invalid site ID must be a valid UUID")
	ErrInvalidURL    = errors.New("invalid URL must be a valid URL")
)

// TODO: add sitemap
type Site struct {
	ID             string    `db:"id"`              // ID uuid идентификатор сайта
	SourceID       string    `db:"source_id"`       // SourceID идентификатор источника к которому относится сайт
	URL            string    `db:"url"`             // URL корневой адрес сайта
	AvailablePages []string  `db:"available_pages"` // url страниц полученных из sitemap
	CreatedAt      time.Time `db:"created_at"`      // CreatedAt время создания сайта
	UpdatedAt      time.Time `db:"updated_at"`      // UpdatedAt время последнего обновления сайта
}

func NewSite(sourceID, siteURL string) (*Site, error) {
	if sourceID == "" {
		return nil, fmt.Errorf("failed to create site: %w", ErrInvalidSiteID)
	}
	if siteURL == "" {
		return nil, fmt.Errorf("failed to create site: %w", ErrInvalidURL)
	}
	_, err := uuid.Parse(sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", ErrInvalidSiteID)
	}
	_, err = url.Parse(siteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w, %w", ErrInvalidURL, err)
	}
	site := &Site{
		ID:       uuid.NewString(),
		SourceID: sourceID,
		URL:      siteURL,
	}
	return site, nil
}

type Page struct {
	ID            string            `db:"id"`            // ID uuid идентификатор страницы
	SiteID        string            `db:"site_id"`       // SiteID идентификатор сайта к которому относится страница
	URL           string            `db:"url"`           // URL адрес страницы
	Metadata      map[string]string `db:"metadata"`      // Metadata метаданные страницы (needs JSONB in Postgres)
	RawObjectID   string            `db:"raw_object_id"` // RawObjectID идентификатор объекта в S3 с необработанным содержанием страницы
	Raw           string            `db:"-" json:"-"`    // Raw необработанное содержание страницы
	Content       string            `db:"content"`       // Content текстовое содержание страницы
	OutgoingPages []string          `db:"outgoing"`      // OutgoingPages список UUID страниц на которые ссылается текущая страница
	CreatedAt     time.Time         `db:"created_at"`    // CreatedAt время создания страницы
	UpdatedAt     time.Time         `db:"updated_at"`    // UpdatedAt время последнего обновления страницы
}

// NewPage конструктор для новой страницы, перед сохранением в хранилище
func NewPage(siteID string, pageURL string) (*Page, error) {
	if siteID == "" {
		return nil, fmt.Errorf("failed to create page: %w", ErrInvalidSiteID)
	}
	if pageURL == "" {
		return nil, fmt.Errorf("failed to create page: %w", ErrInvalidURL)
	}
	_, err := uuid.Parse(siteID)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", ErrInvalidSiteID)
	}
	_, err = url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w, %w", ErrInvalidURL, err)
	}
	page := &Page{
		ID:            uuid.NewString(),
		SiteID:        siteID,
		URL:           pageURL,
		OutgoingPages: make([]string, 0),
		Metadata:      make(map[string]string),
		CreatedAt:     time.Now(),
	}
	return page, nil
}
