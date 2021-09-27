package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	DB          *mongo.Database
	MongoClient *mongo.Client
)

const UserCollection = "user"
const UserProfileCollection = "user_profile"
const VerificationCollection = "verification"

type UserDoc struct {
	ID      primitive.ObjectID `bson:"_id"`
	Email   string             `bson:"email"`
	Mobile  string             `bson:"mobile"`
	Active  bool               `bson:"active"`
	Created time.Time          `bson:"created"`
}

type UserProfileDoc struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        primitive.ObjectID `bson:"user_id"`
	FirstName     string             `bson:"first_name"`
	LastName      string             `bson:"last_name"`
	Gender        string             `bson:"gender"`
	ProfilePicURI string             `bson:"profile_pic_uri"`
	Created       time.Time          `bson:"created"`
}

type VerificationDoc struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        primitive.ObjectID `bson:"user_id"`
	EmailOrMobile string             `bson:"email_or_mobile"`
	Code          int                `bson:"code"`
}
