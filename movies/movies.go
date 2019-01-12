package movies

import (
	"github.com/pkg/errors"
)

func buildMoviesFacade() MoviesFacade {
	return &moviesFacade{repository: newInMemoryRepository()}
}

type moviesFacade struct {
	repository moviesRepository
}

func (f *moviesFacade) ListGenre() {
	panic("implement me")
}

func (f *moviesFacade) Add(createMovie CreateMovie) (CreatedMovieDTO, error) {
	movie, err := f.repository.Save(
		Movie{Title: createMovie.Title, MinimalAge: createMovie.MinimalAge, Year: createMovie.Year},
	)
	if err != nil {
		return CreatedMovieDTO{}, errors.WithMessage(err, "error adding movie")
	}
	return f.createdMovieDTO(movie.ID), nil
}

func (f *moviesFacade) Get(MovieID MovieID) (MovieDTO, error) {
	movie, err := f.repository.Find(MovieID)
	if err != nil {
		return MovieDTO{}, errors.WithMessage(err, "error adding movie")
	}
	return f.movieDTO(movie), nil
}

func (f *moviesFacade) createdMovieDTO(Movie MovieID) CreatedMovieDTO {
	return CreatedMovieDTO{int(Movie)}
}
func (f *moviesFacade) movieDTO(movie Movie) MovieDTO {
	return MovieDTO{ID: movie.ID, Title: movie.Title, Year: movie.Year, MinimalAge: movie.MinimalAge}
}
