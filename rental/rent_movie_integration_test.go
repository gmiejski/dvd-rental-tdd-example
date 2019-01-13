package rental

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		fmt.Printf("Skipping integration tests")
		return
	}
	os.Exit(m.Run())
}

func TestRentingSingleMovieIT(t *testing.T) {
	// given
	clearDB()
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1})
	facade := BuildIntegrationTestFacade(usersFacade, moviesFacade)

	// when
	err := facade.Rent(userID, movieID)

	// then
	require.NoError(t, err)
	require.ElementsMatch(t, []int{movieID}, rentedMoviesIDs(facade, userID))
}

func TestReturningAllRentedMoviesIT(t *testing.T) { // TODO move together with HTTP
	// given
	clearDB()
	usersFacade := users.NewFacadeStub([]users.UserDTO{adult})
	moviesFacade := movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
	facade := BuildIntegrationTestFacade(usersFacade, moviesFacade)
	err := facade.Rent(userID, movieID)
	require.NoError(t, err)
	err = facade.Rent(userID, movieID2)
	require.NoError(t, err)
	// when // TODO fix this :D
	// then
	require.NoError(t, err)
	require.ElementsMatch(t, []int{movieID, movieID2}, rentedMoviesIDs(facade, userID))
}

func clearDB() {
	db, err := sql.Open("postgres", TestConfig().postgresDSN)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE rented_movies")
	if err != nil {
		panic(err.Error())
	}
}
