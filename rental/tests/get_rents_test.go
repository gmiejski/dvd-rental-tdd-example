package domain_crud

import (
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestErrorWhenGettingRentsOfNotExistingUser(t *testing.T) {
	// given
	facade := currentFacadeBuilder(users.NewFacadeStub([]users.UserDTO{}), movies.NewFacadeStub([]movies.MovieDTO{}), noFeesFacade, 2)

	// when
	rents, err := facade.GetRented(userID)

	// then
	require.Error(t, err)
	require.IsType(t, users.UserNotFound{}, errors.Cause(err))
	require.Empty(t, rents.Movies)
}

func TestReturnEmptyRentsIfUserHasNotRentedAnythingYet(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub(
		[]users.UserDTO{adult},
	)
	facade := currentFacadeBuilder(usersFacade, movies.NewFacadeStub([]movies.MovieDTO{}), noFeesFacade, 2)

	// when
	rents, err := facade.GetRented(userID)

	// then
	require.NoError(t, err)
	require.Empty(t, rents.Movies)
}
