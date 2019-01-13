package main

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
)

func main() {

	usersFacade := users.BuildUsersFacade()
	moviesFacade := movies.BuildMoviesFacade()
	feesFacade := fees.NewFacadeStub()

	_ = rental.BuildFacade(usersFacade, moviesFacade, &feesFacade, rental.ProdConfig())

}
