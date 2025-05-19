package service

type (
	OCR interface {
		ProcessImage(filePath string) (string, error)
	}
	PDF interface {
		ConvertPDFToImages(filePath string) ([]string, error)
	}
)
