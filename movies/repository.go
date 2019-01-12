package movies

import (
	"github.com/pkg/errors"
	"sort"
	"sync"
)

type genreFindResult struct {
	movies       []movie
	lastOffset   int
	totalResults int
}

type moviesRepository interface {
	Save(movie movie) (movie, error)
	Find(id movieID) (movie, error)
	FindByGenre(genre string, after int, count int) (genreFindResult, error)
}

func newInMemoryRepository() moviesRepository {
	return &moviesInMemoryRepository{data: make(map[movieID]movie), lock: sync.Mutex{}, nextID: 1}
}

type moviesInMemoryRepository struct {
	data   map[movieID]movie
	lock   sync.Mutex
	nextID int
}

func (r *moviesInMemoryRepository) Save(movie movie) (movie, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	movie.ID = movieID(r.nextID)
	r.data[movie.ID] = movie
	r.nextID += 1
	return movie, nil
}

func (r *moviesInMemoryRepository) Find(movieID movieID) (movie, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for id, movie := range r.data {
		if id == movieID {
			return movie, nil
		}
	}
	return movie{}, errors.Errorf("movie not found %d", movieID)
}

func (r *moviesInMemoryRepository) FindByGenre(genre string, after int, count int) (genreFindResult, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	genreMovies := r.allWithGenre(genre)
	result := make([]movie, 0)
	lastOffset := -1
	for _, genreMovie := range genreMovies {
		if int(genreMovie.ID) > after && len(result) < count {
			result = append(result, genreMovie)
			lastOffset = int(genreMovie.ID)
		}
	}

	return genreFindResult{
		movies:       result,
		totalResults: len(genreMovies),
		lastOffset:   lastOffset,
	}, nil
}

func (r *moviesInMemoryRepository) allWithGenre(genre string) []movie {
	genreMovies := make([]movie, 0)
	for _, movie := range r.data {
		if movie.Genre == genre {
			genreMovies = append(genreMovies, movie)
		}
	}
	sort.Sort(movieIDSorter(genreMovies))
	return genreMovies
}

// NameSorter sorts planets by name.
type movieIDSorter []movie

func (a movieIDSorter) Len() int           { return len(a) }
func (a movieIDSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a movieIDSorter) Less(i, j int) bool { return a[i].ID < a[j].ID }
