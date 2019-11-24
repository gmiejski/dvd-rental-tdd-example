package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gmiejski/dvd-rental-tdd-example/src/rental_es"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"time"
)

type mongoRepository struct {
	collection *mongo.Collection
	events     map[string]func(json string) rental_es.Event
}

type savedEvent struct {
	UserID    int    `json:"user",bson:"user"`
	EventName string `json:"event",bson:"event"`
	Data      string `json:"data",bson:"data"`
}

func (r *mongoRepository) Save(userID int, eventsToSave []rental_es.Event) error {
	for _, event := range eventsToSave {
		jsonValue, err := json.Marshal(event)

		if err != nil {
			return err
		}
		_, err = r.collection.InsertOne(
			context.Background(),
			bson.D{
				{"user", userID},
				{"event", eventName(event)},
				{"data", string(jsonValue)}})

		if err != nil {
			return err
		}
	}
	return nil
}

func eventName(event rental_es.Event) string {
	return reflect.TypeOf(event).Name()
}

func (r *mongoRepository) Get(user int) (*rental_es.UserRents, error) {
	filter := bson.M{"user": user}

	cur, err := r.collection.Find(context.Background(), filter)

	if err != nil {
		return nil, err
	}
	var events []rental_es.Event
	for cur.Next(context.Background()) {
		var saved savedEvent
		elem := &bson.D{}
		if err := cur.Decode(elem); err != nil {
			return nil, err
		}
		saved.UserID = int(elem.Map()["user"].(int64))
		saved.EventName = elem.Map()["event"].(string)
		saved.Data = elem.Map()["data"].(string)

		decodedEvent, err := r.decodeEvents(saved)
		if err != nil {
			return nil, err
		}
		events = append(events, decodedEvent)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	rentals := rental_es.NewUserRents(user)
	rentals.Apply(events)
	return &rentals, nil
}

func (r *mongoRepository) decodeEvents(event savedEvent) (rental_es.Event, error) {
	switch event.EventName {
	case "MovieRentedEvent":
		result := rental_es.MovieRentedEvent{}
		decoder := json.NewDecoder(bytes.NewBuffer([]byte(event.Data)))
		err := decoder.Decode(&result)
		if err != nil {
			return nil, errors.Errorf("error decoding event: %s", event.EventName)
		}
		return result, nil
	case "MovieReturnedEvent":
		result := rental_es.MovieReturnedEvent{}
		decoder := json.NewDecoder(bytes.NewBuffer([]byte(event.Data)))
		err := decoder.Decode(&result)
		if err != nil {
			return nil, errors.Errorf("error decoding event: %s", event.EventName)
		}
		return result, nil
	}
	return nil, errors.Errorf("Cannot find event decoder for name: %s", event.EventName)
}

func getEventsList() map[string]func(string) rental_es.Event {
	eventBuilders := map[string]func(string) rental_es.Event{
		eventName(rental_es.MovieRentedEvent{}): func(jsonString string) rental_es.Event {
			e := rental_es.MovieRentedEvent{}
			d := json.NewDecoder(bytes.NewBuffer([]byte(jsonString)))
			d.Decode(&e)
			return e
		},
		eventName(rental_es.MovieReturnedEvent{}): func(jsonString string) rental_es.Event {
			e := rental_es.MovieReturnedEvent{}
			d := json.NewDecoder(bytes.NewBuffer([]byte(jsonString)))
			d.Decode(&e)
			return e
		},
	}
	return eventBuilders
}

func NewMongoRepository(address string) rental_es.Repository {

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(address))

	if err != nil {
		panic(err.Error())
	}

	collection := client.Database("dvd-rental").Collection("rentals")

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err.Error())
	}

	events := getEventsList() // TODO events used anywhere?

	return &mongoRepository{collection: collection, events: events}
}
