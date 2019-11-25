package main

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental/api"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	usersFacade := users.Build(users.NewInMemoryRepository()) // TODO use SQL implementation
	moviesFacade := movies.Build()
	feesFacade := fees.Build()
	config := rental_crud.ProdConfig(rental.StandardConfig())

	repository := infrastructure.NewPostgresRepository(config.PostgresDSN)
	rentalFacade := rental_crud.Build(usersFacade, moviesFacade, feesFacade, repository, config)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status": "OK"}`))
	})
	err := api.SetupHandlers(router, rentalFacade)
	if err != nil {
		panic(err.Error())
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
