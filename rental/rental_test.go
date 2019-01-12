package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const userID = 10
const movieID = 1272

func TestReturningRentedMovies(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub(
		[]users.UserDTO{{ID: userID, Age: 25, Name: "Greg"}},
	)
	moviesFacade := movies.NewFacadeStub(
		[]movies.MovieDTO{{ID: movieID, Title: "something", Year: 2000, MinimalAge: 0, Genre: "horror"}},
	)
	facade := buildTestFacade(usersFacade, moviesFacade)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.NoError(t, err)
	rentedMovies, err := facade.GetRented(userID)
	require.NoError(t, err)
	assert.EqualValues(t, []int{movieID}, getMoviesIDs(rentedMovies.Movies))
}

func TestReturnErrorWhenUserNotFound(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{})
	facade := buildTestFacade(usersFacade, moviesFacade)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, users.UserNotFound{}, errors.Cause(err))
}

func TestReturnErrorWhenMovieNotFound(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub(
		[]users.UserDTO{{ID: userID, Age: 25, Name: "Greg"}},
	)
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{})
	facade := buildTestFacade(usersFacade, moviesFacade)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, movies.MovieNotFound{}, errors.Cause(err))
}

func TestErrorWhenGettingRentsOfNotExistingUser(t *testing.T) {
	// given
	facade := buildTestFacade(users.NewFacadeStub([]users.UserDTO{}), movies.NewFacadeStub([]movies.MovieDTO{}))

	// when
	rents, err := facade.GetRented(userID)

	// then
	require.Error(t, err)
	require.IsType(t, users.UserNotFound{}, errors.Cause(err))
	require.Empty(t, rents.Movies)
}

func TestReturnEmptyRentsIfUserHasNotRentedAnythingYet(t *testing.T) { // TODO split into files
	// given
	usersFacade := users.NewFacadeStub(
		[]users.UserDTO{{ID: userID, Age: 25, Name: "Greg"}},
	)
	facade := buildTestFacade(usersFacade, movies.NewFacadeStub([]movies.MovieDTO{}))

	// when
	rents, err := facade.GetRented(userID)

	// then
	require.NoError(t, err)
	require.Empty(t, rents.Movies)
}

func getMoviesIDs(rentedMovies []RentedMovieDTO) []int {
	movieIDs := make([]int, 0)
	for _, movie := range rentedMovies {
		movieIDs = append(movieIDs, movie.MovieID)
	}
	return movieIDs
}
