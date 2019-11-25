package rental_crud

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"os"
)

func Build(
	usersFacade users.Facade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
	repository Repository,
	config Config,
) rental.Facade {
	return &facade{
		users:      usersFacade,
		movies:     moviesFacade,
		fees:       feesFacade,
		repository: repository,
		config:     config,
	}
}

type Config struct {
	rental.Config
	PostgresDSN string
}

func IntegrationTestConfig(config rental.Config) Config {
	return Config{
		PostgresDSN: ensureEnv("POSTGRES_DSN"),
		Config:      config,
	}
}

func TestConfig(config rental.Config) Config {
	return Config{
		Config: config,
	}
}

func ProdConfig(config rental.Config) Config {
	return Config{
		PostgresDSN: ensureEnv("POSTGRES_DSN"),
		Config:      config,
	}
}

func ensureEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if value == "" || !exists {
		panic(fmt.Sprintf("Env not found: %s", name))
	}
	return value
}
