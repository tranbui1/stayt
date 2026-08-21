package models

import (
	"time"
)

type Thoughts struct {
	ID         int
	ContentUrl string
	UserID     int // ID of the user who sent it
	ReceiverID int // ID of the receiving user
	CreatedAt  time.Time
	MediaType  string // voice memo, picture, etc.
}
