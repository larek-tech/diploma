package img

type (
	ocr interface {
		Process(filePath string) (string, error)
	}
)
