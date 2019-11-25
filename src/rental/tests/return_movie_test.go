package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReturningRentedMovies(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	facade := currentFacadeBuilder(usersFacade, moviesFacade, noFeesFacade)
	err := facade.Rent(userID, movieID)
	require.NoError(t, err)
	err = facade.Rent(userID, movieID2)
	require.NoError(t, err)

	// when
	rents, err := facade.GetRented(userID)

	// then
	require.NoError(t, err)
	assert.ElementsMatch(t, []int{movieID, movieID2}, getMoviesIDs(rents.Movies))

	// when
	err = facade.Return(userID, movieID)

	// then
	require.NoError(t, err)
	assert.EqualValues(t, []int{movieID2}, rentedMoviesIDs(facade, userID))
}

func TestErrorReturningMovieAsNotExistingUser(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{})
	facade := currentFacadeBuilder(usersFacade, movies.NewFacadeStub([]movies.MovieDTO{}), noFeesFacade)
	// when
	err := facade.Return(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, users.UserNotFound{}, errors.Cause(err))
}

func TestErrorReturningMovieNotRentedPreviously(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1})
	facade := currentFacadeBuilder(usersFacade, moviesFacade, noFeesFacade)
	// when
	err := facade.Return(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, rental.MovieIsNotRented{}, errors.Cause(err))
}
