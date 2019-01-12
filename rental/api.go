package rental

import "time"

type RentedMoviesDTO struct {
	Movies []RentedMovieDTO
}

type RentedMovieDTO struct {
	MovieID  int
	RentedAt time.Time
	ReturnAt time.Time
}

type RentalFacade interface {
	Rent(userID int, movieID int) error
	GetRented(userID int) (RentedMoviesDTO, error)
}
