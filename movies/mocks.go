package movies

type facadeStub struct {
	movies map[int]MovieDTO
}

func (f *facadeStub) Get(movieID MovieID) (MovieDTO, error) {
	for _, movie := range f.movies {
		if movie.ID == movieID {
			return movie, nil
		}
	}
	return MovieDTO{}, MovieNotFound{movieID: int(movieID)}
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
