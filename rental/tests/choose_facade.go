package domain_crud

import (
	"context"
	"database/sql"
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_common"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_crud"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_es"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"time"
)

type createFacadeFunc = func(users users.UsersFacade, movies movies.Facade) domain_common.RentalFacade // TODO

var facadeBuilder = buildEventSourcedTestFacade

func buildInMemoryCrudTestFacade(
	users users.UsersFacade,
	movies movies.Facade,
	fees fees.Facade,
	maximumRentedMovies int,
) domain_common.RentalFacade {
	config := domain_crud.TestConfig()
	config.MaxRentedMoviesCount = maximumRentedMovies

	return domain_crud.BuildFacade(users, movies, fees, domain_crud.NewInMemoryRepository(), config)
}

func buildPostgresCrudTestFacade(
	users users.UsersFacade,
	movies movies.Facade,
	fees fees.Facade,
	maximumRentedMovies int,
) domain_common.RentalFacade {
	config := domain_crud.TestConfig()
	config.MaxRentedMoviesCount = maximumRentedMovies
	clearPostgresDB(config)

	return domain_crud.BuildFacade(users, movies, fees, infrastructure.NewPostgresRepository(config.PostgresDSN), config)
}

func buildEventSourcedTestFacade(
	users users.UsersFacade,
	movies movies.Facade,
	fees fees.Facade,
	maximumRentedMovies int,
) domain_common.RentalFacade {
	clearMongoDB() // TODO make a proper config

	return domain_es.BuildFacade(users, movies, fees, infrastructure.NewMongoRepository(), maximumRentedMovies)
}

func clearMongoDB() {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017") // TODO ensure env
	if err != nil {
		panic(err.Error())
	}
	collection := client.Database("dvd-rental").Collection("rentals")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	if err != nil {
		panic(err.Error())
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err.Error())
	}

	err = collection.Drop(context.Background())
	if err != nil {
		panic(err.Error())
	}
}

func clearPostgresDB(config domain_crud.Config) {
	db, err := sql.Open("postgres", config.PostgresDSN)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE rented_movies")
	if err != nil {
		panic(err.Error())
	}
}
