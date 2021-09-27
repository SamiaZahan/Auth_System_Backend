package repository

import (
	"context"
	"time"

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
		Created:       time.Now(),
	})
	ID = res.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (v *Verification) GetByIDAndCode(ID string, code int) (vDoc VerificationDoc, err error) {
	vID, _ := primitive.ObjectIDFromHex(ID)
	col := DB.Collection(VerificationCollection)
	err = col.
		FindOne(v.Ctx, VerificationDoc{ID: vID, Code: code}).
		Decode(vDoc)
	return
}
