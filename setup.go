package main

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
)

func Build( // TODO move it to rental module
	usersFacade users.Facade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
) rental.RentalFacade {
	config := rental_crud.ProdConfig()

	repository := infrastructure.NewPostgresRepository(config.PostgresDSN)

	return rental_crud.BuildFacade(usersFacade, moviesFacade, feesFacade, repository, config)
}
