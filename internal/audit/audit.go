package audit

import (
	"encoding/json"
	"time"
)

// UpdateItem ..
type UpdateItem struct {
	Field       string
	OriginValue interface{}
	NewValue    interface{}
}

// Content ..
type Content []UpdateItem

// String ..
func (c *Content) String() (string, error) {
	bytes, err := json.Marshal(c)
	return string(bytes), err
}

// Struct ..
func (c *Content) Struct(sc string) (*Content, error) {
	err := json.Unmarshal([]byte(sc), c)
	return c, err
}

// Rsp ..
type Rsp struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id,omitempty"`
	Type      int       `json:"type,omitempty"`
	Username  string    `json:"username,omitempty"`
	Content   Content   `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
