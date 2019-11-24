package users

type facadeStub struct {
	users map[int]UserDTO
}

func (f *facadeStub) Create(user CreateUser) (CreatedUserDTO, error) {
	panic("implement me")
}

func (f *facadeStub) Find(userID int) (UserDTO, error) {
	for _, user := range f.users {
		if int(user.ID) == userID {
			return user, nil
		}
	}
	return UserDTO{}, UserNotFound{userID: int(userID)}
}

func NewFacadeStub(stubbedUsers []UserDTO) Facade {
	usersById := make(map[int]UserDTO)
	for _, user := range stubbedUsers {
		usersById[int(user.ID)] = user
	}
	return &facadeStub{users: usersById}
}
