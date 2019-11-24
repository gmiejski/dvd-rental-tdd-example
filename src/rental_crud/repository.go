package rental_crud

import (
	"sync"
)

type Repository interface {
	Save(rents UserRents) error
	Get(user int) (*UserRents, error)
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{data: make(map[int]UserRents), lock: sync.Mutex{}}
}

type inMemoryRepository struct {
	data map[int]UserRents
	lock sync.Mutex
}

func (r *inMemoryRepository) Save(rentsToSave UserRents) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[rentsToSave.UserID] = rentsToSave
	return nil
}

func (r *inMemoryRepository) Get(user int) (*UserRents, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, rents := range r.data {
		if rents.UserID == user {
			return &rents, nil
		}
	}
	return nil, nil
}
