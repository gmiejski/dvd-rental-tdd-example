package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
)

func SetupProdRentals(
	usersFacade users.UsersFacade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
) domain.RentalFacade {
	config := domain.ProdConfig()

	repository := infrastructure.NewPostgresRepository(config.PostgresDSN)

	return domain.BuildFacade(usersFacade, moviesFacade, feesFacade, repository, config)
}
