package rental

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gmiejski/dvd-rental-tdd-example/src/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/src/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud"
	rental_crud_infra "github.com/gmiejski/dvd-rental-tdd-example/src/rental_crud/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_es"
	rental_es_infra "github.com/gmiejski/dvd-rental-tdd-example/src/rental_es/infrastructure"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

var currentFacadeBuilder = buildInMemoryCrudTestFacade

func buildInMemoryCrudTestFacade(users users.Facade, movies movies.Facade, fees fees.Facade, options ...rental.RentalOption) rental.Facade {
	basicConfig := buildRentalConfig(options)
	config := rental_crud.TestConfig(basicConfig)
	return rental_crud.Build(users, movies, fees, rental_crud.NewInMemoryRepository(), config)
}

func buildPostgresCrudTestFacade(users users.Facade, movies movies.Facade, fees fees.Facade, options ...rental.RentalOption) rental.Facade {
	rentalConfig := buildRentalConfig(options)
	config := rental_crud.IntegrationTestConfig(rentalConfig)
	clearPostgresDB(config)
	return rental_crud.Build(users, movies, fees, rental_crud_infra.NewPostgresRepository(config.PostgresDSN), config)
}

func buildEventSourcedTestFacade(users users.Facade, movies movies.Facade, fees fees.Facade, options ...rental.RentalOption) rental.Facade {
	clearMongoDB()
	rentalConfig := buildRentalConfig(options)
	config := rental_es.NewConfig(rentalConfig)

	return rental_es.Build(users, movies, fees, rental_es_infra.NewMongoRepository(config.MongoDB), config)
}

func buildRentalConfig(rentalOptions []rental.RentalOption) rental.Config {
	config := rental.StandardConfig()
	for _, x := range rentalOptions {
		x(&config)
	}
	return config
}

func clearMongoDB() {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(ensureEnv("MONGODB")))

	if err != nil {
		panic(err.Error())
	}
	collection := client.Database("dvd-rental").Collection("rentals")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err.Error())
	}

	err = collection.Drop(context.Background())
	if err != nil {
		panic(err.Error())
	}
}

func clearPostgresDB(config rental_crud.Config) {
	db, err := sql.Open("postgres", config.PostgresDSN)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec("TRUNCATE TABLE rented_movies")
	if err != nil {
		panic(err.Error())
	}
}

func ensureEnv(name string) string {
	value, exists := os.LookupEnv(name)
	if value == "" || !exists {
		panic(fmt.Sprintf("Env not found: %s", name))
	}
	return value
}
