package repository

import (
	"context"
	"database/sql"
	"image-processor/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, img *models.Image) error {
	if err := img.Validate(); err != nil {
		return err
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO images (id, status, original_path) VALUES ($1,$2,$3)`,
		img.ID, img.Status, img.OriginalPath,
	)
	return err
}

func (r *Repo) Get(ctx context.Context, id string) (*models.Image, error) {
	var img models.Image
	var processed sql.NullString

	err := r.db.QueryRow(ctx,
		`SELECT id, status, original_path, processed_path FROM images WHERE id=$1`,
		id,
	).Scan(&img.ID, &img.Status, &img.OriginalPath, &processed)

	if err != nil {
		return nil, err
	}

	if processed.Valid {
		img.ProcessedPath = processed.String
	}

	return &img, nil
}

func (r *Repo) UpdateStatus(ctx context.Context, id, status, path string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE images SET status=$2, processed_path=$3, updated_at=now() WHERE id=$1`,
		id, status, path,
	)
	return err
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM images WHERE id=$1`, id)
	return err
}

func (r *Repo) GetAll(ctx context.Context) ([]*models.Image, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, status, original_path, processed_path FROM images ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*models.Image

	for rows.Next() {
		var img models.Image
		var processed sql.NullString

		err := rows.Scan(&img.ID, &img.Status, &img.OriginalPath, &processed)
		if err != nil {
			return nil, err
		}

		if processed.Valid {
			img.ProcessedPath = processed.String
		}

		images = append(images, &img)
	}

	return images, nil
}
