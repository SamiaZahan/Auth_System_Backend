package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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

// GetInfoByCountryPrefix Get Country Code
func (a *Auth) GetInfoByCountryPrefix(countryPrefix string) (phoneNumberMap *PhoneNumberMapDoc, err error) {
	col := DB.Collection(PhoneNumberMapCollection)
	err = col.FindOne(a.Ctx, bson.D{{Key: "country_prefix", Value: countryPrefix}}).Decode(&phoneNumberMap)
	return
}

// GetUserByEmailOrMobile get user  by email or mobile
func (a *Auth) GetUserByEmailOrMobile(emailOrMobile string) (user *UserDoc, err error) {
	col := DB.Collection(UserCollection)
	err = col.FindOne(
		a.Ctx,
		bson.D{{"$or", bson.A{
			bson.D{{Key: "email", Value: emailOrMobile}},
			bson.D{{Key: "mobile", Value: emailOrMobile}},
		},
		}},
	).Decode(&user)
	return
}

func (a *Auth) HashPassword(password string) (hashedPass string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return
	}
	return string(hash), nil
}

func (a *Auth) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Error(err)
		return false
	}

	return true
}

func (a *Auth) GetUserByEmail(email string) (user *UserDoc, err error) {
	col := DB.Collection(UserCollection)
	err = col.FindOne(a.Ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)
	return
}

func (a *Auth) GetUserByMobile(mobile string) (user *UserDoc, err error) {
	col := DB.Collection(UserCollection)
	err = col.FindOne(a.Ctx, bson.D{{Key: "mobile", Value: mobile}}).Decode(&user)
	return
}

func (a *Auth) CreateUser(email string, password string) (ID string, err error) {
	col := DB.Collection(UserCollection)
	res, err := col.InsertOne(a.Ctx, UserDoc{
		ID:       primitive.NewObjectID(),
		Email:    email,
		Password: password,
		Active:   false,
		Created:  time.Now(),
		Updated:  time.Now(),
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
		Updated:   time.Now(),
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

func (a *Auth) ActivateUserByEmail(email string) (err error) {
	col := DB.Collection(UserCollection)
	_, err = col.UpdateOne(
		a.Ctx,
		bson.M{"email": bson.M{"$eq": email}},
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

func (a *Auth) SetUserMobileByEmail(email string, mobile string) (err error) {
	col := DB.Collection(UserCollection)
	_, err = col.UpdateOne(
		a.Ctx,
		bson.M{"email": bson.M{"$eq": email}},
		bson.M{"$set": bson.M{"mobile": mobile}},
	)
	return
}
