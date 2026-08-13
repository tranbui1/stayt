package models

import (
	"time"
)

type FriendShip struct {
	ID          int       // ID of friendship
	UserID      int       // User id of the sender
	FriendID    int       // User id of the receiver
	CreatedAt   time.Time // Date of friendship formed
	Status      string    // "accepted", "rejected", "pending"
	RequestedBy int       // Contains ID of the user who sent the request
}
