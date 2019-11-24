package rental_es

import (
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"os"
)

type Config struct {
	MongoDB   string
	MaxRented int
}

func BuildFacade(
	usersFacade users.Facade,
	moviesFacade movies.Facade,
	feesFacade fees.Facade,
	repository Repository,
	config Config,
) rental.RentalFacade {
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
