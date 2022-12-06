package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// used to be table

type Membership struct {
	ID              primitive.ObjectID `bson:"_id"`
	Membership_type string             `json:"membership_type" validate:"eq=Executive|eq=Business|eq=Gold Star"`
	Created_at      time.Time          `json:"created_at"`
	Updated_at      time.Time          `json:"updated_at"`
	Membership_id   string             `json:"membership_id"`
}
