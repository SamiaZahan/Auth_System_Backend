package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type Auth struct{}

func (a *Auth) CreateUserIndex() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := DB.Collection(UserCollection)
	indexModel := mongo.IndexModel{
		Options: options.Index().SetUnique(true),
		Keys:    bsonx.MDoc{"email": bsonx.Int32(-1)},
	}

	_, err = col.Indexes().CreateOne(ctx, indexModel)

	if err != nil {
		return
	}

	return
}

func (a *Auth) GetUserByEmail(email string) (user *User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := DB.Collection(UserCollection)
	err = col.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)

	if err != nil {
		return
	}

	return
}

func (a *Auth) CreateUser(email string) (ID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := DB.Collection(UserCollection)
	res, err := col.InsertOne(ctx, User{
		Email:   email,
		Active:  false,
		Created: time.Now(),
	})

	if err != nil {
		return
	}

	ID = res.InsertedID.(primitive.ObjectID).String()
	return
}

func (a *Auth) CreateUserProfile(userID string, firstName string, lastName string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	UserID, _ := primitive.ObjectIDFromHex(userID)
	col := DB.Collection(UserProfileCollection)
	_, err = col.InsertOne(ctx, UserProfile{
		UserId:    UserID,
		FirstName: firstName,
		LastName:  lastName,
		Created:   time.Now(),
	})

	if err != nil {
		return
	}

	return
}
