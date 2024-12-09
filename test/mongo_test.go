package test

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im/models"
	"testing"
)

func TestFindOne(t *testing.T) {
	uri := "mongodb://mongo:mongo@localhost:27017/?timeoutMS=5000"
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer client.Disconnect(context.TODO())

	database := client.Database("im")

	user := new(models.User)

	err = database.Collection("user").FindOne(context.TODO(), bson.D{}).Decode(&user)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(user)
}
