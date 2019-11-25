package rental_es

import (
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"github.com/pkg/errors"
)

type eventSourcedFacade struct {
	users      users.Facade
	movies     movies.Facade
	fees       fees.Facade
	repository Repository
	config     Config
}

func (f *eventSourcedFacade) Rent(userID int, movieID int) error {
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
		return errors.Wrapf(err, "Error getting fees for user: %d", userID)
	}
	if len(userFees.Fees) > 0 {
		return errors.Wrapf(
			rental.UnpaidFees{UserID: userID, Movies: userFees.OverrentMovieIDs()},
			"error renting movie %d",
			movieID,
		)
	}

	events, err := userRents.rentMovie(movie)
	if err != nil {
		return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
	}
	err = f.repository.Save(userID, events)
	return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
}

func (f *eventSourcedFacade) GetRented(userID int) (rental.RentedMoviesDTO, error) { // TODO rename to Rents
	if _, err := f.users.Find(userID); err != nil {
		return rental.RentedMoviesDTO{}, errors.Wrapf(err, "Error getting user: %d", userID)
	}
	rents, err := f.getUserRents(userID)
	if err != nil {
		return rental.RentedMoviesDTO{}, errors.WithMessagef(err, "Error getting rented movies for user %d", userID)
	}
	return toDTO(rents), nil
}

func (f *eventSourcedFacade) Return(userID int, movieID int) error {
	userRents, err := f.getUserRents(userID)
	if err != nil {
		return errors.Wrapf(err, "Error getting rented movies for user: %d", userID)
	}

	events, err := userRents.returnBack(movieID)
	if err != nil {
		return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
	}
	err = f.repository.Save(userID, events)
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

func (f *eventSourcedFacade) getUserRents(userID int) (UserRents, error) {
	if _, err := f.users.Find(userID); err != nil {
		return UserRents{}, errors.Wrapf(err, "Error getting user: %d", userID)
	}

	userRents, err := f.repository.Get(userID)
	if err != nil {
		return UserRents{}, err
	}
	if userRents == nil {
		return NewUserRents(userID), nil
	}
	return *userRents, nil
}
