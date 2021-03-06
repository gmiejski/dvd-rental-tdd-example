package users

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func buildTestFacade() Facade {
	return Build(NewInMemoryRepository())
}

func TestFindingUserAfterCreation(t *testing.T) {
	// given
	usersFacade := buildTestFacade()
	createdUser, err := usersFacade.Create(CreateUserCommand{Name: "Gabriel", Age: 18})
	require.NoError(t, err)

	// when
	user, err := usersFacade.Find(createdUser.ID)

	// then
	require.NoError(t, err)
	require.EqualValues(t, UserDTO{ID: createdUser.ID, Name: "Gabriel", Age: 18}, user)
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
