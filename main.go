package main

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental/api"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	usersFacade := users.Build()
	moviesFacade := movies.Build()
	feesFacade := fees.Build()

	rentalFacade := Build(usersFacade, moviesFacade, &feesFacade)

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
