package domain_es

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
)

func BuildFacade(
	usersFacade users.UsersFacade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
	repository Repository,
	maxRented int,
) domain_common.RentalFacade {
	return &eventSourcedFacade{
		users:           usersFacade,
		movies:          moviesFacade,
		fees:            feesFacade,
		repository:      repository,
		maxRentedMovies: maxRented,
	}
}
