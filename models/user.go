package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"im/config"
)

type User struct {
	Identity string `bson:"_id,omitempty"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Nickname string `bson:"nickname"`
	Sex      uint8  `bson:"sex"`
	Email    string `bson:"email"`
	Avatar   string `bson:"avatar"`
	CreateAt int64  `bson:"create_at"`
	UpdateAt int64  `bson:"update_at"`
}

func (u *User) Collection() string {
	return "user"
}

func GetUserByUsernamePassword(username, password string) (*User, error) {
	user := User{}
	err := config.Mongo.Collection(user.Collection()).FindOne(context.TODO(), bson.D{{Key: "username", Value: username}, {Key: "password", Value: password}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByIdentity(identity string) (*User, error) {
	user := User{}

	objId, err2 := primitive.ObjectIDFromHex(identity)
	if err2 != nil {
		return nil, err2
	}

	err := config.Mongo.Collection(user.Collection()).FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserCountByUsername(username string) (int64, error) {
	return config.Mongo.Collection((&User{}).Collection()).CountDocuments(context.TODO(), bson.D{{Key: "username", Value: username}})
}

func CreateOneUser(user *User) error {
	_, err := config.Mongo.Collection(user.Collection()).InsertOne(context.TODO(), user)
	return err
}
