package users

type user struct {
	ID   userID
	Name string
	Age  int
}

type userID int
