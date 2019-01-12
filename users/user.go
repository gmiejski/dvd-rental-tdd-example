package users

type User struct {
	ID   UserID
	Name string
	Age  int
}

type UserID int
