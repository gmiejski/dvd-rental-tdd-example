package integration_tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental/api"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
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
var feesFacade = fees.Build()
var config = rental_crud.IntegrationTestConfig(rental.StandardConfig())
var facade = rental_crud.Build(usersFacade, moviesFacade, feesFacade, infrastructure.NewPostgresRepository(config.PostgresDSN), config)

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

func getRentedMovies(userID int) rental.RentedMoviesDTO {
	rs, err := http.Get(fmt.Sprintf("http://localhost:8000/users/%d/rentals", userID))
	panicOnError(err)
	rentedMovies := rental.RentedMoviesDTO{}
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
	db, err := sql.Open("postgres", rental_crud.IntegrationTestConfig(rental.StandardConfig()).PostgresDSN)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE rented_movies")
	if err != nil {
		panic(err.Error())
	}
}

func getMoviesIDs(rentedMovies []rental.RentedMovieDTO) []int {
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
