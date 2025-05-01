package site

import (
	"time"
)

type Site struct {
	ID             string    `db:"id"`              // ID uuid идентификатор сайта
	SourceID       string    `db:"source_id"`       // SourceID идентификатор источника к которому относится сайт
	URL            string    `db:"url"`             // URL корневой адрес сайта
	AvailablePages []string  `db:"available_pages"` // Might need special handling depending on your DB schema
	CreatedAt      time.Time `db:"created_at"`      // CreatedAt время создания сайта
	UpdatedAt      time.Time `db:"updated_at"`      // UpdatedAt время последнего обновления сайта
}

type Page struct {
	ID            string            `db:"id"`         // ID uuid идентификатор страницы
	SiteID        string            `db:"site_id"`    // SiteID идентификатор сайта к которому относится страница
	URL           string            `db:"url"`        // URL адрес страницы
	Metadata      map[string]string `db:"metadata"`   // Metadata метаданные страницы (needs JSONB in Postgres)
	Raw           string            `db:"raw"`        // Raw необработанное содержание страницы
	Content       string            `db:"content"`    // Content текстовое содержание страницы
	OutgoingPages []string          `db:"outgoing"`   // OutgoingPages список UUID страниц на которые ссылается текущая страница
	CreatedAt     time.Time         `db:"created_at"` // CreatedAt время создания страницы
	UpdatedAt     time.Time         `db:"updated_at"` // UpdatedAt время последнего обновления страницы
}
