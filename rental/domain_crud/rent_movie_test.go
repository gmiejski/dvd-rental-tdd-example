package domain_crud

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
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

func TestRentingSingleMovie(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1})
	facade := BuildUnitTestFacade(usersFacade, moviesFacade)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.NoError(t, err)
	require.ElementsMatch(t, []int{movieID}, rentedMoviesIDs(facade, userID))
}

func TestReturnErrorWhenRentingAsNotExistingUser(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{})
	facade := BuildUnitTestFacade(usersFacade, moviesFacade)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.Error(t, err)
	require.IsType(t, users.UserNotFound{}, errors.Cause(err))
}

func TestReturnErrorWhenRentingNotExistingMovie(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{})
	facade := BuildUnitTestFacade(usersFacade, moviesFacade)

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
	facade := BuildUnitTestFacade(usersFacade, moviesFacade)
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
	facade := BuildUnitTestFacade(usersFacade, moviesFacade, withConfig(Config{MaxRentedMoviesCount: 1}))
	err := facade.Rent(userID, movieID)
	require.NoError(t, err)

	// when
	err = facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, domain_common.MaximumMoviesRented{}, errors.Cause(err))
}

func TestErrorWhenUserHasUnpaidFees(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	feesFacade := fees.NewFacadeStub()
	feesFacade.AddFee(userID, movieID, time.Now(), time.Now().Add(time.Hour), 100.00)
	facade := BuildUnitTestFacade(usersFacade, moviesFacade, withFeesFacade(&feesFacade))

	// when
	err := facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, domain_common.UnpaidFees{}, errors.Cause(err))
}

func rentedMoviesIDs(facade domain_common.RentalFacade, userID int) []int {
	rentedMovies, err := facade.GetRented(userID)
	if err != nil {
		panic(err.Error())
	}
	return getMoviesIDs(rentedMovies.Movies)
}

func getMoviesIDs(rentedMovies []domain_common.RentedMovieDTO) []int {
	movieIDs := make([]int, 0)
	for _, movie := range rentedMovies {
		movieIDs = append(movieIDs, movie.MovieID)
	}
	return movieIDs
}
