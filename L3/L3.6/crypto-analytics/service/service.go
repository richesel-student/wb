package service

import (
	"crypto-analytics/bybit"
	"crypto-analytics/model"
	"crypto-analytics/repository"
	"fmt"
)

type Service struct {
	repo  *repository.Repo
	bybit *bybit.Client
}

func NewService(r *repository.Repo, b *bybit.Client) *Service {
	return &Service{repo: r, bybit: b}
}

func (s *Service) CreateTransaction(t model.Transaction) error {
	price := s.bybit.GetPrice()
	if price == 0 {
		return fmt.Errorf("failed to get price")
	}

	t.Price = price
	t.Quantity = t.AmountUSD / price
	t.Fee = t.AmountUSD * 0.001

	return s.repo.Create(t)
}
func (s *Service) GetAll(from, to string) ([]model.Transaction, error) {
	return s.repo.GetAll(from, to)
}

func (s *Service) GetAnalytics(from, to string) (repository.Analytics, error) {
	return s.repo.GetAnalytics(from, to)
}

func (s *Service) Update(id int, t model.Transaction) error {
	return s.repo.Update(id, t)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *Service) GetChart(from, to string) ([]repository.ChartPoint, error) {
	return s.repo.GetChart(from, to)
}
