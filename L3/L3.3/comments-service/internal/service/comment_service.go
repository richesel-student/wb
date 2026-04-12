package service

import (
	"comments-service/internal/model"
	"comments-service/internal/repository"
	"sort"
	"strings"
	"time"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(r repository.CommentRepository) *CommentService {
	return &CommentService{repo: r}
}

func (s *CommentService) Create(text string, parentID *int64) (*model.Comment, error) {
	comment := &model.Comment{
		Text:      text,
		ParentID:  parentID,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(comment)
	return comment, err
}

func (s *CommentService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *CommentService) GetTree(parentID *int64, limit, offset int, search string) ([]*model.Comment, error) {
	all, _ := s.repo.GetAll()

	// фильтр по тексту
	var filtered []*model.Comment
	for _, c := range all {
		if search == "" || strings.Contains(strings.ToLower(c.Text), strings.ToLower(search)) {
			filtered = append(filtered, c)
		}
	}

	// сортировка по дате
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
	})

	tree := buildTree(filtered, parentID)

	// пагинация
	if offset > len(tree) {
		return []*model.Comment{}, nil
	}

	end := offset + limit
	if end > len(tree) {
		end = len(tree)
	}

	return tree[offset:end], nil
}

func buildTree(comments []*model.Comment, parentID *int64) []*model.Comment {
	var res []*model.Comment

	for _, c := range comments {
		if (c.ParentID == nil && parentID == nil) ||
			(c.ParentID != nil && parentID != nil && *c.ParentID == *parentID) {

			children := buildTree(comments, &c.ID)
			c.Children = children
			res = append(res, c)
		}
	}
	return res
}
