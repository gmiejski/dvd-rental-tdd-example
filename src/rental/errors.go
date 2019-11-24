package rental

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
