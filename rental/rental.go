package rental

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/pkg/errors"
	"time"
)

type UnpaidFees struct {
	userID int
	movies []int
}

func (err UnpaidFees) Error() string {
	return fmt.Sprintf("User %d has unpaid fees for movies %v", err.userID, err.movies)
}

type MovieIsNotRented struct {
	userID  int
	movieID int
}

func (err MovieIsNotRented) Error() string {
	return fmt.Sprintf("Movie %d not rented by user %d", err.movieID, err.userID)
}

type MaximumMoviesRented struct {
	userID int
	max    int
}

func (err MaximumMoviesRented) Error() string {
	return fmt.Sprintf("User %d cannot rent more than %d mvoies", err.userID, err.max)
}

type UserRents struct {
	userID       int
	age          int
	rentedMovies []rentedMovie
}

func (r *UserRents) rentMovie(movie movies.MovieDTO) error { // TODO add movies anti corruption layer
	now := time.Now()
	if r.isMovieRented(int(movie.ID)) {
		return errors.Errorf("User %d already rented movie %d", r.userID, movie.ID)
	}
	r.rentedMovies = append(r.rentedMovies, rentedMovie{movieID: int(movie.ID), rentedAt: now, returnAt: now.Add(time.Hour * 24 * 3)})
	return nil
}

func (r *UserRents) returnBack(movieID int) error {
	if !r.isMovieRented(movieID) {
		return errors.Wrapf(MovieIsNotRented{r.userID, movieID}, "error returning movie")
	}
	var rentsAfterReturning []rentedMovie
	for _, rentedMovie := range r.rentedMovies {
		if rentedMovie.movieID != movieID {
			rentsAfterReturning = append(rentsAfterReturning, rentedMovie)
		}
	}
	r.rentedMovies = rentsAfterReturning
	return nil
}

func (r *UserRents) isMovieRented(movieID int) bool {
	for _, movie := range r.rentedMovies {
		if movie.movieID == movieID {
			return true
		}
	}
	return false
}
func (r *UserRents) rentedCount() int {
	return len(r.rentedMovies)
}

type rentedMovie struct {
	movieID  int
	rentedAt time.Time
	returnAt time.Time
}

type facade struct {
	users      users.UsersFacade
	movies     movies.Facade
	repository repository
	config     config
	fees       fees.Facade
}

func (f *facade) Rent(userID int, movieID int) error {
	if _, err := f.users.Get(userID); err != nil {
		return errors.Wrapf(err, "Error getting user: %d", userID) // TODO move to f.getUserRents
	}
	movie, err := f.movies.Get(movies.MovieID(movieID))
	if err != nil {
		return errors.Wrapf(err, "Error finding movie: %d", movieID)
	}

	userRents, err := f.getUserRents(userID)
	if err != nil {
		return errors.Wrapf(err, "Error getting rented movies for user: %d", userID)
	}

	if userRents.rentedCount() >= f.config.maxRentedMoviesCount {
		return errors.Wrapf(
			MaximumMoviesRented{userID: userID, max: f.config.maxRentedMoviesCount},
			"error renting movie %d by user %d", movieID, userID,
		)
	}

	if fees, _ := f.fees.GetFees(userID); len(fees.Fees) > 0 {
		return errors.Wrapf(
			UnpaidFees{userID: userID, movies: fees.OverrentMovieIDs()},
			"error renting movie %d",
			movieID,
		)
	}

	err = userRents.rentMovie(movie)
	if err != nil {
		return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
	}
	err = f.repository.Save(userRents)
	return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
}

func newUserRents(userID int) UserRents {
	return UserRents{userID: userID, rentedMovies: []rentedMovie{}}
}

func (f *facade) GetRented(userID int) (RentedMoviesDTO, error) { // TODO rename to Rents
	if _, err := f.users.Get(userID); err != nil {
		return RentedMoviesDTO{}, errors.Wrapf(err, "Error getting user: %d", userID)
	}
	rents, err := f.getUserRents(userID)
	if err != nil {
		return RentedMoviesDTO{}, errors.WithMessagef(err, "Error getting rented movies for user %d", userID)
	}
	return toDTO(rents), nil
}

func (f *facade) Return(userID int, movieID int) error {
	if _, err := f.users.Get(userID); err != nil {
		return errors.Wrapf(err, "Error getting user: %d", userID)
	}

	userRents, err := f.getUserRents(userID)
	if err != nil {
		return errors.Wrapf(err, "Error getting rented movies for user: %d", userID)
	}

	err = userRents.returnBack(movieID)
	if err != nil {
		return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
	}
	err = f.repository.Save(userRents)
	return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
}

func toDTO(rents UserRents) RentedMoviesDTO {
	rentedMovies := make([]RentedMovieDTO, 0)
	for _, movie := range rents.rentedMovies {
		rentedMovies = append(rentedMovies, toMovieDTO(movie))
	}
	return RentedMoviesDTO{Movies: rentedMovies}
}

func toMovieDTO(movie rentedMovie) RentedMovieDTO {
	return RentedMovieDTO{MovieID: movie.movieID, RentedAt: movie.rentedAt, ReturnAt: movie.returnAt}
}

func (f *facade) getUser(userID int) bool {
	_, err := f.users.Get(userID)
	return err != nil

}
func (f *facade) getUserRents(userID int) (UserRents, error) {
	userRents, err := f.repository.Get(userID)
	if err != nil {
		return UserRents{}, err
	}
	if userRents == nil {
		return newUserRents(userID), nil
	}
	return *userRents, nil
}

type testOptionFacade = func(*facade)

var withConfig = func(c config) testOptionFacade {
	return func(f *facade) {
		f.config = c
	}
}

var withFeesFacade = func(feesFacade fees.Facade) testOptionFacade {
	return func(f *facade) {
		f.fees = feesFacade
	}
}

func buildTestFacade(users users.UsersFacade, movies movies.Facade, options ...testOptionFacade) RentalFacade {
	fees := fees.NewFacadeStub()
	baseTestFacade := &facade{
		users:      users,
		movies:     movies,
		fees:       &fees,
		repository: newInMemoryRepository(),
		config:     config{maxRentedMoviesCount: 10},
	}

	for _, option := range options {
		option(baseTestFacade)
	}

	return baseTestFacade
}
