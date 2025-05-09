package service

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseHTML извлекает основной текстовый контент из HTML-документа, удаляя элементы,
// не несущие ценности для основного содержания страницы.
//
// Функция принимает io.ReadSeeker с HTML-контентом, очищает его от следующих элементов:
//   - <script> (JavaScript-скрипты)
//   - <nav> (навигационные панели)
//   - элементы с классами или id, содержащими "nav", "menu", "navbar"
//   - <a> (гиперссылки)
//
// После удаления ненужных элементов функция возвращает очищенный текст из <body>,
// либо, если <body> отсутствует, — весь текст документа.
//
// Возвращает строку с основным текстовым содержимым и ошибку, если что-то пошло не так.
//
// Пример использования:
//
//	text, err := ParseHTML(reader)
//	if err != nil {
//	    // обработка ошибки
//	}
func ParseHTML(content io.ReadSeeker) (string, error) {
	_, err := content.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek HTML content: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	doc.Find("nav").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	doc.Find("[class*='nav'], [id*='nav'], [class*='menu'], [id*='menu'], [class*='navbar'], [id*='navbar']").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	var sb strings.Builder

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		sb.WriteString(strings.Join(strings.Fields(s.Text()), " "))
	})
	text := sb.String()
	if text == "" {
		text = strings.Join(strings.Fields(doc.Text()), " ")
	}
	return text, nil
}
