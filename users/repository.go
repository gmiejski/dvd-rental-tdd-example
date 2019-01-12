package users

import (
	"github.com/pkg/errors"
	"sync"
)

type usersRepository interface {
	Save(user user) (user, error)
	Find(id userID) (user, error)
}

func newInMemoryRepository() usersRepository {
	return &usersInMemoryRepository{data: make(map[userID]user), lock: sync.Mutex{}, nextID: 1}
}

type usersInMemoryRepository struct {
	data   map[userID]user
	lock   sync.Mutex
	nextID int
}

func (r *usersInMemoryRepository) Save(user user) (user, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	user.ID = userID(r.nextID)
	r.data[user.ID] = user
	r.nextID += 1
	return user, nil
}

func (r *usersInMemoryRepository) Find(userID userID) (user, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for id, user := range r.data {
		if id == userID {
			return user, nil
		}
	}
	return user{}, errors.Errorf("User not found %d", userID)
}
