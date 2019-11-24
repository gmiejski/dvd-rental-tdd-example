package users

import "fmt"

type CreateUser struct {
	Name string
	Age  int
}

type UserDTO struct {
	ID   int
	Name string
	Age  int
}

type CreatedUserDTO struct {
	ID int
}

type Facade interface {
	Create(user CreateUser) (CreatedUserDTO, error)
	Find(user int) (UserDTO, error)
}

type UserNotFound struct {
	userID int
}

func (err UserNotFound) Error() string {
	return fmt.Sprintf("User not found: %d", err.userID)
}
