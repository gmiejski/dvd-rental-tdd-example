package users

import (
	"sync"
)

type usersRepository interface {
	Save(user user) (user, error)
	Find(id int) (*user, error)
}

func newInMemoryRepository() usersRepository {
	return &usersInMemoryRepository{data: make(map[int]user), lock: sync.Mutex{}, nextID: 1}
}

type usersInMemoryRepository struct {
	data   map[int]user
	lock   sync.Mutex
	nextID int
}

func (r *usersInMemoryRepository) Save(user user) (user, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	user.ID = r.nextID
	r.data[user.ID] = user
	r.nextID += 1
	return user, nil
}

func (r *usersInMemoryRepository) Find(userID int) (*user, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for id, user := range r.data {
		if id == userID {
			return &user, nil
		}
	}
	return nil, nil
}
