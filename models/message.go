package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im/config"
)

type Message struct {
	ID           string `bson:"_id,omitempty"`
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
}

func (m *Message) Collection() string {
	return "message"
}

func InsertOneMessage(msg *Message) error {
	_, err := config.Mongo.Collection(msg.Collection()).InsertOne(context.Background(), msg)
	return err
}

func FindMessageListByRoomIdentity(roomIdentity string, size, page int64) (data []*Message, err error) {
	skip := size * (page - 1)
	limit := size

	cur, err := config.Mongo.Collection((&Message{}).Collection()).Find(context.Background(), bson.M{"room_identity": roomIdentity}, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  bson.M{"create_at": -1},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		m := new(Message)
		err := cur.Decode(m)
		if err != nil {
			return nil, err
		}
		data = append(data, m)
	}

	return data, nil
}
