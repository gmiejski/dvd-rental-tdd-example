package movies

import (
	"github.com/pkg/errors"
)

func Build() Facade {
	return &moviesFacade{repository: newInMemoryRepository()}
}

type moviesFacade struct {
	repository moviesRepository
}

func (f *moviesFacade) Add(createMovie CreateMovie) (CreatedMovieDTO, error) {
	movie, err := f.repository.Save(
		movie{Title: createMovie.Title, MinimalAge: createMovie.MinimalAge, Year: createMovie.Year, Genre: createMovie.Genre},
	)
	if err != nil {
		return CreatedMovieDTO{}, errors.WithMessage(err, "error adding movie")
	}
	return f.createdMovieDTO(movie.ID), nil
}

func (f *moviesFacade) Get(movieID MovieID) (MovieDTO, error) {
	movie, err := f.repository.Find(movieID)
	if err != nil {
		return MovieDTO{}, MovieNotFound{movieID: int(movieID)}
	}
	return f.movieDTO(movie), nil
}

func (f *moviesFacade) ListGenre(request GenreListingRequest) (ListingDTO, error) {
	page, err := f.repository.FindByGenre(request.Genre, request.CursorOffset, request.Limit)
	if err != nil {
		return ListingDTO{}, err
	}

	return ListingDTO{
		Movies:       f.toMoviesDTO(page.movies),
		PageInfo:     f.createPageInfo(page),
		TotalResults: page.totalResults,
	}, nil
}

func (f *moviesFacade) toMoviesDTO(movies []movie) []MovieDTO {
	result := make([]MovieDTO, 0)
	for _, movie := range movies {
		result = append(result, f.movieDTO(movie))
	}
	return result
}

func (f *moviesFacade) createdMovieDTO(movie MovieID) CreatedMovieDTO {
	return CreatedMovieDTO{int(movie)}
}

func (f *moviesFacade) movieDTO(movie movie) MovieDTO {
	return MovieDTO{ID: movie.ID, Title: movie.Title, Year: movie.Year, MinimalAge: movie.MinimalAge, Genre: movie.Genre}
}

func (f *moviesFacade) createPageInfo(page genreFindResult) PageInfo {
	return PageInfo{
		HasNextPage: page.lastOffset != page.totalResults && len(page.movies) > 0,
		LastCursor:  page.lastOffset,
	}
}
