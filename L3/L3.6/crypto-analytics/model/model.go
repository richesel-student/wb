package model

import (
	"errors"
	"time"
)

type Transaction struct {
	ID        int       `json:"id"`
	AssetName string    `json:"asset_name"`
	Type      string    `json:"type"`
	AmountUSD float64   `json:"amount_usd"`
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Fee       float64   `json:"fee"`
	Timestamp time.Time `json:"timestamp"`
}

func (t *Transaction) Validate() error {
	if t.AmountUSD <= 0 {
		return errors.New("amount_usd must be > 0")
	}
	if t.Type == "" {
		return errors.New("type required")
	}
	return nil
}
