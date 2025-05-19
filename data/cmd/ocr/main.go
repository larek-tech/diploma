package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/larek-tech/diploma/data/internal/domain/document/service/img"
	"github.com/larek-tech/diploma/data/internal/domain/document/service/pdf"
	"github.com/larek-tech/diploma/data/internal/infrastructure/ocr"
	"github.com/otiai10/gosseract/v2"
)

const (
	testPng       = "bin/mts.png"
	testPdf       = "bin/mts.pdf"
	testPdfResult = "bin/pdf.txt"
	testPngResult = "bin/png.txt"
)

func main() {
	client := gosseract.NewClient()
	client.Languages = []string{"rus", "eng"}
	defer client.Close()
	c := ocr.New(client)

	pdfService := pdf.New(c)
	imgService := img.New(c)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		file, err := os.Open(testPng)
		if err != nil {
			panic("failed to open png")
		}
		defer file.Close()
		text, err := imgService.Parse(file)
		if err != nil {
			panic("failed to parse document" + err.Error())
		}
		resultFile, err := os.Create(testPngResult)
		if err != nil {
			panic(err)
		}
		defer resultFile.Close()

		_, err = resultFile.WriteString(text)
		if err != nil {
			panic(err)
		}
		fmt.Println("finished processing image")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		file, err := os.Open(testPdf)
		if err != nil {
			panic("failed to open pdf")
		}
		defer file.Close()

		text, err := pdfService.Parse(file)
		if err != nil {
			panic("failed to parse document" + err.Error())
		}
		resultFile, err := os.Create(testPdfResult)
		if err != nil {
			panic(err)
		}
		defer resultFile.Close()

		_, err = resultFile.WriteString(text)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
