package users

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrievingAddedUser(t *testing.T) {
	// given
	usersFacade := buildUsersFacade()
	createdUser, err := usersFacade.Add(CreateUser{Name: "Gabriel", Age: 18})
	require.NoError(t, err)

	// when
	user, err := usersFacade.Get(userID(createdUser.ID))

	// then
	require.NoError(t, err)
	assert.EqualValues(t, UserDTO{ID: 1, Name: "Gabriel", Age: 18}, user)
}

func TestErrorWhenUserNotFound(t *testing.T) {
	// given
	usersFacade := buildUsersFacade()

	// when
	_, err := usersFacade.Get(10)

	// then
	require.Error(t, err)
}
