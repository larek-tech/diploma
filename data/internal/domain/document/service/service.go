package service

import (
	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/html"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/img"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/markdown"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/pdf"
)

type Service struct {
	documentStorage documentStorage
	chunkStorage    chunkStorage
	questionStorage questionStorage
	questionService questionService
	parsers         map[document.FileExtension]parser
	embedder        embedder
	trManager       trManager
}

func New(
	documentStorage documentStorage,
	chunkStorage chunkStorage,
	questionStorage questionStorage,
	questionService questionService,
	embedder embedder,
	ocr ocr,
	trManager trManager,
) *Service {
	img := img.New(ocr)

	return &Service{
		documentStorage: documentStorage,
		chunkStorage:    chunkStorage,
		questionStorage: questionStorage,
		questionService: questionService,
		parsers: map[document.FileExtension]parser{
			document.HTML: html.New(),
			document.MD:   markdown.New(),
			document.PNG:  img,
			document.JPEG: img,
			document.PDF:  pdf.New(ocr),
		},
		embedder:  embedder,
		trManager: trManager,
	}
}
