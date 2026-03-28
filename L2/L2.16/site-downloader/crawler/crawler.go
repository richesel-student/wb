package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

// Crawler хранит состояние обхода:
// список посещённых URL, ограничения и настройки
type Crawler struct {
	BaseHost string          // домен, за пределы которого не выходим
	Visited  map[string]bool // уже обработанные ссылки
	Mu       sync.Mutex      // защита Visited при параллельной работе
	Wg       sync.WaitGroup  // ожидание завершения всех задач
	MaxDepth int             // максимальная глубина обхода

	Semaphore chan struct{}   // ограничение количества одновременных запросов
	Robots    map[string]bool // запрещённые пути из robots.txt
}

func NewCrawler(startURL string, depth int) *Crawler {
	u, err := url.Parse(startURL)
	if err != nil {
		panic(err) // некорректный URL — дальнейшая работа невозможна
	}

	c := &Crawler{
		BaseHost:  u.Host,
		Visited:   make(map[string]bool),
		MaxDepth:  depth,
		Semaphore: make(chan struct{}, 10), // ограничиваем параллельные загрузки
	}

	// загружаем robots.txt при инициализации
	c.Robots = c.LoadRobots(startURL)

	return c
}

// Start запускает обход с начального URL
// и ждёт завершения всех горутин
func (c *Crawler) Start(startURL string) {
	c.Wg.Add(1)
	go c.crawl(startURL, 0)
	c.Wg.Wait()
}

// NormalizeURL приводит URL к единому виду,
// чтобы избежать дубликатов
func NormalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	u.Fragment = "" // якорь не влияет на содержимое страницы

	if strings.HasSuffix(u.Path, "/") && len(u.Path) > 1 {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	return u.String()
}

// LoadRobots загружает robots.txt и сохраняет запрещённые пути
// реализовано в упрощённом виде (без User-Agent)
func (c *Crawler) LoadRobots(startURL string) map[string]bool {
	robots := make(map[string]bool)

	u, err := url.Parse(startURL)
	if err != nil {
		return robots
	}

	robotsURL := u.Scheme + "://" + u.Host + "/robots.txt"

	body, _, err := Download(robotsURL)
	if err != nil {
		return robots // если не удалось загрузить — считаем всё доступным
	}

	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Disallow:") {
			path := strings.TrimSpace(strings.TrimPrefix(line, "Disallow:"))
			if path != "" {
				robots[path] = true
			}
		}
	}

	return robots
}

// IsAllowed проверяет, разрешён ли доступ к URL
func (c *Crawler) IsAllowed(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}

	for disallowed := range c.Robots {
		if strings.HasPrefix(u.Path, disallowed) {
			return false
		}
	}

	return true
}

// crawl — основная функция обхода
// скачивает страницу, сохраняет её и обрабатывает ссылки
func (c *Crawler) crawl(link string, depth int) {
	defer c.Wg.Done()

	// ограничение глубины
	if depth > c.MaxDepth {
		return
	}

	// проверка robots.txt
	if !c.IsAllowed(link) {
		return
	}

	// ограничение количества параллельных запросов
	c.Semaphore <- struct{}{}
	defer func() { <-c.Semaphore }()

	link = NormalizeURL(link)

	// проверка на повторную обработку
	c.Mu.Lock()
	if c.Visited[link] {
		c.Mu.Unlock()
		return
	}
	c.Visited[link] = true
	c.Mu.Unlock()

	fmt.Println("Downloading:", link)

	body, contentType, err := Download(link)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// HTML требует разбора и извлечения ссылок
	if IsHTML(contentType) {

		newHTML, pages, resources := RewriteHTML(string(body), link)

		localPath := SaveFile(link, []byte(newHTML))
		fmt.Println("Saved page:", localPath)

		// сначала обрабатываем ресурсы (css, js, изображения)
		for _, res := range resources {
			c.Wg.Add(1)
			go c.crawl(res, depth+1)
		}

		// затем переходим по ссылкам на другие страницы
		for _, p := range pages {
			parsed, err := url.Parse(p)
			if err != nil {
				continue
			}

			// остаёмся в пределах исходного домена
			if parsed.Host != "" && parsed.Host != c.BaseHost {
				continue
			}

			c.Wg.Add(1)
			go c.crawl(p, depth+1)
		}

	} else if strings.Contains(contentType, "text/css") {

		// CSS может содержать дополнительные ссылки (например, на шрифты)
		localPath := SaveFile(link, body)
		fmt.Println("Saved CSS:", localPath)

		resources := ExtractCSSResources(string(body), link)

		for _, res := range resources {
			c.Wg.Add(1)
			go c.crawl(res, depth+1)
		}

	} else {
		// остальные типы файлов сохраняем без дополнительной обработки
		localPath := SaveFile(link, body)
		fmt.Println("Saved resource:", localPath)
	}
}
