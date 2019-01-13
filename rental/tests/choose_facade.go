package domain_crud

import (
	"database/sql"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
)

type createFacadeFunc = func(users users.UsersFacade, movies movies.Facade) domain_common.RentalFacade // TODO

var facadeBuilder = buildPostgresCrudTestFacade

func buildInMemoryCrudTestFacade(users users.UsersFacade, movies movies.Facade) domain_common.RentalFacade {
	return domain_crud.BuildUnitTestFacade(users, movies)
}

func buildPostgresCrudTestFacade(users users.UsersFacade, movies movies.Facade) domain_common.RentalFacade {
	feesStub := fees.NewFacadeStub()
	config := domain_crud.TestConfig()
	clearPostgresDB(config)
	return domain_crud.BuildFacade(users, movies, &feesStub, infrastructure.NewPostgresRepository(config.PostgresDSN), config)
}

func clearPostgresDB(config domain_crud.Config) {
	db, err := sql.Open("postgres", config.PostgresDSN)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE rented_movies")
	if err != nil {
		panic(err.Error())
	}
}
