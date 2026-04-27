package repository

import (
	"crypto-analytics/model"
	"database/sql"
)

type Analytics struct {
	Sum    float64 `json:"sum"`
	Avg    float64 `json:"avg"`
	Count  int     `json:"count"`
	Median float64 `json:"median"`
	P90    float64 `json:"p90"`
}
type ChartPoint struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db}
}

// CREATE
func (r *Repo) Create(t model.Transaction) error {
	_, err := r.db.Exec(`
	INSERT INTO transactions 
	(asset_name, type, amount_usd, price, quantity, fee)
	VALUES ($1,$2,$3,$4,$5,$6)
	`, t.AssetName, t.Type, t.AmountUSD, t.Price, t.Quantity, t.Fee)

	return err
}

// READ (с фильтром по датам)
func (r *Repo) GetAll(from, to string) ([]model.Transaction, error) {
	query := `
	SELECT id, asset_name, type, amount_usd, price, quantity, fee, timestamp
	FROM transactions
	`

	var rows *sql.Rows
	var err error

	if from != "" && to != "" {
		query += " WHERE timestamp BETWEEN $1 AND $2"
		rows, err = r.db.Query(query, from, to)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Transaction

	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(
			&t.ID,
			&t.AssetName,
			&t.Type,
			&t.AmountUSD,
			&t.Price,
			&t.Quantity,
			&t.Fee,
			&t.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, t)
	}

	return list, nil
}

// UPDATE
func (r *Repo) Update(id int, t model.Transaction) error {
	_, err := r.db.Exec(`
		UPDATE transactions 
		SET amount_usd=$1 
		WHERE id=$2
	`, t.AmountUSD, id)
	return err
}

// DELETE
func (r *Repo) Delete(id int) error {
	_, err := r.db.Exec(
		"DELETE FROM transactions WHERE id=$1",
		id,
	)
	return err
}

// ANALYTICS (с фильтром)
func (r *Repo) GetAnalytics(from, to string) (Analytics, error) {
	var a Analytics

	query := `
	SELECT 
	COALESCE(SUM(amount_usd),0),
	COALESCE(AVG(amount_usd),0),
	COUNT(*),
	COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY amount_usd),0),
	COALESCE(percentile_cont(0.9) WITHIN GROUP (ORDER BY amount_usd),0)
	FROM transactions
	`

	if from != "" && to != "" {
		query += " WHERE timestamp BETWEEN $1 AND $2"
		err := r.db.QueryRow(query, from, to).
			Scan(&a.Sum, &a.Avg, &a.Count, &a.Median, &a.P90)
		return a, err
	}

	err := r.db.QueryRow(query).
		Scan(&a.Sum, &a.Avg, &a.Count, &a.Median, &a.P90)

	return a, err

}

func (r *Repo) GetChart(from, to string) ([]ChartPoint, error) {
	query := `
	SELECT 
		DATE(timestamp) as day,
		SUM(amount_usd)
	FROM transactions
	`

	var rows *sql.Rows
	var err error

	if from != "" && to != "" {
		query += " WHERE timestamp::date BETWEEN $1 AND $2 GROUP BY day ORDER BY day"
		rows, err = r.db.Query(query, from, to)
	} else {
		query += " GROUP BY day ORDER BY day"
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ChartPoint

	for rows.Next() {
		var p ChartPoint
		rows.Scan(&p.Date, &p.Total)
		result = append(result, p)
	}

	return result, nil
}
