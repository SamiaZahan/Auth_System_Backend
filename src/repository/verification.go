package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Verification struct {
	Ctx context.Context
}

func (v *Verification) Create(emailOrMobile string, code int, userId string) (ID string, err error) {
	UserID, _ := primitive.ObjectIDFromHex(userId)
	col := DB.Collection(VerificationCollection)
	res, err := col.InsertOne(v.Ctx, VerificationDoc{
		ID:            primitive.NewObjectID(),
		UserID:        UserID,
		EmailOrMobile: emailOrMobile,
		Code:          code,
	})
	ID = res.InsertedID.(primitive.ObjectID).Hex()
	return
}
