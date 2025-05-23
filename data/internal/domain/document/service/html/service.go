package html

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/larek-tech/diploma/data/internal/domain/document"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

// Parse извлекает основной текстовый контент из HTML-документа, удаляя элементы,
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
//	s := html.New()
//	text, err := s.Parse(reader)
//	if err != nil {
//	    // обработка ошибки
//	}
func (s Service) Parse(content io.ReadSeeker) (string, error) {
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
	return document.CleanUTF8(text), nil
}

func (s Service) STDParse(content io.ReadSeeker) (string, error) {
	_, err := content.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek HTML content: %w", err)
	}
	data, err := io.ReadAll(content)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML: %w", err)
	}
	html := string(data)

	// Remove <script>...</script>
	reScript := regexp.MustCompile(`(?is)<script.*?>.*?</script>`)
	html = reScript.ReplaceAllString(html, "")

	// Remove <nav>...</nav>
	reNav := regexp.MustCompile(`(?is)<nav.*?>.*?</nav>`)
	html = reNav.ReplaceAllString(html, "")

	// Remove elements with class/id containing nav, menu, navbar
	reClassId := regexp.MustCompile(`(?is)<[a-z0-9]+\s+[^>]*(class|id)\s*=\s*["'][^"']*(nav|menu|navbar)[^"']*["'][^>]*>.*?</[a-z0-9]+>`)
	html = reClassId.ReplaceAllString(html, "")

	// Remove <a>...</a>
	reA := regexp.MustCompile(`(?is)<a.*?>.*?</a>`)
	html = reA.ReplaceAllString(html, "")

	// Extract text from <body>
	reBody := regexp.MustCompile(`(?is)<body.*?>(.*?)</body>`)
	matches := reBody.FindStringSubmatch(html)
	var text string
	if len(matches) > 1 {
		text = matches[1]
	} else {
		text = html
	}

	// Remove all tags
	reTags := regexp.MustCompile(`(?is)<.*?>`)
	text = reTags.ReplaceAllString(text, "")

	// Normalize whitespace
	text = strings.Join(strings.Fields(text), " ")

	return document.CleanUTF8(text), nil
}
