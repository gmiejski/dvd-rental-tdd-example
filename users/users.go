package users

import (
	"github.com/pkg/errors"
)

func buildUsersFacade() UsersFacade {
	return &usersFacade{repository: newInMemoryRepository()}
}

type usersFacade struct {
	repository usersRepository
}

func (f *usersFacade) Add(createUser CreateUser) (CreatedUserDTO, error) {
	user, err := f.repository.Save(User{Name: createUser.Name, Age: createUser.Age})
	if err != nil {
		return CreatedUserDTO{}, errors.WithMessage(err, "error adding user")
	}
	return f.createdUserDTO(user.ID), nil
}

func (f *usersFacade) Get(userID UserID) (UserDTO, error) {
	user, err := f.repository.Find(userID)
	if err != nil {
		return UserDTO{}, errors.WithMessage(err, "error adding user")
	}
	return f.userDTO(user), nil
}

func (f *usersFacade) createdUserDTO(user UserID) CreatedUserDTO {
	return CreatedUserDTO{int(user)}
}
func (f *usersFacade) userDTO(user User) UserDTO {
	return UserDTO{ID: user.ID, Name: user.Name, Age: user.Age}
}
