package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"im/config"
)

type UserRoom struct {
	ID           string `bson:"_id,omitempty"`
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType     int    `bson:"room_type"`
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
}

func (ur *UserRoom) Collection() string {
	return "user_room"
}

func GetUserRoomByUserIdAndRoomId(userId, roomId string) (*UserRoom, error) {
	ur := UserRoom{}
	err := config.Mongo.Collection(ur.Collection()).FindOne(context.TODO(), bson.D{{Key: "user_identity", Value: userId}, {Key: "room_identity", Value: roomId}}).Decode(&ur)
	if err != nil {
		return nil, err
	}

	return &ur, nil
}

func GetUserRoomsByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	urs := make([]*UserRoom, 0)
	cursor, err := config.Mongo.Collection((&UserRoom{}).Collection()).Find(context.TODO(), bson.D{{Key: "room_identity", Value: roomIdentity}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		ur := UserRoom{}
		err := cursor.Decode(&ur)
		if err != nil {
			return nil, err
		}
		urs = append(urs, &ur)
	}

	return urs, nil
}

func GetUserRoomsByUserIdentityRoomType(userIdentity string, rooType int) ([]*UserRoom, error) {
	urs := make([]*UserRoom, 0)
	cursor, err := config.Mongo.Collection((&UserRoom{}).Collection()).Find(context.TODO(), bson.D{{Key: "user_identity", Value: userIdentity}, {Key: "room_type", Value: rooType}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		ur := UserRoom{}
		err := cursor.Decode(&ur)
		if err != nil {
			return nil, err
		}
		urs = append(urs, &ur)
	}

	return urs, nil
}

func GetUserRoomCountByUserIdentityRoomIdentity(userIdentity string, roomIdentitys []string, rooType int) (int64, error) {
	return config.Mongo.Collection((&UserRoom{}).Collection()).CountDocuments(context.TODO(), bson.M{"user_identity": userIdentity, "room_identity": bson.M{"$in": roomIdentitys}, "room_type": rooType})
}

func GetRoomIdentityCountByUserIdentityRoomIdentity(userIdentity string, roomIdentitys []string, rooType int) (string, error) {
	ur := &UserRoom{}
	err := config.Mongo.Collection(ur.Collection()).FindOne(context.TODO(), bson.M{"user_identity": userIdentity, "room_identity": bson.M{"$in": roomIdentitys}, "room_type": rooType}).Decode(ur)
	return ur.RoomIdentity, err
}

func CreateUserRoom(ur *UserRoom) error {
	_, err := config.Mongo.Collection(ur.Collection()).InsertOne(context.TODO(), ur)
	return err
}

func DeleteUserRoomByUserIdentityAndRoomIdentity(userIdentity, roomIdentity string) error {
	_, err := config.Mongo.Collection((&UserRoom{}).Collection()).DeleteOne(context.TODO(), bson.M{
		"user_identity": userIdentity,
		"room_identity": roomIdentity,
	})

	return err
}
