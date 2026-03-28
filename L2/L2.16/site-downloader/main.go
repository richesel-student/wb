package main

import (
	"flag"
	"fmt"

	"project/crawler"
)

func main() {
	// параметры запуска из командной строки
	url := flag.String("url", "https://example.com", "start url")
	depth := flag.Int("depth", 2, "recursion depth")

	flag.Parse()

	fmt.Println("Start crawling:", *url)

	// создаём и запускаем краулер
	c := crawler.NewCrawler(*url, *depth)
	c.Start(*url)
}
