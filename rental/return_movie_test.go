package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReturningRentedMovies(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	facade := buildTestFacade(usersFacade, moviesFacade)
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
	facade := buildTestFacade(usersFacade, movies.NewFacadeStub([]movies.MovieDTO{}))

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
	facade := buildTestFacade(usersFacade, moviesFacade)

	// when
	err := facade.Return(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, MovieIsNotRented{}, errors.Cause(err))
}

func rentedMoviesIDs(facade RentalFacade, userID int) []int {
	rentedMovies, err := facade.GetRented(userID)
	if err != nil {
		panic(err.Error())
	}
	return getMoviesIDs(rentedMovies.Movies)
}
