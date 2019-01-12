package users

type CreateUser struct {
	Name string
	Age  int
}

type UserDTO struct {
	ID   userID
	Name string
	Age  int
}

type CreatedUserDTO struct {
	ID int
}

type UsersFacade interface {
	Add(user CreateUser) (CreatedUserDTO, error)
	Get(user userID) (UserDTO, error)
}
