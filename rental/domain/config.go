package domain

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"os"
)

type Config struct {
	MaxRentedMoviesCount int
	PostgresDSN          string
}

type testOptionFacade = func(*facade)

var withConfig = func(c Config) testOptionFacade {
	return func(f *facade) {
		f.config = c
	}
}

var withFeesFacade = func(feesFacade fees.Facade) testOptionFacade {
	return func(f *facade) {
		f.fees = feesFacade
	}
}

func buildTestFacade(users users.UsersFacade, movies movies.Facade, options ...testOptionFacade) RentalFacade {
	feesStub := fees.NewFacadeStub()
	baseTestFacade := &facade{
		users:      users,
		movies:     movies,
		fees:       &feesStub,
		repository: newInMemoryRepository(),
		config:     Config{MaxRentedMoviesCount: 10},
	}

	for _, option := range options {
		option(baseTestFacade)
	}

	return baseTestFacade
}

func TestConfig() Config {
	return Config{
		PostgresDSN:          "postgresql://postgres:postgres@localhost:5432/dvd-rental-tdd-example?sslmode=disable", // TODO use docker host
		MaxRentedMoviesCount: 10,
	}
}

func ProdConfig() Config {
	return Config{
		//PostgresDSN:          ensureEnv("POSTGRES_DSN"),// TODO use env dsn
		PostgresDSN:          "postgresql://postgres:postgres@localhost:5432/dvd-rental-tdd-example?sslmode=disable", // use env dsn
		MaxRentedMoviesCount: 10,
	}
}

func ensureEnv(name string) string { // TODO switch to envs
	value, exists := os.LookupEnv(name)
	if value != "" || !exists {
		panic(fmt.Sprintf("Env not found: %s", name))
	}
	return value
}

func BuildFacade(
	usersFacade users.UsersFacade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
	repository Repository,
	config Config,
) RentalFacade {
	return &facade{
		users:      usersFacade,
		movies:     moviesFacade,
		fees:       feesFacade,
		repository: repository,
		config:     config,
	}
}