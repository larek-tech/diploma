package ocr

type (
	tesseract interface {
		SetImage(filePath string) error
		Text() (string, error)
		Close() error
	}
)
