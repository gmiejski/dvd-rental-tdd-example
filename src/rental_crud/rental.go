package rental_crud

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"github.com/pkg/errors"
	"time"
)

type UserRents struct { // TODO hide inside module
	UserID       int
	RentedMovies []RentedMovie
}

var rentingTime = time.Hour * 24 * 3 // TODO should be taken from elsewhere

func (r *UserRents) rentMovie(movie movies.MovieDTO) error { // TODO add movies anti corruption layer
	if r.isMovieRented(int(movie.ID)) {
		return errors.Errorf("User %d already rented movie %d", r.UserID, movie.ID)
	}
	rentedAt := time.Now()
	returnAt := rentedAt.Add(rentingTime)
	r.RentedMovies = append(r.RentedMovies, RentedMovie{MovieID: int(movie.ID), RentedAt: rentedAt, ReturnAt: returnAt})
	return nil
}

func (r *UserRents) returnBack(movieID int) error {
	if !r.isMovieRented(movieID) {
		return errors.Wrapf(rental.MovieIsNotRented{UserID: r.UserID, MovieID: movieID}, "error returning movie")
	}
	var rentsAfterReturning []RentedMovie
	for _, rentedMovie := range r.RentedMovies {
		if rentedMovie.MovieID != movieID {
			rentsAfterReturning = append(rentsAfterReturning, rentedMovie)
		}
	}
	r.RentedMovies = rentsAfterReturning
	return nil
}

func (r *UserRents) isMovieRented(movieID int) bool {
	for _, movie := range r.RentedMovies {
		if movie.MovieID == movieID {
			return true
		}
	}
	return false
}
func (r *UserRents) rentedCount() int {
	return len(r.RentedMovies)
}

type RentedMovie struct {
	MovieID  int
	RentedAt time.Time
	ReturnAt time.Time
}

type facade struct {
	users      users.Facade
	movies     movies.Facade
	repository Repository
	config     Config
	fees       fees.Facade
}

func (f *facade) Rent(userID int, movieID int) error {
	userRents, err := f.getUserRents(userID)
	if err != nil {
		return errors.Wrapf(err, "Error getting rented movies for user: %d", userID)
	}
	movie, err := f.movies.Get(movies.MovieID(movieID))
	if err != nil {
		return errors.Wrapf(err, "Error finding movie: %d", movieID)
	}

	if userRents.rentedCount() >= f.config.MaxRentedMoviesCount {
		return errors.Wrapf(
			rental.MaximumMoviesRented{UserID: userID, Max: f.config.MaxRentedMoviesCount},
			"error renting movie %d by user %d", movieID, userID,
		)
	}

	userFees, err := f.fees.GetFees(userID)
	if err != nil {
		return errors.Wrapf(err, "Error checking user fees: %d", userID)
	}
	if len(userFees.Fees) > 0 {
		return errors.Wrapf(
			rental.UnpaidFees{UserID: userID, Movies: userFees.OverrentMovieIDs()},
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
	return UserRents{UserID: userID, RentedMovies: []RentedMovie{}}
}

func (f *facade) GetRented(userID int) (rental.RentedMoviesDTO, error) { // TODO rename to Rents
	if _, err := f.users.Find(userID); err != nil {
		return rental.RentedMoviesDTO{}, errors.Wrapf(err, "Error getting user: %d", userID)
	}
	rents, err := f.getUserRents(userID)
	if err != nil {
		return rental.RentedMoviesDTO{}, errors.WithMessagef(err, "Error getting rented movies for user %d", userID)
	}
	return toDTO(rents), nil
}

func (f *facade) Return(userID int, movieID int) error {
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

func toDTO(rents UserRents) rental.RentedMoviesDTO {
	rentedMovies := make([]rental.RentedMovieDTO, 0)
	for _, movie := range rents.RentedMovies {
		rentedMovies = append(rentedMovies, toMovieDTO(movie))
	}
	return rental.RentedMoviesDTO{Movies: rentedMovies}
}

func toMovieDTO(movie RentedMovie) rental.RentedMovieDTO {
	return rental.RentedMovieDTO{MovieID: movie.MovieID, RentedAt: movie.RentedAt, ReturnAt: movie.ReturnAt}
}

func (f *facade) getUserRents(userID int) (UserRents, error) {
	if _, err := f.users.Find(userID); err != nil {
		return UserRents{}, errors.Wrapf(err, "Error getting user: %d", userID)
	}

	userRents, err := f.repository.Get(userID)
	if err != nil {
		return UserRents{}, err
	}
	if userRents == nil {
		return newUserRents(userID), nil
	}
	return *userRents, nil
}
