package domain_crud

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
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

func BuildUnitTestFacade(users users.UsersFacade, movies movies.Facade, options ...testOptionFacade) domain_common.RentalFacade {
	feesStub := fees.NewFacadeStub()
	baseTestFacade := &facade{
		users:      users,
		movies:     movies,
		fees:       &feesStub,
		repository: NewInMemoryRepository(),
		config:     Config{MaxRentedMoviesCount: 10},
	}

	for _, option := range options {
		option(baseTestFacade)
	}

	return baseTestFacade
}

func IntegrationTestConfig() Config {
	return Config{
		PostgresDSN:          ensureEnv("POSTGRES_DSN"),
		MaxRentedMoviesCount: 10,
	}
}

func TestConfig() Config {
	return Config{
		MaxRentedMoviesCount: 10,
	}
}

func ProdConfig() Config {
	return Config{
		PostgresDSN:          ensureEnv("POSTGRES_DSN"),
		MaxRentedMoviesCount: 10,
	}
}

func ensureEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if value == "" || !exists {
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
) domain_common.RentalFacade {
	return &facade{
		users:      usersFacade,
		movies:     moviesFacade,
		fees:       feesFacade,
		repository: repository,
		config:     config,
	}
}
