package domain

import "fmt"

type UnpaidFees struct {
	userID int
	movies []int
}

func (err UnpaidFees) Error() string {
	return fmt.Sprintf("User %d has unpaid fees for movies %v", err.userID, err.movies)
}

type MovieIsNotRented struct {
	userID  int
	movieID int
}

func (err MovieIsNotRented) Error() string {
	return fmt.Sprintf("Movie %d not rented by user %d", err.movieID, err.userID)
}

type MaximumMoviesRented struct {
	userID int
	max    int
}

func (err MaximumMoviesRented) Error() string {
	return fmt.Sprintf("User %d cannot rent more than %d mvoies", err.userID, err.max)
}
