package crawler

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// преобразует URL в путь на диске
func URLToLocalPath(rawURL string) string {
	u, _ := url.Parse(rawURL)

	path := u.Path

	// если это "директория" — сохраняем как index.html
	if path == "" || strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	return filepath.Join(u.Host, path)
}

// сохраняет файл и создаёт нужные директории
func SaveFile(rawURL string, data []byte) string {
	localPath := filepath.Join("site", URLToLocalPath(rawURL))

	os.MkdirAll(filepath.Dir(localPath), os.ModePerm)

	err := os.WriteFile(localPath, data, 0644)
	if err != nil {
		println("write error:", err.Error())
	}

	return localPath
}
