package document

type FileExtension string

const (
	PDF  FileExtension = ".pdf"
	PNG  FileExtension = ".png"
	JPEG FileExtension = ".jpeg"
	TXT  FileExtension = ".txt"
	MD   FileExtension = ".md"
	CSV  FileExtension = ".csv"
	HTML FileExtension = ".html"
)

var FileExtensionMap = map[string]FileExtension{
	".pdf":  PDF,
	".png":  PNG,
	".jpeg": JPEG,
	".jpg":  JPEG, // common alias for JPEG
	".txt":  TXT,
	".md":   MD,
	".csv":  CSV,
	".html": HTML,
}
