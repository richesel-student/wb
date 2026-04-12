package service

import (
	"comments-service/internal/repository"
	"testing"
)

func setup() *CommentService {
	repo := repository.NewMemoryRepository()
	return NewCommentService(repo)
}

func TestCreateComment(t *testing.T) {
	s := setup()

	c, err := s.Create("hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	if c.ID == 0 {
		t.Fatal("expected ID to be set")
	}

	if c.Text != "hello" {
		t.Fatal("text mismatch")
	}
}

func TestTreeStructure(t *testing.T) {
	s := setup()

	root, _ := s.Create("root", nil)
	child, _ := s.Create("child", &root.ID)

	tree, _ := s.GetTree(nil, 10, 0, "")

	if len(tree) != 1 {
		t.Fatalf("expected 1 root, got %d", len(tree))
	}

	if len(tree[0].Children) != 1 {
		t.Fatal("expected child comment")
	}

	if tree[0].Children[0].ID != child.ID {
		t.Fatal("child mismatch")
	}
}

func TestDeleteRecursive(t *testing.T) {
	s := setup()

	root, _ := s.Create("root", nil)
	_, _ = s.Create("child", &root.ID)

	_ = s.Delete(root.ID)

	tree, _ := s.GetTree(nil, 10, 0, "")

	if len(tree) != 0 {
		t.Fatal("expected empty tree after delete")
	}
}

func TestSearch(t *testing.T) {
	s := setup()

	s.Create("hello world", nil)
	s.Create("golang", nil)

	res, _ := s.GetTree(nil, 10, 0, "hello")

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
}

func TestDeepTree(t *testing.T) {
	s := setup()

	var parentID *int64

	for i := 0; i < 5; i++ {
		c, _ := s.Create("lvl", parentID)
		parentID = &c.ID
	}

	tree, _ := s.GetTree(nil, 10, 0, "")

	current := tree[0]
	depth := 1

	for len(current.Children) > 0 {
		current = current.Children[0]
		depth++
	}

	if depth != 5 {
		t.Fatalf("expected depth 5, got %d", depth)
	}
}
