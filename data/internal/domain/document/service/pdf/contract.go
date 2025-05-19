package pdf

type (
	ocr interface {
		Process(filePath string) (string, error)
	}
)
