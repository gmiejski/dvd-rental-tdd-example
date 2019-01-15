package domain_common

import "fmt"

type UnpaidFees struct {
	UserID int
	Movies []int
}

func (err UnpaidFees) Error() string {
	return fmt.Sprintf("User %d has unpaid fees for Movies %v", err.UserID, err.Movies)
}

type MovieIsNotRented struct {
	UserID  int
	MovieID int
}

func (err MovieIsNotRented) Error() string {
	return fmt.Sprintf("Movie %d not rented by user %d", err.MovieID, err.UserID)
}

type MaximumMoviesRented struct {
	UserID int
	Max    int
}

func (err MaximumMoviesRented) Error() string {
	return fmt.Sprintf("User %d cannot rent more than %d mvoies", err.UserID, err.Max)
}

type TooYoungError struct {
	UserID        int
	MovieID       int
	MovieAgeLimit int
}

func (err TooYoungError) Error() string {
	return fmt.Sprintf("User %d too young for movie %d with required %d age", err.UserID, err.MovieID, err.MovieAgeLimit)
}
