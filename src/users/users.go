package users

import (
	"github.com/pkg/errors"
)

func Build() Facade {
	return &usersFacade{repository: newInMemoryRepository()}
}

type usersFacade struct {
	repository usersRepository
}

func (f *usersFacade) Create(createUser CreateUser) (CreatedUserDTO, error) {
	user, err := f.repository.Save(user{Name: createUser.Name, Age: createUser.Age})
	if err != nil {
		return CreatedUserDTO{}, errors.WithMessage(err, "error adding user")
	}
	return f.createdUserDTO(user.ID), nil
}

func (f *usersFacade) Find(userID int) (UserDTO, error) {
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

type user struct {
	ID   int
	Name string
	Age  int
}
