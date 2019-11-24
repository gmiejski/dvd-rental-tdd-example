package users

import "fmt"

type Facade interface {
	Create(user CreateUserCommand) (CreatedUserDTO, error)
	Find(user int) (UserDTO, error)
}

type CreateUserCommand struct {
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

type UserNotFound struct {
	userID int
}

func (err UserNotFound) Error() string {
	return fmt.Sprintf("User not found: %d", err.userID)
}
