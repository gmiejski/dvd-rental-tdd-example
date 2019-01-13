package domain_crud

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/pkg/errors"
	"time"
)

type UserRents struct {
	UserID       int
	RentedMovies []RentedMovie
}

func (r *UserRents) rentMovie(movie movies.MovieDTO) error { // TODO add movies anti corruption layer
	now := time.Now()
	if r.isMovieRented(int(movie.ID)) {
		return errors.Errorf("User %d already rented movie %d", r.UserID, movie.ID)
	}
	r.RentedMovies = append(r.RentedMovies, RentedMovie{MovieID: int(movie.ID), RentedAt: now, ReturnAt: now.Add(time.Hour * 24 * 3)})
	return nil
}

func (r *UserRents) returnBack(movieID int) error {
	if !r.isMovieRented(movieID) {
		return errors.Wrapf(domain_common.MovieIsNotRented{r.UserID, movieID}, "error returning movie")
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
	users      users.UsersFacade
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
			domain_common.MaximumMoviesRented{UserID: userID, Max: f.config.MaxRentedMoviesCount},
			"error renting movie %d by user %d", movieID, userID,
		)
	}

	if fees, _ := f.fees.GetFees(userID); len(fees.Fees) > 0 {
		return errors.Wrapf(
			domain_common.UnpaidFees{UserID: userID, Movies: fees.OverrentMovieIDs()},
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

func (f *facade) GetRented(userID int) (domain_common.RentedMoviesDTO, error) { // TODO rename to Rents
	if _, err := f.users.Get(userID); err != nil {
		return domain_common.RentedMoviesDTO{}, errors.Wrapf(err, "Error getting user: %d", userID)
	}
	rents, err := f.getUserRents(userID)
	if err != nil {
		return domain_common.RentedMoviesDTO{}, errors.WithMessagef(err, "Error getting rented movies for user %d", userID)
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

func toDTO(rents UserRents) domain_common.RentedMoviesDTO {
	rentedMovies := make([]domain_common.RentedMovieDTO, 0)
	for _, movie := range rents.RentedMovies {
		rentedMovies = append(rentedMovies, toMovieDTO(movie))
	}
	return domain_common.RentedMoviesDTO{Movies: rentedMovies}
}

func toMovieDTO(movie RentedMovie) domain_common.RentedMovieDTO {
	return domain_common.RentedMovieDTO{MovieID: movie.MovieID, RentedAt: movie.RentedAt, ReturnAt: movie.ReturnAt}
}

func (f *facade) getUserRents(userID int) (UserRents, error) {
	if _, err := f.users.Get(userID); err != nil {
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
