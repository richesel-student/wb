package repository

import (
	"sync"

	"comments-service/internal/model"
)

type CommentRepository interface {
	Create(comment *model.Comment) error
	GetAll() ([]*model.Comment, error)
	Delete(id int64) error
}

type memoryRepository struct {
	mu       sync.RWMutex
	comments map[int64]*model.Comment
	nextID   int64
}

func NewMemoryRepository() CommentRepository {
	return &memoryRepository{
		comments: make(map[int64]*model.Comment),
		nextID:   1,
	}
}

func (r *memoryRepository) Create(comment *model.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	comment.ID = r.nextID
	r.nextID++

	r.comments[comment.ID] = comment
	return nil
}

func (r *memoryRepository) GetAll() ([]*model.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var res []*model.Comment
	for _, c := range r.comments {
		res = append(res, c)
	}
	return res, nil
}

func (r *memoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.deleteRecursive(id)
	return nil
}

func (r *memoryRepository) deleteRecursive(id int64) {
	for _, c := range r.comments {
		if c.ParentID != nil && *c.ParentID == id {
			r.deleteRecursive(c.ID)
		}
	}
	delete(r.comments, id)
}
