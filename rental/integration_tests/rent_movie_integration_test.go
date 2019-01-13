package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/api"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const userID = 10
const movieID = 1272
const movieID2 = 1273

var adult = users.UserDTO{ID: userID, Age: 25, Name: "Greg"}

var movie1 = movies.MovieDTO{ID: movieID, Title: "something", Year: 2000, MinimalAge: 0, Genre: "horror"}
var movie2 = movies.MovieDTO{ID: movieID2, Title: "family fun", Year: 2010, MinimalAge: 0, Genre: "family"}

var usersFacade = users.NewFacadeStub([]users.UserDTO{adult})
var moviesFacade = movies.NewFacadeStub([]movies.MovieDTO{movie1, movie2})
var facade = BuildIntegrationTestFacade(usersFacade, moviesFacade)

const testServerAddress = ":8000"

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		fmt.Printf("Skipping integration tests")
		return
	}

	go runServer()
	waitForServiceRunning()
	os.Exit(m.Run())
}

func TestReturningAllRentedMoviesIT(t *testing.T) {
	// given
	clearDB()
	rentMovie(userID, movieID)
	rentMovie(userID, movieID2)

	// when
	rented := getRentedMovies(userID)
	// then

	require.ElementsMatch(t, []int{movieID, movieID2}, getMoviesIDs(rented.Movies))
}

func rentMovie(userID int, movie int) {
	request := api.RentMovieRequest{MovieID: movie}
	data, err := json.Marshal(request)
	panicOnError(err)
	rs, err := http.Post(fmt.Sprintf("http://localhost:8000/users/%d/rentals", userID), "application/json", bytes.NewBuffer(data))
	panicOnError(err)
	if rs.StatusCode != http.StatusOK {
		panic(errors.Errorf("Wrong status code: %d", rs.StatusCode))
	}
}

func getRentedMovies(userID int) domain_common.RentedMoviesDTO {
	rs, err := http.Get(fmt.Sprintf("http://localhost:8000/users/%d/rentals", userID))
	panicOnError(err)
	rentedMovies := domain_common.RentedMoviesDTO{}
	d := json.NewDecoder(rs.Body)
	err = d.Decode(&rentedMovies)
	panicOnError(err)
	return rentedMovies
}

func panicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func clearDB() {
	db, err := sql.Open("postgres", domain_crud.TestConfig().PostgresDSN)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE rented_movies")
	if err != nil {
		panic(err.Error())
	}
}

func BuildIntegrationTestFacade(usersFacade users.UsersFacade, moviesFacade movies.Facade) domain_common.RentalFacade {
	feesStub := fees.NewFacadeStub()

	config := domain_crud.TestConfig()

	return domain_crud.BuildFacade(usersFacade, moviesFacade, &feesStub, infrastructure.NewPostgresRepository(config.PostgresDSN), config)
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

func waitForServiceRunning() {
	healthy := Eventually(serviceRunning)
	if !healthy {
		panic("Unable to connect to server")
	}
}

func serviceRunning() bool {
	var url = "http://localhost:8000/"
	resp, err := http.Get(url)
	return err == nil && resp.StatusCode == http.StatusOK
}
func runServer() {
	router := mux.NewRouter()
	err := api.SetupHandlers(router, facade)
	if err != nil {
		panic(err.Error())
	}
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(testServerAddress, router)
}

func Eventually(functionToCheck func() bool) bool {
	timer := time.NewTimer(5 * time.Second)
	ticker := time.NewTicker(time.Millisecond * 20)
	checkPassed := make(chan bool)
	defer timer.Stop()
	defer ticker.Stop()

	go func() {
		checkPassed <- functionToCheck()
	}()
	for {
		select {
		case <-timer.C:
			return false
		case result := <-checkPassed:
			if result {
				return true
			}
		case <-ticker.C:
			go func() {
				checkPassed <- functionToCheck()
			}()
		}
	}
}
