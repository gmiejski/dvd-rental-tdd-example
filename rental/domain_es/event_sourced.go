package domain_es

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
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

type eventSourcedFacade struct {
	users  users.UsersFacade
	movies movies.Facade
	//repository Repository
	//config     Config
	fees fees.Facade
}

func (*eventSourcedFacade) Rent(userID int, movieID int) error {
	panic("implement me")
}

func (*eventSourcedFacade) GetRented(userID int) (domain_common.RentedMoviesDTO, error) {
	panic("implement me")
}

func (*eventSourcedFacade) Return(userID int, movieID int) error {
	panic("implement me")
}
