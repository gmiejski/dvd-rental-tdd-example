package domain_es

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"os"
)

type Config struct {
	MongoDB   string
	MaxRented int
}

func BuildFacade(
	usersFacade users.UsersFacade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
	repository Repository,
	config Config,
) domain_common.RentalFacade {
	return &eventSourcedFacade{
		users:           usersFacade,
		movies:          moviesFacade,
		fees:            feesFacade,
		repository:      repository,
		maxRentedMovies: config.MaxRented,
	}
}

func NewConfig() Config {
	return Config{
		MongoDB:   ensureEnv("MONGODB"),
		MaxRented: 10,
	}
}

func ensureEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if value == "" || !exists {
		panic(fmt.Sprintf("Env not found: %s", name))
	}
	return value
}
