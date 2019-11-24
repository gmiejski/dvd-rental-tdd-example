package users

import (
	"database/sql"
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users/infrastructure"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func buildIntegrationTestFacade() users.Facade {
	return users.Build(infrastructure.NewPostgresRepository(ensureEnv("POSTGRES_DSN")))
}

func TestRetrievingAddedUser(t *testing.T) {
	// given
	clearDB()
	usersFacade := buildIntegrationTestFacade()
	createdUser, err := usersFacade.Create(users.CreateUser{Name: "Gabriel", Age: 18})
	require.NoError(t, err)

	// when
	user, err := usersFacade.Find(createdUser.ID)

	// then
	require.NoError(t, err)
	assert.EqualValues(t, users.UserDTO{ID: createdUser.ID, Name: "Gabriel", Age: 18}, user)
}

func TestErrorWhenUserNotFound(t *testing.T) {
	// given
	clearDB()
	usersFacade := buildIntegrationTestFacade()

	// when
	_, err := usersFacade.Find(10)

	// then
	require.Error(t, err)
	require.IsType(t, users.UserNotFound{}, errors.Cause(err))
}

func ensureEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if value == "" || !exists {
		panic(fmt.Sprintf("Env not found: %s", name))
	}
	return value
}

func clearDB() {
	db, err := sql.Open("postgres", ensureEnv("POSTGRES_DSN"))
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE users")
	if err != nil {
		panic(err.Error())
	}
}
