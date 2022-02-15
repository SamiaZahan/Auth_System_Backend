package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DB          *mongo.Database
	MongoClient *mongo.Client
)

const UserCollection = "user"
const UserProfileCollection = "user_profile"

const VerificationCollection = "verification"
const PhoneNumberMapCollection = "phone_number_map"

type UserDoc struct {
	ID             primitive.ObjectID `bson:"_id"`
	Email          string             `bson:"email"`
	Password       string             `bson:"password"`
	Mobile         string             `bson:"mobile"`
	Active         bool               `bson:"active"`
	EmailVerified  bool               `bson:"email_verified"`
	MobileVerified bool               `bson:"mobile_verified"`
	ExistingUser   bool               `bson:"existing_user"`
	Created        time.Time          `bson:"created"`
	Updated        time.Time          `bson:"updated"`
}

type UserProfileDoc struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        primitive.ObjectID `bson:"user_id"`
	FirstName     string             `bson:"first_name"`
	LastName      string             `bson:"last_name"`
	Gender        string             `bson:"gender"`
	ProfilePicURI string             `bson:"profile_pic_uri"`
	Created       time.Time          `bson:"created"`
	Updated       time.Time          `bson:"updated"`
}

type VerificationDoc struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        primitive.ObjectID `bson:"user_id"`
	EmailOrMobile string             `bson:"email_or_mobile"`
	Code          int                `bson:"code"`
	Created       time.Time          `bson:"created"`
	Updated       time.Time          `bson:"updated"`
}

type PhoneNumberMapDoc struct {
	CountryName   string `bson:"country_name"`
	CountryCode   string `bson:"country_code"`
	CountryPrefix string `bson:"country_prefix"`
	Active        bool   `bson:"active"`
}
