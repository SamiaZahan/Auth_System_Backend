package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

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
		Updated:       time.Now(),
	})
	ID = res.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (v *Verification) GetByID(ID string) (vDoc VerificationDoc, err error) {
	vID, _ := primitive.ObjectIDFromHex(ID)
	col := DB.Collection(VerificationCollection)
	err = col.
		FindOne(
			v.Ctx,
			bson.M{"_id": bson.M{"$eq": vID}}).
		Decode(&vDoc)
	return
}

func (v *Verification) GetByIDAndCode(ID string, code int) (vDoc VerificationDoc, err error) {
	vID, _ := primitive.ObjectIDFromHex(ID)
	col := DB.Collection(VerificationCollection)
	err = col.
		FindOne(
			v.Ctx,
			bson.M{
				"_id":  bson.M{"$eq": vID},
				"code": bson.M{"$eq": code},
			}).
		Decode(&vDoc)
	return
}

func (v *Verification) GetByEmailOrMobileAndCode(emailOrMobile string, code int) (vDoc VerificationDoc, err error) {
	col := DB.Collection(VerificationCollection)
	err = col.
		FindOne(
			v.Ctx,
			bson.M{
				"email_or_mobile": bson.M{"$eq": emailOrMobile},
				"code":            bson.M{"$eq": code},
			}).
		Decode(&vDoc)
	return
}

func (v *Verification) DeleteByID(ID string) (err error) {
	vID, _ := primitive.ObjectIDFromHex(ID)
	col := DB.Collection(VerificationCollection)
	_, err = col.DeleteOne(v.Ctx, bson.M{"_id": bson.M{"$eq": vID}})
	return
}
