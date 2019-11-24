package domain_es

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental/domain_common"
	"github.com/pkg/errors"
	"time"
)

type UserRents struct {
	UserID       int
	RentedMovies []RentedMovie
}

type RentedMovie struct {
	MovieID  int
	RentedAt time.Time
	ReturnAt time.Time
}

func NewUserRents(userID int) UserRents {
	return UserRents{UserID: userID, RentedMovies: []RentedMovie{}}
}

func (r *UserRents) rentedCount() int {
	return len(r.RentedMovies)
}

func (r *UserRents) rentMovie(movie movies.MovieDTO) ([]Event, error) {
	now := time.Now()
	if r.isMovieRented(int(movie.ID)) {
		return nil, errors.Errorf("User %d already rented movie %d", r.UserID, movie.ID)
	}
	eventsProduced := []Event{MovieRentedEvent{UserID: r.UserID, MovieID: int(movie.ID), RentedAt: now}}
	defer r.Apply(eventsProduced)
	return eventsProduced, nil
}

func (r *UserRents) returnBack(movieID int) ([]Event, error) {
	if !r.isMovieRented(movieID) {
		return nil, errors.Wrapf(domain_common.MovieIsNotRented{r.UserID, movieID}, "error returning movie")
	}
	eventsProduced := []Event{MovieReturnedEvent{UserID: r.UserID, MovieID: movieID}}
	defer r.Apply(eventsProduced)
	return eventsProduced, nil
}

func (r *UserRents) Apply(events []Event) {
	for _, event := range events {
		switch v := event.(type) {
		case MovieRentedEvent:
			e := event.(MovieRentedEvent)
			r.UserID = e.UserID
			r.RentedMovies = append(r.RentedMovies, RentedMovie{MovieID: e.MovieID, RentedAt: e.RentedAt})
		case MovieReturnedEvent:
			e := event.(MovieReturnedEvent)
			r.UserID = e.UserID
			var afterRemove []RentedMovie
			for _, rented := range r.RentedMovies {
				if rented.MovieID != e.MovieID {
					afterRemove = append(afterRemove, rented)
				}
			}
			r.RentedMovies = afterRemove
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
}

func (r *UserRents) isMovieRented(movieID int) bool {
	for _, movie := range r.RentedMovies {
		if movie.MovieID == movieID {
			return true
		}
	}
	return false
}
