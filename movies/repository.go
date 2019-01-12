package movies

import (
	"github.com/pkg/errors"
	"sync"
)

type moviesRepository interface {
	Save(Movie Movie) (Movie, error)
	Find(id MovieID) (Movie, error)
	//FindByGenre(id MovieID) (MovieListingDTO, error)
}

func newInMemoryRepository() moviesRepository {
	return &moviesInMemoryRepository{data: make(map[MovieID]Movie), lock: sync.Mutex{}, nextID: 1}
}

type moviesInMemoryRepository struct {
	data   map[MovieID]Movie
	lock   sync.Mutex
	nextID int
}

func (r *moviesInMemoryRepository) Save(Movie Movie) (Movie, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	Movie.ID = MovieID(r.nextID)
	r.data[Movie.ID] = Movie
	r.nextID += 1
	return Movie, nil
}

func (r *moviesInMemoryRepository) Find(MovieID MovieID) (Movie, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for id, Movie := range r.data {
		if id == MovieID {
			return Movie, nil
		}
	}
	return Movie{}, errors.Errorf("Movie not found %d", MovieID)
}
