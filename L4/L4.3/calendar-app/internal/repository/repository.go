package repository

import (
	"context"
	"sync"

	"calendar-app/internal/models"
)

type Repository struct {
	mu          sync.RWMutex
	events      map[string]*models.Event // active events
	archived    map[string]*models.Event // archived events
}

func NewRepository() *Repository {
	return &Repository{
		events:   make(map[string]*models.Event),
		archived: make(map[string]*models.Event),
	}
}

func (r *Repository) Create(ctx context.Context, event *models.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mu.Lock()
		defer r.mu.Unlock()

		r.events[event.ID] = event
		return nil
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*models.Event, bool, error) {
	select {
	case <-ctx.Done():
		return nil, false, ctx.Err()
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()

		// сначала проверяем активные события
		if event, ok := r.events[id]; ok {
			return event, true, nil
		}

		// затем архивные
		if event, ok := r.archived[id]; ok {
			return event, true, nil
		}

		return nil, false, nil
	}
}

func (r *Repository) GetAll(ctx context.Context) ([]*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()

		events := make([]*models.Event, 0, len(r.events))
		for _, event := range r.events {
			events = append(events, event)
		}

		return events, nil
	}
}

func (r *Repository) Update(ctx context.Context, id string, updatedEvent *models.Event) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		r.mu.Lock()
		defer r.mu.Unlock()

		if _, ok := r.events[id]; ok {
			updatedEvent.ID = id
			r.events[id] = updatedEvent
			return true, nil
		}

		return false, nil
	}
}

func (r *Repository) Delete(ctx context.Context, id string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		r.mu.Lock()
		defer r.mu.Unlock()

		if _, ok := r.events[id]; ok {
			delete(r.events, id)
			return true, nil
		}

		if _, ok := r.archived[id]; ok {
			delete(r.archived, id)
			return true, nil
		}

		return false, nil
	}
}

func (r *Repository) Archive(ctx context.Context, id string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		r.mu.Lock()
		defer r.mu.Unlock()

		if event, ok := r.events[id]; ok {
			event.Archived = true
			r.archived[id] = event
			delete(r.events, id)
			return true, nil
		}

		return false, nil
	}
}

func (r *Repository) GetActive(ctx context.Context) ([]*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()

		events := make([]*models.Event, 0, len(r.events))
		for _, event := range r.events {
			events = append(events, event)
		}

		return events, nil
	}
}

func (r *Repository) GetArchived(ctx context.Context) ([]*models.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mu.RLock()
		defer r.mu.RUnlock()

		events := make([]*models.Event, 0, len(r.archived))
		for _, event := range r.archived {
			events = append(events, event)
		}

		return events, nil
	}
}