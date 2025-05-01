package crawler

import (
	"fmt"
	"net/url"
	"strings"
)

// isSameDomain определяет, принадлежат ли два URL-адреса к одному базовому домену
// Примеры:
// - example.com и app.example.com => true (один базовый домен)
// - example.com и example.org => false (разные домены)
// - sub.example.com и other.example.com => true (один базовый домен)
func isSameDomain(url1, url2 string) (bool, error) {
	// Разбираем первый URL
	u1, err := url.Parse(url1)
	if err != nil {
		return false, fmt.Errorf("failed to parse first URL %s: %w", url1, err)
	}

	// Разбираем второй URL
	u2, err := url.Parse(url2)
	if err != nil {
		return false, fmt.Errorf("failed to parse second URL %s: %w", url2, err)
	}

	// Извлекаем хост из каждого URL
	host1 := u1.Hostname()
	host2 := u2.Hostname()

	// Разделяем имена хостов на части
	parts1 := strings.Split(host1, ".")
	parts2 := strings.Split(host2, ".")

	// Проверяем, что оба хоста имеют хотя бы 2 части
	if len(parts1) < 2 || len(parts2) < 2 {
		return false, fmt.Errorf("invalid hostname format: %s or %s", host1, host2)
	}

	// Обрабатываем особые случаи, такие как .co.uk, .com.au и т.д.
	domainPartsToConsider := 2

	// Проверяем известные многокомпонентные TLD
	if len(parts1) >= 3 && len(parts2) >= 3 {
		lastTwoParts1 := strings.Join(parts1[len(parts1)-2:], ".")
		lastTwoParts2 := strings.Join(parts2[len(parts2)-2:], ".")

		// Список известных многокомпонентных TLD
		multiPartTLDs := []string{"co.uk", "com.au", "co.jp", "co.nz", "org.uk", "gov.uk"}

		for _, tld := range multiPartTLDs {
			if lastTwoParts1 == tld || lastTwoParts2 == tld {
				domainPartsToConsider = 3
				break
			}
		}
	}

	// Получаем части домена, включая TLD
	var domain1, domain2 string
	if len(parts1) >= domainPartsToConsider {
		domain1 = strings.Join(parts1[len(parts1)-domainPartsToConsider:], ".")
	} else {
		domain1 = host1
	}

	if len(parts2) >= domainPartsToConsider {
		domain2 = strings.Join(parts2[len(parts2)-domainPartsToConsider:], ".")
	} else {
		domain2 = host2
	}

	return domain1 == domain2, nil
}

func cleanURL(rawURL string, removeQueryParams bool) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	parsedURL.Fragment = ""

	if removeQueryParams {
		parsedURL.RawQuery = ""
	}

	cleanedURL := parsedURL.String()

	if strings.HasSuffix(cleanedURL, "/") {
		cleanedURL = cleanedURL[:len(cleanedURL)-1]
	}

	return cleanedURL, nil
}
