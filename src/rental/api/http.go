package api

import (
	"encoding/json"
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type RentMovieRequest struct {
	MovieID int `json:"movieID"`
}

func SetupHandlers(router *mux.Router, rentalFacade rental.RentalFacade) error {
	router.Handle("/users/{user}/rentals", handlers.MethodHandler{
		http.MethodPost: rentMovieHandler(rentalFacade),
		http.MethodGet:  getRentedMoviesHandler(rentalFacade),
	})
	return nil
}

func rentMovieHandler(facade rental.RentalFacade) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userID, _ := mux.Vars(request)["user"]

		user, err := strconv.Atoi(userID)

		if err != nil {
			fmt.Printf(errors.Wrap(err, "UserID must be integer").Error())

			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		rq := RentMovieRequest{}

		decoder := json.NewDecoder(request.Body)

		err = decoder.Decode(&rq)
		if err != nil {
			fmt.Printf(errors.Wrap(err, "Error decoding request").Error())
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = facade.Rent(user, rq.MovieID)
		if err != nil {
			fmt.Printf(err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func getRentedMoviesHandler(facade rental.RentalFacade) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userID, _ := mux.Vars(request)["user"]

		user, err := strconv.Atoi(userID)
		if err != nil {
			writeErrorAndLog(writer, errors.Wrap(err, "UserID must be integer"), http.StatusBadRequest)
			return
		}

		rentals, err := facade.GetRented(user)
		if err != nil {
			writeErrorAndLog(writer, err, http.StatusInternalServerError)
			return
		}

		rs, err := json.Marshal(rentals)
		if err != nil {
			writeErrorAndLog(writer, errors.Wrap(err, "Error encoding response"), http.StatusInternalServerError)
			return
		}

		writer.Write(rs)
	}
}

type httpError struct {
	Message string `json:"message"`
}

func writeErrorAndLog(writer http.ResponseWriter, err error, status int) {
	errorMessage, err := json.Marshal(httpError{Message: err.Error()})
	if err != nil {
		fmt.Printf(errors.Wrap(err, "Error encoding response").Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(status)
	writer.Write(errorMessage)
}
