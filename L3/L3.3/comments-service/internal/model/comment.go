package model

import "time"

type Comment struct {
	ID        int64      `json:"id"`
	ParentID  *int64     `json:"parent_id,omitempty"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	Children  []*Comment `json:"children,omitempty"`
}
