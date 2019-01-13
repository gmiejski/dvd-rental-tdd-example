package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
)

func SetupProdRentals(
	usersFacade users.UsersFacade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
) domain_common.RentalFacade {
	config := domain_crud.ProdConfig()

	repository := infrastructure.NewPostgresRepository(config.PostgresDSN)

	return domain_crud.BuildFacade(usersFacade, moviesFacade, feesFacade, repository, config)
}
