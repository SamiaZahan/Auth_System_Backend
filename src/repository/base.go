package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	DB *mongo.Database
)

const UserCollection = "user"
const UserProfileCollection = "user_profile"

type User struct {
	ID      primitive.ObjectID `bson:"_id"`
	Email   string             `bson:"email"`
	Mobile  string             `bson:"mobile"`
	Active  bool               `bson:"active"`
	Created time.Time          `bson:"created"`
}

type UserProfile struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserId        primitive.ObjectID `bson:"user_id"`
	FirstName     string             `bson:"first_name"`
	LastName      string             `bson:"last_name"`
	Gender        string             `bson:"gender"`
	ProfilePicURI string             `bson:"profile_pic_uri"`
	Created       time.Time          `bson:"created"`
}
