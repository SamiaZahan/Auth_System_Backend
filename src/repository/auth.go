package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Auth struct {
	Ctx context.Context
}

func (a *Auth) CreateUserIndex() (err error) {
	col := DB.Collection(UserCollection)
	_, err = col.Indexes().CreateOne(a.Ctx, mongo.IndexModel{
		Options: options.Index().SetUnique(true),
		Keys:    bsonx.MDoc{"email": bsonx.Int32(-1)},
	})
	return
}

func (a *Auth) GetUserByID(ID string) (user *UserDoc, err error) {
	UserID, _ := primitive.ObjectIDFromHex(ID)
	col := DB.Collection(UserCollection)
	err = col.FindOne(a.Ctx, bson.M{"_id": bson.M{"$eq": UserID}}).Decode(&user)
	return
}

func (a *Auth) GetUserProfileByID(ID string) (user *UserProfileDoc, err error) {
	UserID, _ := primitive.ObjectIDFromHex(ID)
	col := DB.Collection(UserProfileCollection)
	err = col.FindOne(a.Ctx, bson.M{"user_id": bson.M{"$eq": UserID}}).Decode(&user)
	return
}

func (a *Auth) GetUserByEmail(email string) (user *UserDoc, err error) {
	col := DB.Collection(UserCollection)
	err = col.FindOne(a.Ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)
	return
}

func (a *Auth) CreateUser(email string) (ID string, err error) {
	col := DB.Collection(UserCollection)
	res, err := col.InsertOne(a.Ctx, UserDoc{
		ID:      primitive.NewObjectID(),
		Email:   email,
		Active:  false,
		Created: time.Now(),
	})

	if err != nil {
		return
	}

	ID = res.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (a *Auth) CreateUserProfile(userID string, firstName string, lastName string) (err error) {
	UserID, _ := primitive.ObjectIDFromHex(userID)
	col := DB.Collection(UserProfileCollection)
	_, err = col.InsertOne(a.Ctx, UserProfileDoc{
		ID:        primitive.NewObjectID(),
		UserID:    UserID,
		FirstName: firstName,
		LastName:  lastName,
		Created:   time.Now(),
	})
	return
}

func (a *Auth) ActivateUserByID(userID string) (err error) {
	UserID, _ := primitive.ObjectIDFromHex(userID)
	col := DB.Collection(UserCollection)
	_, err = col.UpdateOne(
		a.Ctx,
		bson.M{"_id": bson.M{"$eq": UserID}},
		bson.M{"$set": bson.M{"active": true}},
	)
	return
}

func (a *Auth) SetUserMobileByID(userID string, mobile string) (err error) {
	UserID, _ := primitive.ObjectIDFromHex(userID)
	col := DB.Collection(UserCollection)
	_, err = col.UpdateOne(
		a.Ctx,
		bson.M{"_id": bson.M{"$eq": UserID}},
		bson.M{"$set": bson.M{"mobile": mobile}},
	)
	return
}
