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

func TestRentingSingleMovie(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1})
	facade := facadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)

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
	facade := facadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)

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
	facade := facadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)

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
	facade := facadeBuilder(usersFacade, moviesFacade, noFeesFacade, 2)
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
	facade := facadeBuilder(usersFacade, moviesFacade, noFeesFacade, 1)
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
	facade := facadeBuilder(usersFacade, moviesFacade, &feesFacade, 2)

	// when
	err := facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, domain_common.UnpaidFees{}, errors.Cause(err))
}
