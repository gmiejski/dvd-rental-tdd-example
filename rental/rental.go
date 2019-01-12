package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/pkg/errors"
	"time"
)

type UserRents struct {
	userID       int
	age          int
	rentedMovies []rentedMovie
}

func (r *UserRents) rentMovie(movie movies.MovieDTO) error {
	now := time.Now()
	r.rentedMovies = append(r.rentedMovies, rentedMovie{movieID: int(movie.ID), rentedAt: now, returnAt: now.Add(time.Hour * 24 * 3)})
	return nil
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
}

func (f *facade) Rent(userID int, movieID int) error {
	if _, err := f.users.Get(userID); err != nil {
		return errors.Wrapf(err, "Error getting user: %d", userID)
	}
	movie, err := f.movies.Get(movies.MovieID(movieID))
	if err != nil {
		return errors.Wrapf(err, "Error finding movie: %d", movieID)
	}

	userRents, err := f.getUserRents(userID)
	if err != nil {
		return errors.Wrapf(err, "Error getting rented movies for user: %d", userID)
	}

	err = userRents.rentMovie(movie)
	if err != nil {
		return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
	}
	err = f.repository.Save(userRents)
	return errors.WithMessagef(err, "error renting movie %d by user %d", movieID, userID)
}

func newUserRents(userID int) UserRents {
	return UserRents{userID: userID}
}

func (f *facade) GetRented(userID int) (RentedMoviesDTO, error) {
	rents, err := f.repository.Get(userID)
	if err != nil {
		return RentedMoviesDTO{}, errors.WithMessagef(err, "Error getting rented movies for user %d", userID)
	}
	return toDTO(rents), nil
}

func toDTO(rents *UserRents) RentedMoviesDTO {
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

func buildTestFacade(users users.UsersFacade, movies movies.Facade) RentalFacade {
	return &facade{users: users, movies: movies, repository: newInMemoryRepository()}
}
