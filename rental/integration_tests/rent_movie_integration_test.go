package integration_tests

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const userID = 10
const movieID = 1272
const movieID2 = 1273

var adult = users.UserDTO{ID: userID, Age: 25, Name: "Greg"}

var movie1 = movies.MovieDTO{ID: movieID, Title: "something", Year: 2000, MinimalAge: 0, Genre: "horror"}
var movie2 = movies.MovieDTO{ID: movieID2, Title: "family fun", Year: 2010, MinimalAge: 0, Genre: "family"}

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
	db, err := sql.Open("postgres", domain.TestConfig().PostgresDSN)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE rented_movies")
	if err != nil {
		panic(err.Error())
	}
}

func BuildIntegrationTestFacade(usersFacade users.UsersFacade, moviesFacade movies.Facade) domain.RentalFacade {
	feesStub := fees.NewFacadeStub()

	config := domain.TestConfig()

	return domain.BuildFacade(usersFacade, moviesFacade, &feesStub, infrastructure.NewPostgresRepository(config.PostgresDSN), config)
}

func rentedMoviesIDs(facade domain.RentalFacade, userID int) []int {
	rentedMovies, err := facade.GetRented(userID)
	if err != nil {
		panic(err.Error())
	}
	return getMoviesIDs(rentedMovies.Movies)
}

func getMoviesIDs(rentedMovies []domain.RentedMovieDTO) []int {
	movieIDs := make([]int, 0)
	for _, movie := range rentedMovies {
		movieIDs = append(movieIDs, movie.MovieID)
	}
	return movieIDs
}
