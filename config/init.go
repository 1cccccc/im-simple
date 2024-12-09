package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	username = "mongo"
	password = "mongo"
	host     = "localhost"
	port     = 27017
	database = "im"
	timeout  = 5000

	Mongo = InitMongoDB()
	Redis = InitRedis()
)

func InitMongoDB() *mongo.Database {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?timeoutMS=%d", username, password, host, port, timeout)
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	database := client.Database(database)

	return database
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
