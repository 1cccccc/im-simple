package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"im/config"
)

type Room struct {
	ID           string `bson:"_id,omitempty"`
	Number       int    `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	UserIdentity string `bson:"user_identity"`
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
}

func (r *Room) Collection() string {
	return "room"
}

func CreateRoom(room *Room) error {
	_, err := config.Mongo.Collection(room.Collection()).InsertOne(context.TODO(), room)
	return err
}

func DeleteRoom(roomIdentity string) error {
	_, err := config.Mongo.Collection((&Room{}).Collection()).DeleteOne(context.TODO(), bson.M{"_id": roomIdentity})
	return err
}
