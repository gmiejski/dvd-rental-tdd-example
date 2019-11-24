package domain_es

import "time"

type Event interface {
}

type MovieRentedEvent struct {
	// TODO add returnAt
	UserID   int
	MovieID  int
	RentedAt time.Time
}

type MovieReturnedEvent struct {
	UserID  int
	MovieID int
}
