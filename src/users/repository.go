package users

import (
	"sync"
)

type Repository interface {
	Save(user User) (User, error)
	Find(id int) (*User, error)
}

func NewInMemoryRepository() Repository {
	return &usersInMemoryRepository{data: make(map[int]User), lock: sync.Mutex{}, nextID: 1}
}

type usersInMemoryRepository struct {
	data   map[int]User
	lock   sync.Mutex
	nextID int
}

func (r *usersInMemoryRepository) Save(user User) (User, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	user.ID = r.nextID
	r.data[user.ID] = user
	r.nextID += 1
	return user, nil
}

func (r *usersInMemoryRepository) Find(userID int) (*User, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for id, user := range r.data {
		if id == userID {
			return &user, nil
		}
	}
	return nil, nil
}
