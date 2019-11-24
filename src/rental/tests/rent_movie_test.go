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

func TestRentingSingleMovie(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1})
	facade := currentFacadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.NoError(t, err)
	require.ElementsMatch(t, []int{movieID}, rentedMoviesIDs(facade, userID))
}

func TestErrorWhenRentingAsNotExistingUser(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{})
	facade := currentFacadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, users.UserNotFound{}, errors.Cause(err))
}

func TestErrorWhenRentingNotExistingMovie(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{})
	facade := currentFacadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, movies.MovieNotFound{}, errors.Cause(err))
}

func TestErrorWhenRentingSameMovieTwice(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1})
	facade := currentFacadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)
	err := facade.Rent(userID, movieID)
	require.NoError(t, err)

	// when trying to rent same movie second time
	err = facade.Rent(userID, movieID)

	// then
	require.Error(t, err)
}

func TestCannotRentMoreMoviesThanMaximum(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	facade := currentFacadeBuilder(usersFacade, moviesFacade, noFeesFacade, 1)
	err := facade.Rent(userID, movieID)
	require.NoError(t, err)

	// when
	err = facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, rental.MaximumMoviesRented{}, errors.Cause(err))
}

func TestErrorWhenUserHasUnpaidFees(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	feesFacade := fees.Build()
	feesFacade.AddFee(userID, movieID, time.Now(), time.Now().Add(time.Hour), 100.00)
	facade := currentFacadeBuilder(usersFacade, moviesFacade, &feesFacade, 2)

	// when
	err := facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, rental.UnpaidFees{}, errors.Cause(err))
}
