package rental_crud

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const userID = 10
const movieID = 1272
const movieID2 = 1273

var adult = users.UserDTO{ID: userID, Age: 25, Name: "Greg"}

var movie1 = movies.MovieDTO{ID: movieID, Title: "something", Year: 2000, MinimalAge: 0, Genre: "horror"}
var movie2 = movies.MovieDTO{ID: movieID2, Title: "family fun", Year: 2010, MinimalAge: 0, Genre: "family"}

func TestErrorWhenUserHasUnpaidFees(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	feesFacade := fees.Build()
	feesFacade.AddFee(userID, movieID, time.Now(), time.Now().Add(time.Hour), 100.00)
	facade := BuildUnitTestFacade(usersFacade, moviesFacade, withFeesFacade(&feesFacade))

	// when
	err := facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, rental.UnpaidFees{}, errors.Cause(err))
}

func TestCannotRentMoreMoviesThanMaximum(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	facade := BuildUnitTestFacade(usersFacade, moviesFacade, withConfig(Config{MaxRentedMoviesCount: 1}))
	err := facade.Rent(userID, movieID)
	require.NoError(t, err)

	// when
	err = facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, rental.MaximumMoviesRented{}, errors.Cause(err))
}
