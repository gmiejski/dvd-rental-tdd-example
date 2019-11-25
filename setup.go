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
) rental.Facade {
	config := rental_crud.ProdConfig(rental.StandardConfig())

	repository := infrastructure.NewPostgresRepository(config.PostgresDSN)

	return rental_crud.Build(usersFacade, moviesFacade, feesFacade, repository, config)
}

// TODO update docs, because some things must have changed
