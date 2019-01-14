package domain_es

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/pkg/errors"
)

type eventSourcedFacade struct {
	users           users.UsersFacade
	movies          movies.Facade
	fees            fees.Facade
	repository      Repository
	maxRentedMovies int
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

	if userRents.rentedCount() >= f.maxRentedMovies {
		return errors.Wrapf(
			domain_common.MaximumMoviesRented{UserID: userID, Max: f.maxRentedMovies},
			"error renting movie %d by user %d", movieID, userID,
		)
	}

	dto, e := f.fees.GetFees(userID)
	if userFees, _ := dto, e; len(userFees.Fees) > 0 {
		return errors.Wrapf(
			domain_common.UnpaidFees{UserID: userID, Movies: userFees.OverrentMovieIDs()},
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

func (f *eventSourcedFacade) GetRented(userID int) (domain_common.RentedMoviesDTO, error) { // TODO rename to Rents
	if _, err := f.users.Get(userID); err != nil {
		return domain_common.RentedMoviesDTO{}, errors.Wrapf(err, "Error getting user: %d", userID)
	}
	rents, err := f.getUserRents(userID)
	if err != nil {
		return domain_common.RentedMoviesDTO{}, errors.WithMessagef(err, "Error getting rented movies for user %d", userID)
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

func (f *eventSourcedFacade) getUserRents(userID int) (UserRents, error) {
	if _, err := f.users.Get(userID); err != nil {
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
