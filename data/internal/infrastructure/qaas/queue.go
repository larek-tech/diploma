package qaas

// Queue ограничение доступных очередей используемых в data сервисах
type Queue string

const (
	ParseSiteQueue       Queue = "web_parse_site"
	ParsePageQueue       Queue = "web_parse_page"
	ParsePageResultQueue Queue = "web_parse_page_result" // replacement for WebResult

	EmbedResultQueue Queue = "document_embed_result"
)
