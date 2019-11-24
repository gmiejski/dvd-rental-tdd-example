package users

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func buildTestFacade() Facade {
	return Build(NewInMemoryRepository())
}

func TestRetrievingAddedUser(t *testing.T) {
	// given
	usersFacade := buildTestFacade()
	createdUser, err := usersFacade.Create(CreateUser{Name: "Gabriel", Age: 18})
	require.NoError(t, err)

	// when
	user, err := usersFacade.Find(createdUser.ID)

	// then
	require.NoError(t, err)
	assert.EqualValues(t, UserDTO{ID: 1, Name: "Gabriel", Age: 18}, user)
}

func TestErrorWhenUserNotFound(t *testing.T) {
	// given
	usersFacade := buildTestFacade()

	// when
	_, err := usersFacade.Find(10)

	// then
	require.Error(t, err)
	require.IsType(t, UserNotFound{}, errors.Cause(err))
}
