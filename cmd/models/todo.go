package models

import (
	"time"
)

type Todo struct {
	ID        *int      `json:"id,omitempty" bson:"_id,omitempty"`
	Task      *string   `json:"task"`
	Completed *bool     `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
