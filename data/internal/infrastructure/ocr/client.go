package ocr

type job struct {
	imgPath string
	resCh   chan result
}

type result struct {
	text string
	err  error
}

type OCR struct {
	t     tesseract
	jobCh chan job
}

func New(t tesseract) *OCR {
	jobCh := make(chan job)
	ocr := &OCR{
		t:     t,
		jobCh: jobCh,
	}
	go ocr.Run()
	return ocr
}

func (o *OCR) Run() {
	for j := range o.jobCh {
		err := o.t.SetImage(j.imgPath)
		if err != nil {
			j.resCh <- result{"", err}
			continue
		}
		text, err := o.t.Text()
		j.resCh <- result{text, err}
	}
}

func (o *OCR) Process(filePath string) (string, error) {
	resCh := make(chan result)
	o.jobCh <- job{imgPath: filePath, resCh: resCh}
	res := <-resCh
	return res.text, res.err
}
