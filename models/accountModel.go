package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// will be similar to order item
// so create a user and then an account

type Account struct {
	ID         primitive.ObjectID `bson:"_id"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	User_id    string             `json:"user_id" validate:"required"` // like food_id
	Account_id string             `json:"account_id"`
}
