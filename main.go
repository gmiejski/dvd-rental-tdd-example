package main

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/api"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	usersFacade := users.BuildUsersFacade()
	moviesFacade := movies.BuildMoviesFacade()
	feesFacade := fees.NewFacadeStub()

	rentalFacade := rental.SetupProdRentals(usersFacade, moviesFacade, &feesFacade)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "OK"}`))
	})
	err := api.SetupHandlers(router, rentalFacade)
	if err != nil {
		panic(err.Error())
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
