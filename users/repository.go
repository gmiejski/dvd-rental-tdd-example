package users

import (
	"github.com/pkg/errors"
	"sync"
)

type usersRepository interface {
	Save(user User) (User, error)
	Find(id UserID) (User, error)
}

func newInMemoryRepository() usersRepository {
	return &usersInMemoryRepository{data: make(map[UserID]User), lock: sync.Mutex{}, nextID: 1}
}

type usersInMemoryRepository struct {
	data   map[UserID]User
	lock   sync.Mutex
	nextID int
}

func (r *usersInMemoryRepository) Save(user User) (User, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	user.ID = UserID(r.nextID)
	r.data[user.ID] = user
	r.nextID += 1
	return user, nil
}

func (r *usersInMemoryRepository) Find(userID UserID) (User, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for id, user := range r.data {
		if id == userID {
			return user, nil
		}
	}
	return User{}, errors.Errorf("User not found %d", userID)
}
