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
	require.IsType(t, errors.Cause(err), users.UserNotFound{})
}

func getMoviesIDs(rentedMovies []RentedMovieDTO) []int {
	movieIDs := make([]int, 0)
	for _, movie := range rentedMovies {
		movieIDs = append(movieIDs, movie.MovieID)
	}
	return movieIDs
}
