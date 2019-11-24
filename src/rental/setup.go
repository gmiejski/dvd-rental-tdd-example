package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental/domain_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
)

func Build(
	usersFacade users.Facade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
) domain_common.RentalFacade {
	config := domain_crud.ProdConfig()

	repository := infrastructure.NewPostgresRepository(config.PostgresDSN)

	return domain_crud.BuildFacade(usersFacade, moviesFacade, feesFacade, repository, config)
}
