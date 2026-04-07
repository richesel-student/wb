package repository

import (
	"database/sql"
	"log"
	"shortener/internal/domain"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) Save(u domain.URL) error {
	_, err := r.db.Exec(
		"INSERT INTO urls (original, short) VALUES (?, ?)",
		u.Original, u.Short,
	)
	return err
}

func (r *Repo) Get(short string) (string, error) {
	var original string
	err := r.db.QueryRow(
		"SELECT original FROM urls WHERE short = ?",
		short,
	).Scan(&original)

	return original, err
}

func (r *Repo) SaveClick(c domain.Click) error {
	_, err := r.db.Exec(
		"INSERT INTO clicks (short, user_agent) VALUES (?, ?)",
		c.Short, c.UserAgent,
	)
	return err
}

func (r *Repo) GetStats(short string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// total
	var total int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM clicks WHERE short = ?",
		short,
	).Scan(&total)
	if err != nil {
		return nil, err
	}
	result["total"] = total

	// by day
	rowsDay, err := r.db.Query(`
		SELECT DATE(created_at), COUNT(*)
		FROM clicks
		WHERE short = ?
		GROUP BY DATE(created_at)
	`, short)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rowsDay.Close(); err != nil {
			log.Printf("failed to close rowsDay: %v", err)
		}
	}()

	var byDay []map[string]interface{}
	for rowsDay.Next() {
		var date string
		var count int
		if err := rowsDay.Scan(&date, &count); err != nil {
			return nil, err
		}

		byDay = append(byDay, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}
	result["by_day"] = byDay

	// by month
	rowsMonth, err := r.db.Query(`
		SELECT strftime('%Y-%m', created_at), COUNT(*)
		FROM clicks
		WHERE short = ?
		GROUP BY 1
	`, short)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rowsMonth.Close(); err != nil {
			log.Printf("failed to close rowsMonth: %v", err)
		}
	}()

	var byMonth []map[string]interface{}
	for rowsMonth.Next() {
		var month string
		var count int
		if err := rowsMonth.Scan(&month, &count); err != nil {
			return nil, err
		}

		byMonth = append(byMonth, map[string]interface{}{
			"month": month,
			"count": count,
		})
	}
	result["by_month"] = byMonth

	// by user agent
	rowsUA, err := r.db.Query(`
		SELECT user_agent, COUNT(*)
		FROM clicks
		WHERE short = ?
		GROUP BY user_agent
	`, short)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rowsUA.Close(); err != nil {
			log.Printf("failed to close rowsUA: %v", err)
		}
	}()

	var byUA []map[string]interface{}
	for rowsUA.Next() {
		var ua string
		var count int
		if err := rowsUA.Scan(&ua, &count); err != nil {
			return nil, err
		}

		byUA = append(byUA, map[string]interface{}{
			"user_agent": ua,
			"count":      count,
		})
	}
	result["by_ua"] = byUA

	return result, nil
}
