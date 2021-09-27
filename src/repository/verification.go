package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Verification struct {
	Ctx context.Context
}

func (v *Verification) Create(emailOrMobile string, code int, userId string) (err error) {
	UserID, _ := primitive.ObjectIDFromHex(userId)
	col := DB.Collection(VerificationCollection)
	_, err = col.InsertOne(v.Ctx, VerificationDoc{
		ID:            primitive.NewObjectID(),
		UserID:        UserID,
		EmailOrMobile: emailOrMobile,
		Code:          code,
	})
	return
}
