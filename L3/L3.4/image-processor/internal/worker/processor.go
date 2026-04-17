package worker

import (
	"bytes"
	"context"
	"image-processor/internal/models"
	"log"

	"github.com/disintegration/imaging"
)

// Интерфейсы для зависимостей
type Repo interface {
	Get(ctx context.Context, id string) (*models.Image, error)
	UpdateStatus(ctx context.Context, id, status, path string) error
}

type Storage interface {
	Download(ctx context.Context, key string) ([]byte, error)
	UploadBytes(ctx context.Context, key string, data []byte) error
}

type Reader interface {
	Read(ctx context.Context) (string, error)
}

// Структура процессора
type Processor struct {
	repo Repo
	st   Storage
}

// Конструктор
func NewProcessor(repo Repo, st Storage) *Processor {
	return &Processor{
		repo: repo,
		st:   st,
	}
}

// Start запускает цикл прослушивания очереди
func (p *Processor) Start(ctx context.Context, queue Reader) error {
	log.Println("Worker started...")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			id, err := queue.Read(ctx)
			if err != nil {
				log.Printf("Error reading from queue: %v", err)
				continue
			}

			log.Printf("Processing image ID: %s", id)
			if err := p.handle(ctx, id); err != nil {
				log.Printf("Error processing %s: %v", id, err)
			}
		}
	}
}

// handle — внутренняя логика обработки изображения
func (p *Processor) handle(ctx context.Context, id string) error {
	// 1. Получаем запись из базы
	dbImg, err := p.repo.Get(ctx, id)
	if err != nil {
		_ = p.repo.UpdateStatus(ctx, id, string(models.StatusFailed), "")
		return err
	}

	_ = p.repo.UpdateStatus(ctx, id, string(models.StatusProcessing), "")

	// 2. Скачиваем данные из S3
	data, err := p.st.Download(ctx, dbImg.OriginalPath)
	if err != nil {
		_ = p.repo.UpdateStatus(ctx, id, string(models.StatusFailed), "")
		return err
	}

	// 3. Декодируем изображение
	srcImg, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		_ = p.repo.UpdateStatus(ctx, id, string(models.StatusFailed), "")
		return err
	}

	// 4. Ресайз (основная копия)
	processed := imaging.Resize(srcImg, 800, 0, imaging.Lanczos)
	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, processed, imaging.JPEG); err != nil {
		return err
	}

	processedKey := "processed/" + id + ".jpg"
	if err := p.st.UploadBytes(ctx, processedKey, buf.Bytes()); err != nil {
		return err
	}

	// 5. Создание миниатюры
	thumb := imaging.Thumbnail(srcImg, 150, 150, imaging.Lanczos)
	thumbBuf := new(bytes.Buffer)
	if err := imaging.Encode(thumbBuf, thumb, imaging.JPEG); err != nil {
		return err
	}

	thumbKey := "thumbnails/" + id + ".jpg"
	if err := p.st.UploadBytes(ctx, thumbKey, thumbBuf.Bytes()); err != nil {
		return err
	}

	// 6. Успешное завершение
	return p.repo.UpdateStatus(ctx, id, string(models.StatusDone), processedKey)
}
