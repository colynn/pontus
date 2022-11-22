package db

import (
	"time"
)

// BaseModel ..
type BaseModel struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}
