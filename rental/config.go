package rental

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"os"
)

type config struct {
	maxRentedMoviesCount int
	postgresDSN          string
}

type testOptionFacade = func(*facade)

var withConfig = func(c config) testOptionFacade {
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
		config:     config{maxRentedMoviesCount: 10},
	}

	for _, option := range options {
		option(baseTestFacade)
	}

	return baseTestFacade
}

func TestConfig() config {
	return config{
		postgresDSN:          "postgresql://postgres:postgres@localhost:5432/dvd-rental-tdd-example?sslmode=disable", // TODO usedocker host
		maxRentedMoviesCount: 10,
	}
}

func ProdConfig() config {
	return config{
		//postgresDSN:          ensureEnv("POSTGRES_DSN"),// TODO use env dsn
		postgresDSN:          "postgresql://postgres:postgres@localhost:5432/dvd-rental-tdd-example?sslmode=disable", // use env dsn
		maxRentedMoviesCount: 10,
	}
}

func BuildFacade(
	usersFacade users.UsersFacade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
	config config,
) RentalFacade {
	return &facade{
		users:      usersFacade,
		movies:     moviesFacade,
		fees:       feesFacade,
		repository: NewPostgresRepository(config.postgresDSN),
		config:     config,
	}
}

func BuildIntegrationTestFacade(usersFacade users.UsersFacade, moviesFacade movies.Facade) RentalFacade {
	feesStub := fees.NewFacadeStub()
	return BuildFacade(usersFacade, moviesFacade, &feesStub, TestConfig())
}

func ensureEnv(name string) string { // TODO switch to envs
	value, exists := os.LookupEnv(name)
	if value != "" || !exists {
		panic(fmt.Sprintf("Env not found: %s", name))
	}
	return value
}
