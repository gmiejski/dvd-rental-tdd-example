package movies

import "github.com/pkg/errors"

type facadeStub struct {
	movies map[int]MovieDTO
}

func (f *facadeStub) Get(movieID MovieID) (MovieDTO, error) {
	for _, movie := range f.movies {
		if movie.ID == movieID {
			return movie, nil
		}
	}
	return MovieDTO{}, errors.New("movie not found") // TODO make specific error
}

func (f *facadeStub) ListGenre(request GenreListingRequest) (ListingDTO, error) {
	panic("implement me")
}

func (f *facadeStub) Add(user CreateMovie) (CreatedMovieDTO, error) {
	panic("implement me")
}

func NewFacadeStub(stubbedUsers []MovieDTO) Facade {
	usersById := make(map[int]MovieDTO)
	for _, user := range stubbedUsers {
		usersById[int(user.ID)] = user
	}

	return &facadeStub{movies: usersById}
}
