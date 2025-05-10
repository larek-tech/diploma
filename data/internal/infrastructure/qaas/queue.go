package qaas

// Queue ограничение доступных очередей используемых в data сервисах
type Queue string

const (
	ParseSiteQueue Queue = "web_parse_site"

	ParsePageQueue       Queue = "web_parse_page"
	ParsePageResultQueue Queue = "web_parse_page_result" // replacement for WebResult

	ParseS3Queue       Queue = "web_parse_s3"        // job for parsing s3 bucket
	ParseS3ResultQueue Queue = "web_parse_s3_result" // job for parsing s3 bucket result

	ParseSiteStatusQueue Queue = "web_parse_site_status" // job for collecting parsing status
	EmbedResultQueue     Queue = "document_embed_result"
)
