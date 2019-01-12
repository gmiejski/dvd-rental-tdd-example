package rental

import (
	"sync"
)

type repository interface {
	Save(rents UserRents) error
	Get(user int) (*UserRents, error)
}

func newInMemoryRepository() repository {
	return &inMemoryRepository{data: make(map[int]UserRents), lock: sync.Mutex{}}
}

type inMemoryRepository struct {
	data map[int]UserRents
	lock sync.Mutex
}

func (r *inMemoryRepository) Save(rentsToSave UserRents) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[rentsToSave.userID] = rentsToSave
	return nil
}

func (r *inMemoryRepository) Get(user int) (*UserRents, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, rents := range r.data {
		if rents.userID == user {
			return &rents, nil
		}
	}
	return nil, nil
}
