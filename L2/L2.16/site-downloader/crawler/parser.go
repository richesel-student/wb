package crawler

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// отбрасываем ссылки, которые не ведут на реальные страницы или ресурсы
func IsValidLink(link string) bool {
	if strings.HasPrefix(link, "mailto:") ||
		strings.HasPrefix(link, "javascript:") ||
		strings.HasPrefix(link, "tel:") ||
		strings.HasPrefix(link, "#") {
		return false
	}
	return true
}

// считаем страницами только ссылки вида <a href="...">
func IsPageLink(tag, attr string) bool {
	return tag == "a" && attr == "href"
}

// RewriteHTML:
// - разбирает HTML
// - находит ссылки
// - делает их абсолютными
// - заменяет на локальные пути
func RewriteHTML(htmlStr, base string) (string, []string, []string) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr, nil, nil
	}

	var resources []string
	var pages []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, attr := range n.Attr {

				if attr.Key == "href" || attr.Key == "src" {

					if !IsValidLink(attr.Val) {
						continue
					}

					abs := ResolveURL(attr.Val, base)
					if abs == "" {
						continue
					}

					abs = NormalizeURL(abs)

					// заменяем ссылку на путь в локальной файловой системе
					local := URLToLocalPath(abs)
					n.Attr[i].Val = local

					if IsPageLink(n.Data, attr.Key) {
						pages = append(pages, abs)
					} else {
						resources = append(resources, abs)
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	var buf strings.Builder
	html.Render(&buf, doc)

	return buf.String(), pages, resources
}

// извлекаем ресурсы из CSS (url(...) и @import)
func ExtractCSSResources(css string, base string) []string {
	var result []string

	re := regexp.MustCompile(`url\((.*?)\)`)
	matches := re.FindAllStringSubmatch(css, -1)

	for _, m := range matches {
		link := strings.Trim(m[1], `"'`)
		if !IsValidLink(link) {
			continue
		}

		abs := ResolveURL(link, base)
		if abs != "" {
			result = append(result, abs)
		}
	}

	re2 := regexp.MustCompile(`@import\s+["'](.*?)["']`)
	matches2 := re2.FindAllStringSubmatch(css, -1)

	for _, m := range matches2 {
		abs := ResolveURL(m[1], base)
		if abs != "" {
			result = append(result, abs)
		}
	}

	return result
}

// преобразует относительный URL в абсолютный
func ResolveURL(href, base string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseURL, _ := url.Parse(base)
	return baseURL.ResolveReference(u).String()
}

// проверяем, является ли ответ HTML
func IsHTML(contentType string) bool {
	return strings.Contains(contentType, "text/html")
}
