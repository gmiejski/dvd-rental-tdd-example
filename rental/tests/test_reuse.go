package domain_crud

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
)

const userID = 10
const movieID = 1272
const movieID2 = 1273

var adult = users.UserDTO{ID: userID, Age: 25, Name: "Greg"}

var movie1 = movies.MovieDTO{ID: movieID, Title: "something", Year: 2000, MinimalAge: 0, Genre: "horror"}
var movie2 = movies.MovieDTO{ID: movieID2, Title: "family fun", Year: 2010, MinimalAge: 0, Genre: "family"}

func getMoviesIDs(rentedMovies []domain_common.RentedMovieDTO) []int {
	movieIDs := make([]int, 0)
	for _, movie := range rentedMovies {
		movieIDs = append(movieIDs, movie.MovieID)
	}
	return movieIDs
}

func rentedMoviesIDs(facade domain_common.RentalFacade, userID int) []int {
	rentedMovies, err := facade.GetRented(userID)
	if err != nil {
		panic(err.Error())
	}
	return getMoviesIDs(rentedMovies.Movies)
}

var noFeesFacade = func() fees.Facade { feesStub := fees.NewFacadeStub(); return &feesStub }()
