package users

import (
	"github.com/pkg/errors"
)

func BuildUsersFacade() UsersFacade {
	return &usersFacade{repository: newInMemoryRepository()}
}

type usersFacade struct {
	repository usersRepository
}

func (f *usersFacade) Add(createUser CreateUser) (CreatedUserDTO, error) {
	user, err := f.repository.Save(user{Name: createUser.Name, Age: createUser.Age})
	if err != nil {
		return CreatedUserDTO{}, errors.WithMessage(err, "error adding user")
	}
	return f.createdUserDTO(user.ID), nil
}

func (f *usersFacade) Get(userID int) (UserDTO, error) {
	user, err := f.repository.Find(userID)
	if err != nil {
		return UserDTO{}, errors.WithMessage(err, "error adding user")
	}
	if user == nil {
		return UserDTO{}, UserNotFound{userID: int(userID)}
	}
	return f.userDTO(user), nil
}

func (f *usersFacade) createdUserDTO(user int) CreatedUserDTO {
	return CreatedUserDTO{int(user)}
}
func (f *usersFacade) userDTO(user *user) UserDTO {
	return UserDTO{ID: user.ID, Name: user.Name, Age: user.Age}
}
