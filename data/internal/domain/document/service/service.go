package service

import (
	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/html"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/img"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/markdown"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/pdf"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/txt"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	documentStorage documentStorage
	chunkStorage    chunkStorage
	questionStorage questionStorage
	questionService questionService
	parsers         map[document.FileExtension]parser
	embedder        embedder
	trManager       trManager
	tracer          trace.Tracer
}

func New(
	documentStorage documentStorage,
	chunkStorage chunkStorage,
	questionStorage questionStorage,
	questionService questionService,
	embedder embedder,
	ocr ocr,
	trManager trManager,
	tracer trace.Tracer,
) *Service {
	imgProcessor := img.New(ocr)

	return &Service{
		documentStorage: documentStorage,
		chunkStorage:    chunkStorage,
		questionStorage: questionStorage,
		questionService: questionService,
		parsers: map[document.FileExtension]parser{
			document.HTML: html.New(),
			document.MD:   markdown.New(),
			document.PNG:  imgProcessor,
			document.JPEG: imgProcessor,
			document.PDF:  pdf.New(ocr),
			document.TXT:  txt.New(),
		},
		embedder:  embedder,
		trManager: trManager,
		tracer:    tracer,
	}
}
