package models

import (
	"time"
)

type Thoughts struct {
	ID         int       `json:"id"`
	ContentUrl string    `json:"content_url"`
	UserID     int       `json:"user_id"`     // ID of the user who sent it
	ReceiverID int       `json:"receiver_id"` // ID of the receiving user
	CreatedAt  time.Time `json:"created_at"`
	MediaType  string    `json:"media_type"` // voice memo, picture, etc.
}
