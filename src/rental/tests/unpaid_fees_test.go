package rental

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

func TestErrorWhenUserHasUnpaidFees(t *testing.T) {
	// given
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	feesFacade := fees.Build()
	feesFacade.AddFee(userID, movieID, time.Now(), time.Now().Add(time.Hour), 100.00)
	facade := currentFacadeBuilder(usersFacade, moviesFacade, feesFacade)

	// when
	err := facade.Rent(userID, movieID2)

	// then
	require.Error(t, err)
	require.IsType(t, rental.UnpaidFees{}, errors.Cause(err))
}
