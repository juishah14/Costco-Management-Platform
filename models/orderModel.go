package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID            primitive.ObjectID `bson:"_id"`
	Order_Date    time.Time          `json:"order_date"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	Order_id      string             `json:"order_id"`
	Membership_id *string            `json:"membership_id" validate:"required"` // used to be table_id
}
