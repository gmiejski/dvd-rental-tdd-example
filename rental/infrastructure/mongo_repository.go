package infrastructure

import (
	"context"
	"github.com/gmiejski/dvd-rental-tdd-example/rental/domain_crud"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"time"
)

type mongoRepository struct {
	collection *mongo.Collection
}

func (r *mongoRepository) Save(rents domain_crud.UserRents) error {
	//res, err := r.collection.InsertOne(context.Background(), bson.D{{"user_id",10}, {"movie_id",10}})
	return nil
}

func (r *mongoRepository) Get(user int) (*domain_crud.UserRents, error) {
	panic("implement me")
}

func NewMongoRepository() domain_crud.Repository {

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")
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
	return &mongoRepository{collection: collection}
}
