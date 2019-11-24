package rental

import "time"

type Facade interface {
	Rent(userID int, movieID int) error
	GetRented(userID int) (RentedMoviesDTO, error)
	Return(userID int, movieID int) error
}

type RentedMoviesDTO struct {
	Movies []RentedMovieDTO
}

type RentedMovieDTO struct {
	MovieID  int
	RentedAt time.Time
	ReturnAt time.Time
}
